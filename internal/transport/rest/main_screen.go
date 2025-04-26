package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"redirectServer/internal/domain"
	"redirectServer/internal/service"
	"redirectServer/internal/transport/dto/resp"
	"redirectServer/pkg"
)

type MainScreenHandler interface {
	MainScreen(*gin.Context)
}

type mainScreenHandler struct {
	linkService   service.LinkService
	renderService service.RenderService
}

func NewMainScreenHandler(l service.LinkService, r service.RenderService) MainScreenHandler {
	return &mainScreenHandler{linkService: l, renderService: r}
}

// MainScreen godoc
//	@Summary		Main web window
//	@Description	web window with button on application
//	@Tags			html
//	@Accept			json
//	@Produce		json
//	@Param			User-Agent	header		string	true	"Юзер агент пользователя. ex: Android"
//	@Param			linkId		path		string	true	"Идентификатор (NanoID)"
//	@Success		200			{object}	string
//	@Failure		400			{object}	resp.ErrorDTO	"Bad request"
//	@Failure		500			{object}	resp.ErrorDTO	"Internal server error"
//	@Router			/{linkId} [get]
func (h *mainScreenHandler) MainScreen(ctx *gin.Context) {
	linkId, exists := ctx.Params.Get("linkId")
	if !exists {
		ctx.JSON(http.StatusBadRequest, resp.NewErrorDTO("linkId is required"))
		return
	}

	link, err := h.linkService.LinkTap(ctx, linkId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, resp.NewErrorDTO(err.Error()))
		return
	}

	ua := domain.ParseUserAgent(ctx.GetHeader(pkg.UserAgent))

	html, err := h.renderService.RenderMain(ctx, link, ua)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.NewErrorDTO(err.Error()))
		return
	}

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
