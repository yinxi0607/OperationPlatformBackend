package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"operation-platform/handlers"
	"strings"
)

func InitUsersRouter(router *gin.Engine) {
	usersRouter := router.Group("/users")
	usersRouter.GET("/login", handlers.GetUserLogin)
	usersRouter.GET("/callback", handlers.GetUserLoginCallback)
	usersRouter.POST("/logout", handlers.PostUserLogout)
	usersRouter.GET("/info/:userId", handlers.GetUserInfo)
	usersRouter.GET("/list", handlers.GetAllUsers)
}

func TokenAuthMiddleware1() gin.HandlerFunc {
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
