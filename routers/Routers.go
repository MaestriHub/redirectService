package routers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"redirectServer/clientData"
	"redirectServer/models"
	"strings"

	"io"

	"github.com/lib/pq"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitRouters(db *gorm.DB) {
	DB = db
	http.HandleFunc("/", ServeHTML)
	http.HandleFunc("/PC", CollectDataPC)
	http.HandleFunc("/Mobile", CollectDataMobile)
	http.HandleFunc("/createDirectURL", CreateDirectURL)
	http.HandleFunc("/findRequester", FindRequester)

}
func CreateDirectURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var particalDirectURL models.ParticalDirectURL
	err = json.Unmarshal(body, &particalDirectURL)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := gonanoid.New(8)
	if err != nil {
		log.Fatal(err)
	}
	//TODO: сделать валидатор на URLEvent
	directURL := models.DirectURL{ID: id, Payload: particalDirectURL.Payload, URLEvent: particalDirectURL.URLEvent}

	if err := DB.Create(&directURL).Error; err != nil {
		http.Error(w, "Failed to create direct URL", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "URL created successfully")
}

func FindRequester(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var inputRequester models.ParticalRequester
	err = json.Unmarshal(body, &inputRequester)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	var query *gorm.DB

	if inputRequester.Memory != nil {
		query = query.Where("memory = ?", *inputRequester.Memory)
	} else {
		query = query.Where("memory IS NULL")
	}

	if inputRequester.Cores != nil {
		query = query.Where("cores = ?", *inputRequester.Cores)
	} else {
		query = query.Where("cores IS NULL")
	}

	if inputRequester.VendorRender != nil {
		query = query.Where("vendor_render = ?", *inputRequester.VendorRender)
	} else {
		query = query.Where("vendor_render IS NULL")
	}
	conditions := map[string]interface{}{
		"ip":              strings.Split(r.RemoteAddr, ":")[0],
		"platform":        inputRequester.Platform,
		"version":         inputRequester.Version,
		"language":        inputRequester.Language,
		"languages":       pq.StringArray(inputRequester.Languages),
		"screen_width":    inputRequester.ScreenWidth,
		"screen_height":   inputRequester.ScreenHeight,
		"color_depth":     inputRequester.ColorDepth,
		"pixel_ratio":     inputRequester.PixelRatio,
		"viewport_width":  inputRequester.ViewportWidth,
		"viewport_height": inputRequester.ViewportHeight,
		"renderer":        inputRequester.Renderer,
		"time_zone":       inputRequester.TimeZone,
	}
	var foundRequester *models.Requester
	var foundURL *models.DirectURL
	res := query.Where(conditions).Order("createdAt desc").First(&foundRequester)
	if res.Error != nil {
		var foundHistory models.HistoryRequester
		res := DB.Where("requester_id = ?", foundRequester.ID).Order("createdAt desc").First(&foundHistory)
		if res.Error == nil {
			DB.Where("id = ?", foundHistory.DirectURLID).First(&foundURL)
		}
	}
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if isOrganic(foundRequester, inputRequester, inputRequester.UniversalLink, ip, DB) {
		w.WriteHeader(http.StatusOK)
		return
	}

	if isFound(foundRequester, foundURL, inputRequester.UniversalLink, DB) {
		w.WriteHeader(http.StatusOK)
		return
	}
	if isNotFound(foundRequester, foundURL, inputRequester, inputRequester.UniversalLink, ip, DB) {
		w.WriteHeader(http.StatusOK)
		return
	}
	if isFoundUncorrect(foundRequester, foundURL, inputRequester.UniversalLink, DB) {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func isOrganic(
	foundRequester *models.Requester,
	inputRequester models.ParticalRequester,
	inputURL *string,
	inputIP string,
	db *gorm.DB,
) bool {
	if foundRequester == nil && inputURL == nil {
		fullRequester := inputRequester.ToRequester(inputIP, nil, []string{string(models.Organic)})
		db.Create(&fullRequester)
		return true
	}
	return false
}

func isFound(
	foundRequester *models.Requester,
	foundURL *models.DirectURL,
	inputURL *string,
	db *gorm.DB,
) bool {
	if foundRequester != nil && foundURL != nil && foundURL.ParseToURL() == *inputURL {
		foundRequester.Statuses = append(foundRequester.Statuses, string(models.Found))
		db.Save(&foundRequester)
		return true
	}
	return false
}

func isNotFound(
	foundRequester *models.Requester,
	foundURL *models.DirectURL,
	inputRequester models.ParticalRequester,
	inputURL *string,
	inputIP string,
	db *gorm.DB,
) bool {
	if foundRequester == nil && foundURL == nil && inputURL != nil {
		fullRequester := inputRequester.ToRequester(inputIP, nil, []string{string(models.NotFound)})
		db.Create(&fullRequester)
		return true
	}
	return false
}

func isFoundUncorrect(
	foundRequester *models.Requester,
	foundURL *models.DirectURL,
	inputURL *string,
	db *gorm.DB,
) bool {
	if foundRequester != nil && foundURL != nil && foundURL.ParseToURL() != *inputURL {
		foundRequester.Statuses = append(foundRequester.Statuses, string(models.FoundUncorrect))
		db.Save(&foundRequester)
		return true
	}
	return false
}

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing query parameter 'code'", http.StatusBadRequest)
		return
	}
	var directUrl models.DirectURL
	DB.First(&directUrl, "id = ?", code)
	directUrl.Сlicks++
	DB.Save(&directUrl)

	userAgent := r.Header.Get("User-Agent")

	if strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "Android") || strings.Contains(userAgent, "iPad") {
		htmlFile, err := os.ReadFile("static/mobileScreen.html")
		if err != nil {
			http.Error(w, "Error reading HTML file", http.StatusInternalServerError)
			return
		}
		assembledPayload, err := directUrl.GetPayload(DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directUrl.ParseToURL(), -1)
		modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", assembledPayload.Name, -1)
		modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Title}}", assembledPayload.Event, -1)
		modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", assembledPayload.Description, -1)
		modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", code, -1)

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(modifiedHTML))
	} else {
		http.ServeFile(w, r, "static/webScreen.html")
	}

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

	//w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Printf("Получены данные клиента: %+v\n", data)
	fmt.Fprintf(w, "Данные полученыфы")
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
	var requesterInfo models.Requester = data.ToRequester()
	ipParts := strings.Split(r.RemoteAddr, ":")
	if len(ipParts) == 2 {
		requesterInfo.IP = ipParts[0]
		//equesterInfo.Port = ipParts[1]
	} else {
		fmt.Println("Некорректный формат строки")
	}

	conditions := map[string]interface{}{
		"ip":         requesterInfo.IP,
		"user_agent": requesterInfo.UserAgent,
		"platform":   requesterInfo.Platform,
		"version":    requesterInfo.Version,
		"language":   requesterInfo.Language,
		"languages":  pq.StringArray(requesterInfo.Languages),
		"cores":      requesterInfo.Cores,
		//"memory":          requesterInfo.Memory,
		"screen_width":    requesterInfo.ScreenWidth,
		"screen_height":   requesterInfo.ScreenHeight,
		"color_depth":     requesterInfo.ColorDepth,
		"pixel_ratio":     requesterInfo.PixelRatio,
		"viewport_width":  requesterInfo.ViewportWidth,
		"viewport_height": requesterInfo.ViewportHeight,
		"renderer":        requesterInfo.Renderer,
		"vendor_render":   requesterInfo.VendorRender,
		"time_zone":       requesterInfo.TimeZone,
	}
	query := DB.Where(conditions)
	if requesterInfo.Memory != nil {
		query = query.Where("memory = ?", *requesterInfo.Memory)
	} else {
		query = query.Where("memory IS NULL")
	}
	var existingRequester models.Requester

	if query.First(&existingRequester).Error != nil {
		requesterInfo.Statuses = pq.StringArray([]string{string(models.Linked)})
		if err := DB.Create(&requesterInfo).Error; err != nil {
			http.Error(w, "Failed to create requester info", http.StatusInternalServerError)
			return
		}
		DB.Create(
			&models.HistoryRequester{
				RequesterID: requesterInfo.ID,
				Port:        ipParts[1],
				DirectURLID: data.DirectURLID,
			},
		)
	} else {
		DB.Create(
			&models.HistoryRequester{
				RequesterID: existingRequester.ID,
				Port:        ipParts[1],
				DirectURLID: data.DirectURLID,
			},
		)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Printf("Получены данные клиента: %+v\n", data)
	fmt.Fprintf(w, "Данные получены")
}
