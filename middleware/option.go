package middleware

import "github.com/gin-gonic/gin"

func IsFunctionEnabled(functionKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO 站点配置相关的功能启用状态
		c.Next()
	}
}
