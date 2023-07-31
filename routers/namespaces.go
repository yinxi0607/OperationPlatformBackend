package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/handlers"
)

func InitNamespacesRouter(router *gin.RouterGroup) {
	namespacesRouter := router.Group("/namespaces")
	namespacesRouter.GET("/", handlers.GetAllNamespaces)
}
