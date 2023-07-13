package routers

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(router *gin.Engine) {
	InitDeploymentRouter(router)
	//initNamespaceRouter()
	//initPodRouter()
	//initNodeRouter()
	//initIngressRouter()
}
