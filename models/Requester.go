package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type RequesterStatus string

const (
	Linked         RequesterStatus = "linked"
	Found          RequesterStatus = "found"
	FoundUncorrect RequesterStatus = "found_uncorrect"
	NotFound       RequesterStatus = "not_found"
	Organic        RequesterStatus = "organic"
)

type Requester struct {
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
	Statuses       pq.StringArray `gorm:"type:text[]"`
}

type HistoryRequester struct {
	gorm.Model
	Port        string
	RequesterID uint
	DirectURLID string
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
	Renderer       string
	VendorRender   *string
	TimeZone       string  `json:"timeZone"`
	UniversalLink  *string `json:"universalLink"`
}

func (p *ParticalRequester) ToRequester(ip string, userAgent *string, statuses []string) *Requester {
	return &Requester{
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
		Statuses:       pq.StringArray(statuses),
	}
}
