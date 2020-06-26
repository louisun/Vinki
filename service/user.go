package service

import (
	"github.com/gin-gonic/gin"
	"github.com/louisun/vinki/models"
	"github.com/louisun/vinki/pkg/serializer"
	"github.com/louisun/vinki/pkg/session"
)

const INVITATION_CODE = "VINKI_BY_RENZO"

// UserRegisterService 用户注册服务
type UserRegisterService struct {
	UserName       string `form:"userName" json:"userName" binding:"required,email"`
	NickName       string `form:"nickName" json:"nickName" binding:"required"`
	Password       string `form:"password" json:"password" binding:"required,min=4,max=64"`
	InvitationCode string `form:"invitationCode" json:"invitationCode" binding:"required"`
}

// UserLoginService 用户登录服务
type UserLoginService struct {
	UserName string `form:"userName" json:"userName" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=4,max=64"`
}

// UserLogoutService 用户登出服务
type UserLogoutService struct{}

// GetApplications
type GetApplicationsService struct{}

// ActivateUserService 激活用户服务
type ActivateUserService struct {
	UserID uint64   `form:"userID" json:"userID" binding:"required"`
	Repos  []string `form:"repos" json:"repos" binding:"required"`
}
type RejectUserService struct {
	UserID uint64 `form:"userID" json:"userID" binding:"required"`
}

type BanUserService struct {
	UserID uint64 `form:"userID" json:"userID" binding:"required"`
}

// ApplyForActivateService 用户申请激活服务
type ApplyForActivateService struct {
	Message string `form:"message" json:"message"`
}

// UserResetService 密码重设服务
type UserResetService struct {
	Password    string `form:"password" json:"password" binding:"required,min=4,max=64"`
	NewPassword string `form:"newPassword" json:"newPassword" binding:"required,min=4,max=64"`
}

// Register 用户注册
func (service *UserRegisterService) Register(c *gin.Context) serializer.Response {
	user := models.User{
		Email:    service.UserName,
		NickName: service.NickName,
	}
	_ = user.SetPassword(service.Password)
	user.Status = models.STATUS_NOT_ACTIVE
	if err := models.CreateUser(&user); err != nil {
		return serializer.CreateDBErrorResponse("", err)
	}
	return serializer.CreateSuccessResponse("", "注册成功")
}

// Login 用户登录
func (service *UserLoginService) Login(c *gin.Context) serializer.Response {
	user, err := models.GetUserByEmail(service.UserName)
	if err != nil {
		return serializer.CreateErrorResponse(401, "用户邮箱或密码错误", err)
	}
	if passwordCorrect, err := user.CheckPassword(service.Password); !passwordCorrect {
		return serializer.CreateErrorResponse(serializer.CodeUnauthorized, "用户邮箱或密码错误", err)
	}
	if user.Status == models.STATUS_BANNED {
		return serializer.CreateErrorResponse(serializer.CodeForbidden, "该用户已被封禁", err)
	}
	session.SetSession(c, map[string]interface{}{
		"user_id": user.ID,
	})
	return serializer.CreateUserResponse(&user)
}

// Logout 用户登出
func (service *UserLogoutService) Logout(c *gin.Context) serializer.Response {
	session.DeleteSession(c, "user_id")
	return serializer.CreateSuccessResponse("", "登出成功")
}

// ResetPassword 重置用户密码
func (service *UserResetService) ResetPassword(c *gin.Context, user *models.User) serializer.Response {
	// 验证旧密码
	if passwordCorrect, err := user.CheckPassword(service.Password); !passwordCorrect {
		return serializer.CreateErrorResponse(serializer.CodeUnauthorized, "当前用户密码错误，无法重置密码", err)
	}
	// 设置新密码
	if err := user.SetPassword(service.NewPassword); err != nil {
		return serializer.CreateErrorResponse(200, "重置密码失败", err)
	}
	if err := models.UpdateUser(user.ID, map[string]interface{}{"password": user.Password}); err != nil {
		return serializer.CreateDBErrorResponse("重置密码失败", err)
	}

	return serializer.CreateSuccessResponse("", "重置密码成功")
}

// ApplyForActivate 用户向管理员申请激活
func (service *ApplyForActivateService) ApplyForActivate(c *gin.Context, user *models.User) serializer.Response {
	err := models.UpdateUser(user.ID, map[string]interface{}{"status": models.STATUS_APPLYING, "apply_message": service.Message})
	if err != nil {
		return serializer.CreateDBErrorResponse("用户申请激活权限失败", err)
	}
	return serializer.CreateSuccessResponse("", "已向管理员申请激活，请耐心等待")
}

// GetApplications 管理员获取用户申请列表
func (service *GetApplicationsService) GetApplications() serializer.Response {
	applyInfos, err := models.GetApplyingUserInfo()
	if err != nil {
		return serializer.CreateDBErrorResponse("获取用户激活申请列表失败", err)
	}
	return serializer.CreateSuccessResponse(applyInfos, "")
}

// ActivateUser 管理员激活用户：授予指定仓库访问权限
func (service *ActivateUserService) ActivateUser() serializer.Response {
	user, err := models.GetUserByID(service.UserID)
	if err != nil {
		return serializer.CreateDBErrorResponse("获取用户失败", err)
	}
	if user.Status != models.STATUS_APPLYING {
		return serializer.CreateErrorResponse(serializer.CodeConditionNotMeet, "激活用户权限失败：该用户非申请状态", nil)
	}
	err = models.UpdateUserAllowedRepos(user.ID, service.Repos)
	if err != nil {
		return serializer.CreateDBErrorResponse("激活用户权限失败", err)
	}
	return serializer.CreateSuccessResponse("", "激活用户权限成功")
}

// RejectUser 管理员拒绝用户申请：取消申请状态
func (service *RejectUserService) RejectUser() serializer.Response {
	user, err := models.GetUserByID(service.UserID)
	if err != nil {
		return serializer.CreateDBErrorResponse("获取用户失败", err)
	}
	if user.Status != models.STATUS_APPLYING {
		return serializer.CreateErrorResponse(serializer.CodeConditionNotMeet, "拒绝用户失败：该用户非申请状态", nil)
	}
	err = models.UpdateUser(user.ID, map[string]interface{}{"status": models.STATUS_NOT_ACTIVE, "apply_message": ""})
	if err != nil {
		return serializer.CreateDBErrorResponse("拒绝用户失败", err)
	}
	return serializer.CreateSuccessResponse("", "拒绝用户成功")
}

// BanUser 管理员封禁用户
func (service *BanUserService) BanUser() serializer.Response {
	user, err := models.GetUserByID(service.UserID)
	if err != nil {
		return serializer.CreateDBErrorResponse("获取用户失败", err)
	}
	err = models.UpdateUser(user.ID, map[string]interface{}{"status": models.STATUS_BANNED, "apply_message": ""})
	if err != nil {
		return serializer.CreateDBErrorResponse("封禁用户失败", err)
	}
	return serializer.CreateSuccessResponse("", "封禁用户成功")
}
