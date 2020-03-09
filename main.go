package main

import (
	"flag"
	"fmt"

	"github.com/jinzhu/configor"
	"github.com/vinki/db"
	"github.com/vinki/routers"
	"github.com/vinki/utils"
)

func main() {
	// 1. Load Configs
	configFile := flag.String("c", "./conf/config.yml", "configuration file")
	flag.Parse()
	configor.Load(&utils.Config, *configFile)

	// 2. Init Database
	db.InitDatabase()
	defer db.CloseDB()

	// 3. Load Routes
	r := routers.LoadRoutes()
	r.Run(fmt.Sprintf(":%d", utils.Config.Server.Port))
}
