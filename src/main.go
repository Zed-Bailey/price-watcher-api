package main

import (
	"RestService/src/controller"
	"RestService/src/model"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	// https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/
	// https://github.com/Depado/gin-auth-example/blob/master/main.go
	r := gin.Default()
	// initalize db
	model.SetupDatabase()

	//setup session store
	r.Use(sessions.Sessions("application_session", sessions.NewCookieStore([]byte("token"))))

	r.POST("/signup", controller.Signup)
	r.GET("/login", controller.Login)

	private := r.Group("/private")
	// setup endpoints that can only be accessed with a bearer token
	private.Use(AuthorizedEndpoint)
	{
		private.GET("/logout", controller.Logout)

		private.GET("/items", controller.GetItems)
		private.POST("/items", controller.CreateItem)
		private.DELETE("/items/:id", controller.DeleteItem)

	}

	r.Run()
}

func AuthorizedEndpoint(c *gin.Context) {
	token := c.Request.Header.Get("Bearer")
	session := sessions.Default(c)
	user := session.Get(token)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}
