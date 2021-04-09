package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/louisun/vinki/models"
	"github.com/louisun/vinki/pkg/serializer"
	"github.com/louisun/vinki/service"
)

// GetCurrentUserFromCtx 获取当前用户
func GetCurrentUserFromCtx(c *gin.Context) *models.User {
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*models.User); ok {
			return u
		}
	}

	return nil
}

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	var s service.UserRegisterService
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(200, serializer.CreateParamErrorResponse(err))
	} else {
		res := s.Register(c)
		c.JSON(200, res)
	}
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var s service.UserLoginService
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(200, serializer.CreateParamErrorResponse(err))
	} else {
		res := s.Login(c)
		c.JSON(200, res)
	}
}

// UserLogout 用户登出
func UserLogout(c *gin.Context) {
	var s service.UserLogoutService
	res := s.Logout(c)
	c.JSON(200, res)
}

// UserPasswordReset 用户重置密码
func UserResetPassword(c *gin.Context) {
	var s service.UserResetService
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(200, serializer.CreateParamErrorResponse(err))
	} else {
		var user *models.User
		userCtx, ok := c.Get("user")
		if !ok {
			c.JSON(200, serializer.GetUnauthorizedResponse())
			return
		}
		user = userCtx.(*models.User)
		// 校验并重置密码
		res := s.ResetPassword(c, user)
		c.JSON(200, res)
	}
}

// GetApplications 管理员查看用户申请
func GetApplications(c *gin.Context) {
	var s service.GetApplicationsService
	res := s.GetApplications()
	c.JSON(200, res)
}

// ActivateUser 激活用户
func ActivateUser(c *gin.Context) {
	var s service.ActivateUserService
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(200, serializer.CreateParamErrorResponse(err))
	} else {
		res := s.ActivateUser()
		c.JSON(200, res)
	}
}

// RejectUserApplication 拒绝用户申请
func RejectUserApplication(c *gin.Context) {
	var s service.RejectUserService
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(200, serializer.CreateParamErrorResponse(err))
	} else {
		res := s.RejectUser()
		c.JSON(200, res)
	}
}

// BanUser 封禁用户
func BanUser(c *gin.Context) {
	var s service.BanUserService
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(200, serializer.CreateParamErrorResponse(err))
	} else {
		res := s.BanUser()
		c.JSON(200, res)
	}
}

// ApplyForActivate 向管理员申请激活
func ApplyForActivate(c *gin.Context) {
	// 实际就是修改自己的 user 状态和申请 message
	var s service.ApplyForActivateService
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(200, serializer.CreateParamErrorResponse(err))
	} else {
		var user *models.User
		userCtx, ok := c.Get("user")
		if !ok {
			c.JSON(200, serializer.GetUnauthorizedResponse())
			return
		}
		user = userCtx.(*models.User)
		res := s.ApplyForActivate(c, user)
		c.JSON(200, res)
	}
}

// SetCurrentRepo 设置当前的仓库
func SetCurrentRepo(c *gin.Context) {
	var s service.UserSetCurrentRepo
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(200, serializer.CreateParamErrorResponse(err))
	} else {
		var user *models.User
		userCtx, ok := c.Get("user")
		if !ok {
			c.JSON(200, serializer.GetUnauthorizedResponse())
			return
		}

		user = userCtx.(*models.User)
		res := s.SetCurrentRepo(user.ID)

		c.JSON(200, res)
	}
}

// GetCurrentRepo 获取当前仓库
func GetCurrentRepo(c *gin.Context) {
	var user *models.User
	userCtx, ok := c.Get("user")
	if !ok {
		c.JSON(200, serializer.GetUnauthorizedResponse())
		return
	}

	user = userCtx.(*models.User)
	res := service.GetCurrentRepo(user.ID)

	c.JSON(200, res)
}
