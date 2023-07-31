package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/handlers"
)

func InitConfigmapRouter(router *gin.RouterGroup) {
	configmapRouter := router.Group("/configmaps")
	configmapRouter.GET("/", handlers.GetAllNSConfigmaps)
	configmapRouter.GET("/:namespace", handlers.GetAllConfigmaps)
	configmapRouter.POST("/", handlers.PostConfigmap)
	configmapRouter.PUT("/", handlers.PutConfigmap)
	configmapRouter.DELETE("/", handlers.DeleteConfigmap)
	configmapRouter.GET("/:namespace/:configmapName", handlers.GetConfigmap)
}
