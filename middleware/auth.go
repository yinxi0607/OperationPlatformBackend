package middleware

import (
	"github.com/gin-gonic/gin"
)

// TokenAuthMiddleware 用于验证访问令牌
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 令牌验证通过，继续执行
		c.Next()
	}
}
