package service

import (
	"github.com/gin-gonic/gin"
	repo "redirectServer/internal/database/repo"
	"redirectServer/internal/domain"
	"redirectServer/pkg"
)

type FingerprintService interface {
	Create(ctx *gin.Context, fp *domain.Fingerprint) *pkg.ErrorS
	Find(ctx *gin.Context, fp *domain.FingerprintFields) (*domain.DirectLink, *pkg.ErrorS)
}

type fingerprintService struct {
	repoFp   repo.FingerprintRepo
	repoLink repo.LinkRepo
}

func NewFingerprintService(fp repo.FingerprintRepo, l repo.LinkRepo) FingerprintService {
	return &fingerprintService{repoFp: fp, repoLink: l}
}

func (f fingerprintService) Create(ctx *gin.Context, fp *domain.Fingerprint) *pkg.ErrorS {
	if err := fp.Validate(); err != nil {
		return pkg.NewBadRequestError(err.Error())
	}

	if err := f.repoFp.Create(ctx, fp); err != nil {
		return pkg.NewInternalServerError(err.Error())
	}

	return nil
}

func (f fingerprintService) Find(ctx *gin.Context, fpFields *domain.FingerprintFields) (*domain.DirectLink, *pkg.ErrorS) {
	matchedFp, err := f.repoFp.Find(ctx, fpFields)
	if err != nil {
		return nil, pkg.NewNotFoundError(err.Error())
	}

	link, err := f.repoLink.Find(ctx, matchedFp.LinkId)
	if err != nil {
		return nil, pkg.NewNotFoundError(err.Error())
	}

	return link, nil
}
