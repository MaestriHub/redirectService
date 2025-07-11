package repository

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"redirectServer/internal/database/models"
	"redirectServer/internal/domain"
)

type LinkRepo interface {
	Create(ctx *gin.Context, link *domain.DirectLink) error
	Update(ctx *gin.Context, link *domain.DirectLink) error
	Find(ctx *gin.Context, id domain.NanoID) (*domain.DirectLink, error)
}

type linkRepo struct {
	db *gorm.DB
}

func NewLinkRepo(db *gorm.DB) LinkRepo {
	return &linkRepo{db: db}
}

func (l linkRepo) Create(ctx *gin.Context, link *domain.DirectLink) error {
	dbLink, err := models.NewDirectLinkDB(link)
	if err != nil {
		return err
	}

	if err := l.db.Create(&dbLink).Error; err != nil {
		return fmt.Errorf("db create link: %w", err)
	}
	return nil
}

func (l linkRepo) Update(ctx *gin.Context, link *domain.DirectLink) error {
	dbLink, err := models.NewDirectLinkDB(link)
	if err != nil {
		return err
	}

	if err := l.db.
		Model(&models.DirectLink{}).
		Where("nano_id = ?", link.NanoId).
		Limit(1).
		Order("create_at asc").
		Updates(&dbLink).Error; err != nil {
		return fmt.Errorf("db update link: %w", err)
	}

	return nil
}

func (l linkRepo) Find(ctx *gin.Context, id domain.NanoID) (*domain.DirectLink, error) {
	var link models.DirectLink
	if err := l.db.Where("nano_id = ?", id).First(&link).Error; err != nil {
		return nil, fmt.Errorf("db find link: %w", err)
	}

	event, err := domain.NewEvent(link.Event, link.Payload)
	if err != nil {
		return nil, fmt.Errorf("db create event: %w", err)
	}

	domainLink := &domain.DirectLink{
		NanoId: link.NanoId,
		Clicks: link.Clicks,
		Event:  event,
	}

	return domainLink, nil
}
