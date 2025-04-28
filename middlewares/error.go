package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
	"redirectServer/internal/transport/dto/resp"
)

// ErrorHandlerMiddleware По сути недостижимо. Если случайно просрем ошибку то это отработает.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf("Error occurred: %v", err)
			}

			c.JSON(500, resp.NewErrorDTO("Internal Server Error"))
		}
	}
}
