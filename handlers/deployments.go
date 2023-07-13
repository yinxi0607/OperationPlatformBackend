package handlers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/services"
)

var deploymentServices *services.DeploymentService

func init() {
	deploymentServices = services.NewDeploymentService()
}

func GetAllDeployment(c *gin.Context) {
	deploymentServices.GetAllDeployment(c)
}

func GetDeploymentPods(c *gin.Context) {
	deploymentServices.GetDeploymentPods(c)
}
