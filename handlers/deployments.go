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

func GetAllNSDeployment(c *gin.Context) {
	deploymentServices.GetAllNSDeployment(c)
}

func GetDeploymentPods(c *gin.Context) {
	deploymentServices.GetDeploymentPods(c)
}

func PostDeployment(c *gin.Context) {
	deploymentServices.PostDeployment(c)
}

func PutDeployment(c *gin.Context) {
	deploymentServices.PutDeployment(c)
}

func DeleteDeployment(c *gin.Context) {
	deploymentServices.DeleteDeployment(c)
}
