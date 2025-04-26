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
	Language       string          `json:"language"`
	Languages      pkg.StringArray `json:"languages"`
	Cores          *int            `json:"cores"`
	Memory         *int            `json:"memory"`
	ScreenWidth    int             `json:"screenWidth"`
	ScreenHeight   int             `json:"screenHeight"`
	ColorDepth     int             `json:"colorDepth"`
	PixelRatio     float64         `json:"pixelRatio"`
	ViewportWidth  int             `json:"viewportWidth"`
	ViewportHeight int             `json:"viewportHeight"`
	TimeZone       string          `json:"timeZone"`
}

func (fp *Fingerprint) Validate() error {
	return nil
}
