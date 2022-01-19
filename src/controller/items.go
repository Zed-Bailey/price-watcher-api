package controller

import (
	"RestService/src/model"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getUser(c *gin.Context) model.User {
	token := c.Request.Header.Get("Bearer")
	session := sessions.Default(c)
	// gets the user id associated with the token, converts it to a string
	id := session.Get(token).(uint)

	var user model.User
	findUser := model.User{Model: gorm.Model{ID: id}}
	//TODO check result for error
	model.DB.First(&user, findUser)
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
	model.DB.Model(&user).Related(&user.Products)
	c.JSON(http.StatusOK, gin.H{"data": user.Products})

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

	item := model.Product{Url: input.Url, ItemName: input.ItemName, LastChecked: "", NextCheck: "", CurrentPrice: 0}
	// add the new item to the user via association
	result := model.DB.Model(&user).Association("Products").Append(&item)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}
