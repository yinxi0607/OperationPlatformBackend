package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/services"
)

func InitDeploymentRouter(router *gin.Engine) {
	deploymentServices := services.NewDefault()

	router.GET("/pods/:namespace/:deployment", func(c *gin.Context) {
		deploymentServices.GetDeploymentPods(c)
	})
}
