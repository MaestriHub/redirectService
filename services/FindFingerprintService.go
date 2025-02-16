package services

import (
	"redirectServer/models"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func FindFingerprint(data models.Fingerprint, db *gorm.DB) *models.Fingerprint {
	conditions := map[string]interface{}{
		"ip":        data.IP,
		"platform":  data.Platform,
		"version":   data.Version,
		"language":  data.Language,
		"languages": pq.StringArray(data.Languages),
		"cores":     data.Cores,
		//"memory":          data.Memory,
		"screen_width":    data.ScreenWidth,
		"screen_height":   data.ScreenHeight,
		"color_depth":     data.ColorDepth,
		"pixel_ratio":     data.PixelRatio,
		"viewport_width":  data.ViewportWidth,
		"viewport_height": data.ViewportHeight,
		"renderer":        data.Renderer,
		"vendor_render":   data.VendorRender,
		"time_zone":       data.TimeZone,
	}

	query := db.Where(conditions)
	if data.Memory != nil {
		query = query.Where("memory = ?", *data.Memory)
	} else {
		query = query.Where("memory IS NULL")
	}

	if data.UserAgent != nil {
		query = query.Where("user_agent = ?", *data.UserAgent)
	}
	var existingFingerprint models.Fingerprint
	query.First(&existingFingerprint)
	return &existingFingerprint

}
