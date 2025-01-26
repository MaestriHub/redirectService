package models

import (
	"time"

	"gorm.io/gorm"
)

type DirectURL struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	URL       string
	Сlicks    int `gorm:"default:0"`
}
