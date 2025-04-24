package domain

import (
	"redirectServer/pkg"
)

type Fingerprint struct {
	FingerprintFields
	LinkId NanoID
}

type FingerprintFields struct {
	IP             string
	UserAgent      string
	Platform       string
	Version        string
	Language       string
	Languages      pkg.StringArray `gorm:"type:text[]" `
	Cores          *int            `gorm:"default:null"`
	Memory         *int            `gorm:"default:null"`
	ScreenWidth    int
	ScreenHeight   int
	ColorDepth     int
	PixelRatio     float64
	ViewportWidth  int
	ViewportHeight int
	Renderer       string
	VendorRender   *string
	TimeZone       string
}

func (fp *Fingerprint) Validate() error {
	return nil
}
