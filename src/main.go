package main

import (
	"RestService/src/controller"
	"RestService/src/jobs"
	"RestService/src/logger"
	"RestService/src/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

/*

TODO add/update code documentation
improvments,
- calculate percentage change between prices, this will require updating the model
- batch update on a per user basis and run that in individual go routines rather then per item
- how to handle link no longer existing?
*/

func main() {

	// setup custom logger
	logger.SetupLogger()
	logger.Log.Info().Msg("Initalized Logger")

	// setting up cron job
	// https://pkg.go.dev/github.com/robfig/cron
	cj := cron.New()
	cj.AddFunc("@daily", jobs.CheckSites)
	defer cj.Stop()
	logger.Log.Info().Msg("Initalized cron jobs")

	// initalize db
	model.SetupDatabase()
	logger.Log.Info().Msg("Initalized Database")

	// initalize router
	router := SetupRouter()
	logger.Log.Info().Msg("Initalized Router")

	// start cron jobs and router
	cj.Start()
	router.Run()
}

// Setup the router and api routes
func SetupRouter() *gin.Engine {
	// setup gin api routes
	// https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/
	// https://github.com/Depado/gin-auth-example/blob/master/main.go
	r := gin.Default()

	r.Use(CORSMiddleware())

	//setup session store
	store := sessions.NewCookieStore([]byte("userID"))
	store.Options(sessions.Options{MaxAge: 60 * 60 * 12})

	// add session store to the router
	r.Use(sessions.Sessions("token", store))

	r.POST("/signup", controller.Signup)
	r.POST("/login", controller.Login)

	private := r.Group("/private")
	// setup endpoints that require an authorized user, uses the AuthorizedEndpoint middleware to
	// terminate requests from invalid users
	private.Use(AuthorizedEndpoint)
	{
		private.GET("/logout", controller.Logout)
		private.GET("/items", controller.GetItems)
		private.POST("/items", controller.CreateItem)
		private.DELETE("/items/:id", controller.DeleteItem)
		private.PATCH("/items/:id", controller.UpdateItem)
	}

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

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

	session := sessions.Default(c)
	// get userID stored in the session
	user := session.Get("userID")
	logger.Log.Debug().
		Str("token", token.Value).
		Interface("user", user).
		Msg("Authorized endpoint")

	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}
