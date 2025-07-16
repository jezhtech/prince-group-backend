package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jezhtech/prince-group-backend/config"
	"github.com/jezhtech/prince-group-backend/models"
)

// JWTMiddleware verifies JWT tokens from the backend OTP system
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// Get JWT secret
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "your-secret-key-change-in-production"
		}

		// Parse and verify token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Set user info in context
			c.Set("user_id", claims["user_id"])
			c.Set("email", claims["email"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

// CombinedAuthMiddleware tries JWT first, then Firebase if JWT fails
func CombinedAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// Try JWT first
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "your-secret-key-change-in-production"
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err == nil && token.Valid {
			// JWT is valid
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				c.Set("user_id", claims["user_id"])
				c.Set("email", claims["email"])
				c.Set("role", claims["role"])
				c.Set("auth_type", "jwt")
				c.Next()
				return
			}
		}

		// Try Firebase if JWT failed
		firebaseToken, err := config.FirebaseAuth.VerifyIDToken(context.Background(), tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("firebaseId", firebaseToken.UID)
		c.Set("auth_type", "firebase")
		c.Next()
	}
}

func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		token, err := config.FirebaseAuth.VerifyIDToken(context.Background(), tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("firebaseId", token.UID)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		token, err := config.FirebaseAuth.VerifyIDToken(context.Background(), tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("firebaseId", token.UID)
		user, err := models.GetUserByFirebaseId(token.UID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
