package models

// User 用户
type User struct {
	ID       uint64 `gorm:"primary_key"`
	Email    string `gorm:"type:varchar(100);unique_index"`
	Name     string `gorm:"size:30"`
	Password string `json:"-"`
	IsAdmin  bool
}

func (User) TableName() string {
	return "users"
}

// 根据 ID 获取 User
func GetUserByID(ID uint) (User, error) {
	var user User
	result := DB.First(&user, ID)
	return user, result.Error
}

// 根据 Email 获取 User
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
