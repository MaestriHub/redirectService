package service

import (
	"redirectServer/internal/domain"
)

type UrlGeneratorService interface {
	Generate(link *domain.DirectLink) string
}

type urlGeneratorService struct {
	urlPrefix string
}

func NewUrlGeneratorService(urlPrefix string) UrlGeneratorService {
	return urlGeneratorService{urlPrefix: urlPrefix}
}

func (u urlGeneratorService) Generate(link *domain.DirectLink) string {
	return u.urlPrefix + link.NanoId
}
