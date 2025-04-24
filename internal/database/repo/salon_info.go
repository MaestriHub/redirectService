package repository

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"redirectServer/internal/domain"
)

type SalonInfoRepo interface {
	GetInfo(ctx *gin.Context, salonId uuid.UUID) (*domain.Salon, error)
}

type salonInfoRepo struct {
	db *gorm.DB
}

func NewSalonInfoRepo(db *gorm.DB) SalonInfoRepo {
	return &salonInfoRepo{db: db}
}

func (s *salonInfoRepo) GetInfo(ctx *gin.Context, salonId uuid.UUID) (*domain.Salon, error) {
	var salon domain.Salon
	if err := s.db.
		Table("salons").
		Where("id = ?", salonId).
		First(&salon).Error; err != nil {
		return nil, fmt.Errorf("SalonInfoRepo.GetInfo: %w", err)
	}
	return &salon, nil
}
