package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/config"
	"github.com/jezhtech/prince-group-backend/helper"
	"github.com/jezhtech/prince-group-backend/models"
)

func GetAllUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	fmt.Println(users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
	})
}

func GetUserData(c *gin.Context) {
	firebaseID := c.GetString("firebaseId")

	user, err := models.GetUserByFirebaseId(firebaseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(200, gin.H{
		"user": user,
	})
}

func GetUser(c *gin.Context) {
	userID := c.GetString("uid")

	var user models.User

	err := config.DB.Where("user_id = ?", userID).First(&user).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(200, gin.H{
		"user": user,
	})
}

func CreateUser(c *gin.Context) {
	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Generate unique UserID with retry logic
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		user.UserID = helper.GenerateUserID()

		// Check if UserID already exists
		var existingUser models.User
		err := config.DB.Where("user_id = ?", user.UserID).First(&existingUser).Error

		if err != nil {
			// UserID doesn't exist, we can use it
			break
		}

		// If we've tried maxRetries times, return an error
		if i == maxRetries-1 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate unique user ID after multiple attempts"})
			return
		}
	}

	err = config.DB.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

func UpdateUser(c *gin.Context) {
	firebaseID := c.GetString("firebaseId")

	var user models.User

	err := config.DB.Where("firebase_id = ?", firebaseID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	err = c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Explicitly set role to 'user' to prevent someone from sending role as admin and get easy get the admin access
	user.Role = "user"

	err = config.DB.Save(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

func DeleteUser(c *gin.Context) {
	userID := c.GetString("uid")

	var user models.User

	err := config.DB.Where("userId = ?", userID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	err = config.DB.Delete(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "User deleted successfully",
	})
}
