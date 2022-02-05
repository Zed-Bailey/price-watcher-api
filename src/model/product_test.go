package model_test

import (
	"RestService/src/model"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ProductTestSuite struct {
	suite.Suite
	testUser model.User
}

// initalize some variables for the suite
func (suite *ProductTestSuite) SetupSuite() {
	suite.testUser = model.User{
		Email:    "test@test.com",
		Password: "test",
	}
	model.SetupDatabaseForTesting()
	err := model.CreateNewUser(&suite.testUser)
	suite.Nil(err, "Failed to add user to database, err:%v\n", err)
}

// func (suite *ProductTestSuite) TestAddProduct() {
func (suite *ProductTestSuite) BeforeTest(suiteName, testName string) {
	product := model.Product{
		ItemName:     "testProduct",
		Url:          "https://test.com",
		CurrentPrice: 0,
		LastChecked:  "",
		NextCheck:    "",
	}
	err := suite.testUser.AddProduct(product)
	suite.Nil(err, "Failed to add product! error: %v\n", err)
}

func (suite *ProductTestSuite) AfterTest(suiteName, testName string) {
	// model.DB.Model(&model.Product{}).Unscoped().Delete(model.Product{})
	// model.DB.Unscoped().Delete(model.Product{})
	// model.DB.Exec("DELETE FROM products")
	model.DB.Where("1=1").Unscoped().Delete(&model.Product{})
	// for sime reason the TestDeleteProduct fails if this line isn't here?
	suite.testUser.GetAllProducts()
	// fmt.Printf("num products after delete: %v\n", len(prod))
}

func (suite *ProductTestSuite) TestGetProduct() {
	products, err := suite.testUser.GetAllProducts()
	suite.Nil(err, "error getting products for user %v\nerror: %v\n", suite.testUser, err)
	suite.Equalf(1, len(products), "The number of products for the example user should be 1, got: %v\n", len(products))
}

//
type updateProductData struct {
	Url      string
	ItemName string
}

func (suite *ProductTestSuite) TestUpdateProduct() {
	products, err := suite.testUser.GetAllProducts()
	suite.Nil(err, "error getting products for user %v\nerror: %v\n", suite.testUser, err)
	suite.Equalf(1, len(products), "The number of products for the example user should be 1, got: %v\n", len(products))

	productToUpdate := products[0]
	// define the columns to update and the new value
	newInfo := updateProductData{
		Url:      "NewName",
		ItemName: "https://newUrl.com",
	}
	err = suite.testUser.UpdateProduct(&productToUpdate, newInfo)
	suite.Nil(err, "failed to update product: %v\n", err)

	// refetch the products from the database
	products, _ = suite.testUser.GetAllProducts()
	suite.Equalf(newInfo.ItemName, products[0].ItemName, "ItemName column did not update")
	suite.Equalf(newInfo.Url, products[0].Url, "Url column did not update")

}

func (suite *ProductTestSuite) TestDeleteProduct() {
	// get products
	products, err := suite.testUser.GetAllProducts()
	suite.Nil(err, "error getting products for user %v\nerror: %v\n", suite.testUser, err)
	suite.Equalf(1, len(products), "The number of products for the example user should be 1, got: %v\n", len(products))

	productToDelete := products[0]
	// delete product
	err = suite.testUser.DeleteProduct(productToDelete)
	suite.Nil(err, "error deleting product: %v\n", productToDelete)

	products, _ = suite.testUser.GetAllProducts()

	suite.Equalf(0, len(products), "Product was not deleted properly, got %v but wanted 0", len(products))

}

func TestProductSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}
