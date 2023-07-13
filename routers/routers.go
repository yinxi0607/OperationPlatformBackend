package routers

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(router *gin.Engine) {
	InitDeploymentRouter(router)
	InitNamespacesRouter(router)
	InitPodsRouter(router)
	//initNodeRouter()
	//initIngressRouter()
}
