package main

import (
	"flag"
	"fmt"

	"github.com/louisun/vinki/bootstrap"
	"github.com/louisun/vinki/pkg/conf"
	"github.com/louisun/vinki/pkg/utils"
	"github.com/louisun/vinki/routers"
)

func initConfig() {
	var confPath string

	flag.StringVar(&confPath, "c", "./conf/config.yml", "configuration file")
	flag.Parse()
	bootstrap.Init(confPath)
}

func main() {
	// 初始化配置
	initConfig()
	// 初始化路由
	engine := routers.InitRouter()

	utils.Log().Infof("Listening: %d", conf.GlobalConfig.System.Port)

	if err := engine.Run(fmt.Sprintf(":%d", conf.GlobalConfig.System.Port)); err != nil {
		utils.Log().Errorf("无法启动服务端口 [%d]: %s", conf.GlobalConfig.System.Port, err)
	}
}
