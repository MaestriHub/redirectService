package routers

import (
	"net/http"
	"os"
	"redirectServer/model"
	"strings"

	"github.com/gin-gonic/gin"
)

func ServeRouter(app *gin.Engine) {
	app.GET("/:code", serveHTML)
}

func serveHTML(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.String(http.StatusBadRequest, "Missing query parameter 'code'")
		return
	}
	var directLink model.DirectLink
	result := DB.First(&directLink, "id = ?", code)
	if result.Error != nil {
		c.AbortWithError(http.StatusNotFound, result.Error)
		return
	}
	directLink.Сlicks++
	DB.Save(&directLink)

	userAgent := c.Request.Header.Get("User-Agent")
	var response string
	var err error
	switch {
	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(model.MasterInviteToSalon):
		response, err = getMasterToSalonInviteIOS(directLink)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

	case strings.Contains(userAgent, "Android") && directLink.Event == string(model.MasterInviteToSalon):
		response, err = getMasterToSalonInviteAndroid(directLink)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(model.SalonInvite):
		response, err = getSalonInviteIOS(directLink)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

	case strings.Contains(userAgent, "Android") && directLink.Event == string(model.SalonInvite):
		response, err = getSalonInviteAndroid(directLink)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(model.CustomerInvite):
		response, err = getCustomerInviteIOS(directLink)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	case strings.Contains(userAgent, "Android") && directLink.Event == string(model.CustomerInvite):
		response, err = getCustomerInviteAndroid(directLink)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	case (strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad")) && directLink.Event == string(model.EmployeerInvite):
		response, err = getEmployeerInviteIOS(directLink)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	case strings.Contains(userAgent, "Android") && directLink.Event == string(model.EmployeerInvite):
		response, err = getEmployeerInviteAndroid(directLink)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

	default:
		if directLink.Event == string(model.MasterInviteToSalon) {
			response, err = getMasterToSalonInviteWeb(directLink)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		} else if directLink.Event == string(model.SalonInvite) {
			response, err = getSalonInviteWeb(directLink)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		} else if directLink.Event == string(model.CustomerInvite) {
			response, err = getCustomerInviteWeb(directLink)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		} else if directLink.Event == string(model.EmployeerInvite) {
			response, err = getEmployeerInviteWeb(directLink)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		} else {
			c.String(http.StatusInternalServerError, "Unsupported device")
		}
	}
	c.Data(http.StatusOK, "text/html", []byte(response))
}

func getMasterToSalonInviteIOS(
	directLink model.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/MasterToSalon/IOS/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadMasterToSalon(DB)
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", name, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", description, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getMasterToSalonInviteAndroid(
	directLink model.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/MasterToSalon/Android/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadMasterToSalon(DB)
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", name, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", description, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getMasterToSalonInviteWeb(
	directLink model.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/MasterToSalon/Web/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadMasterToSalon(DB)
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", name, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", description, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getSalonInviteIOS(
	directLink model.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/SalonInvite/IOS/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadSalon(DB)
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", name, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", description, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getSalonInviteAndroid(
	directLink model.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/SalonInvite/Android/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadSalon(DB)
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", name, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", description, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getSalonInviteWeb(
	directLink model.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/SalonInvite/Web/Screen.html")
	if err != nil {
		return "", err
	}
	name, description, err := directLink.GetPayloadSalon(DB)
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Name}}", name, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.Description}}", description, -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getCustomerInviteIOS(
	directLink model.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/CustomerInvite/IOS/Screen.html")
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getCustomerInviteAndroid(
	directLink model.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/CustomerInvite/Android/Screen.html")
	if err != nil {
		return "", err
	}

	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getCustomerInviteWeb(
	directLink model.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/CustomerInvite/Web/Screen.html")
	if err != nil {
		return "", err
	}

	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getEmployeerInviteIOS(
	directLink model.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/EmployeerInvite/IOS/Screen.html")
	if err != nil {
		return "", err
	}
	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getEmployeerInviteAndroid(
	directLink model.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/EmployeerInvite/Android/Screen.html")
	if err != nil {
		return "", err
	}

	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}

func getEmployeerInviteWeb(
	directLink model.DirectLink,
) (string, error) {
	htmlFile, err := os.ReadFile("static/EmployeerInvite/Web/Screen.html")
	if err != nil {
		return "", err
	}

	modifiedHTML := strings.Replace(string(htmlFile), "{{.DynamicUniversalLink}}", directLink.ParseToURL(), -1)
	modifiedHTML = strings.Replace(string(modifiedHTML), "{{.LinkID}}", directLink.ID, -1)

	return modifiedHTML, nil
}
