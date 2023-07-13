package handlers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/services"
)

var namespacesServices *services.NamespacesService

func init() {
	namespacesServices = services.NewNamespacesService()
}

func GetAllNamespaces(c *gin.Context) {
	namespacesServices.GetAllNamespaces(c)
}
