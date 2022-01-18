package controller

import (
	"RestService/src/model"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getUser(c *gin.Context) model.User {
	token := c.Request.Header.Get("Bearer")
	session := sessions.Default(c)
	// gets the user id associated with the token, converts it to a string
	id := session.Get(token).(string)

	var user model.User
	//TODO check result for error
	model.DB.First(&user, model.User{UserID: id})
	return user
}

func GetItems(c *gin.Context) {
	// token := c.Request.Header.Get("Bearer")
	// session := sessions.Default(c)
	// // gets the user id associated with the token, converts it to a string
	// id := session.Get(token).(string)

	// var user model.User
	// //TODO check result for error
	// model.DB.First(&user, model.User{UserID: id})
	user := getUser(c)
	model.DB.Model(&user).Related(&user.WatchingProducts)
	c.JSON(http.StatusOK, gin.H{"data": user.WatchingProducts})

}

type ItemInput struct {
	Url      string `json:"url"`
	ItemName string `json:"item_name"`
}

func CreateItem(c *gin.Context) {
	var input ItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get the user
	user := getUser(c)
	// fetch the current items
	currentItems := user.WatchingProducts
	// appened the enw item
	item := model.WatchingItem{Url: input.Url, ItemName: input.ItemName, LastChecked: "", NextCheck: "", CurrentPrice: 0}
	currentItems = append(currentItems, item)
	// update the user in the db
	model.DB.Model(&user).Updates(model.User{WatchingProducts: currentItems})
	c.JSON(http.StatusOK, gin.H{"data": item})
}
