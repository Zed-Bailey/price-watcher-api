package model

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model

	UserID       uint
	Url          string  `json:"url"`
	ItemName     string  `json:"item_name"`
	LastChecked  string  `json:"last_check"`
	NextCheck    string  `json:"next_check"`
	CurrentPrice float64 `json:"curr_price"`
}

/**************************************
*			  	REPOSITORY FUNCTIONS 		 		*
***************************************/

// Fetches all products that the user has added
func (u *User) GetAllProducts() ([]Product, error) {
	// products can be directly returned through the attached User interface 'u *User'
	// it would simplify this function to one line,
	// but i may confuse myself how the function returns in the future
	result := DB.Model(&u).Related(&u.Products)
	return u.Products, result.Error
}

// fetches a single product
func (u *User) GetProduct(id uint) (Product, error) {
	var findProduct Product
	err := DB.Where("id = ?", id).First(&findProduct).Error
	return findProduct, err
}

// Add a product to the user
func (user *User) AddProduct(p Product) error {
	// add the new item to the user via association
	return DB.Model(&user).Association("Products").Append(&p).Error
}

func (u *User) DeleteProduct(p Product) error {
	// remove association between product and user
	err := DB.Model(&u).Association("Products").Delete(&p).Error
	if err != nil {
		return err
	}

	// NOTE: at the moment the following line is doing a soft delete, so the object is still in the db
	// BUT it cant be queried, so calling GetItems will not return the object in the response
	return DB.Delete(&p).Error
}

// update a product with new data
// 'product' is the product to update
// 'with' is an interface with values to update the product columns to
func (u *User) UpdateProduct(product *Product, with interface{}) error {
	result := DB.Model(&product).Update(with)
	return result.Error
}
