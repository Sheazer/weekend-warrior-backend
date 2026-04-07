package handlers

import (
	"net/http"

	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/models"
	"github.com/gin-gonic/gin"
)

//Тут мы делаем фильтрацию по категориям и дате и отдаем обратно
func GetActivities(c *gin.Context) {
	category := c.Query("category")
	date := c.Query("date")

	var activities []models.Activity
	query := db.DB

	// Динамически добавляем фильтры в SQL запрос
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if date != "" {
		query = query.Where("date = ?", date)
	}

	query.Find(&activities)
	c.JSON(http.StatusOK, activities)
}

//Тут мы принимаем JSON от фронтенда и сохраняем его в базу.
func CreateActivity(c *gin.Context) {
	var newActivity models.Activity

	if err := c.ShouldBindJSON(&newActivity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем статус по умолчанию
	newActivity.Status = "active"

	// Сохраняем в SQLite
	result := db.DB.Create(&newActivity)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать событие"})
		return
	}

	c.JSON(http.StatusCreated, newActivity)
}