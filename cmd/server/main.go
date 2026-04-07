package main

import (
	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Подключаем базу
	db.InitDB()

	api := gin.Default()

	api.GET("/user/:id", handlers.GetUserByIDHandler)
	api.GET("/activities", handlers.GetActivities)
    api.POST("/activities", handlers.CreateActivity)

	api.Run(":8080")
}