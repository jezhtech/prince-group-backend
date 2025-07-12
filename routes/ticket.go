package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/controllers"
	"github.com/jezhtech/prince-group-backend/middleware"
)

func TicketRoutes(router *gin.RouterGroup) {
	ticketRouter := router.Group("/ticket")

	ticketRouter.GET("/:id", middleware.UserMiddleware(), controllers.GetTicket)
	ticketRouter.GET("/", controllers.GetAllTickets)
	ticketRouter.POST("/", middleware.AdminMiddleware(), controllers.CreateTicket)
	ticketRouter.PUT("/:id", middleware.AdminMiddleware(), controllers.UpdateTicket)
	ticketRouter.DELETE("/:id", middleware.AdminMiddleware(), controllers.DeleteTicket)
}
