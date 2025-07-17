package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/controllers"
	"github.com/jezhtech/prince-group-backend/middleware"
)

func ClientRoutes(router *gin.RouterGroup) {
	clientRouter := router.Group("/client")

	// Client-specific booking routes
	clientRouter.GET("/bookings/paginated", middleware.ClientMiddleware(), controllers.GetClientBookingsPaginated)

	// Add more client-specific routes here as needed
	// clientRouter.GET("/dashboard", middleware.ClientMiddleware(), controllers.GetClientDashboard)
	// clientRouter.GET("/profile", middleware.ClientMiddleware(), controllers.GetClientProfile)
	// clientRouter.PUT("/profile", middleware.ClientMiddleware(), controllers.UpdateClientProfile)
}
