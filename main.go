package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jezhtech/prince-group-backend/config"
	"github.com/jezhtech/prince-group-backend/routes"
)

func main() {
	router := gin.Default()
	config.InitDatabase()
	config.InitFirebase()
	InitAutoMigrate()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type", "X-Google-Access-Token"}
	router.Use(cors.New(corsConfig))

	routes.AppRouter(router)

	router.Run(":8000")
}
