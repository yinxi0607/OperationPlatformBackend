package handlers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/services"
)

var configmapServices *services.ConfigmapService

func init() {
	configmapServices = services.NewConfigmapService()
}

func GetAllConfigmaps(c *gin.Context) {
	configmapServices.GetAllConfigmaps(c)
}

func GetConfigmap(c *gin.Context) {
	configmapServices.GetConfigmap(c)
}

func PostConfigmap(c *gin.Context) {
	configmapServices.PostConfigmap(c)
}

func PutConfigmap(c *gin.Context) {
	configmapServices.PutConfigmap(c)
}

func DeleteConfigmap(c *gin.Context) {
	configmapServices.DeleteConfigmap(c)
}
