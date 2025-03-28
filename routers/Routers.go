package routers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"redirectServer/clientData"
	"redirectServer/models"
	"redirectServer/models/payload"
	"redirectServer/services"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

var DB *gorm.DB
var logger *slog.Logger

func InitRouters(db *gorm.DB) {
	DB = db
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Routers initialized")
	http.HandleFunc("/", ServeHTML)
	http.HandleFunc("/collect/pc", CollectDataPC)
	http.HandleFunc("/collect/mobile", CollectDataMobile)
	http.HandleFunc("/create/salon", CreateSalonInvite)
	http.HandleFunc("/create/employeer", CreateEmployeerInvite)
	http.HandleFunc("/create/customer", CreateCustomerInvite)
	http.HandleFunc("/create/master-to-salon", CreateMasterToSalonInvite)
	// TODO: хуевый нейминг, один с фингерпринтом, другой без
	http.HandleFunc("/find/without-link", FindFingerprint)
	http.HandleFunc("/find/with-link", GetDirectLinkPayload)

}

func CreateEmployeerInvite(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payloadObject payload.Employeer
	err = json.Unmarshal(body, &payloadObject)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := gonanoid.New(8)
	if err != nil {
		http.Error(w, "Failed to create nanoId", http.StatusInternalServerError)
		return
	}
	directLink := models.DirectLink{ID: id, Event: string(models.EmployeerInvite)}
	directLink.SetPayload(payloadObject)

	if err := DB.Create(&directLink).Error; err != nil {
		http.Error(w, "Failed to create direct URL", http.StatusInternalServerError)
		return
	}
	//TODO: подумать про fmt.Fprintf
	logger.Info("URL created successfully")
}

func CreateMasterToSalonInvite(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payloadObject payload.MasterToSalon
	err = json.Unmarshal(body, &payloadObject)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := gonanoid.New(8)
	if err != nil {
		http.Error(w, "Failed to create nanoId", http.StatusInternalServerError)
		return
	}
	directLink := models.DirectLink{ID: id, Event: string(models.MasterInviteToSalon)}
	directLink.SetPayload(payloadObject)

	if err := DB.Create(&directLink).Error; err != nil {
		http.Error(w, "Failed to create direct URL", http.StatusInternalServerError)
		return
	}
	logger.Info("URL created successfully")
}

func CreateSalonInvite(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payloadObject payload.Salon
	err = json.Unmarshal(body, &payloadObject)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := gonanoid.New(8)
	if err != nil {
		http.Error(w, "Failed to create nanoId", http.StatusInternalServerError)
		return
	}
	directLink := models.DirectLink{ID: id, Event: string(models.SalonInvite)}
	directLink.SetPayload(payloadObject)

	if err := DB.Create(&directLink).Error; err != nil {
		http.Error(w, "Failed to create direct URL", http.StatusInternalServerError)
		return
	}
	logger.Info("URL created successfully")
}

func CreateCustomerInvite(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payloadObject payload.Customer
	err = json.Unmarshal(body, &payloadObject)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := gonanoid.New(8)
	if err != nil {
		http.Error(w, "Failed to create nanoId", http.StatusInternalServerError)
		return
	}
	directLink := models.DirectLink{ID: id, Event: string(models.CustomerInvite)}
	directLink.SetPayload(payloadObject)

	if err := DB.Create(&directLink).Error; err != nil {
		http.Error(w, "Failed to create direct URL", http.StatusInternalServerError)
		return
	}
	logger.Info("URL created successfully")
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

func getMasterToSalonInviteIOS(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/MasterToSalon/IOS/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadMasterToSalon(DB)
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", name, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", description, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getMasterToSalonInviteAndroid(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/MasterToSalon/Android/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadMasterToSalon(DB)
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", name, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", description, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getMasterToSalonInviteWeb(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/MasterToSalon/Web/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadMasterToSalon(DB)
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

func getSalonInviteWeb(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/SalonInvite/Web/Screen.html")
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

func getCustomerInviteWeb(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/CustomerInvite/Web/Screen.html")
	if err != nil {
		return "", err
	}

	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getEmployeerInviteIOS(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/EmployeerInvite/IOS/Screen.html")
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getEmployeerInviteAndroid(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/EmployeerInvite/Android/Screen.html")
	if err != nil {
		return "", err
	}

	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getEmployeerInviteWeb(
	directLink models.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/EmployeerInvite/Web/Screen.html")
	if err != nil {
		return "", err
	}

	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" {
		http.Error(w, "Missing query parameter 'code'", http.StatusBadRequest)
		return
	}
	var directLink models.DirectLink
	result := DB.First(&directLink, "id = ?", code)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	directLink.Сlicks++
	DB.Save(&directLink)

	userAgent := r.Header.Get("User-Agent")
	var response string
	var err error
	switch {
	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(models.MasterInviteToSalon):
		response, err = getMasterToSalonInviteIOS(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case strings.Contains(userAgent, "Android") && directLink.Event == string(models.MasterInviteToSalon):
		response, err = getMasterToSalonInviteAndroid(directLink)
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
	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(models.EmployeerInvite):
		response, err = getEmployeerInviteIOS(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case strings.Contains(userAgent, "Android") && directLink.Event == string(models.EmployeerInvite):
		response, err = getEmployeerInviteAndroid(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		if directLink.Event == string(models.MasterInviteToSalon) {
			response, err = getMasterToSalonInviteWeb(directLink)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if directLink.Event == string(models.SalonInvite) {
			response, err = getSalonInviteWeb(directLink)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if directLink.Event == string(models.CustomerInvite) {
			response, err = getCustomerInviteWeb(directLink)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if directLink.Event == string(models.EmployeerInvite) {
			response, err = getEmployeerInviteWeb(directLink)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Unsupported device", http.StatusInternalServerError)
		}
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(response))
}

func CollectDataPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST-запросы поддерживаются", http.StatusMethodNotAllowed)
		return
	}

	var data clientData.PC
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Ошибка декодирования данных", http.StatusBadRequest)
		return
	}
	var fingerprintInfo models.Fingerprint = data.ToFingerprint()
	ipParts := strings.Split(r.RemoteAddr, ":")
	if len(ipParts) == 2 {
		fingerprintInfo.IP = ipParts[0]
	} else {
		logger.Info("Некорректный формат строки")
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
	logger.Info("Данные получены")
}

func CollectDataMobile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST-запросы поддерживаются", http.StatusMethodNotAllowed)
		return
	}

	var data clientData.Mobile
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
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
	logger.Info("Данные получены")
}
