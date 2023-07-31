package routers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/middleware"
)

func InitRouters(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	api := router.Group("/api/v1")
	InitUsersRouter(api)
	InitHealthRouter(api)
	api.Use(middleware.TokenAuthMiddleware())
	{
		InitDeploymentRouter(api)
		InitNamespacesRouter(api)
		InitPodsRouter(api)
		InitConfigmapRouter(api)
	}

	//initNodeRouter()
	//initIngressRouter()
}
