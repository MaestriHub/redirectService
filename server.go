package main

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

	"github.com/joho/godotenv"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	godotenv.Load()
	db = initDB()
	migrate(db)
	routers()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	err := http.ListenAndServeTLS(":8081", "server.crt", "server.key", nil)
	log.Println("Сервер запущен на http://localhost:8081")

	if err != nil {
		log.Fatal(err)
	}
}

func migrate(db *gorm.DB) {
	db.Migrator().DropTable(&models.DirectURL{}, &models.Requester{})
	if db.AutoMigrate(&models.DirectURL{}, &models.Requester{}) != nil {
		log.Fatal("Failed to migrate database")
	}
	directURL := models.DirectURL{ID: "YSg6UgcF", URL: "tg://resolve?domain=vitalik_shevtsov&text=ass&profile"}
	db.Create(&directURL)

}

func initDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_CONNECT")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных: ", err)
	}

	return db
}

func routers() {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/PC", collectDataPC)
	http.HandleFunc("/Mobile", collectDataMobile)
	http.HandleFunc("/createDirectURL", createDirectURL)
	http.HandleFunc("/findRequester", findRequester)

}

func findRequester(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var requester models.ParticalRequester
	err = json.Unmarshal(body, &requester)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	var query *gorm.DB
	var response models.DirectURL
	if requester.UniversalLink != nil {
		res := db.Where(&models.DirectURL{URL: *requester.UniversalLink}).Order("createdAt desc").First(&response)
		if res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusNotFound)
			return
		}
		query.Where("direct_url_id = ?", response.ID)
	}

	if requester.Memory != nil {
		query.Where("memory = ?", *requester.Memory)
	}

	if requester.Cores != nil {
		query.Where("cores = ?", *requester.Cores)
	}
	ip := strings.Split(r.RemoteAddr, ":")[0]
	query.Where(`ip = ? 
		AND platform = ? 
		AND version = ? 
		AND language = ? 
		AND languages = ?
		AND screenWidth = ?
		AND screenHeight = ?
		AND colorDepth = ?
		AND pixelRatio = ?
		AND viewportWidth = ?
		AND viewportHeight = ?
		AND timeZone = ?`,
		ip,
		requester.Platform,
		requester.Version,
		requester.Language,
		requester.Languages,
		requester.ScreenWidth,
		requester.ScreenHeight,
		requester.ColorDepth,
		requester.PixelRatio,
		requester.ViewportWidth,
		requester.ViewportHeight,
		requester.TimeZone)
	var installedRequester models.Requester
	res := query.Order("createdAt desc").First(&installedRequester)
	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusNotFound)
		return
	}

	installedRequester.IsInstalled = true
	db.Save(&installedRequester)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if response == (models.DirectURL{}) {
		res := db.Where(&models.DirectURL{URL: *requester.UniversalLink}).Order("createdAt desc").First(&response)
		if res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusNotFound)
			return
		}
	}

	json.NewEncoder(w).Encode(response.URL)
}

func createDirectURL(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query().Get("universalLink")
	if queryParam == "" {
		http.Error(w, "Missing query parameter 'universalLink'", http.StatusBadRequest)
		return
	}
	id, err := gonanoid.New(8)
	if err != nil {
		log.Fatal(err)
	}
	directURL := models.DirectURL{ID: id, URL: queryParam}

	if err := db.Create(&directURL).Error; err != nil {
		http.Error(w, "Failed to create direct URL", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "URL created successfully: %s", directURL.URL)
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing query parameter 'code'", http.StatusBadRequest)
		return
	}
	var directUrl models.DirectURL
	db.First(&directUrl, "id = ?", code)
	directUrl.Сlicks++
	db.Save(&directUrl)

	userAgent := r.Header.Get("User-Agent")

	if strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "Android") || strings.Contains(userAgent, "iPad") {
		htmlFile, err := os.ReadFile("static/mobileScreen.html")
		if err != nil {
			http.Error(w, "Error reading HTML file", http.StatusInternalServerError)
			return
		}

		modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directUrl.URL, -1)
		modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", code, -1)

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(modifiedHTML))
	} else {
		http.ServeFile(w, r, "static/webScreen.html")
	}

}

func collectDataPC(w http.ResponseWriter, r *http.Request) {
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

func collectDataMobile(w http.ResponseWriter, r *http.Request) {
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
		requesterInfo.Port = ipParts[1]
	} else {
		fmt.Println("Некорректный формат строки")
	}

	if err := db.Create(&requesterInfo).Error; err != nil {
		http.Error(w, "Failed to create requester info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Printf("Получены данные клиента: %+v\n", data)
	fmt.Fprintf(w, "Данные получены")
}
