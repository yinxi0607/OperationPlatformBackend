package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/handlers"
)

func InitDeploymentRouter(router *gin.Engine) {
	deploymentRouter := router.Group("/deployments")
	deploymentRouter.GET("/:namespace", handlers.GetAllDeployment)
	deploymentRouter.GET("/:namespace/:deployment", handlers.GetDeploymentPods)
	deploymentRouter.POST("/:namespace", handlers.PostDeployment)
}
