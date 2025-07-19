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
	clientRouter.GET("/bookings/stats", middleware.ClientMiddleware(), controllers.GetClientBookingsStats)
}
