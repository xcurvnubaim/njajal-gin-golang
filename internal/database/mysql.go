package database

import (
	"github.com/xcurvnubaim/njajal-gin-golang/internal/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	dsn := configs.Config.DatabaseURL + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
