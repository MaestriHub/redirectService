package test

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"redirectServer/internal/domain"
)

type FingerprintBuilder struct {
	Fingerprint *domain.Fingerprint
}

func NewFingerprintBuilder() *FingerprintBuilder {
	userAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 18_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Mobile/15E148 Safari/604.1"
	cores := 4
	memory := 0
	return &FingerprintBuilder{
		Fingerprint: &domain.Fingerprint{
			IP:             "192.0.2.1",
			UserAgent:      &userAgent,
			Platform:       "iPhone",
			Version:        "18.2.1",
			Language:       "ru",
			Languages:      []string{"ru"},
			ScreenWidth:    414,
			ScreenHeight:   896,
			ColorDepth:     24,
			PixelRatio:     2,
			ViewportWidth:  414,
			ViewportHeight: 714,
			Renderer:       "Apple GPU",
			TimeZone:       "Europe/Moscow",
			Cores:          &cores,
			Memory:         &memory,
			DirectLinkID:   "none",
		},
	}
}

func (rb *FingerprintBuilder) SetUserAgent(userAgent *string) *FingerprintBuilder {
	rb.Fingerprint.UserAgent = userAgent
	return rb
}

func (rb *FingerprintBuilder) SetDirectLink(directLink domain.DirectLink) *FingerprintBuilder {
	rb.Fingerprint.DirectLinkID = directLink.ID
	return rb
}

func (rb *FingerprintBuilder) SetPlatform(platform string) *FingerprintBuilder {
	rb.Fingerprint.Platform = platform
	return rb
}

func (rb *FingerprintBuilder) SetVersion(version string) *FingerprintBuilder {
	rb.Fingerprint.Version = version
	return rb
}

func (rb *FingerprintBuilder) SetLanguage(language string) *FingerprintBuilder {
	rb.Fingerprint.Language = language
	return rb
}

func (rb *FingerprintBuilder) SetLanguages(languages pq.StringArray) *FingerprintBuilder {
	rb.Fingerprint.Languages = languages
	return rb
}

func (rb *FingerprintBuilder) SetCores(cores int) *FingerprintBuilder {
	rb.Fingerprint.Cores = &cores
	return rb
}

func (rb *FingerprintBuilder) SetMemory(memory int) *FingerprintBuilder {
	rb.Fingerprint.Memory = &memory
	return rb
}

func (rb *FingerprintBuilder) SetScreenWidth(width int) *FingerprintBuilder {
	rb.Fingerprint.ScreenWidth = width
	return rb
}

func (rb *FingerprintBuilder) SetScreenHeight(height int) *FingerprintBuilder {
	rb.Fingerprint.ScreenHeight = height
	return rb
}

func (rb *FingerprintBuilder) SetColorDepth(depth int) *FingerprintBuilder {
	rb.Fingerprint.ColorDepth = depth
	return rb
}

func (rb *FingerprintBuilder) SetPixelRatio(ratio float64) *FingerprintBuilder {
	rb.Fingerprint.PixelRatio = ratio
	return rb
}

func (rb *FingerprintBuilder) SetViewportWidth(width int) *FingerprintBuilder {
	rb.Fingerprint.ViewportWidth = width
	return rb
}

func (rb *FingerprintBuilder) SetViewportHeight(height int) *FingerprintBuilder {
	rb.Fingerprint.ViewportHeight = height
	return rb
}

func (rb *FingerprintBuilder) SetRenderer(renderer string) *FingerprintBuilder {
	rb.Fingerprint.Renderer = renderer
	return rb
}

func (rb *FingerprintBuilder) SetVendorRender(vendorRender *string) *FingerprintBuilder {
	rb.Fingerprint.VendorRender = vendorRender
	return rb
}

func (rb *FingerprintBuilder) SetTimeZone(timeZone string) *FingerprintBuilder {
	rb.Fingerprint.TimeZone = timeZone
	return rb
}

func (rb *FingerprintBuilder) SetStatuses(directLink string) *FingerprintBuilder {
	rb.Fingerprint.DirectLinkID = directLink
	return rb
}

// Завершающий метод для создания и сохранения Fingerprint в базе данных
func (rb *FingerprintBuilder) Build(db *gorm.DB) *domain.Fingerprint {
	// Сохраняем Fingerprint в базу данных перед возвратом
	if rb.Fingerprint.DirectLinkID == "none" {
		directLink := NewDirectLinkBuilder().Build(db)
		rb.Fingerprint.DirectLinkID = directLink.ID
	}
	if err := db.Create(rb.Fingerprint).Error; err != nil {
		return nil
	}
	return rb.Fingerprint
}
