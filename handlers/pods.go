package handlers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/services"
)

var podsServices *services.PodsService

func init() {
	podsServices = services.NewPodsService()
}

func GetPodInfo(c *gin.Context) {
	podsServices.GetPodInfo(c)
}

func GetAllPods(c *gin.Context) {
	podsServices.GetAllPods(c)
}

func GetAllNSPods(c *gin.Context) {
	podsServices.GetAllNSPods(c)
}
