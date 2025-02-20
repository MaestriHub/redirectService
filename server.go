package main

import (
	"os"
	"redirectServer/source"
)

func main() {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbPort := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("DATABASE_NAME")
	dbUser := os.Getenv("DATABASE_NAME")

	dataConnect := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=GMT"
	source.StartApp(dataConnect)
}
