package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"redirectServer/internal/database/repo"
	"redirectServer/internal/domain"
)

type LinkService interface {
	CreateInvite(ctx *gin.Context, link *domain.DirectLink) error
	LinkTap(ctx *gin.Context, nanoId domain.NanoID) (*domain.DirectLink, error)
	Find(ctx *gin.Context, nanoId domain.NanoID) (*domain.DirectLink, error)
}

type linkService struct {
	repo repository.LinkRepo
}

func NewLinkService(repo repository.LinkRepo) LinkService {
	return &linkService{repo: repo}
}

func (l linkService) CreateInvite(ctx *gin.Context, link *domain.DirectLink) error {
	if err := link.Validate(); err != nil {
		return fmt.Errorf("validate link: %w", err)
	}

	if err := l.repo.Create(ctx, link); err != nil {
		return fmt.Errorf("create link: %w", err)
	}

	return nil
}

func (l linkService) LinkTap(ctx *gin.Context, nanoId domain.NanoID) (*domain.DirectLink, error) {
	link, err := l.repo.Find(ctx, nanoId)
	if err != nil {
		return nil, fmt.Errorf("not found link: %w", err)
	}

	link.IncClicks()

	if err = l.repo.Update(ctx, link); err != nil {
		return nil, fmt.Errorf("update link: %w", err)
	}

	return link, nil
}

func (l linkService) Find(ctx *gin.Context, nanoId domain.NanoID) (*domain.DirectLink, error) {
	link, err := l.repo.Find(ctx, nanoId)
	if err != nil {
		return nil, fmt.Errorf("not found link: %w", err)
	}
	return link, nil
}
