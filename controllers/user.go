package controllers

import (
	"net/http"
	"strconv"

	"auth.jwt.api/database"
	"auth.jwt.api/models"
	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {
	userId := c.MustGet("id").(string)
	var user models.User
	id, err := strconv.Atoi(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	database.DB.First(&user, id)
	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "data": user})
}
