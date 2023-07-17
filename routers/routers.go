package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func InitRouters(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	InitHealthRouter(router)
	router.Use(TokenAuthMiddleware())
	InitDeploymentRouter(router)
	InitNamespacesRouter(router)
	InitPodsRouter(router)
	InitConfigmapRouter(router)
	//initNodeRouter()
	//initIngressRouter()
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API token required"})
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		// 这里你可以验证 token 的合法性，如查询数据库或缓存等
		if token != "CJ9eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}
