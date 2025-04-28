package service

import (
	"github.com/gin-gonic/gin"
	"redirectServer/internal/database/repo"
	"redirectServer/internal/domain"
	"redirectServer/pkg"
)

type LinkService interface {
	CreateInvite(ctx *gin.Context, link *domain.DirectLink) *pkg.ErrorS
	LinkTap(ctx *gin.Context, nanoId domain.NanoID) (*domain.DirectLink, *pkg.ErrorS)
	Find(ctx *gin.Context, nanoId domain.NanoID) (*domain.DirectLink, *pkg.ErrorS)
}

type linkService struct {
	repo repository.LinkRepo
}

func NewLinkService(repo repository.LinkRepo) LinkService {
	return &linkService{repo: repo}
}

func (l linkService) CreateInvite(ctx *gin.Context, link *domain.DirectLink) *pkg.ErrorS {
	if err := link.Validate(); err != nil {
		return pkg.NewBadRequestError(err.Error())
	}

	if err := l.repo.Create(ctx, link); err != nil {
		return pkg.NewInternalServerError(err.Error())
	}

	return nil
}

func (l linkService) LinkTap(ctx *gin.Context, nanoId domain.NanoID) (*domain.DirectLink, *pkg.ErrorS) {
	link, err := l.repo.Find(ctx, nanoId)
	if err != nil {
		return nil, pkg.NewNotFoundError(err.Error())
	}

	link.IncClicks()

	if err = l.repo.Update(ctx, link); err != nil {
		return nil, pkg.NewInternalServerError(err.Error())
	}

	return link, nil
}

func (l linkService) Find(ctx *gin.Context, nanoId domain.NanoID) (*domain.DirectLink, *pkg.ErrorS) {
	link, err := l.repo.Find(ctx, nanoId)
	if err != nil {
		return nil, pkg.NewNotFoundError(err.Error())
	}
	return link, nil
}
