package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email    string
	Password string
	Products []Product `json:"items" gorm:"foreignKey:UserID"`
}
