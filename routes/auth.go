package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/controllers"
)

func AuthRoutes(router *gin.RouterGroup) {
	authRouter := router.Group("/auth")

	authRouter.POST("/send-otp", controllers.SendOTP)
	authRouter.POST("/verify-otp", controllers.VerifyOTP)
	authRouter.POST("/resend-otp", controllers.ResendOTP)
}
