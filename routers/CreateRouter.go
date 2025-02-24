package routers

import (
	"net/http"
	"redirectServer/model"
	"redirectServer/model/payload"

	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func CreateRouter(app *gin.Engine) {

	create := app.Group("/create")
	{
		create.POST("salon", createSalonInvite)
		create.POST("employeer", createEmployeerInvite)
		create.POST("customer", createCustomerInvite)
		create.POST("master-to-salon", createMasterToSalonInvite)
	}
}

func createEmployeerInvite(c *gin.Context) {

	payloadObject := payload.Employeer{}
	if err := c.ShouldBindJSON(&payloadObject); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid json format")
	}
	id, err := gonanoid.New(8)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	directLink := model.DirectLink{ID: id, Event: string(model.EmployeerInvite)}
	directLink.SetPayload(payloadObject)

	if err := DB.Create(&directLink).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, directLink.ParseToURL())
}

func createMasterToSalonInvite(c *gin.Context) {
	payloadObject := payload.MasterToSalon{}
	if err := c.ShouldBindJSON(&payloadObject); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid json format")
	}

	id, err := gonanoid.New(8)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	directLink := model.DirectLink{ID: id, Event: string(model.MasterInviteToSalon)}
	directLink.SetPayload(payloadObject)

	if err := DB.Create(&directLink).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, directLink.ParseToURL())
}

func createSalonInvite(c *gin.Context) {
	payloadObject := payload.Salon{}
	if err := c.ShouldBindJSON(&payloadObject); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid json format")
	}

	id, err := gonanoid.New(8)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	directLink := model.DirectLink{ID: id, Event: string(model.SalonInvite)}
	directLink.SetPayload(payloadObject)

	if err := DB.Create(&directLink).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, directLink.ParseToURL())
}

func createCustomerInvite(c *gin.Context) {
	payloadObject := payload.Customer{}
	if err := c.ShouldBindJSON(&payloadObject); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid json format")
	}

	id, err := gonanoid.New(8)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	directLink := model.DirectLink{ID: id, Event: string(model.CustomerInvite)}
	directLink.SetPayload(payloadObject)

	if err := DB.Create(&directLink).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, directLink.ParseToURL())
}
