package models

import (
	"redirectServer/internal/domain"
)

type DirectLink struct {
	BaseUUIDModel
	*domain.DirectLink
}

func NewDirectLinkDB(link *domain.DirectLink) *DirectLink {
	return &DirectLink{
		DirectLink: link,
	}
}
