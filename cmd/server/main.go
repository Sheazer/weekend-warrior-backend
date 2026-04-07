package main

import (
	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Подключаем базу
	db.InitDB()

	r := gin.Default()

	r.GET("/user/:id", handlers.GetUserByIDHandler)
    // Сюда добавим POST для создания юзера чуть позже

	r.Run(":8080")
}