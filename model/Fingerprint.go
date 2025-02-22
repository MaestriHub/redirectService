package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// TODO: кноака отправить ссыфлку.
// погуглить насчет action действия, если нет то просто скопировать в будер обмена
type Fingerprint struct {
	gorm.Model
	IP             string
	UserAgent      *string
	Platform       string
	Version        string
	Language       string
	Languages      pq.StringArray `gorm:"type:text[]" `
	Cores          *int           `gorm:"default:null"`
	Memory         *int           `gorm:"default:null"`
	ScreenWidth    int
	ScreenHeight   int
	ColorDepth     int
	PixelRatio     float64
	ViewportWidth  int
	ViewportHeight int
	Renderer       string
	VendorRender   *string
	TimeZone       string
	DirectLinkID   string
}
