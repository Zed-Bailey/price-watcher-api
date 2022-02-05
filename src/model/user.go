package model

import (
	"errors"

	"github.com/jinzhu/gorm"
)

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

// creates a new user, will return an error if the email is already in use
func CreateNewUser(newUser *User) error {
	// check if the email already exists in the database
	var existing User
	result := DB.First(&existing, User{Email: newUser.Email})
	// if a row was returned then the email is already in use
	if result.RowsAffected != 0 {
		return errors.New("email already in use")
	}

	return DB.Create(&newUser).Error
}

// Checks the database for a user with matching email/password combination
// returns true if there is a match false otherwise
// if true is returned then the returned User struct will contain the users information
func CheckUserDoesExist(email string, pass string) (bool, User) {

	var loggedInUser User

	// select the first user with matching email/password combo
	result := DB.Where(&User{Email: email, Password: pass}).First(&loggedInUser)

	// check if a record was returned
	return !errors.Is(result.Error, gorm.ErrRecordNotFound), loggedInUser
}
