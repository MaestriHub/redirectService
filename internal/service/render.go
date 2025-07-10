package service

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"redirectServer/configs"
	"redirectServer/internal/database/repo"
	"redirectServer/internal/domain"
	"redirectServer/pkg"
)

type RenderService interface {
	RenderMain(ctx *gin.Context, link *domain.DirectLink, ua domain.UserAgent) (string, *pkg.ErrorS)
}

type renderService struct {
	cfg       configs.AppStoreLinksConfig
	repoSalon repository.SalonInfoRepo
}

func NewRenderService(cfg *configs.AppStoreLinksConfig, s repository.SalonInfoRepo) RenderService {
	return &renderService{cfg: *cfg, repoSalon: s}
}

type placeholder = string

const (
	hAppStoreLink         placeholder = "{{.AppStoreLink}}"
	hName                 placeholder = "{{.Name}}"
	hDescription          placeholder = "{{.Description}}"
	hDynamicUniversalLink placeholder = "{{.DynamicUniversalLink}}"
	hLinkID               placeholder = "{{.LinkID}}"
)

func (s *renderService) RenderMain(ctx *gin.Context, link *domain.DirectLink, ua domain.UserAgent) (string, *pkg.ErrorS) {
	event := link.Event

	path := getHTMLPath(event, ua)

	htmlFile, err := os.ReadFile(path)
	if err != nil {
		return "", pkg.NewInternalServerError(fmt.Errorf("read file %s: %w", path, err).Error())
	}

	mHTML := string(htmlFile)

	switch ua {
	case domain.IOS:
		mHTML = strings.Replace(mHTML, hAppStoreLink, s.cfg.IOSLink, -1)
	case domain.ANDROID:
		mHTML = strings.Replace(mHTML, hAppStoreLink, s.cfg.AndroidLink, -1)
	default:
	}

	switch e := event.(type) {
	case *domain.EmployeeInviteEvent:
		salon, err := s.repoSalon.GetInfo(ctx, e.SalonId)
		if err != nil {
			log.Printf("не нашли салон по айдишнику %v", e.SalonId)
			salon = domain.NewDummySalon(e.SalonId)
		}
		mHTML = strings.Replace(mHTML, hName, salon.Name, -1)
		mHTML = strings.Replace(mHTML, hDescription, salon.Description, -1)

	case *domain.SalonInviteEvent:
		salon, err := s.repoSalon.GetInfo(ctx, e.SalonId)
		if err != nil {
			log.Printf("не нашли салон по айдишнику %v", e.SalonId)
			salon = domain.NewDummySalon(e.SalonId)
		}
		mHTML = strings.Replace(mHTML, hName, salon.Name, -1)
		mHTML = strings.Replace(mHTML, hDescription, salon.Description, -1)

	case *domain.ClientInviteEvent:
		salon, err := s.repoSalon.GetInfo(ctx, e.SalonId)
		if err != nil {
			log.Printf("не нашли салон по айдишнику %v", e.SalonId)
			salon = domain.NewDummySalon(e.SalonId)
		}
		mHTML = strings.Replace(mHTML, hName, salon.Name, -1)
		mHTML = strings.Replace(mHTML, hDescription, salon.Description, -1)
	}

	mHTML = injectAppStoreLink(mHTML, s.cfg, ua)
	mHTML = strings.Replace(mHTML, hDynamicUniversalLink, link.ToURL(), -1)
	mHTML = strings.Replace(mHTML, hLinkID, link.NanoId, -1)

	return mHTML, nil
}

func getHTMLPath(event domain.Event, ua domain.UserAgent) string {
	return fmt.Sprintf("static/%s/%s/Screen.html", event.GetType(), ua)
}

func injectAppStoreLink(html string, c configs.AppStoreLinksConfig, ua domain.UserAgent) string {
	switch ua {
	case domain.IOS:
		return strings.Replace(html, hAppStoreLink, c.IOSLink, -1)
	case domain.ANDROID:
		return strings.Replace(html, hAppStoreLink, c.AndroidLink, -1)
	default:
		return html
	}
}
