package main

import (
	"RestService/src/controller"
	"RestService/src/jobs"
	"RestService/src/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

/*

TODO prevent user from logging in multiple times and generating multiple tokens, should return the same token as whats in the session

*/

func main() {

	// setting up cron job
	// https://pkg.go.dev/github.com/robfig/cron
	cj := cron.New()
	cj.AddFunc("@daily", jobs.CheckSites)
	defer cj.Stop()

	// setup gin api routes
	// https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/
	// https://github.com/Depado/gin-auth-example/blob/master/main.go
	r := gin.Default()
	// initalize db
	model.SetupDatabase()
	r.Use(CORSMiddleware())
	//setup session store
	r.Use(sessions.Sessions("application_session", sessions.NewCookieStore([]byte("token"))))

	r.POST("/signup", controller.Signup)
	r.POST("/login", controller.Login)

	private := r.Group("/private")
	// setup endpoints that can only be accessed with a bearer token
	private.Use(AuthorizedEndpoint)
	{
		private.GET("/logout", controller.Logout)
		private.GET("/items", controller.GetItems)
		private.POST("/items", controller.CreateItem)
		private.DELETE("/items/:id", controller.DeleteItem)
		private.PATCH("/items/:id", controller.UpdateItem)
	}
	// start cron jobs and router
	cj.Start()
	r.Run()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func AuthorizedEndpoint(c *gin.Context) {
	// get token from request header
	token, err := c.Request.Cookie("token")
	if errors.Is(err, http.ErrNoCookie) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no cookie named token sent with request!"})
		return
	}

	// token := c.Request.Header.Get("Bearer")
	session := sessions.Default(c)
	// get token value
	user := session.Get(token.Value)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}
