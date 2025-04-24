package models

import (
	"gorm.io/gorm"
	"redirectServer/internal/domain"
)

type FingerprintDB struct {
	gorm.Model
	*domain.Fingerprint
}

func NewFingerprintDB(fp *domain.Fingerprint) *FingerprintDB {
	return &FingerprintDB{
		Fingerprint: fp,
	}
}
