package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/controllers"
	"github.com/jezhtech/prince-group-backend/middleware"
)

func YouTubeRoutes(router *gin.RouterGroup) {
	youtubeRouter := router.Group("/youtube")

	youtubeRouter.GET("/check-subscription", middleware.UserMiddleware(), controllers.CheckYouTubeSubscription)
}
