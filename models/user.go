package models

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/louisun/vinki/pkg/utils"
)

const (
	// 未激活
	STATUS_NOT_ACTIVE = iota
	// 申请中
	STATUS_APPLYING
	// 激活
	STATUS_ACTIVE
	// 禁用
	STATUS_BANNED
)

// User 用户
type User struct {
	ID           uint64 `gorm:"primary_key"`
	Email        string `gorm:"type:varchar(100);unique_index"`
	NickName     string `gorm:"size:30"`
	Password     string `json:"-"`
	IsAdmin      bool
	Repos        string   `json:"-"`
	RepoNames    []string `gorm:"-" json:"repos"`
	ApplyMessage string   `json:"-"`
	Status       int
}

type UserApplyInfo struct {
	ID           uint64
	NickName     string
	ApplyMessage string
}

func (User) TableName() string {
	return "users"
}

// GetUserByID 根据 ID 获取 User
func GetUserByID(ID uint64) (User, error) {
	var user User
	result := DB.First(&user, ID)
	return user, result.Error
}

// GetAvailableUserByID 用ID获取可登录用户
func GetAvailableUserByID(ID interface{}) (User, error) {
	var user User
	result := DB.Where("status != ?", STATUS_BANNED).First(&user, ID)
	return user, result.Error
}

// GetUserByEmail 根据 Email 获取 User
func GetUserByEmail(email string) (User, error) {
	var user User
	result := DB.Where("email = ?", email).First(&user)
	return user, result.Error
}

// CreateUser 创建 User
func CreateUser(user *User) error {
	result := DB.Create(user)
	return result.Error
}

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	// 随机 salt 值：16位
	salt := utils.RandString(16)
	// 计算密码和 salt 组合的 SHA1 摘要
	hash := sha1.New()
	_, err := hash.Write([]byte(password + salt))
	bs := hex.EncodeToString(hash.Sum(nil))
	if err != nil {
		return err
	}
	// 设置密码为 salt 值和摘要的组合
	user.Password = salt + ":" + bs
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) (bool, error) {
	passwordEncrpt := strings.Split(user.Password, ":")
	if len(passwordEncrpt) != 2 {
		return false, errors.New("unknown password type")
	}
	// 生成摘要，判断密码是否匹配
	hash := sha1.New()
	_, err := hash.Write([]byte(password + passwordEncrpt[0]))
	bs := hex.EncodeToString(hash.Sum(nil))
	if err != nil {
		return false, err
	}
	// 判断哈希是否一致
	return bs == passwordEncrpt[1], nil
}

// SetStatus 设置用户状态
func SetStatus(userID uint64, status int) {
	DB.Model(&User{}).Where("id = ?", userID).Update("status", status)
}

// Update 更新用户
func UpdateUser(userID uint64, val map[string]interface{}) error {
	return DB.Model(&User{}).Where("id = ?", userID).Updates(val).Error
}

// UpdateUserAllowedRepos 授予用户可访问的仓库列表
func UpdateUserAllowedRepos(userID uint64, repos []string) error {
	b, err := json.Marshal(&repos)
	if err != nil {
		return err
	}
	result := DB.Model(&User{}).Update(map[string]interface{}{"id": userID, "repos": string(b), "status": STATUS_ACTIVE})
	return result.Error
}

// AfterFind 钩子：反序列化 Repo 列表
func (user *User) AfterFind() (err error) {
	if user.Repos != "" {
		err = json.Unmarshal([]byte(user.Repos), &user.RepoNames)
	}
	return err
}

// GetApplyingUserInfo 获取申请的用户信息
func GetApplyingUserInfo() ([]UserApplyInfo, error) {
	var users []UserApplyInfo
	result := DB.Model(&User{}).Where("status = ?", STATUS_APPLYING).Select("id, nick_name, apply_message").Scan(&users)
	return users, result.Error
}
