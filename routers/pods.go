package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/services"
)

func InitPodsRouter(router *gin.Engine) {
	podsServices := services.NewPodsService()
	podsRouter := router.Group("/pods")
	podsRouter.GET("/:namespace/:pod", func(c *gin.Context) {
		podsServices.GetPodInfo(c)
	})
	podsRouter.GET("/:namespace", func(c *gin.Context) {
		podsServices.GetAllPods(c)
	})
}
