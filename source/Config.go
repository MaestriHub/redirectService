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

// TODO: передать коннект
func StartApp() {
	godotenv.Load()
	DB = initDB()
	migrate(DB)
	routers.InitRouters(DB)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	err := http.ListenAndServeTLS(":8081", "server.crt", "server.key", nil)
	log.Println("Сервер запущен на http://localhost:8081")
	//TODO: fmt заменить на логи slog
	if err != nil {
		log.Fatal(err)
	}
}

func migrate(db *gorm.DB) {
	//db.Migrator().DropTable(&models.DirectLink{}, &models.Fingerprint{})
	if db.AutoMigrate(&models.DirectLink{}, &models.Fingerprint{}) != nil {
		log.Fatal("Failed to migrate database")
	}

	// parsedUUID, _ := uuid.Parse("c6acaff7-29a5-4c60-b8b8-4be02503bd8b")
	// payload1 := payload.Salon{
	// 	ID: parsedUUID,
	// }
	// var directLink = models.DirectLink{
	// 	ID:    "YSg6UgcF",
	// 	Event: string(models.SalonInvite),
	// }
	// directLink.SetPayload(payload1)
	// db.Create(&directLink)

	// parsedUUID, _ = uuid.Parse("53bb0f86-a94e-4302-8a07-ea0b083d3bde")
	// payload2 := payload.MasterToSalon{
	// 	EmployeeId: parsedUUID,
	// }
	// directLink = models.DirectLink{
	// 	ID:    "YSg6Ugcf",
	// 	Event: string(models.MasterInviteToSalon),
	// }
	// directLink.SetPayload(payload2)
	// db.Create(&directLink)

	// parsedUUID, _ = uuid.Parse("c6acaff7-29a5-4c60-b8b8-4be02503bd8b")
	// payload3 := payload.Customer{
	// 	ID: parsedUUID,
	// }
	// directLink = models.DirectLink{
	// 	ID:    "YSg6Ugc",
	// 	Event: string(models.CustomerInvite),
	// }
	// directLink.SetPayload(payload3)
	// db.Create(&directLink)

	// parsedUUID, _ = uuid.Parse("c6acaff7-29a5-4c60-b8b8-4be02503bd8b")
	// parsedUUID2, _ := uuid.Parse("53bb0f86-a94e-4302-8a07-ea0b083d3bde")
	// payload4 := payload.Employeer{
	// 	ID:      parsedUUID2,
	// 	SalonId: parsedUUID,
	// }
	// directLink = models.DirectLink{
	// 	ID:    "YSg6UgC",
	// 	Event: string(models.EmployeerInvite),
	// }
	// directLink.SetPayload(payload4)
	// db.Create(&directLink)

}

func initDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}
	dbHost := os.Getenv("DATABASE_HOST")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbPort := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("DATABASE_NAME")
	dbUser := os.Getenv("DATABASE_NAME")

	dataConnect := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=GMT"
	//TODO: подумать над env мб у go вообще другие приколы
	db, err := gorm.Open(postgres.Open(dataConnect), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных: ", err)
	}

	return db
}
