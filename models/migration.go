package models

import (
	"github.com/vinki/pkg/conf"
	"github.com/vinki/pkg/utils"
)

// migration 初始化数据表
func migration() {
	utils.Log().Info("正在初始化数据表...")
	if conf.GlobalConfig.Database.Type == "mysql" {
		DB = DB.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	// 自动迁移模式
	DB.AutoMigrate(&Repo{}, &Tag{}, &Article{}, &User{})
	addAdmin()

	// 创建管理员
	addAdmin()
	utils.Log().Info("数据表初始化完成")

}

func addAdmin() {

}
