package handlers

import (
	"net/http"

	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// Тут мы делаем фильтрацию по категориям и дате и отдаем обратно
func GetActivities(c *gin.Context) {
	category := c.Query("category")
	date := c.Query("date")
	status := c.Query("status") // принимаем параметр status из URL

	var activities []models.Activity
	query := db.DB

	// Фильтр по статусу (если передан)
	if status != "" {
		if status == "all" {
			// Не добавляем фильтр, показываем всё
		} else if status == "active" || status == "finished" || status == "cancelled" {
			query = query.Where("status = ?", status)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid status. Allowed: active, finished, cancelled, all",
			})
			return
		}
	} else {
		// По умолчанию показываем только активные
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

	query.Find(&activities)
	c.JSON(http.StatusOK, activities)
}

// Тут мы принимаем JSON от фронтенда и сохраняем его в базу.
func CreateActivity(c *gin.Context) {
	var newActivity models.Activity

	if err := c.ShouldBindJSON(&newActivity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
