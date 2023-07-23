package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/handlers"
)

func InitDeploymentRouter(router *gin.Engine) {
	deploymentRouter := router.Group("/deployments")
	deploymentRouter.GET("/:namespace", handlers.GetAllDeployment)
	deploymentRouter.GET("/", handlers.GetAllNSDeployment)
	deploymentRouter.GET("/:namespace/:deployment/pods", handlers.GetDeploymentPods)
	deploymentRouter.GET("/:namespace/:deployment/info", handlers.GetDeploymentPods)
	deploymentRouter.POST("/", handlers.PostDeployment)
	deploymentRouter.PUT("/", handlers.PutDeployment)
	deploymentRouter.DELETE("/", handlers.DeleteDeployment)
}
