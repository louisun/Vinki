package controllers

import "github.com/gin-gonic/gin"

// UserLogin 用户登录
func Login(c *gin.Context) {
	// 用户名、密码、验证码
	//var service user.UserLoginService
	//if err := c.ShouldBindJSON(&service); err == nil {
	//	res := service.Login(c)
	//	c.JSON(200, res)
	//} else {
	//	c.JSON(200, ErrorResponse(err))
	//}
}

// UserRegister 用户注册
func Register(c *gin.Context) {
	//var service user.UserRegisterService
	//if err := c.ShouldBindJSON(&service); err == nil {
	//	res := service.Register(c)
	//	c.JSON(200, res)
	//} else {
	//	c.JSON(200, ErrorResponse(err))
	//}
}

// SetVisibility 设置用户对 repo, tag, article 的可见性
func SetVisibility(c *gin.Context) {

}
