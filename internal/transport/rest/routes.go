package rest

import (
	"github.com/gin-gonic/gin"
)

func MapLinkRoutes(e *gin.Engine, h LinkHandler) {
	createLink := e.Group("/link")
	{
		createLink.POST("/salon", h.CreateInviteSalon)
		createLink.POST("/employee", h.CreateInviteEmployee)
		createLink.POST("/client", h.CreateInviteClient)
	}
}

func MapFPRoutes(e *gin.Engine, h FingerprintHandler) {
	fp := e.Group("/fingerprint")
	{
		fp.POST("/:linkId", h.Create)
		// MARK: Gin not support optional params
		fp.POST("/find/:linkId", h.Find)
		fp.POST("/find", h.Find)
	}
}

func MapMainScreenRoutes(e *gin.Engine, h MainScreenHandler) {
	e.GET("/:linkId", h.MainScreen)
}
