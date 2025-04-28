package migrations

import (
	"fmt"

	"gorm.io/gorm"
	models2 "redirectServer/internal/database/models"
)

type Migration1 struct {
	db *gorm.DB
}

func NewMigration1(db *gorm.DB) *Migration1 {
	return &Migration1{db: db}
}

func (m *Migration1) Up() error {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Migrator().AutoMigrate(&models2.DirectLink{}, &models2.Fingerprint{})
		if err != nil {
			return fmt.Errorf("failed to UP CreateFingerPrint_Migration: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (m *Migration1) Down() error {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Migrator().DropTable(&models2.DirectLink{}, &models2.Fingerprint{})
		if err != nil {
			return fmt.Errorf("failed to DOWN CreateFingerPrint_Migration: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
