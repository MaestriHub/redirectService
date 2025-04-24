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
	var fp models.FingerprintDB
	if err := f.db.Where("ip = ?", fpFields.IP).First(&fp).Error; err != nil {
		return nil, fmt.Errorf("db find fingerprint: %w", err)
	}
	return fp.Fingerprint, nil
}
