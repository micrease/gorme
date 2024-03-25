package tests

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var gDB *gorm.DB

func GetDB() *gorm.DB {
	if gDB == nil {
		dsn := "root:123456@tcp(127.0.0.1:3306)/gorme?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{

			Logger: logger.Default.LogMode(logger.Info)})
		if err != nil {
			log.Fatalln("连接数据库失败")
		}

		gDB = db
	}
	return gDB
}
