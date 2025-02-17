package routers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"redirectServer/clientData"
	"redirectServer/models"
	"redirectServer/services"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitRouters(db *gorm.DB) {
	DB = db

	http.HandleFunc("/", ServeHTML)
	http.HandleFunc("/collect/pc", CollectDataPC)
	http.HandleFunc("/collect/mobile", CollectDataMobile)
	http.HandleFunc("/create/salon", CreateSalonInvite)
	http.HandleFunc("/create/employeer", CreateEmployeerInvite)
	http.HandleFunc("/create/customer", CreateCustomerInvite)
	http.HandleFunc("/find/without-link", FindFingerprint)
	http.HandleFunc("/find/with-link", GetDirectLinkPayload)

}

// TODO:сделать три метода на каждый ивент
func CreateEmployeerInvite(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var particalDirectLink models.ParticalDirectLink
	err = json.Unmarshal(body, &particalDirectLink)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := gonanoid.New(8)
	if err != nil {
		log.Fatal(err)
	}
	directLink := models.DirectLink{ID: id, Payload: particalDirectLink.Payload, Event: string(models.EmployeerInvite)}

	if err := DB.Create(&directLink).Error; err != nil {
		http.Error(w, "Failed to create direct URL", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "URL created successfully")
}

func CreateSalonInvite(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var particalDirectLink models.ParticalDirectLink
	err = json.Unmarshal(body, &particalDirectLink)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := gonanoid.New(8)
	if err != nil {
		log.Fatal(err)
	}
	directLink := models.DirectLink{ID: id, Payload: particalDirectLink.Payload, Event: string(models.SalonInvite)}

	if err := DB.Create(&directLink).Error; err != nil {
		http.Error(w, "Failed to create direct URL", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "URL created successfully")
}

func CreateCustomerInvite(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var particalDirectLink models.ParticalDirectLink
	err = json.Unmarshal(body, &particalDirectLink)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := gonanoid.New(8)
	if err != nil {
		log.Fatal(err)
	}
	directLink := models.DirectLink{ID: id, Payload: particalDirectLink.Payload, Event: string(models.CustomerInvite)}

	if err := DB.Create(&directLink).Error; err != nil {
		http.Error(w, "Failed to create direct URL", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "URL created successfully")
}

func FindFingerprint(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var inputFingerprint models.ParticalFingerprint
	err = json.Unmarshal(body, &inputFingerprint)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	ip := strings.Split(r.RemoteAddr, ":")[0]
	fingerprint := inputFingerprint.ToFingerprint(ip, nil)
	findFingerprint := services.FindFingerprint(*fingerprint, DB)
	if findFingerprint == nil && inputFingerprint.UniversalLink == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if findFingerprint != nil && inputFingerprint.UniversalLink == nil {
		var DirectLink models.DirectLink
		DB.First(&DirectLink, "id = ?", findFingerprint.DirectLinkID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DirectLink)
		return
	}
	if inputFingerprint.UniversalLink != nil {
		DirectLink, err := models.ParseURL(*inputFingerprint.UniversalLink, DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DirectLink)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
func GetDirectLinkPayload(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing query parameter 'code'", http.StatusBadRequest)
		return
	}
	var directLink models.DirectLink
	if DB.First(&directLink, "id = ?", code).Error != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(directLink)

}

func getEmployeeInviteIOS(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/EmployeeInvite/IOS/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadEmployee(DB)
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", name, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", description, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getEmployeeInviteAndroid(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/EmployeeInvite/Android/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadEmployee(DB)
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", name, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", description, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getSalonInviteIOS(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/SalonInvite/IOS/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadSalon(DB)
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", name, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", description, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getSalonInviteAndroid(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/SalonInvite/Android/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadSalon(DB)
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", name, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", description, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getCustomerInviteIOS(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/CustomerInvite/IOS/Screen.html")
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getCustomerInviteAndroid(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/CustomerInvite/Android/Screen.html")
	if err != nil {
		return "", err
	}

	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getWebInvite(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/webScreen.html")
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing query parameter 'code'", http.StatusBadRequest)
		return
	}
	var directLink models.DirectLink
	DB.First(&directLink, "id = ?", code)
	directLink.Сlicks++
	DB.Save(&directLink)

	userAgent := r.Header.Get("User-Agent")
	var response string
	var err error
	switch {
	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(models.EmployeerInvite):
		response, err = getEmployeeInviteIOS(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case strings.Contains(userAgent, "Android") && directLink.Event == string(models.EmployeerInvite):
		response, err = getEmployeeInviteAndroid(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(models.SalonInvite):
		response, err = getSalonInviteIOS(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case strings.Contains(userAgent, "Android") && directLink.Event == string(models.SalonInvite):
		response, err = getSalonInviteAndroid(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(models.CustomerInvite):
		response, err = getCustomerInviteIOS(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case strings.Contains(userAgent, "Android") && directLink.Event == string(models.CustomerInvite):
		response, err = getCustomerInviteAndroid(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		response, err = getWebInvite(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(response))
}

func CollectDataPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Неверный метод:", r.Method)
		http.Error(w, "Только POST-запросы поддерживаются", http.StatusMethodNotAllowed)
		return
	}

	var data clientData.PC
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("Ошибка декодирования данных:", err)
		http.Error(w, "Ошибка декодирования данных", http.StatusBadRequest)
		return
	}
	var fingerprintInfo models.Fingerprint = data.ToFingerprint()
	ipParts := strings.Split(r.RemoteAddr, ":")
	if len(ipParts) == 2 {
		fingerprintInfo.IP = ipParts[0]
	} else {
		fmt.Println("Некорректный формат строки")
	}
	var existingFingerprint *models.Fingerprint = services.FindFingerprint(fingerprintInfo, DB)
	if existingFingerprint != nil {
		if err := DB.Create(&fingerprintInfo).Error; err != nil {
			http.Error(w, "Failed to create fingerprint info", http.StatusInternalServerError)
			return
		}
		DB.Delete(&existingFingerprint)
	} else {
		DB.Create(&fingerprintInfo)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Printf("Получены данные клиента: %+v\n", data)
	fmt.Fprintf(w, "Данные получены")
}

func CollectDataMobile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Неверный метод:", r.Method)
		http.Error(w, "Только POST-запросы поддерживаются", http.StatusMethodNotAllowed)
		return
	}

	var data clientData.Mobile
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("Ошибка декодирования данных:", err)
		http.Error(w, "Ошибка декодирования данных", http.StatusBadRequest)
		return
	}
	var fingerprintInfo models.Fingerprint = data.ToFingerprint()
	ipParts := strings.Split(r.RemoteAddr, ":")
	if len(ipParts) == 2 {
		fingerprintInfo.IP = ipParts[0]
	} else {
		fmt.Println("Некорректный формат строки")
	}
	var existingFingerprint *models.Fingerprint = services.FindFingerprint(fingerprintInfo, DB)
	if existingFingerprint != nil {
		if err := DB.Create(&fingerprintInfo).Error; err != nil {
			http.Error(w, "Failed to create fingerprint info", http.StatusInternalServerError)
			return
		}
		DB.Delete(&existingFingerprint)
	} else {
		DB.Create(&fingerprintInfo)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Printf("Получены данные клиента: %+v\n", data)
	fmt.Fprintf(w, "Данные получены")
}
