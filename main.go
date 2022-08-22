package main

import (
	"net/http"

	"auth.jwt.api/controllers"
	"auth.jwt.api/database"
	"auth.jwt.api/middleware"
	"github.com/gin-gonic/gin"
)

func userRegister(c *gin.Context) {
}

func init() {
	database.Setup()
}

func htmlMain(c *gin.Context) {
	htmlIndex := `<html>
	<body>
		<a href="/google">Google Log In</a>
	</body>
	</html>`

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(htmlIndex))
	return
}

func main() {
	r := gin.Default()

	r.GET("/", htmlMain)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/google", controllers.GoogleAuth)
	r.GET("/google/callback", controllers.GoogleAuthCallback)

	auth := r.Group("/user", middleware.AuthMiddleware())

	{
		auth.GET("/profile", controllers.Profile)
	}

	r.Run(":8080")
}
