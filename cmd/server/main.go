package main

import (
	_ "github.com/Erzhan/weekend-warrior-backend/docs"
	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/handlers"
	"github.com/Erzhan/weekend-warrior-backend/internal/middleware"
	"github.com/gin-contrib/cors"
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
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // Для локальной разработки — идеально
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"} // 🔥 Явно разрешаем Authorization

	api.Use(cors.New(config))

	// --- 🔴 ПУБЛИЧНЫЕ РОУТЫ (Доступны абсолютно всем) ---
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api.POST("/register", handlers.RegisterHandler)
	api.POST("/login", handlers.LoginHandler)
	
	// Главную страницу и просмотр ивентов оставляем открытыми, чтобы неавторизованные гости тоже видели список
	api.GET("/activities", handlers.GetActivities)
	api.GET("/user/:id", handlers.GetUserByIDHandler)
	api.GET("/users/:id/feedback", handlers.GetOrganizerFeedbackHandler)
	api.GET("/activities/:id/reviews", handlers.GetActivityReviewsHandler)

	// --- 🔒 ЗАЩИЩЕННЫЕ РОУТЫ (Только для тех, у кого есть JWT-токен в Headers) ---
	// Создаем изолированную группу роутов
	protected := api.Group("/api")
	
	// Подключаем Middleware авторизации к этой группе
	protected.Use(middleware.AuthMiddleware())
	{
		// Все роуты внутри этих фигурных скобок автоматически требуют токен!
		protected.POST("/activities", handlers.CreateActivity)
		protected.POST("/activities/:id/join", handlers.JoinActivityHandler)
		protected.PUT("/activities/:id/participants/:user_id/approve", handlers.ApproveParticipantHandler)
		protected.DELETE("/activities/:id/participants/:user_id/reject", handlers.RejectParticipantHandler)
		protected.PATCH("/activities/:id/status", handlers.UpdateActivityStatusHandler)
		protected.POST("/activities/:id/review", handlers.CreateReviewHandler)
		protected.GET("/activities/:id/chat", handlers.GetActivityMessages)
		protected.POST("/activities/:id/chat", handlers.CreateMessage)
	}   

	api.Run(":8080")
}