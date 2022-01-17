package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	r.Run()
}
