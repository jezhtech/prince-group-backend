package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/controllers"
	"github.com/jezhtech/prince-group-backend/middleware"
)

func UserRoutes(router *gin.RouterGroup) {
	userRouter := router.Group("/user")

	userRouter.GET("/", middleware.CombinedAuthMiddleware(), controllers.GetUserData)
	userRouter.POST("/", middleware.CombinedAuthMiddleware(), controllers.CreateUser)
	userRouter.GET("/:id", middleware.CombinedAuthMiddleware(), controllers.GetUser)
	userRouter.PUT("/", middleware.CombinedAuthMiddleware(), controllers.UpdateUser)

	userRouter.GET("/all", middleware.AdminMiddleware(), controllers.GetAllUsers)
	userRouter.DELETE("/:id", middleware.AdminMiddleware(), controllers.DeleteUser)
}
