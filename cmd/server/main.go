package main

import (
	"time" // time

	_ "github.com/Erzhan/weekend-warrior-backend/docs" // Важно: импорт сгенерированных доков
	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/handlers"
	"github.com/gin-contrib/cors" // cors import
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // импорт swagger
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

	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api.GET("/user/:id", handlers.GetUserByIDHandler)
	api.GET("/activities", handlers.GetActivities)
	api.POST("/activities", handlers.CreateActivity)
	api.POST("/activities/:id/join", handlers.JoinActivityHandler)
	api.PUT("/activities/:id/participants/:user_id/approve", handlers.ApproveParticipantHandler)
	api.DELETE("/activities/:id/participants/:user_id/reject", handlers.RejectParticipantHandler)
	api.PATCH("/activities/:id/status", handlers.UpdateActivityStatusHandler)
	api.GET("/users/:id/feedback", handlers.GetOrganizerFeedbackHandler)
	api.POST("/activities/:id/review", handlers.CreateReviewHandler)
	api.GET("/activities/:id/reviews", handlers.GetActivityReviewsHandler)
	api.GET("/activities/:id/chat", handlers.GetActivityMessages)
	api.POST("/activities/:id/chat", handlers.CreateMessage)
	api.Run(":8080")
}
