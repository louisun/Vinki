package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/louisun/vinki/model"
	"github.com/louisun/vinki/pkg/conf"
	"github.com/louisun/vinki/pkg/utils"
)

func Init(path string) {
	InitApplication()
	// 加载配置文件
	conf.Init(path)
	// 设置 gin 模式
	if !conf.GlobalConfig.System.Debug {
		utils.Log().Info("gin 当前为 Release 模式")
		gin.SetMode(gin.ReleaseMode)
	} else {
		utils.Log().Info("gin 当前为 Test 模式")
		gin.SetMode(gin.TestMode)
	}
	// 初始化 Redis
	// cache.Init()

	// 初始化数据库
	model.Init()

	// 初始化静态文件系统 StaticFS
	InitStatic()

	// 初始化通用鉴权器
	//auth.Init()
}
