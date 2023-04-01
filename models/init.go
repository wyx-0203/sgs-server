package models

import (
	"github.com/wyx-0203/sgs-server/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open(mysql.Open(global.MYSQL_DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Personal{})
}
