package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"redirectServer/internal/domain"
	"redirectServer/internal/service"
	"redirectServer/internal/transport/dto/params"
	"redirectServer/internal/transport/dto/resp"
	"redirectServer/pkg"
)

type FingerprintHandler interface {
	Create(*gin.Context)
	Find(*gin.Context)
}

type fingerprintHandler struct {
	linkService        service.LinkService
	fingerprintService service.FingerprintService
}

func NewFingerprintHandler(l service.LinkService, f service.FingerprintService) FingerprintHandler {
	return &fingerprintHandler{linkService: l, fingerprintService: f}
}

// Create godoc
//	@Summary		Create device fingerprint
//	@Description	Uses for define who go to app with link
//	@Tags			fingerprint
//	@Accept			json
//	@Produce		json
//	@Param			User-Agent	header		string	true	"Юзер агент пользователя. ex: Android"
//	@Param			linkId		path		string	true	"Идентификатор (NanoID)"
//	@Param			request	body		params.Fingerprint	true	"Данные об устройстве"
//	@Success		201
//	@Failure		400		{object}	resp.ErrorDTO	"Bad request"
//	@Failure		500		{object}	resp.ErrorDTO	"Internal server error"
//	@Router			/fingerprint/{linkId} [post]
func (h *fingerprintHandler) Create(ctx *gin.Context) {
	fpDTO := params.Fingerprint{}
	if err := ctx.ShouldBindWith(&fpDTO, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, resp.NewErrorDTO(err.Error()))
		return
	}
	ua := ctx.GetHeader(pkg.UserAgent)

	linkId, ok := ctx.Params.Get("linkId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, resp.NewErrorDTO("linkId is required"))
		return
	}

	fp := fpDTO.ToDomain(ctx.ClientIP(), domain.ParseUserAgent(ua), linkId)
	if err := h.fingerprintService.Create(ctx, fp); err != nil {
		ctx.JSON(err.Status, resp.NewErrorDTO(err.Message))
		return
	}

	ctx.Status(http.StatusCreated)
}

// Find godoc
//	@Summary		Find device by fingerprint
//	@Description	we want to associate direct link with new user in app by fingerprint
//	@Tags			fingerprint
//	@Accept			json
//	@Produce		json
//	@Param			User-Agent	header		string	true	"Юзер агент пользователя. ex: Android"
//	@Param			linkId		path		string	false	"Идентификатор (NanoID)"
//	@Param			request	body		params.Fingerprint	true	"Данные об устройстве"
//	@Success		200		{object}	resp.DirectLinkDTO
//	@Failure		400		{object}	resp.ErrorDTO	"Bad request"
//	@Failure		500		{object}	resp.ErrorDTO	"Internal server error"
//	@Router			/fingerprint/find/{linkId} [post]
func (h *fingerprintHandler) Find(ctx *gin.Context) {
	var link *domain.DirectLink
	var domainErr *pkg.ErrorS
	if linkId, ok := ctx.Params.Get("linkId"); ok {
		link, domainErr = h.FindByLink(ctx, linkId)
	} else {
		link, domainErr = h.FindByFP(ctx)
	}

	if domainErr != nil {
		ctx.JSON(domainErr.Status, resp.NewErrorDTO(domainErr.Error()))
		return
	}

	dto, parsed := resp.NewDirectLinkDTO(*link)
	if parsed != nil {
		ctx.JSON(parsed.Status, resp.NewErrorDTO(parsed.Error()))
	}
	ctx.JSON(http.StatusOK, dto)
}

func (h *fingerprintHandler) FindByLink(ctx *gin.Context, linkId domain.NanoID) (*domain.DirectLink, *pkg.ErrorS) {
	link, err := h.linkService.Find(ctx, linkId)
	return link, err
}

func (h *fingerprintHandler) FindByFP(ctx *gin.Context) (*domain.DirectLink, *pkg.ErrorS) {
	fpDTO := params.Fingerprint{}
	if err := ctx.ShouldBindWith(&fpDTO, binding.JSON); err != nil {
		return nil, pkg.NewBadRequestError("bad DTO")
	}

	ua := ctx.GetHeader(pkg.UserAgent)
	fpFields := fpDTO.ToFields(ctx.ClientIP(), domain.ParseUserAgent(ua))

	link, err := h.fingerprintService.Find(ctx, fpFields)
	return link, err
}
