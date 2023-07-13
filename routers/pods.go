package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/handlers"
)

func InitPodsRouter(router *gin.Engine) {
	podsRouter := router.Group("/pods")
	podsRouter.GET("/:namespace/:pod", handlers.GetPodInfo)
	podsRouter.GET("/:namespace", handlers.GetAllPods)
}
