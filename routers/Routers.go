package routers

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB
var logger *slog.Logger

func InitRouters(db *gorm.DB, app *gin.Engine) {
	DB = db
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	app.With(CollectRouter)
	app.With(CreateRouter)
	app.With(FindRouter)
	app.With(ServeRouter)
	app.StaticFS("/static", http.Dir("static"))
	logger.Info("Routers initialized")

}
