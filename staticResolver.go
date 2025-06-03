package redirectService

import (
	"embed"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed static
var staticFolder embed.FS

func EmbedStatic(r *gin.Engine) {
	fs, err := static.EmbedFolder(staticFolder, "static")
	if err != nil {
		panic(err)
	}
	r.Use(static.Serve("/", fs))
}
