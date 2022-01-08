package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"models"
)

var Instance *gorm.DB

func Connect(path string) {
	db, err := gorm.Open(mysql.Open(path), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// set the pointer instance to the connected database
	Instance = db

	//func (db *DB) AutoMigrate(dst ...interface{}) error
	db.AutoMigrate(&models.Profile{})
}
