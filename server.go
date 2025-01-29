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
	requesterInfo.IP = r.RemoteAddr
	if err := db.Create(&requesterInfo).Error; err != nil {
		http.Error(w, "Failed to create requester info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Printf("Получены данные клиента: %+v\n", data)
	fmt.Fprintf(w, "Данные получены")
}
