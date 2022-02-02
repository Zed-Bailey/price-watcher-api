package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email    string
	Password string
	Products []Product `json:"items" gorm:"foreignKey:UserID"`
}

/**************************************
*			  	REPOSITORY FUNCTIONS 		 		*
***************************************/

// Find a user from their ID
func FindUser(id uint) (User, error) {
	// define the out interface
	var user User

	// create an interface to match our search to
	searchFor := User{Model: gorm.Model{ID: id}}

	// query db for the first user to match
	result := DB.First(&user, searchFor)

	return user, result.Error
}

// ceate a new user
func CreateNewUser(newUser User) error {
	return nil
}
