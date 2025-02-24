package routers

import (
	"net/http"
	"redirectServer/dto"
	"redirectServer/model"
	"redirectServer/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func CollectRouter(app *gin.Engine) {

	collect := app.Group("/collect")
	{
		collect.POST("", collectData)
	}
}

func collectData(c *gin.Context) {
	data := dto.FingerprintJS{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid json format")
	}
	var fingerprintInfo model.Fingerprint = data.ToFingerprint()
	fingerprintInfo.IP = strings.Split(c.Request.RemoteAddr, ":")[0]

	var existingFingerprint *model.Fingerprint = services.FindFingerprint(fingerprintInfo, DB)
	if existingFingerprint != nil {
		if err := DB.Create(&fingerprintInfo).Error; err != nil {
			c.String(http.StatusInternalServerError, "Failed to create fingerprint info")
			return
		}
		DB.Delete(&existingFingerprint)
	} else {
		DB.Create(&fingerprintInfo)
	}

	c.String(http.StatusOK, "ok")
}
