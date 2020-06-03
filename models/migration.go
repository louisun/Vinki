package models

import (
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"github.com/louisun/vinki/pkg/conf"
	"github.com/louisun/vinki/pkg/utils"
)

// migration 初始化数据表
func migration() {
	utils.Log().Info("正在初始化数据表...")
	if conf.GlobalConfig.Database.Type == "mysql" {
		DB = DB.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	// 自动迁移模式
	DB.AutoMigrate(&Repo{}, &Tag{}, &Article{}, &User{})

	// 创建管理员
	addAdmin()
	utils.Log().Info("数据表初始化完成")

}

func addAdmin() {
	_, err := GetUserByID(1)

	if gorm.IsRecordNotFoundError(err) {
		password := "vinkipass"
		adminUser := User{
			Email:    "admin@vinki.org",
			NickName: "Renzo",
			IsAdmin:  true,
			Status:   STATUS_ACTIVE,
		}
		if err = adminUser.SetPassword("vinkipass"); err != nil {
			utils.Log().Panic("无法设置管理员密码, %s", err)
		}
		if err = DB.Create(&adminUser).Error; err != nil {
			utils.Log().Panic("无法创建管理员用户, %s", err)
		}

		c := color.New(color.FgWhite).Add(color.BgBlack).Add(color.Bold)
		utils.Log().Info("初始管理员账号：" + c.Sprint("admin@cloudreve.org"))
		utils.Log().Info("初始管理员密码：" + c.Sprint(password))
	}
}
