package dto

import (
	"redirectServer/model"

	"github.com/lib/pq"
)

type FingerprintIOS struct {
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
	Renderer       string
	VendorRender   *string
	TimeZone       string  `json:"timeZone"`
	UniversalLink  *string `json:"universalLink"`
}

func (p *FingerprintIOS) ToFingerprint(ip string, userAgent *string) *model.Fingerprint {
	return &model.Fingerprint{
		IP:             ip,
		UserAgent:      userAgent,
		Platform:       p.Platform,
		Version:        p.Version,
		Language:       p.Language,
		Languages:      p.Languages,
		Cores:          p.Cores,
		Memory:         p.Memory,
		ScreenWidth:    p.ScreenWidth,
		ScreenHeight:   p.ScreenHeight,
		ColorDepth:     p.ColorDepth,
		PixelRatio:     p.PixelRatio,
		ViewportWidth:  p.ViewportWidth,
		ViewportHeight: p.ViewportHeight,
		Renderer:       p.Renderer,
		VendorRender:   p.VendorRender,
		TimeZone:       p.TimeZone,
	}
}
