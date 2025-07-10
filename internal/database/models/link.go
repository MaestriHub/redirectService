package models

import (
	"net/http"

	"redirectServer/internal/domain"
	"redirectServer/pkg"
)

type DirectLink struct {
	BaseUUIDModel
	NanoId  domain.NanoID
	Event   string
	Payload []byte `gorm:"type:jsonb"`
	Clicks  int    `gorm:"default:0"`
}

func NewDirectLinkDB(link *domain.DirectLink) (*DirectLink, error) {
	payload, err := link.Event.GetPayload()
	if err != nil {
		return nil, pkg.NewErrorS("Ooops", http.StatusInternalServerError)
	}

	return &DirectLink{
		NanoId:  link.NanoId,
		Event:   link.Event.GetType(),
		Payload: payload,
		Clicks:  link.Clicks,
	}, nil
}
