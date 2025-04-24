package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"redirectServer/internal/domain"
	"redirectServer/internal/service"
	"redirectServer/internal/transport/dto/params"
	"redirectServer/internal/transport/dto/resp"
	"redirectServer/pkg"
)

type FingerprintHandler interface {
	Create(ctx *gin.Context)
	Find(ctx *gin.Context)
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
	if err := ctx.ShouldBindBodyWithJSON(&fpDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, resp.NewErrorDTO(err.Error()))
		return
	}
	ua := ctx.GetHeader(pkg.UserAgent)

	linkId, exists := ctx.Params.Get("linkId")
	if !exists {
		ctx.JSON(http.StatusBadRequest, resp.NewErrorDTO("linkId is required"))
		return
	}

	fp := fpDTO.ToDomain(ctx.ClientIP(), domain.ParseUserAgent(ua), linkId)
	if err := h.fingerprintService.Create(ctx, fp); err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.NewErrorDTO(err.Error()))
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
//	@Success		200		{object}	domain.DirectLink
//	@Failure		400		{object}	resp.ErrorDTO	"Bad request"
//	@Failure		500		{object}	resp.ErrorDTO	"Internal server error"
//	@Router			/fingerprint/find/{linkId} [post]
func (h *fingerprintHandler) Find(ctx *gin.Context) {
	fpDTO := params.Fingerprint{}
	if err := ctx.ShouldBindBodyWithJSON(&fpDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, resp.NewErrorDTO(err.Error()))
		return
	}

	if linkId, exists := ctx.Params.Get("linkId"); exists {
		if link, err := h.linkService.Find(ctx, linkId); err == nil {
			ctx.JSON(http.StatusOK, link)
		}
	}

	ua := ctx.GetHeader(pkg.UserAgent)
	fpFields := fpDTO.ToFields(ctx.ClientIP(), domain.ParseUserAgent(ua))

	link, err := h.fingerprintService.Find(ctx, fpFields)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.NewErrorDTO(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, link)
}
