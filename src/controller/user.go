package controller

import (
	"RestService/src/model"
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// https://gorm.io/docs/query.html

// POST /user
// create a new user
// returns a token that can be used to further access the authenticated points of the api
func Signup(c *gin.Context) {
	session := sessions.Default(c)

	email := c.PostForm("email")
	pass := c.PostForm("password")

	if strings.Trim(email, " ") == " " || strings.Trim(pass, " ") == " " {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// check if the email already exists in the database
	var existing model.User
	result := model.DB.First(&existing, model.User{Email: email})

	if result.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use"})
		return
	}
	newUser := model.User{Email: email, Password: pass}
	model.DB.Create(&newUser)

	// generate a new session token
	token := generateToken()
	session.Set(token, newUser.ID)
	// save session
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving user session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "successfully created new user",
		"id":    newUser.ID,
		"token": token,
	})
}

// GET /user
// Login a user with email/password combo using basic auth
// email and password should be joined with a ':' and then encoded with base64
// then attached to an 'Authorization' header
// returns a 412 error code if an email or password field is missing
// returns a 404 error code if the user couldn't be found
// returns a 200 ok code if the user was found, attatched is a token that can be used fro future requests
func Login(c *gin.Context) {

	auth := c.Request.Header.Get("Authorization")
	// base64 decode
	dec, _ := b64.StdEncoding.DecodeString(auth)
	// split on ':'
	split := strings.Split(string(dec), ":")

	// check that there are 2 parts to the split string
	if len(split) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode authorization header"})
		return
	}

	var loggedInUser model.User

	// select the first user with matching email/password combo
	result := model.DB.Where(&model.User{Email: split[0], Password: split[1]}).First(&loggedInUser)

	// // check if a record was returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "No user found with that email/password"})
		return
	}
	token := generateToken()
	// save the token to the session
	session := sessions.Default(c)
	session.Set(token, loggedInUser.ID)
	session.Options(sessions.Options{
		MaxAge: 3600 * 12, // set session to expire in 12 hours
	})
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.SetCookie("token", token, 60*60*12, "/", "localhost", false, false)

	// should return a token that can be used to access authenticated points later on
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": token}})
}

// GET /logout
// logs a user out of the application
//
func Logout(c *gin.Context) {
	token := c.Request.Header.Get("Bearer")

	session := sessions.Default(c)

	// check if there is a valid token in the session
	if user := session.Get(token); user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}

	// remove token and save session
	session.Delete(token)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// generates a random 30 character token, proably not secure but for now it'll do
func generateToken() string {
	b := make([]byte, 30)
	rand.Read(b)
	return hex.EncodeToString(b)
}
