package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/middleware"
)

func InitRouters(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	InitUsersRouter(router)
	InitHealthRouter(router)
	router.Use(middleware.TokenAuthMiddleware())
	InitDeploymentRouter(router)
	InitNamespacesRouter(router)
	InitPodsRouter(router)
	InitConfigmapRouter(router)

	//initNodeRouter()
	//initIngressRouter()
}
