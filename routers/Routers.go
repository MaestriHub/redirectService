package routers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"redirectServer/dto"
	"redirectServer/model"
	"redirectServer/model/payload"
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
	http.HandleFunc("/collect", CollectData)
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
	directLink := model.DirectLink{ID: id, Event: string(model.EmployeerInvite)}
	directLink.SetPayload(payloadObject)

	if err := DB.Create(&directLink).Error; err != nil {
		http.Error(w, "Failed to create direct URL", http.StatusInternalServerError)
		return
	}
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
	directLink := model.DirectLink{ID: id, Event: string(model.MasterInviteToSalon)}
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
	directLink := model.DirectLink{ID: id, Event: string(model.SalonInvite)}
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
	directLink := model.DirectLink{ID: id, Event: string(model.CustomerInvite)}
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

	var inputFingerprint dto.FingerprintIOS
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
		var DirectLink model.DirectLink
		DB.First(&DirectLink, "id = ?", findFingerprint.DirectLinkID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DirectLink)
		return
	}
	if inputFingerprint.UniversalLink != nil {
		DirectLink, err := model.ParseURL(*inputFingerprint.UniversalLink, DB)
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
	var directLink model.DirectLink
	if DB.First(&directLink, "id = ?", code).Error != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(directLink)

}

func getMasterToSalonInviteIOS(
	directLink model.DirectLink,
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
	directLink model.DirectLink,
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
	directLink model.DirectLink,
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
	directLink model.DirectLink,
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
	directLink model.DirectLink,
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
	directLink model.DirectLink,
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
	directLink model.DirectLink,
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
	directLink model.DirectLink,
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
	directLink model.DirectLink,
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
	directLink model.DirectLink,
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
	directLink model.DirectLink,
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
	directLink model.DirectLink,
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
	var directLink model.DirectLink
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
	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(model.MasterInviteToSalon):
		response, err = getMasterToSalonInviteIOS(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case strings.Contains(userAgent, "Android") && directLink.Event == string(model.MasterInviteToSalon):
		response, err = getMasterToSalonInviteAndroid(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(model.SalonInvite):
		response, err = getSalonInviteIOS(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case strings.Contains(userAgent, "Android") && directLink.Event == string(model.SalonInvite):
		response, err = getSalonInviteAndroid(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(model.CustomerInvite):
		response, err = getCustomerInviteIOS(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case strings.Contains(userAgent, "Android") && directLink.Event == string(model.CustomerInvite):
		response, err = getCustomerInviteAndroid(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(model.EmployeerInvite):
		response, err = getEmployeerInviteIOS(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case strings.Contains(userAgent, "Android") && directLink.Event == string(model.EmployeerInvite):
		response, err = getEmployeerInviteAndroid(directLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		if directLink.Event == string(model.MasterInviteToSalon) {
			response, err = getMasterToSalonInviteWeb(directLink)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if directLink.Event == string(model.SalonInvite) {
			response, err = getSalonInviteWeb(directLink)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if directLink.Event == string(model.CustomerInvite) {
			response, err = getCustomerInviteWeb(directLink)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if directLink.Event == string(model.EmployeerInvite) {
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

func CollectData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST-запросы поддерживаются", http.StatusMethodNotAllowed)
		return
	}

	var data dto.FingerprintJS
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Ошибка декодирования данных", http.StatusBadRequest)
		return
	}
	var fingerprintInfo model.Fingerprint = data.ToFingerprint()
	ipParts := strings.Split(r.RemoteAddr, ":")
	if len(ipParts) == 2 {
		fingerprintInfo.IP = ipParts[0]
	} else {
		fmt.Println("Некорректный формат строки")
	}
	var existingFingerprint *model.Fingerprint = services.FindFingerprint(fingerprintInfo, DB)
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
