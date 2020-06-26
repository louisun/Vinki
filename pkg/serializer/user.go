package serializer

import (
	"github.com/louisun/vinki/models"
	"github.com/louisun/vinki/pkg/conf"
	"github.com/louisun/vinki/pkg/utils"
)

type UserDTO struct {
	ID       string `json:"id"`
	Email    string `json:"user_name"`
	NickName string `json:"nickname"`
	IsAdmin  bool   `json:"is_admin"`
	Status   int    `json:"status"`
}

// GetUnauthorizedResponse 检查登录
func GetUnauthorizedResponse() Response {
	return Response{
		Code: CodeUnauthorized,
		Msg:  "未认证",
	}
}

func CreateUserResponse(user *models.User) Response {
	data := UserDTO{
		ID:       utils.GenerateHash(user.ID, utils.UserID, conf.GlobalConfig.System.HashIDSalt),
		Email:    user.Email,
		NickName: user.NickName,
		IsAdmin:  user.IsAdmin,
		Status:   user.Status,
	}
	return CreateSuccessResponse(data, "")
}
