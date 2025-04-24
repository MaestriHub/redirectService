package migrations

import (
	"fmt"

	"gorm.io/gorm"
	models2 "redirectServer/internal/database/models"
)

// create tables fingerprint and directlink

type Migration1 struct {
	db *gorm.DB
}

func NewMigration_1(db *gorm.DB) *Migration1 {
	return &Migration1{db: db}
}

func (m *Migration1) Up() error {
	m.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Migrator().AutoMigrate(&models2.DirectLinkDB{}, &models2.FingerprintDB{})
		if err != nil {
			return fmt.Errorf("failed to UP CreateFingerPrint_Migration: %w", err)
		}
		return nil
	})

	return nil
}

func (m *Migration1) Down() error {
	m.db.Transaction(func(tx *gorm.DB) error {
		err := m.db.Migrator().DropTable(&models2.DirectLinkDB{}, &models2.FingerprintDB{})
		if err != nil {
			return fmt.Errorf("failed to DOWN CreateFingerPrint_Migration: %w", err)
		}
		return nil
	})

	return nil
}
