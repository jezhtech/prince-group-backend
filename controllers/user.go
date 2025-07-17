package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/config"
	"github.com/jezhtech/prince-group-backend/helper"
	"github.com/jezhtech/prince-group-backend/models"
)

func GetAllUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
	})
}

func GetUserData(c *gin.Context) {
	authType := c.GetString("auth_type")

	var user models.User
	var err error

	if authType == "jwt" {
		// JWT authentication
		userID := c.GetString("user_id")
		email := c.GetString("email")

		// Try to get user by ID first, then by email
		if userID != "" {
			if userIDInt, parseErr := strconv.ParseUint(userID, 10, 32); parseErr == nil {
				err = config.DB.Where("id = ?", uint(userIDInt)).First(&user).Error
			}
		}

		// If ID lookup fails or userID is empty, try by email
		if err != nil && email != "" {
			err = config.DB.Where("email = ?", email).First(&user).Error
		}
	} else {
		// Firebase authentication
		firebaseID := c.GetString("firebaseId")
		if firebaseID != "" {
			user, err = models.GetUserByFirebaseId(firebaseID)
		} else {
			err = models.ErrUserNotFound
		}
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{
		"user": user,
	})
}

func GetUser(c *gin.Context) {
	authType := c.GetString("auth_type")

	var user models.User
	var err error

	if authType == "jwt" {
		// JWT authentication - get user by ID or email
		userID := c.GetString("user_id")
		email := c.GetString("email")

		if userID != "" {
			if userIDInt, parseErr := strconv.ParseUint(userID, 10, 32); parseErr == nil {
				err = config.DB.Where("id = ?", uint(userIDInt)).First(&user).Error
			}
		}

		// If ID lookup fails or userID is empty, try by email
		if err != nil && email != "" {
			err = config.DB.Where("email = ?", email).First(&user).Error
		}
	} else {
		// Firebase authentication - get user by firebase_id
		firebaseID := c.GetString("firebaseId")
		if firebaseID != "" {
			user, err = models.GetUserByFirebaseId(firebaseID)
		} else {
			err = models.ErrUserNotFound
		}
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{
		"user": user,
	})
}

func CreateUser(c *gin.Context) {
	authType := c.GetString("auth_type")

	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Set authentication-specific fields
	if authType == "jwt" {
		// For JWT users, we already have the user data from OTP verification
		// Just ensure the user exists in the database
		userID := c.GetString("user_id")
		email := c.GetString("email")

		// Check if user already exists
		var existingUser models.User
		if userID != "" {
			if userIDInt, parseErr := strconv.ParseUint(userID, 10, 32); parseErr == nil {
				err = config.DB.Where("id = ?", uint(userIDInt)).First(&existingUser).Error
			}
		}

		// If ID lookup fails, try by email
		if err != nil && email != "" {
			err = config.DB.Where("email = ?", email).First(&existingUser).Error
		}

		if err == nil {
			// User exists, return it
			c.JSON(200, gin.H{
				"message": "User already exists",
				"user":    existingUser,
			})
			return
		}

		// Set the user ID from JWT claims
		if userID != "" {
			if userIDInt, parseErr := strconv.ParseUint(userID, 10, 32); parseErr == nil {
				user.ID = uint(userIDInt)
			}
		}
	} else {
		// Firebase authentication - generate unique UserID
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
	authType := c.GetString("auth_type")

	var user models.User
	var err error

	if authType == "jwt" {
		// JWT authentication
		userID := c.GetString("user_id")
		email := c.GetString("email")

		// Try to get user by ID first, then by email
		if userID != "" {
			if userIDInt, parseErr := strconv.ParseUint(userID, 10, 32); parseErr == nil {
				err = config.DB.Where("id = ?", uint(userIDInt)).First(&user).Error
			}
		}

		// If ID lookup fails or userID is empty, try by email
		if err != nil && email != "" {
			err = config.DB.Where("email = ?", email).First(&user).Error
		}
	} else {
		// Firebase authentication
		firebaseID := c.GetString("firebaseId")
		if firebaseID != "" {
			user, err = models.GetUserByFirebaseId(firebaseID)
		} else {
			err = models.ErrUserNotFound
		}
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Bind the update data
	var updateData models.User
	err = c.ShouldBindJSON(&updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update only the allowed fields
	user.FullName = updateData.FullName
	user.Mobile = updateData.Mobile
	user.Address = updateData.Address
	user.City = updateData.City
	user.State = updateData.State
	user.Pincode = updateData.Pincode
	user.Aadhaar = updateData.Aadhaar

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
	authType := c.GetString("auth_type")

	var user models.User
	var err error

	if authType == "jwt" {
		// JWT authentication
		userID := c.GetString("user_id")
		email := c.GetString("email")

		// Try to get user by ID first, then by email
		if userID != "" {
			if userIDInt, parseErr := strconv.ParseUint(userID, 10, 32); parseErr == nil {
				err = config.DB.Where("id = ?", uint(userIDInt)).First(&user).Error
			}
		}

		// If ID lookup fails or userID is empty, try by email
		if err != nil && email != "" {
			err = config.DB.Where("email = ?", email).First(&user).Error
		}
	} else {
		// Firebase authentication
		firebaseID := c.GetString("firebaseId")
		if firebaseID != "" {
			user, err = models.GetUserByFirebaseId(firebaseID)
		} else {
			err = models.ErrUserNotFound
		}
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
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
