package source

import (
	"log"
	"net/http"
	"os"
	"redirectServer/models"
	"redirectServer/routers"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func StartApp() {
	godotenv.Load()
	DB = initDB()
	migrate(DB)
	routers.InitRouters(DB)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	err := http.ListenAndServeTLS(":8081", "server.crt", "server.key", nil)
	log.Println("Сервер запущен на http://localhost:8081")

	if err != nil {
		log.Fatal(err)
	}
}

func migrate(db *gorm.DB) {
	db.Migrator().DropTable(&models.DirectURL{}, &models.Requester{}, &models.HistoryRequester{})
	if db.AutoMigrate(&models.DirectURL{}, &models.Requester{}, &models.HistoryRequester{}) != nil {
		log.Fatal("Failed to migrate database")
	}
	var directURL = models.DirectURL{
		ID:       "YSg6UgcF",
		Payload:  "c6acaff7-29a5-4c60-b8b8-4be02503bd8b",
		URLEvent: "SalonInvite",
	}

	db.Create(&directURL)
	directURL = models.DirectURL{
		ID:       "YSg6Ugcf",
		Payload:  "53bb0f86-a94e-4302-8a07-ea0b083d3bde",
		URLEvent: "EmployeerInvite",
	}

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
