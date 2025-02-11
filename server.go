package main

import (
	"redirectServer/source"

	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	source.StartApp()
}
