package common

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initGORM(cfg DBAuth) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/task?charset=utf8mb4&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.IPPort)

	src, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = src
}
