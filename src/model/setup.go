package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// database instance
var DB *gorm.DB

// Conect to the database
func SetupDatabase() {
	database, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to open database!")
	}

	database.AutoMigrate(&User{})
	database.AutoMigrate(&WatchingItem{})
	DB = database
}
