package main

import (
	"backend-crawler/config"
	"backend-crawler/handlers"
	"backend-crawler/middleware"
	"backend-crawler/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	r := gin.Default()
	r.Use(middleware.AuthMiddleware())

	handlers.InitDB()
	go services.StartWorker()

	r.POST("/api/urls", handlers.AddURL)
	r.GET("/api/urls", handlers.ListURLs)
	r.GET("/api/urls/:id", handlers.GetURLDetails)
	r.POST("/api/urls/:id/reanalyze", handlers.ReanalyzeURL)
	r.DELETE("/api/urls", handlers.DeleteURLs)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
