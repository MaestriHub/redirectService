package models

import (
	"fmt"
	"time"

	"net/url"
	"strings"

	"gorm.io/gorm"
)

// подумать о связке directLink + DeviceId
//TODO: Метод по NanoId получить DirectLink этот метод при установленном приложении
//TODO: Метод по фингерпринту + Optional(link.maetry.com/{NanoId}) вернуть DirectLink. при первой установки приложения

type inviteEvent string

const (
	EmployeerInvite inviteEvent = "EmployeerInvite"
	SalonInvite     inviteEvent = "SalonInvite"
	CustomerInvite  inviteEvent = "CustomerInvite"
)

type DirectLink struct { //link.maetry.com/{NanoId
	ID        string // NanoId
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Сlicks    int            `gorm:"default:0"`
	Payload   string         `json:"payload"`
	Event     string         `json:"event"`
}

type ParticalDirectLink struct {
	Payload string
	Event   string
}

func (directLink DirectLink) ParseToURL() string {
	return "https://link.maetry.com/" + directLink.ID
}

func ParseURL(input string, db *gorm.DB) (*DirectLink, error) {
	var directLink DirectLink

	parsedURL, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}
	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("URL doesn't contain asigned parts")
	}

	db.Where("url_event = ? AND payload = ?", parts[0], parts[1]).Order("createdAt desc").First(&directLink)
	if directLink.ID == "" {
		return nil, fmt.Errorf("URL doesn't exist")
	}
	return &directLink, nil
}
func (directLink DirectLink) GetPayloadEmployee(db *gorm.DB) (string, string, error) {
	var rawAnswers struct {
		Name  string
		Title string
	}
	err := db.Raw(`
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
	`, directLink.Payload).Scan(&rawAnswers)
	if err.Error != nil {
		return "", "", err.Error
	}
	return rawAnswers.Name, rawAnswers.Title, nil
}

func (directLink DirectLink) GetPayloadSalon(db *gorm.DB) (string, string, error) {
	var rawAnswer struct {
		Name        string
		Description string
	}
	err := db.Raw("SELECT name, description FROM salons WHERE id = ?", directLink.Payload).Scan(&rawAnswer)
	if err.Error != nil {
		return "", "", err.Error
	}
	return rawAnswer.Name, rawAnswer.Description, nil
}
