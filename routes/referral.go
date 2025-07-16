package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/controllers"
	"github.com/jezhtech/prince-group-backend/middleware"
)

func ReferralRoutes(router *gin.RouterGroup) {
	referralRouter := router.Group("/referral")

	referralRouter.GET("/:id", middleware.UserMiddleware(), controllers.GetReferral)
	referralRouter.GET("/all", middleware.UserMiddleware(), controllers.GetAllReferrals)
	referralRouter.POST("/", middleware.AdminMiddleware(), controllers.CreateReferral)
	referralRouter.PUT("/:id", middleware.AdminMiddleware(), controllers.UpdateReferral)
	referralRouter.DELETE("/:id", middleware.AdminMiddleware(), controllers.DeleteReferral)
	referralRouter.GET("/check-referral", middleware.UserMiddleware(), controllers.CheckReferral)
}
