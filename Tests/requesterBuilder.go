package test

import (
	"redirectServer/models"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type RequesterBuilder struct {
	requester *models.Requester
}

func NewRequesterBuilder() *RequesterBuilder {
	userAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 18_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Mobile/15E148 Safari/604.1"
	cores := 4
	memory := 0
	return &RequesterBuilder{
		requester: &models.Requester{
			IP:             "192.168.0.111",
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
			Statuses:       []string{string(models.Linked)},
		},
	}
}

func (rb *RequesterBuilder) SetUserAgent(userAgent *string) *RequesterBuilder {
	rb.requester.UserAgent = userAgent
	return rb
}

func (rb *RequesterBuilder) SetPlatform(platform string) *RequesterBuilder {
	rb.requester.Platform = platform
	return rb
}

func (rb *RequesterBuilder) SetVersion(version string) *RequesterBuilder {
	rb.requester.Version = version
	return rb
}

func (rb *RequesterBuilder) SetLanguage(language string) *RequesterBuilder {
	rb.requester.Language = language
	return rb
}

func (rb *RequesterBuilder) SetLanguages(languages pq.StringArray) *RequesterBuilder {
	rb.requester.Languages = languages
	return rb
}

func (rb *RequesterBuilder) SetCores(cores int) *RequesterBuilder {
	rb.requester.Cores = &cores
	return rb
}

func (rb *RequesterBuilder) SetMemory(memory int) *RequesterBuilder {
	rb.requester.Memory = &memory
	return rb
}

func (rb *RequesterBuilder) SetScreenWidth(width int) *RequesterBuilder {
	rb.requester.ScreenWidth = width
	return rb
}

func (rb *RequesterBuilder) SetScreenHeight(height int) *RequesterBuilder {
	rb.requester.ScreenHeight = height
	return rb
}

func (rb *RequesterBuilder) SetColorDepth(depth int) *RequesterBuilder {
	rb.requester.ColorDepth = depth
	return rb
}

func (rb *RequesterBuilder) SetPixelRatio(ratio float64) *RequesterBuilder {
	rb.requester.PixelRatio = ratio
	return rb
}

func (rb *RequesterBuilder) SetViewportWidth(width int) *RequesterBuilder {
	rb.requester.ViewportWidth = width
	return rb
}

func (rb *RequesterBuilder) SetViewportHeight(height int) *RequesterBuilder {
	rb.requester.ViewportHeight = height
	return rb
}

func (rb *RequesterBuilder) SetRenderer(renderer string) *RequesterBuilder {
	rb.requester.Renderer = renderer
	return rb
}

func (rb *RequesterBuilder) SetVendorRender(vendorRender *string) *RequesterBuilder {
	rb.requester.VendorRender = vendorRender
	return rb
}

func (rb *RequesterBuilder) SetTimeZone(timeZone string) *RequesterBuilder {
	rb.requester.TimeZone = timeZone
	return rb
}

func (rb *RequesterBuilder) SetStatuses(statuses pq.StringArray) *RequesterBuilder {
	rb.requester.Statuses = statuses
	return rb
}

// Завершающий метод для создания и сохранения Requester в базе данных
func (rb *RequesterBuilder) Build(db *gorm.DB) (*models.Requester, error) {
	// Сохраняем Requester в базу данных перед возвратом
	if err := db.Create(rb.requester).Error; err != nil {
		return nil, err
	}
	return rb.requester, nil
}
