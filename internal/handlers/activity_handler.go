package handlers

import (
	"fmt"
	"net/http"

	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func GetActivities(c *gin.Context) {
    category := c.Query("category")
    date := c.Query("date")
    status := c.Query("status") 
    include := c.Query("include") 
	// 🔥 ДОБАВЬ ЭТИ ЛОГИ: они покажут, что видит сервер в консоли Go
    fmt.Println("=== ВХОДЯЩИЙ ЗАПРОС ===")
    fmt.Println("Полноценный URL:", c.Request.URL.String())
    fmt.Println("Значение include:", include)
    // ===================================

    var activities []models.Activity
    query := db.DB

    // Фильтр по статусу
    if status != "" {
        if status == "all" {
            // Показываем всё
        } else if status == "active" || status == "finished" || status == "cancelled" {
            query = query.Where("status = ?", status)
        } else {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "invalid status. Allowed: active, finished, cancelled, all",
            })
            return
        }
    } else {
        query = query.Where("status = ?", "active")
    }

    // Фильтр по категории
    if category != "" {
        query = query.Where("category = ?", category)
    }

    // Фильтр по дате
    if date != "" {
        query = query.Where("date = ?", date)
    }

    if include == "participants" {
		fmt.Println("🚀 МАГИЯ: Зашли в Preload!")
        query = query.Preload("Participants")
    }

    // Выполняем итоговый запрос
    if err := query.Find(&activities).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch activities"})
        return
    }
    
    c.JSON(http.StatusOK, activities)
}

func CreateActivity(c *gin.Context) {
    var newActivity models.Activity

    if err := c.ShouldBindJSON(&newActivity); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID, ok := c.Get("user_id") // Проверь, как именно называется ключ в твоем AuthMiddleware (userID или user_id)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не идентифицирован"})
        return
    }

    // Приводим тип из интерфейса к uint (так как в gorm.Model ID имеет тип uint)
    if id, ok := userID.(uint); ok {
        newActivity.OrganizerID = id
    } else if idFloat, ok := userID.(float64); ok { // На случай, если JWT парсит числа как float64
        newActivity.OrganizerID = uint(idFloat)
    }

    // Устанавливаем статус по умолчанию
    newActivity.Status = "active"
    newActivity.NeedModeration = true

    // Сохраняем в SQLite
    result := db.DB.Create(&newActivity)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать событие"})
        return
    }

    c.JSON(http.StatusCreated, newActivity)
}
