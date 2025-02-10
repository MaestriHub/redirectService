package models

import (
	"fmt"
	"time"

	"errors"

	"net/url"
	"strings"

	"gorm.io/gorm"
)

type DirectURL struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Сlicks    int            `gorm:"default:0"`
	Payload   string
	URLEvent  string
}

type AssembledPayload struct {
	Event       string
	Name        string
	Description string
}

type ParticalDirectURL struct {
	Payload  string
	URLEvent string
}

func (directURL DirectURL) ParseToURL() string {
	return "https://link.maetry.com/" + directURL.URLEvent + "/" + directURL.Payload
}

func ParseURL(input string, db *gorm.DB) (*DirectURL, error) {
	var directURL DirectURL

	parsedURL, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}
	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("URL doesn't contain asigned parts")
	}

	db.Where("url_event = ? AND payload = ?", parts[0], parts[1]).Order("createdAt desc").First(&directURL)
	if directURL.ID == "" {
		return nil, fmt.Errorf("URL doesn't exist")
	}
	return &directURL, nil
}

func (directURL DirectURL) GetPayload(db *gorm.DB) (AssembledPayload, error) {
	var result AssembledPayload = AssembledPayload{
		Event:       "",
		Name:        "",
		Description: "",
	}
	if directURL.URLEvent == "EmployeerInvite" {
		var rawAnswers struct {
			Name  string
			Title string
		}
		db.Raw(`
		SELECT 
			s.name,
			p.title
		FROM 
			employees e
		JOIN 
			salons s ON e.salon_id = s.id
		JOIN 
			positions p ON e.position_id = p.id
		WHERE
			e.id = ?
		`, directURL.Payload).Scan(&rawAnswers)
		result.Event = "EmployeerInvite"
		result.Name = rawAnswers.Name
		result.Description = rawAnswers.Title
		return result, nil
	}
	if directURL.URLEvent == "CustomerInvite" {
	}
	if directURL.URLEvent == "SalonInvite" {
		var rawAnswer struct {
			Name        string
			Description string
		}
		db.Raw("SELECT name, description FROM salons WHERE id = ?", directURL.Payload).Scan(&rawAnswer)
		result.Event = "SalonInvite"
		result.Name = rawAnswer.Name
		result.Description = rawAnswer.Description
		return result, nil
	}

	return result, errors.New("неизвестый ивент URL")
}
