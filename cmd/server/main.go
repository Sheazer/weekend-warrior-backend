package main

import (
	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // импорт swagger
	// Важно: импорт сгенерированных доков
)

// @title Weekend Warrior API
// @version 1.0
// @description Это сервер для поиска активностей на выходные.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 1. Подключаем базу
	db.InitDB()

	api := gin.Default()

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api.GET("/user/:id", handlers.GetUserByIDHandler)
	api.GET("/activities", handlers.GetActivities)
    api.POST("/activities", handlers.CreateActivity)

	api.Run(":8080")
}