package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"redirectServer/models/payload"
	"strings"
	"time"

	"gorm.io/gorm"
)

type inviteEvent string

const (
	EmployeerInvite     inviteEvent = "EmployeerInvite"
	SalonInvite         inviteEvent = "SalonInvite"
	MasterInviteToSalon inviteEvent = "MasterInviteToSalon"
	CustomerInvite      inviteEvent = "CustomerInvite"
)

type DirectLink struct { //link.maetry.com/{NanoId
	ID        string // NanoId
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt  `gorm:"index"`
	Сlicks    int             `gorm:"default:0"`
	Payload   json.RawMessage `json:"payload" gorm:"type:jsonb"`
	Event     string          `json:"event"`
}

func (directLink DirectLink) ParseToURL() string {
	return "https://link.maetry.com/" + directLink.ID
}

// TODO: parsing dataa go посмотреть либы
func ParseURL(input string, db *gorm.DB) (*DirectLink, error) {
	parsedURL, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	code := strings.TrimPrefix(parsedURL.Path, "/")

	if code == "" {
		return nil, fmt.Errorf("code not found in URL")
	}
	var result DirectLink
	db.Where("id = ?", code).First(&result)
	if result.ID == "" {
		return nil, fmt.Errorf("URL not found")
	}

	return &result, nil
}

// TODO: подумать над нейменгом функции
// TODO: у нас есть тип салона индивидуальный мастер
func (directLink DirectLink) GetPayloadMasterToSalon(db *gorm.DB) (string, string, error) {
	var rawAnswers struct {
		Name  string
		Title string
	}
	var payload payload.MasterToSalon
	err := directLink.ToObject(&payload)
	if err != nil {
		return "", "", err
	}
	result := db.Raw(`
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
	`, payload.EmployeeId).Scan(&rawAnswers)
	if result.Error != nil {
		return "", "", result.Error
	}
	return rawAnswers.Name, rawAnswers.Title, nil
}

func (directLink DirectLink) GetPayloadSalon(db *gorm.DB) (string, string, error) {
	var rawAnswer struct {
		Name        string
		Description string
	}
	var payload payload.Salon
	err := directLink.ToObject(&payload)
	if err != nil {
		return "", "", err
	}
	result := db.Raw("SELECT name, description FROM salons WHERE id = ?", payload.ID).Scan(&rawAnswer)
	if result.Error != nil {
		return "", "", result.Error
	}
	return rawAnswer.Name, rawAnswer.Description, nil
}

// ToObject - преобразует JSON в указанный объект
func (directLink *DirectLink) ToObject(target any) error {
	return json.Unmarshal(directLink.Payload, target)
}

// SetData - сериализует объект в JSON
func (directLink *DirectLink) SetPayload(data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	directLink.Payload = bytes
	return nil
}
