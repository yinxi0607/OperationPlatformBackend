package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/services"
)

func InitNamespacesRouter(router *gin.Engine) {
	namespacesServices := services.NewNamespacesService()
	namespacesRouter := router.Group("/namespaces")
	namespacesRouter.GET("/", func(c *gin.Context) {
		namespacesServices.GetAllNamespaces(c)
	})
}
