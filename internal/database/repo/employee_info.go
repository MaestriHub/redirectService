package repository

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"redirectServer/internal/domain"
)

type EmployeeInfoRepo interface {
	GetInfo(ctx *gin.Context, employeeId uuid.UUID) (*domain.Employee, error)
}

type employeeInfoRepo struct {
	db *gorm.DB
}

func NewEmployeeInfoRepo(db *gorm.DB) EmployeeInfoRepo {
	return &employeeInfoRepo{db: db}
}

func (s *employeeInfoRepo) GetInfo(ctx *gin.Context, employeeId uuid.UUID) (*domain.Employee, error) {
	var employee domain.Employee
	if err := s.db.
		Table("employees").
		Where("id = ?", employeeId).
		First(&employee).Error; err != nil {
		return nil, fmt.Errorf("SalonInfoRepo.GetInfo: %w", err)
	}
	return &employee, nil
}
