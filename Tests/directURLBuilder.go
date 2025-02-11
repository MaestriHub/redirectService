package test

import (
	"redirectServer/models"

	"gorm.io/gorm"
)

type DirectURLBuilder struct {
	directURL *models.DirectURL
}

func NewDirectURLBuilder() *DirectURLBuilder {
	return &DirectURLBuilder{
		directURL: &models.DirectURL{
			ID:       "YSg6Ugcf",
			Payload:  "53bb0f86-a94e-4302-8a07-ea0b083d3bde",
			URLEvent: "EmployeerInvite",
		},
	}
}

func (du *DirectURLBuilder) SetID(id string) *DirectURLBuilder {
	du.directURL.ID = id
	return du
}

func (du *DirectURLBuilder) SetPayload(payload string) *DirectURLBuilder {
	du.directURL.Payload = payload
	return du
}

func (du *DirectURLBuilder) SetURLEvent(urlEvent string) *DirectURLBuilder {
	du.directURL.URLEvent = urlEvent
	return du
}

func (du *DirectURLBuilder) SetСlicks(clicks int) *DirectURLBuilder {
	du.directURL.Сlicks = clicks
	return du
}

func (du *DirectURLBuilder) Build(db *gorm.DB) (*models.DirectURL, error) {
	if err := db.Create(du.directURL).Error; err != nil {
		return nil, err
	}
	return du.directURL, nil
}
