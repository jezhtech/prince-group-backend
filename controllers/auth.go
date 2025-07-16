package controllers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jezhtech/prince-group-backend/config"
	"github.com/jezhtech/prince-group-backend/helper"
	"github.com/jezhtech/prince-group-backend/models"
)

// OTP storage in memory (in production, use Redis or database)
var otpStore = make(map[string]OTPData)

type OTPData struct {
	OTP       string
	Email     string
	ExpiresAt time.Time
	Attempts  int
}

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type SendOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	OTP     string `json:"otp,omitempty"` // Only for development
}

type VerifyOTPResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Token   string      `json:"token,omitempty"`
	User    models.User `json:"user,omitempty"`
}

// Generate a random 6-digit OTP
func generateOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

// SendOTP generates and sends OTP to the provided email
func SendOTP(c *gin.Context) {
	var req SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// Check if user exists
	var user models.User
	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		// User doesn't exist, create a temporary user for OTP verification
		user = models.User{
			Email: req.Email,
			Role:  "user",
		}
	}

	// Generate OTP
	otp := generateOTP()

	// Store OTP with expiration (5 minutes)
	otpStore[req.Email] = OTPData{
		OTP:       otp,
		Email:     req.Email,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Attempts:  0,
	}

	// Send OTP email
	err := helper.SendOTPEmail(req.Email, otp)
	if err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		// Still log the OTP for development
		fmt.Printf("OTP for %s: %s\n", req.Email, otp)
	}

	c.JSON(http.StatusOK, SendOTPResponse{
		Success: true,
		Message: "OTP sent successfully",
		OTP:     otp, // Remove this in production
	})
}

// VerifyOTP verifies the OTP and logs in the user
func VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// Get stored OTP data
	otpData, exists := otpStore[req.Email]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "OTP not found or expired",
		})
		return
	}

	// Check if OTP is expired
	if time.Now().After(otpData.ExpiresAt) {
		delete(otpStore, req.Email)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "OTP has expired",
		})
		return
	}

	// Check attempts
	if otpData.Attempts >= 3 {
		delete(otpStore, req.Email)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Too many failed attempts",
		})
		return
	}

	// Verify OTP
	if otpData.OTP != req.OTP {
		otpData.Attempts++
		otpStore[req.Email] = otpData
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid OTP",
		})
		return
	}

	// OTP is valid, get or create user
	var user models.User
	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		// Create new user
		user = models.User{
			Email: req.Email,
			Role:  "user",
		}
		config.DB.Create(&user)
	}

	// Generate JWT token
	token, err := generateJWTToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate token",
		})
		return
	}

	// Clean up OTP
	delete(otpStore, req.Email)

	c.JSON(http.StatusOK, VerifyOTPResponse{
		Success: true,
		Message: "OTP verified successfully",
		Token:   token,
		User:    user,
	})
}

// ResendOTP resends OTP to the provided email
func ResendOTP(c *gin.Context) {
	var req SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// Generate new OTP
	otp := generateOTP()

	// Store new OTP with expiration (5 minutes)
	otpStore[req.Email] = OTPData{
		OTP:       otp,
		Email:     req.Email,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Attempts:  0,
	}

	// Send OTP email
	err := helper.SendOTPEmail(req.Email, otp)
	if err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		// Still log the OTP for development
		fmt.Printf("New OTP for %s: %s\n", req.Email, otp)
	}

	c.JSON(http.StatusOK, SendOTPResponse{
		Success: true,
		Message: "OTP resent successfully",
		OTP:     otp, // Remove this in production
	})
}

// generateJWTToken generates a JWT token for the user
func generateJWTToken(user models.User) (string, error) {
	// Get JWT secret from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production" // Default for development
	}

	// Create claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}
