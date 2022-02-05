package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// database instance set by one of the SetupDatabase methods
var DB *gorm.DB

// Connect to the database
func SetupDatabase() {
	database, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to open database!")
	}
	database.AutoMigrate(&User{})
	database.AutoMigrate(&Product{})
	DB = database
}

// Sets up an in memory database that can be used when running test suites
func SetupDatabaseForTesting() {
	// https://gorm.io/docs/connecting_to_the_database.html
	database, err := gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic("failed to setup in-memory database!")
	}
	database.AutoMigrate(&User{})
	database.AutoMigrate(&Product{})
	DB = database
}
