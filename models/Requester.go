package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// TODO: кноака отправить ссыфлку.
// погуглить насчет action действия, если нет то просто скопировать в будер обмена
type Fingerprint struct {
	gorm.Model
	//TODO: по-хорошему надо сделать связь в user так как может меняться ip,но пользуется один и тот же человек
	//UserID         uuid
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

type ParticalFingerprint struct {
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

func (p *ParticalFingerprint) ToFingerprint(ip string, userAgent *string) *Fingerprint {
	return &Fingerprint{
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
