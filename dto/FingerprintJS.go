package dto

import (
	"redirectServer/model"

	"github.com/lib/pq"
)

// TODO: Подумать над неймингом и перенести в модель(как минимум DeviceData и тд)
type FingerprintJS struct {
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

func (mobile FingerprintJS) ToFingerprint() model.Fingerprint {
	return model.Fingerprint{
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
