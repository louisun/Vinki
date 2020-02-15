package db

import (
	"fmt"
	"log"

	"github.com/vinki/pkg/utils"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func InitDatabase() {
	var err error
	db, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		utils.Config.Database.Host,
		utils.Config.Database.Port,
		utils.Config.Database.User,
		utils.Config.Database.Password,
		utils.Config.Database.Database),
	)

	if err != nil {
		log.Fatalf("db.Setup err: %v", err)
	}
	db.DB().SetMaxOpenConns(100)
	log.Println("DB init success.")
}

// CloseDB closes database connection
func CloseDB() {
	defer db.Close()
}
