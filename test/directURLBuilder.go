package test

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"redirectServer/internal/domain"
	"redirectServer/internal/repo/models/payload"
)

type DirectLinkBuilder struct {
	directLink *domain.DirectLink
}

func NewDirectLinkBuilder() *DirectLinkBuilder {
	parsedUUID, _ := uuid.Parse("53bb0f86-a94e-4302-8a07-ea0b083d3bde")
	payload := payload.MasterToSalon{
		EmployeeId: parsedUUID,
	}
	directLink := domain.DirectLink{
		ID:    "YSg6Ugcf",
		Event: "EmployeeInvite",
	}
	directLink.SetPayload(payload)

	return &DirectLinkBuilder{
		directLink: &directLink,
	}
}

func (du *DirectLinkBuilder) SetID(id string) *DirectLinkBuilder {
	du.directLink.ID = id
	return du
}

func (du *DirectLinkBuilder) SetPayload(payload payload.MasterToSalon) *DirectLinkBuilder {

	du.directLink.SetPayload(payload)
	return du
}

func (du *DirectLinkBuilder) SetEvent(event string) *DirectLinkBuilder {
	du.directLink.Event = event
	return du
}

func (du *DirectLinkBuilder) SetСlicks(clicks int) *DirectLinkBuilder {
	du.directLink.Clicks = clicks
	return du
}

func (du *DirectLinkBuilder) Build(db *gorm.DB) *domain.DirectLink {
	if err := db.Create(du.directLink).Error; err != nil {
		return nil
	}
	return du.directLink
}
