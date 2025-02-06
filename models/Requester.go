package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Requester struct {
	gorm.Model
	IP             string
	Port           string
	UserAgent      string
	Platform       string
	Version        string
	Language       string
	Languages      pq.StringArray `gorm:"type:text[]" `
	Cores          int            `gorm:"default:null"`
	Memory         int            `gorm:"default:null"`
	ScreenWidth    int
	ScreenHeight   int
	ColorDepth     int
	PixelRatio     float64
	ViewportWidth  int
	ViewportHeight int
	Renderer       string
	VendorRender   string
	TimeZone       string
	DirectURLID    string
	IsInstalled    bool `gorm:"default:false"`
}

type ParticalRequester struct {
	Platform       string         `json:"platform"`
	Version        string         `json:"version"`
	Language       string         `json:"language"`
	Languages      pq.StringArray `json:"languages"`
	Cores          *int           `json:"cores"`
	Memory         *int           `json:"memory"`
	ScreenWidth    int            `json:"screenWidth"`
	ScreenHeight   int            `json:"screenHeight"`
	ColorDepth     int            `json:"colorDepth"`
	PixelRatio     float64        `json:"pixelRatio"`
	ViewportWidth  int            `json:"viewportWidth"`
	ViewportHeight int            `json:"viewportHeight"`
	TimeZone       string         `json:"timeZone"`
	UniversalLink  *string        `json:"universalLink"`
}
