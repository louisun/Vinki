package main

import (
	"fmt"

	"github.com/jinzhu/configor"
	"github.com/vinki/db"
	"github.com/vinki/pkg/utils"
	"github.com/vinki/routers"
)

func main() {
	// 1. Load Config
	configor.Load(&utils.Config, "./conf/config.yml")

	// 2. InitDatabase
	db.InitDatabase()
	defer db.CloseDB()
	// 3. LoadRoutes
	r := routers.LoadRoutes()
	r.Run(fmt.Sprintf(":%d", utils.Config.Server.Port))
}
