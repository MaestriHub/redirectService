package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"redirectServer/internal/domain"
	"redirectServer/internal/service"
	"redirectServer/internal/transport/dto/params"
	"redirectServer/internal/transport/dto/resp"
)

type LinkHandler interface {
	CreateInviteEmployee(*gin.Context)
	CreateInviteClient(*gin.Context)
	CreateInviteSalon(*gin.Context)
}

type linkHandler struct {
	linkService service.LinkService
}

func NewLinkHandler(l service.LinkService) LinkHandler {
	return &linkHandler{linkService: l}
}

// CreateInviteEmployee godoc
//	@Summary		Create an employee invite
//	@Description	Generates a new invite link for an employee
//	@Tags			link
//	@Accept			json
//	@Produce		json
//	@Param			request	body		params.CreateEmployeeInviteLink	true	"Данные сотрудника"
//	@Success		201		{object}	domain.DirectLink
//	@Failure		400		{object}	resp.ErrorDTO	"Bad request"
//	@Failure		500		{object}	resp.ErrorDTO	"Internal server error"
//	@Router			/link/employee [post]
func (h *linkHandler) CreateInviteEmployee(ctx *gin.Context) {
	emp := params.CreateEmployeeInviteLink{}
	if err := ctx.ShouldBindBodyWithJSON(&emp); err != nil {
		ctx.JSON(http.StatusBadRequest, resp.NewErrorDTO(err.Error()))
		return
	}

	event := domain.NewEmployeeInviteEvent(emp.SalonId, emp.EmployeeId)
	link, err := domain.NewDirectLink(*event)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resp.NewErrorDTO(err.Error()))
	}

	if err := h.linkService.CreateInvite(ctx, link); err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.NewErrorDTO(err.Error())) // TODO:
		return
	}

	// TODO: convert to dto
	ctx.JSON(http.StatusCreated, link)
}

// CreateInviteSalon godoc
//	@Summary		Create to salon invite
//	@Description	Generates a new invite link for salon
//	@Tags			link
//	@Accept			json
//	@Produce		json
//	@Param			request	body		params.CreateSalonInviteLink	true	"Данные салона"
//	@Success		201		{object}	domain.DirectLink
//	@Failure		400		{object}	resp.ErrorDTO	"Bad request"
//	@Failure		500		{object}	resp.ErrorDTO	"Internal server error"
//	@Router			/link/salon [post]
func (h *linkHandler) CreateInviteSalon(ctx *gin.Context) {
	salon := params.CreateSalonInviteLink{}
	if err := ctx.ShouldBindBodyWithJSON(&salon); err != nil {
		ctx.JSON(http.StatusBadRequest, resp.NewErrorDTO(err.Error()))
		return
	}

	event := domain.NewSalonInviteEvent(salon.SalonId)
	link, err := domain.NewDirectLink(*event)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resp.NewErrorDTO(err.Error()))
	}

	if err := h.linkService.CreateInvite(ctx, link); err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.NewErrorDTO(err.Error())) // TODO:
		return
	}

	// TODO: convert to dto
	ctx.JSON(http.StatusCreated, link)
}

// CreateInviteClient godoc
//	@Summary		Create to client invite
//	@Description	Generates a new invite link for client
//	@Tags			link
//	@Accept			json
//	@Produce		json
//	@Param			request	body		params.CreateClientInviteLink	true	"Данные клиента"
//	@Success		201		{object}	domain.DirectLink
//	@Failure		400		{object}	resp.ErrorDTO	"Bad request"
//	@Failure		500		{object}	resp.ErrorDTO	"Internal server error"
//	@Router			/link/client [post]
func (h *linkHandler) CreateInviteClient(ctx *gin.Context) {
	customer := params.CreateClientInviteLink{}
	if err := ctx.ShouldBindBodyWithJSON(&customer); err != nil {
		ctx.JSON(http.StatusBadRequest, resp.NewErrorDTO(err.Error()))
		return
	}

	event := domain.NewClientInviteEvent(customer.ClientId, customer.EmployeeId)
	link, err := domain.NewDirectLink(*event)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, resp.NewErrorDTO(err.Error()))
	}

	if err := h.linkService.CreateInvite(ctx, link); err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.NewErrorDTO(err.Error())) // TODO:
		return
	}

	// TODO: convert to dto
	ctx.JSON(http.StatusCreated, link)
}
