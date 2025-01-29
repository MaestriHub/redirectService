package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Requester struct {
	gorm.Model
	IP                     string
	UserAgent              string
	Platform               string
	Version                string
	Language               string
	Languages              pq.StringArray `gorm:"type:text[]" `
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
