package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDB() (*gorm.DB, error) {
	dsn := "root:root@tcp(localhost:3306)/go-blog-posts?parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config())
	if err != nil {
		return nil, err
	}

	return db, nil
}
