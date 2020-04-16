package models

import (
	"fmt"
	"time"

	"github.com/vinki/pkg/conf"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/vinki/pkg/utils"
)

var DB *gorm.DB

func Init() {
	utils.Log().Info("Init DatabaseConfig Connection")

	var (
		db  *gorm.DB
		err error
	)

	if gin.Mode() == gin.TestMode {
		// 测试环境使用内存数据库
		utils.Log().Info(utils.RelativePath(conf.GlobalConfig.Database.DBFile))
		db, err = gorm.Open("sqlite3", utils.RelativePath(conf.GlobalConfig.Database.DBFile))
		//db, err = gorm.Open("sqlite3", ":memory:")
	} else {
		if conf.GlobalConfig.Database.Type == "UNSET" {
			// 未指定数据库时，使用 SQLite
			db, err = gorm.Open("sqlite3", utils.RelativePath(conf.GlobalConfig.Database.DBFile))
		} else {
			// 指定数据库
			db, err = gorm.Open(conf.GlobalConfig.Database.Type, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				conf.GlobalConfig.Database.Host,
				conf.GlobalConfig.Database.Port,
				conf.GlobalConfig.Database.User,
				conf.GlobalConfig.Database.Password,
				conf.GlobalConfig.Database.DBName),
			)
		}
	}

	if err != nil {
		utils.Log().Panicf("连接数据库失败，%s", err)
	}

	if conf.GlobalConfig.System.Debug {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}
	db.DB().SetMaxIdleConns(50)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Second * 30)
	DB = db

	// 数据表迁移
	migration()
}
