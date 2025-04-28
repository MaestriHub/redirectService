package models

import (
	"redirectServer/internal/domain"
)

type Fingerprint struct {
	BaseUUIDModel
	*domain.Fingerprint
}

func NewFingerprintDB(fp *domain.Fingerprint) *Fingerprint {
	return &Fingerprint{
		Fingerprint: fp,
	}
}
