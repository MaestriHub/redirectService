package models

import (
	"gorm.io/gorm"
)

type Requester struct {
	gorm.Model
	IP                     string
	UserAgent              string
	Platform               string
	Language               string
	Languages              []string `gorm:"type:json"`
	CookiesEnabled         bool
	ConnectionType         string
	IsOnline               bool
	Cores                  int `gorm:"default:null"`
	Memory                 int `gorm:"default:null"`
	ScreenWidth            int
	ScreenHeight           int
	ColorDepth             int
	PixelRatio             float64
	ViewportWidth          int
	ViewportHeight         int
	TimeZone               string
	CurrentTime            string
	BatteryLevel           float64 `gorm:"default:null"`
	BatteryCharging        bool    `gorm:"default:null"`
	BatteryChargingTime    float64 `gorm:"default:null"`
	BatteryDischargingTime float64 `gorm:"default:null"`
	DirectURLID            string
}
