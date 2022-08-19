package main

import (
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

func main() {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	auth := r.Group("/user", middleware.AuthMiddleware())

	{
		auth.GET("/profile", controllers.Profile)
	}

	r.Run(":8080")
}
