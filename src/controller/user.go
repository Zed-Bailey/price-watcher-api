package controller

import (
	"RestService/src/logger"
	"RestService/src/model"
	b64 "encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// https://gorm.io/docs/query.html

/**************************************
* 					 	 SIGNUP    							*
***************************************/

// POST /signup
// create a new user
// returns a token that can be used to further access the authenticated points of the api
func Signup(c *gin.Context) {
	session := sessions.Default(c)

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
	// parse the split string and trim the strings
	email := split[0]
	pass := split[1]
	if strings.Trim(email, " ") == " " || strings.Trim(pass, " ") == " " {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	newUser := model.User{Email: email, Password: pass}

	// create the user
	if result := model.CreateNewUser(&newUser); result != nil {
		if result.Error() == "email already in use" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "That email is already in use"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result})
		}
		return
	}

	// save userID in session
	session.Set("userID", newUser.ID)
	// save session
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving user session"})
		return
	}

	// return ok status
	c.JSON(http.StatusOK, gin.H{})
}

/**************************************
* 					 	 LOGIN    							*
***************************************/

// POST /login
// Login a user with email/password combo using basic auth
// email and password should be joined with a ':' and then encoded with base64
// then attached to an 'Authorization' header
// returns a 412 error code if an email or password field is missing
// returns a 404 error code if the user couldn't be found
// returns a 200 ok code if the user was found, attatched is a token that can be used fro future requests
func Login(c *gin.Context) {
	session := sessions.Default(c)
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

	exists, user := model.CheckUserDoesExist(split[0], split[1])
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or password is incorrect"})
		return
	}

	// save the user ID in the session
	session.Set("userID", user.ID)
	logger.Log.Info().Interface("userID", session.Get("userID")).Msg("user now logged in")

	// update session
	session.Save()

	c.JSON(http.StatusOK, gin.H{})
}

/**************************************
* 					 	 LOGOUT    							*
***************************************/

// GET /logout
// logs a user out of the application
func Logout(c *gin.Context) {

	// check that the token cookie was passed along into the session
	_, err := c.Request.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No session cookie passed"})
		return
	}

	session := sessions.Default(c)

	// check if there is a valid user in the session
	if user := session.Get("userID"); user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}

	// remove userID and save session
	session.Clear()

	// this issue solved the problem of the session persisting even once it's been deleted
	// https://github.com/gin-contrib/sessions/issues/89
	session.Options(sessions.Options{Path: "/", MaxAge: -1}) // this sets the cookie with a MaxAge of 0

	logger.Log.Info().
		Interface("userID", session.Get("userID")).
		Msg("cleared the session")

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
