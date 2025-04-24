package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	repo "redirectServer/internal/database/repo"
	"redirectServer/internal/domain"
)

type FingerprintService interface {
	Create(ctx *gin.Context, fp *domain.Fingerprint) error
	Find(ctx *gin.Context, fp *domain.FingerprintFields) (*domain.DirectLink, error)
}

type fingerprintService struct {
	repoFp   repo.FingerprintRepo
	repoLink repo.LinkRepo
}

func NewFingerprintService(fp repo.FingerprintRepo, l repo.LinkRepo) FingerprintService {
	return &fingerprintService{repoFp: fp, repoLink: l}
}

func (f fingerprintService) Create(ctx *gin.Context, fp *domain.Fingerprint) error {
	if err := fp.Validate(); err != nil {
		return fmt.Errorf("validate fingerprint: %w", err)
	}

	if err := f.repoFp.Create(ctx, fp); err != nil {
		return fmt.Errorf("create fingerprint: %w", err)
	}

	return nil
}

func (f fingerprintService) Find(ctx *gin.Context, fpFields *domain.FingerprintFields) (*domain.DirectLink, error) {
	matchedFp, err := f.repoFp.Find(ctx, fpFields)
	if err != nil {
		return nil, fmt.Errorf("find fingerprint: %w", err)
	}

	link, err := f.repoLink.Find(ctx, matchedFp.LinkId)
	if err != nil {
		return nil, fmt.Errorf("find link: %w", err)
	}

	return link, nil
}
