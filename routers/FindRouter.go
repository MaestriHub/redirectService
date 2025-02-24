package routers

import (
	"net/http"
	"redirectServer/dto"
	"redirectServer/model"
	"redirectServer/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func FindRouter(app *gin.Engine) {

	find := app.Group("/find")
	{
		// TODO: хуевый нейминг, один с фингерпринтом, другой без
		find.POST("with-info", findFingerprint)
		find.GET("without-info", getDirectLinkPayload)
	}
}

func findFingerprint(c *gin.Context) {

	inputFingerprint := dto.FingerprintIOS{}
	if err := c.ShouldBindJSON(&inputFingerprint); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid json format")
	}
	ip := strings.Split(c.Request.RemoteAddr, ":")[0]
	fingerprint := inputFingerprint.ToFingerprint(ip, nil)
	findFingerprint := services.FindFingerprint(*fingerprint, DB)
	if findFingerprint == nil && inputFingerprint.UniversalLink == nil {
		c.String(http.StatusNotFound, "fingerPrintNotFound")
		return
	}
	if findFingerprint != nil && inputFingerprint.UniversalLink == nil {
		var directLink model.DirectLink
		DB.First(&directLink, "id = ?", findFingerprint.DirectLinkID)
		c.JSON(http.StatusOK, directLink)
		return
	}
	if inputFingerprint.UniversalLink != nil {
		directLink, err := model.ParseURL(*inputFingerprint.UniversalLink, DB)
		if err != nil {
			c.String(http.StatusBadRequest, "invalid link")
		}
		c.JSON(http.StatusOK, directLink)
		return
	}
	c.String(http.StatusNotFound, "fingerPrintNotFound")
}
func getDirectLinkPayload(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "Missing query parameter 'code'")
		return
	}
	var directLink model.DirectLink
	if DB.First(&directLink, "id = ?", code).Error != nil {
		c.String(http.StatusNotFound, "URL not found")
		return
	}
	c.JSON(http.StatusOK, directLink)

}
