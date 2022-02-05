package controller

import (
	"RestService/src/logger"
	"RestService/src/model"
	"RestService/src/web"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

/*
useful links
https://gorm.io/docs/associations.html#Association-Mode
*/

/**************************************
*					 	   UTILITY  							*
***************************************/

// function to simplify getting a user from session
func getUser(c *gin.Context) (model.User, error) {
	session := sessions.Default(c)
	// gets the user id associated with the token, converts it to a string
	id := session.Get("userID").(uint)
	logger.Log.Info().Uint("id", id).Msg("id fetched from session")
	// var user model.User
	// findUser := model.User{Model: gorm.Model{ID: id}}
	// //TODO check result for error
	// result := model.DB.First(&user, findUser)
	return model.FindUser(id)
}

/**************************************
*					 	 FETCH ITEM 							*
***************************************/

// returns all the items
func GetItems(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	products, err := user.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, products)
}

/**************************************
*						CREATE ITEM 							*
***************************************/

// a struct to bind post request data to
type CreateItemInput struct {
	Url      string `json:"url" binding:"required"`
	ItemName string `json:"item_name" binding:"required"`
}

// POST /private/items
func CreateItem(c *gin.Context) {
	var input CreateItemInput
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
	currPrice, err := web.Fetch(input.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	dateNow := time.Now()
	// add a day
	nextCheck := dateNow.AddDate(0, 0, 1)
	// add values to struct
	item := model.Product{
		Url:          input.Url,
		ItemName:     input.ItemName,
		LastChecked:  dateNow.Format("01/02/2006"),
		NextCheck:    nextCheck.Format("01/02/2006"),
		CurrentPrice: currPrice,
	}

	if err = user.AddProduct(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// add the new item to the user via association
	// result := model.DB.Model(&user).Association("Products").Append(&item)
	// if result.Error != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"data": item})
}

/**************************************
*						DELETE ITEM 							*
***************************************/

// DELETE /private/items/:id
// removes a product with the specified id from the database
func DeleteItem(c *gin.Context) {

	// https://gorm.io/docs/delete.html

	itemId := c.Param("id")
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var productToDelete model.Product

	// find product in db
	if err := model.DB.Where("id = ?", itemId).First(&productToDelete).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product could not be found"})
		return
	}

	// remove association between product and user
	// model.DB.Model(&user).Association("Products").Delete(&productToDelete)

	// NOTE: at the moment the following line is doing a soft delete, so the object is still in the db
	// BUT it cant be queried, so calling GetItems will not return the object in the response
	// model.DB.Delete(&productToDelete)
	err = user.DeleteProduct(productToDelete)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "removed product sucessfully", "product": productToDelete})
}

/**************************************
*						UPDATE ITEM 							*
***************************************/
type UpdateItemInput struct {
	Url      string `json:"url"`
	ItemName string `json:"item_name"`
}

// PATCH /private/items/:id
func UpdateItem(c *gin.Context) {

	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// find product
	var updateProduct model.Product
	if err := model.DB.Where("id = ?", c.Param("id")).First(&updateProduct).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product could not be found"})
		return
	}

	// validate update input
	var input UpdateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update product
	// model.DB.Model(&updateProduct).Update(input)
	if err := user.UpdateProduct(&updateProduct, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "updated product", "product": updateProduct})
}
