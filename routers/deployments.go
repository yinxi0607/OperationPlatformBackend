package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/services"
)

func InitDeploymentRouter(router *gin.Engine) {
	deploymentServices := services.NewDeploymentService()
	deploymentRouter := router.Group("/deployments")
	deploymentRouter.GET("/:namespace/:deployment", func(c *gin.Context) {
		deploymentServices.GetDeploymentPods(c)
	})
}
