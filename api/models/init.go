package models

import (
	"fmt"
	"os"

	// "github.com/wyx-0203/sgs-server/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	var err error
	// db, err = gorm.Open(mysql.Open(global.MYSQL_DSN), &gorm.Config{})
	url := os.Getenv("MYSQL_URL")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, url, dbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Personal{})
}
