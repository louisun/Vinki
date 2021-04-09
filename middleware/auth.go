package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/louisun/vinki/models"
	"github.com/louisun/vinki/pkg/serializer"
	"github.com/louisun/vinki/pkg/utils"
)

// RequireAuth 需要登录
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userCtx, _ := c.Get("user"); userCtx != nil {
			if user, ok := userCtx.(*models.User); ok {
				if user.Status == models.STATUS_BANNED {
					c.JSON(200, serializer.CreateErrorResponse(serializer.CodeForbidden, "该用户已被禁用", nil))
					return
				}

				c.Next()

				return
			}
		}

		c.JSON(200, serializer.GetUnauthorizedResponse())

		c.Abort()
	}
}

// InitCurrentUserIfExists 设置登录用户
func InitCurrentUserIfExists() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// Login 会在 session 中 Set user_id
		uid := session.Get("user_id")
		if uid != nil {
			user, err := models.GetAvailableUserByID(uid)
			if err == nil {
				c.Set("user", &user)
			}
		}

		c.Next()
	}
}

// RequireAdmin 判断用户是否是管理员
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")
		if !user.(*models.User).IsAdmin {
			c.JSON(200, serializer.CreateErrorResponse(serializer.CodeAdminRequired, "非管理员无法操作", nil))
			c.Abort()
		}

		c.Next()
	}
}

// CheckPermission 判断用户是否已激活，确认对应访问权限
func CheckPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx, _ := c.Get("user")
		user := userCtx.(*models.User)
		// 管理员直接允许
		if user.IsAdmin {
			c.Next()
			return
		}

		if user.Status == models.STATUS_NOT_ACTIVE {
			c.JSON(200, serializer.CreateErrorResponse(serializer.CodeActiveRequired, "账号需要激活，请向管理员申请", nil))
			c.Abort()

			return
		}

		if user.Status == models.STATUS_APPLYING {
			c.JSON(200, serializer.CreateErrorResponse(serializer.CodeActiveRequired, "已申请访问，请耐心等待", nil))
			c.Abort()

			return
		}
		// 判断是否在访问操作权限之外的 Repo
		repoName := c.Query("repoName")
		if repoName != "" {
			if !utils.IsInList(user.RepoNames, repoName) {
				c.JSON(200, serializer.CreateErrorResponse(serializer.CodeForbidden, "无权限访问", nil))
				c.Abort()

				return
			}
		}

		c.Next()
	}
}
