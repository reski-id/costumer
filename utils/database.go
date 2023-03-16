package utils

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := "root:@tcp(127.0.0.1:3306)/costumerdb?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
