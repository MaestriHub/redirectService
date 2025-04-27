package repository

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"redirectServer/internal/database/models"
	"redirectServer/internal/domain"
)

type FingerprintRepo interface {
	Create(ctx *gin.Context, fp *domain.Fingerprint) error
	Find(ctx *gin.Context, fp *domain.FingerprintFields) (*domain.Fingerprint, error)
}

type fingerprintRepo struct {
	db *gorm.DB
}

func NewFingerprintRepo(db *gorm.DB) FingerprintRepo {
	return &fingerprintRepo{db: db}
}

func (f fingerprintRepo) Create(ctx *gin.Context, fp *domain.Fingerprint) error {
	dbFp := models.NewFingerprintDB(fp)
	if err := f.db.Create(&dbFp).Error; err != nil {
		return fmt.Errorf("db create finger print create: %w", err)
	}
	return nil
}

func (f fingerprintRepo) Find(ctx *gin.Context, fpFields *domain.FingerprintFields) (*domain.Fingerprint, error) {
	var fp models.Fingerprint

	w := NewFingerprintWeights()

	query := f.db.Table("fingerprints").
		Select(
			`*, (
		        CASE WHEN ip = ? THEN ? ELSE 0 END +
		        CASE WHEN user_agent = ? THEN ? ELSE 0 END +
		        CASE WHEN language = ? THEN ? ELSE 0 END +
		        CASE WHEN cores = ? THEN ? ELSE 0 END +
		        CASE WHEN memory = ? THEN ? ELSE 0 END +
		        CASE WHEN screen_width = ? THEN ? ELSE 0 END +
		        CASE WHEN screen_height = ? THEN ? ELSE 0 END +
		        CASE WHEN color_depth = ? THEN ? ELSE 0 END +
		        CASE WHEN pixel_ratio = ? THEN ? ELSE 0 END +
		        CASE WHEN time_zone = ? THEN ? ELSE 0 END
		    ) as score`,
			fpFields.IP, w.IP,
			fpFields.UserAgent, w.UserAgent,
			fpFields.Language, w.Language,
			fpFields.Cores, w.Cores,
			fpFields.Memory, w.Memory,
			fpFields.ScreenWidth, w.ScreenWidth,
			fpFields.ScreenHeight, w.ScreenHeight,
			fpFields.ColorDepth, w.ColorDepth,
			fpFields.PixelRatio, w.PixelRatio,
			fpFields.TimeZone, w.TimeZone,
		).
		Where(
			`(
		        CASE WHEN ip = ? THEN ? ELSE 0 END +
		        CASE WHEN user_agent = ? THEN ? ELSE 0 END +
		        CASE WHEN language = ? THEN ? ELSE 0 END +
		        CASE WHEN cores = ? THEN ? ELSE 0 END +
		        CASE WHEN memory = ? THEN ? ELSE 0 END +
		        CASE WHEN screen_width = ? THEN ? ELSE 0 END +
		        CASE WHEN screen_height = ? THEN ? ELSE 0 END +
		        CASE WHEN color_depth = ? THEN ? ELSE 0 END +
		        CASE WHEN pixel_ratio = ? THEN ? ELSE 0 END +
		        CASE WHEN time_zone = ? THEN ? ELSE 0 END
		    ) > 50
		`,
			fpFields.IP, w.IP,
			fpFields.UserAgent, w.UserAgent,
			fpFields.Language, w.Language,
			fpFields.Cores, w.Cores,
			fpFields.Memory, w.Memory,
			fpFields.ScreenWidth, w.ScreenWidth,
			fpFields.ScreenHeight, w.ScreenHeight,
			fpFields.ColorDepth, w.ColorDepth,
			fpFields.PixelRatio, w.PixelRatio,
			fpFields.TimeZone, w.TimeZone).
		Order("score DESC")

	if err := query.Scan(&fp).Error; err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}

	if fp.Fingerprint == nil {
		return nil, fmt.Errorf("fingerprint not found")
	}

	return fp.Fingerprint, nil
}

type FingerprintWeights struct {
	IP             uint
	UserAgent      uint
	Language       uint
	Languages      uint
	Cores          uint
	Memory         uint
	ScreenWidth    uint
	ScreenHeight   uint
	ColorDepth     uint
	PixelRatio     uint
	ViewportWidth  uint
	ViewportHeight uint
	TimeZone       uint
}

func NewFingerprintWeights() *FingerprintWeights {
	return &FingerprintWeights{
		IP:           40, // Высокий вес, так как IP часто стабилен
		UserAgent:    25, // Важный параметр для идентификации устройства
		Language:     15, // Язык браузера обычно фиксирован
		Languages:    0,  // Игнорируем, так как может сильно варьироваться
		Cores:        10, // Менее надежный параметр
		Memory:       10, // Менее надежный параметр
		ScreenWidth:  5,  // Зависит от устройства
		ScreenHeight: 5,  // Зависит от устройства
		ColorDepth:   5,  // Редко меняется
		PixelRatio:   5,  // Редко меняется
		TimeZone:     5,  // Зависит от региона
	}
}
