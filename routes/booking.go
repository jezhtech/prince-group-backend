package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/controllers"
	"github.com/jezhtech/prince-group-backend/middleware"
)

func BookingRoutes(router *gin.RouterGroup) {
	bookingRouter := router.Group("/booking")

	bookingRouter.GET("/:id", middleware.UserMiddleware(), controllers.GetBooking)
	bookingRouter.GET("/all", middleware.UserMiddleware(), controllers.GetAllBookings)
	bookingRouter.POST("/", middleware.UserMiddleware(), controllers.CreateBooking)
	bookingRouter.PUT("/:id", middleware.UserMiddleware(), controllers.UpdateBooking)
	bookingRouter.DELETE("/:id", middleware.UserMiddleware(), controllers.DeleteBooking)
}
