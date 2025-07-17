package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/controllers"
	"github.com/jezhtech/prince-group-backend/middleware"
)

func BookingRoutes(router *gin.RouterGroup) {
	bookingRouter := router.Group("/booking")

	bookingRouter.GET("/:bookingNumber", middleware.UserMiddleware(), controllers.GetBooking)
	bookingRouter.GET("/admin/all", middleware.AdminMiddleware(), controllers.GetAllBookings)
	bookingRouter.GET("/admin/paginated", middleware.AdminMiddleware(), controllers.GetAllBookingsPaginated)
	bookingRouter.POST("/", middleware.UserMiddleware(), controllers.CreateBooking)
	bookingRouter.PUT("/:bookingNumber", middleware.UserMiddleware(), controllers.UpdateBooking)
	bookingRouter.DELETE("/:id", middleware.UserMiddleware(), controllers.DeleteBooking)
	bookingRouter.GET("/user", middleware.UserMiddleware(), controllers.GetBookingsByUserId)
	bookingRouter.GET("/check-payment/:bookingNumber", middleware.UserMiddleware(), controllers.CheckPayment)
}
