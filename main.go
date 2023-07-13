package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"operation-platform/routers"
)

func main() {
	// 在Kubernetes集群内部部署时，使用以下代码创建一个in-cluster配置

	// 设置Gin路由
	router := gin.Default()
	routers.InitRouters(router)

	router.GET("/health", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "operation-platform is running"})
	})

	router.Run(":58180")
}
