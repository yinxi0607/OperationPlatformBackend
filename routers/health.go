package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/utils"
)

func InitHealthRouter(router *gin.Engine) {

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, utils.Response{
			Code:    0,
			Message: "success",
			Data:    "operation-platform is running",
		})
	})
}
