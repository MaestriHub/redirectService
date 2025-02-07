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
	Сlicks    int        `gorm:"default:0"`
	Payload   PayloadURL `gorm:"embedded"`
}

type PayloadURL struct {
	Title       string
	Name        string
	Description string
}

type ParticalDirectURL struct {
	URL     string
	Payload PayloadURL
}
