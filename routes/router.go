package routes

import (
	"github.com/gin-gonic/gin"
)

func AppRouter(router *gin.Engine) {
	apiRouter := router.Group("/api/v1")
	apiRouter.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	{
		AuthRoutes(apiRouter)
		UserRoutes(apiRouter)
		BookingRoutes(apiRouter)
		ReferralRoutes(apiRouter)
		TicketRoutes(apiRouter)
		YouTubeRoutes(apiRouter)
		PaymentRoutes(apiRouter)
	}
}
