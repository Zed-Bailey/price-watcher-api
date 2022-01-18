package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email            string
	Password         string
	WatchingProducts []WatchingItem `json:"items" gorm:"foreignKey:UserReference"`
}
