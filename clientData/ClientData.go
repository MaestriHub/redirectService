package clientData

import (
	"redirectServer/models"

	"github.com/lib/pq"
)

type PC struct {
	UserAgent      string   `json:"userAgent"`
	Platform       string   `json:"platform"`
	Language       string   `json:"language"`
	Languages      []string `json:"languages"`
	CookiesEnabled bool     `json:"cookiesEnabled"`
	ConnectionType string   `json:"connectionType"`
	IsOnline       bool     `json:"isOnline"`
	Cores          int      `json:"cores"`
	Memory         int      `json:"memory"`
	ScreenWidth    int      `json:"screenWidth"`
	ScreenHeight   int      `json:"screenHeight"`
	ColorDepth     int      `json:"colorDepth"`
	PixelRatio     float64  `json:"pixelRatio"`
	ViewportWidth  int      `json:"viewportWidth"`
	ViewportHeight int      `json:"viewportHeight"`
	TimeZone       string   `json:"timeZone"`
	CurrentTime    string   `json:"currentTime"`
	DirectLinkID   string   `json:"directLinkID"`
}

type Mobile struct {
	UserAgent      string   `json:"userAgent"`
	Platform       string   `json:"platform"`
	Version        string   `json:"version"`
	Language       string   `json:"language"`
	Languages      []string `json:"languages"`
	Cores          int      `json:"cores"`
	Memory         int      `json:"memory"`
	ScreenWidth    int      `json:"screenWidth"`
	ScreenHeight   int      `json:"screenHeight"`
	ColorDepth     int      `json:"colorDepth"`
	PixelRatio     float64  `json:"pixelRatio"`
	ViewportWidth  int      `json:"viewportWidth"`
	ViewportHeight int      `json:"viewportHeight"`
	Renderer       string   `json:"renderer"`
	VendorRender   string   `json:"vendor"`
	TimeZone       string   `json:"timeZone"`
	DirectLinkID   string   `json:"directLinkID"`
}

func (pc PC) ToFingerprint() models.Fingerprint {
	return models.Fingerprint{
		UserAgent:      &pc.UserAgent,
		Platform:       pc.Platform,
		Language:       pc.Language,
		Languages:      pq.StringArray(pc.Languages),
		Cores:          &pc.Cores,
		Memory:         &pc.Memory,
		ScreenWidth:    pc.ScreenWidth,
		ScreenHeight:   pc.ScreenHeight,
		ColorDepth:     pc.ColorDepth,
		PixelRatio:     pc.PixelRatio,
		ViewportWidth:  pc.ViewportWidth,
		ViewportHeight: pc.ViewportHeight,
		TimeZone:       pc.TimeZone,
		DirectLinkID:   pc.DirectLinkID,
	}
}
func (mobile Mobile) ToFingerprint() models.Fingerprint {
	return models.Fingerprint{
		UserAgent:      &mobile.UserAgent,
		Platform:       mobile.Platform,
		Version:        mobile.Version,
		Language:       mobile.Language,
		Languages:      pq.StringArray(mobile.Languages),
		Cores:          &mobile.Cores,
		Memory:         &mobile.Memory,
		ScreenWidth:    mobile.ScreenWidth,
		ScreenHeight:   mobile.ScreenHeight,
		ColorDepth:     mobile.ColorDepth,
		PixelRatio:     mobile.PixelRatio,
		ViewportWidth:  mobile.ViewportWidth,
		ViewportHeight: mobile.ViewportHeight,
		Renderer:       mobile.Renderer,
		VendorRender:   &mobile.VendorRender,
		TimeZone:       mobile.TimeZone,
		DirectLinkID:   mobile.DirectLinkID,
	}
}
