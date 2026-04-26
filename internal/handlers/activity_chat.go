package handlers

import (
	"net/http"
	"strconv"

	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// GetActivityMessages принимает activity_id из URL и возвращает список сообщений
func GetActivityMessages(c *gin.Context) {
    // Извлекаем ID из параметров пути (например, /activities/:id/messages)
    activityID := c.Param("id")

    var messages []models.Message
    
    // Выполняем запрос к БД с фильтрацией по activity_id
    err := db.DB.Where("activity_id = ?", activityID).
        Order("sent_at asc").
        Find(&messages).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
        return
    }

    c.JSON(http.StatusOK, messages)
}


func CreateMessage(c *gin.Context) {
	activityID, _:= strconv.Atoi(c.Param("id"))

	var newMessage models.Message

	newMessage.ActivityID = activityID


	if err := c.ShouldBindJSON(&newMessage); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

	if err := db.DB.Create(&newMessage).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save message"})
        return
    }
	c.JSON(http.StatusCreated, newMessage)
}