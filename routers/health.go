package routers

import (
	"github.com/gin-gonic/gin"
)

func InitHealthRouter(router *gin.Engine) {

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "operation-platform is running",
			"status":  "ok",
		})
	})
}
