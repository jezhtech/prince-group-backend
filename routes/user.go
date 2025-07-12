package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/controllers"
	"github.com/jezhtech/prince-group-backend/middleware"
)

func UserRoutes(router *gin.RouterGroup) {
	userRouter := router.Group("/user")

	userRouter.GET("/", middleware.UserMiddleware(), controllers.GetUserData)
	userRouter.POST("/", middleware.UserMiddleware(), controllers.CreateUser)
	userRouter.GET("/:id", middleware.UserMiddleware(), controllers.GetUser)
	userRouter.PUT("/", middleware.UserMiddleware(), controllers.UpdateUser)

	userRouter.GET("/all", middleware.AdminMiddleware(), controllers.GetAllUsers)
	userRouter.DELETE("/:id", middleware.AdminMiddleware(), controllers.DeleteUser)
}
