package models

import (
	"gorm.io/gorm"
	"redirectServer/internal/domain"
)

type DirectLinkDB struct {
	gorm.Model
	*domain.DirectLink
}

func NewDirectLinkDB(link *domain.DirectLink) *DirectLinkDB {
	return &DirectLinkDB{
		DirectLink: link,
	}
}
