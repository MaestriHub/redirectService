package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseUUIDModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"omitempty"`
}

func (base *BaseUUIDModel) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID = uuid.New()
	return nil
}
