package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/config"
	"github.com/jezhtech/prince-group-backend/models"
)

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
		fmt.Println("AdminMiddleware")
		tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		token, err := config.FirebaseAuth.VerifyIDToken(context.Background(), tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("firebaseId", token.UID)
		fmt.Println("token.UID", token.UID)
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
