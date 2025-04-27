package params

import (
	"redirectServer/internal/domain"
)

type Fingerprint struct {
	Language     string   `json:"language"`
	Languages    []string `json:"languages"`
	Cores        int      `json:"cores"`
	Memory       int      `json:"memory"`
	ScreenWidth  int      `json:"screenWidth"`
	ScreenHeight int      `json:"screenHeight"`
	ColorDepth   int      `json:"colorDepth"`
	PixelRatio   float64  `json:"pixelRatio"`
	TimeZone     string   `json:"timeZone"`
}

func (p *Fingerprint) ToDomain(ip string, ua domain.UserAgent, linkId domain.NanoID) *domain.Fingerprint {
	return &domain.Fingerprint{
		FingerprintFields: *p.ToFields(ip, ua),
		LinkId:            linkId,
	}
}

func (p *Fingerprint) ToFields(ip string, ua domain.UserAgent) *domain.FingerprintFields {
	return &domain.FingerprintFields{
		IP:           ip,
		UserAgent:    ua,
		Language:     p.Language,
		Languages:    p.Languages,
		Cores:        p.Cores,
		Memory:       p.Memory,
		ScreenWidth:  p.ScreenWidth,
		ScreenHeight: p.ScreenHeight,
		ColorDepth:   p.ColorDepth,
		PixelRatio:   p.PixelRatio,
		TimeZone:     p.TimeZone,
	}
}
