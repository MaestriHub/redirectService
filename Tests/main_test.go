package test

import (
	"fmt"
	"log"
	"net/http"
	"redirectServer/models"
	"redirectServer/source"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func TestMain(m *testing.M) {
	fmt.Println("Setting up the database...")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	DB = initTestDB()
	migrateTest(DB)
	source.DB = DB
	m.Run()

}

func initTestDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=GMT"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных: ", err)
	}

	return db
}

func migrateTest(db *gorm.DB) {
	db.Migrator().DropTable(&models.DirectURL{}, &models.Requester{}, &models.HistoryRequester{})
	if db.AutoMigrate(&models.DirectURL{}, &models.Requester{}, &models.HistoryRequester{}) != nil {
		log.Fatal("Failed to migrate database")
	}

}
