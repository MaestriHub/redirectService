package test

import (
	"redirectServer/models"

	"gorm.io/gorm"
)

type DirectLinkBuilder struct {
	directLink *models.DirectLink
}

func NewDirectLinkBuilder() *DirectLinkBuilder {
	return &DirectLinkBuilder{
		directLink: &models.DirectLink{
			ID:      "YSg6Ugcf",
			Payload: "53bb0f86-a94e-4302-8a07-ea0b083d3bde",
			Event:   "EmployeerInvite",
		},
	}
}

func (du *DirectLinkBuilder) SetID(id string) *DirectLinkBuilder {
	du.directLink.ID = id
	return du
}

func (du *DirectLinkBuilder) SetPayload(payload string) *DirectLinkBuilder {
	du.directLink.Payload = payload
	return du
}

func (du *DirectLinkBuilder) SetEvent(event string) *DirectLinkBuilder {
	du.directLink.Event = event
	return du
}

func (du *DirectLinkBuilder) SetСlicks(clicks int) *DirectLinkBuilder {
	du.directLink.Сlicks = clicks
	return du
}

func (du *DirectLinkBuilder) Build(db *gorm.DB) *models.DirectLink {
	if err := db.Create(du.directLink).Error; err != nil {
		return nil
	}
	return du.directLink
}
