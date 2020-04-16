package main

import (
	"flag"
	"fmt"

	"github.com/vinki/bootstrap"
	"github.com/vinki/pkg/conf"
	"github.com/vinki/pkg/utils"
	"github.com/vinki/routers"
)

func init() {
	var confPath string
	flag.StringVar(&confPath, "c", "/Users/louisun/.vinki/code/config.yml", "configuration file")
	flag.Parse()
	bootstrap.Init(confPath)
}

func main() {
	// 初始化路由
	engine := routers.InitRouter()
	utils.Log().Infof("Listening: %d", conf.GlobalConfig.System.Port)
	if err := engine.Run(fmt.Sprintf(":%d", conf.GlobalConfig.System.Port)); err != nil {
		utils.Log().Errorf("无法启动服务端口 [%d]: %s", conf.GlobalConfig.System.Port, err)
	}
}
