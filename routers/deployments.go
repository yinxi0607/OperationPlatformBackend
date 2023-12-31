package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/handlers"
)

func InitDeploymentRouter(router *gin.RouterGroup) {
	deploymentRouter := router.Group("/deployments")
	deploymentRouter.GET("/:namespace", handlers.GetAllDeployment)
	deploymentRouter.GET("/", handlers.GetAllNSDeployment)
	deploymentRouter.GET("/:namespace/:deployment/pods", handlers.GetDeploymentPods)
	deploymentRouter.GET("/:namespace/:deployment/info", handlers.GetDeploymentInfo)
	deploymentRouter.POST("/", handlers.PostDeployment)
	deploymentRouter.PUT("/", handlers.PutDeployment)
	deploymentRouter.DELETE("/", handlers.DeleteDeployment)
}
