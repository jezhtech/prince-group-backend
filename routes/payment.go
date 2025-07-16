package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/controllers"
	"github.com/jezhtech/prince-group-backend/middleware"
)

func PaymentRoutes(router *gin.RouterGroup) {
	paymentRouter := router.Group("/payment")

	// Protected routes (require authentication)
	paymentRouter.POST("/links", middleware.UserMiddleware(), controllers.CreatePaymentLink)
	paymentRouter.GET("/status/:linkId", middleware.UserMiddleware(), controllers.CheckPaymentStatus)
	paymentRouter.GET("/history", middleware.UserMiddleware(), controllers.GetPaymentHistory)
	paymentRouter.POST("/send-email/:bookingNumber", middleware.UserMiddleware(), controllers.SendPaymentConfirmationEmail)

	// Public routes (no authentication required)
	paymentRouter.POST("/webhook", controllers.PaymentWebhook)
	paymentRouter.GET("/callback", controllers.PaymentCallback)

	// Test endpoint (for development only)
	paymentRouter.GET("/test-email", controllers.TestEmailEndpoint)
}
