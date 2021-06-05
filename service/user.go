package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/louisun/vinki/models"
	"github.com/louisun/vinki/pkg/serializer"
	"github.com/louisun/vinki/pkg/session"
)

const InvitationCode = "VINKI_BY_RENZO"

// UserRegisterRequest 用户注册服务
type UserRegisterRequest struct {
	UserName       string `form:"userName" json:"userName" binding:"required,email"`
	NickName       string `form:"nickName" json:"nickName" binding:"required"`
	Password       string `form:"password" json:"password" binding:"required,min=4,max=64"`
	InvitationCode string `form:"invitationCode" json:"invitationCode" binding:"required"`
}

// UserLoginRequest 用户登录服务
type UserLoginRequest struct {
	UserName string `form:"userName" json:"userName" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=4,max=64"`
}

// UserLogoutRequest 用户登出服务
type UserLogoutRequest struct{}

// GetApplications
type GetApplicationsRequest struct{}

// ActivateUserRequest 激活用户服务
type ActivateUserRequest struct {
	UserID uint64   `form:"userID" json:"userID" binding:"required"`
	Repos  []string `form:"repos" json:"repos" binding:"required"`
}

type RejectUserRequest struct {
	UserID uint64 `form:"userID" json:"userID" binding:"required"`
}

type BanUserRequest struct {
	UserID uint64 `form:"userID" json:"userID" binding:"required"`
}

// ApplyForActivateRequest 用户申请激活服务
type ApplyForActivateRequest struct {
	Message string `form:"message" json:"message"`
}

// UserResetRequest 密码重设服务
type UserResetRequest struct {
	Password    string `form:"password" json:"password" binding:"required,min=4,max=64"`
	NewPassword string `form:"newPassword" json:"newPassword" binding:"required,min=4,max=64"`
}

// UserSetCurrentRepo 设置当前仓库
type UserSetCurrentRepo struct {
	CurrentRepo string `form:"currentRepo" json:"currentRepo" binding:"required"`
}

// Register 用户注册
func (service *UserRegisterRequest) Register(c *gin.Context) serializer.Response {
	if service.InvitationCode != InvitationCode {
		return serializer.CreateParamErrorResponse(errors.New("invitation code incorrect"))
	}
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
func (service *UserLoginRequest) Login(c *gin.Context) serializer.Response {
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
func (service *UserLogoutRequest) Logout(c *gin.Context) serializer.Response {
	session.DeleteSession(c, "user_id")
	return serializer.CreateSuccessResponse("", "登出成功")
}

// ResetPassword 重置用户密码
func (service *UserResetRequest) ResetPassword(c *gin.Context, user *models.User) serializer.Response {
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
func (service *ApplyForActivateRequest) ApplyForActivate(c *gin.Context, user *models.User) serializer.Response {
	err := models.UpdateUser(user.ID, map[string]interface{}{"status": models.STATUS_APPLYING, "apply_message": service.Message})
	if err != nil {
		return serializer.CreateDBErrorResponse("用户申请激活权限失败", err)
	}

	return serializer.CreateSuccessResponse("", "已向管理员申请激活，请耐心等待")
}

// GetApplications 管理员获取用户申请列表
func (service *GetApplicationsRequest) GetApplications() serializer.Response {
	applyInfos, err := models.GetApplyingUserInfo()
	if err != nil {
		return serializer.CreateDBErrorResponse("获取用户激活申请列表失败", err)
	}

	return serializer.CreateSuccessResponse(applyInfos, "")
}

// ActivateUser 管理员激活用户：授予指定仓库访问权限
func (service *ActivateUserRequest) ActivateUser() serializer.Response {
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
func (service *RejectUserRequest) RejectUser() serializer.Response {
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
func (service *BanUserRequest) BanUser() serializer.Response {
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

func (service *UserSetCurrentRepo) SetCurrentRepo(userID uint64) serializer.Response {
	err := models.SetCurrentRepo(userID, service.CurrentRepo)
	if err != nil {
		return serializer.CreateInternalErrorResponse("设置当前仓库失败", err)
	}

	return serializer.CreateSuccessResponse("", "设置当前仓库成功")

}

func GetCurrentRepo(userID uint64) serializer.Response {
	repo, err := models.GetCurrentRepo(userID)
	if err != nil {
		return serializer.CreateInternalErrorResponse("获取当前仓库失败", err)
	}

	return serializer.CreateSuccessResponse(repo, "获取当前仓库成功")
}
