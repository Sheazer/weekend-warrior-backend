package handlers

import (
	"net/http"
	"strconv"

	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/models"
	"github.com/gin-gonic/gin"
)


func GetUserByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Ищем пользователя
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Получаем активности, где пользователь организатор
	var organizedActivities []models.Activity
	db.DB.Where("organizer_id = ?", userID).Find(&organizedActivities)

	// Получаем активности, где пользователь участник
	var participations []models.Participant
	db.DB.Where("user_id = ? AND status = ?", userID, "joined").
		Preload("Activity").
		Find(&participations)

	// Собираем активности участника
	var joinedActivities []models.Activity
	for _, p := range participations {
		joinedActivities = append(joinedActivities, p.Activity)
	}

	c.JSON(http.StatusOK, gin.H{
		"user":                 user,
		"organized_activities": organizedActivities,
		"joined_activities":    joinedActivities,
	})
}
