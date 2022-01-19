package controller

import (
	"RestService/src/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

/*
useful links
https://gorm.io/docs/associations.html#Association-Mode
*/

func getUser(c *gin.Context) (model.User, error) {
	token := c.Request.Header.Get("Bearer")
	session := sessions.Default(c)
	// gets the user id associated with the token, converts it to a string
	id := session.Get(token).(uint)

	var user model.User
	findUser := model.User{Model: gorm.Model{ID: id}}
	//TODO check result for error
	result := model.DB.First(&user, findUser)
	return user, result.Error
}

func GetItems(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

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
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	item := model.Product{Url: input.Url, ItemName: input.ItemName, LastChecked: "", NextCheck: "", CurrentPrice: 0}
	// add the new item to the user via association
	result := model.DB.Model(&user).Association("Products").Append(&item)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// DELETE /private/items/:id
// removes a product with the specified id from the database
func DeleteItem(c *gin.Context) {

	// https://gorm.io/docs/delete.html

	itemIdString := c.Param("id")
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	itemId, err := strconv.ParseUint(itemIdString, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var productToDelete model.Product

	// find product in db
	model.DB.Where(&model.Product{Model: gorm.Model{ID: uint(itemId)}}).First(&productToDelete)

	// remove association between product and user
	model.DB.Model(&user).Association("Products").Delete(&productToDelete)

	// NOTE: at the moment the following line is doing a soft delete, so the object is still in the db
	// BUT it cant be queried, so calling GetItems will not return the object in the response
	model.DB.Delete(&productToDelete)

	c.JSON(http.StatusOK, gin.H{"data": "removed product sucessfully", "product": productToDelete})
}
