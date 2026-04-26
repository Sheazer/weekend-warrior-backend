package handlers

import (
	"net/http"
	"strconv"

	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type StatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active cancelled finished"`
}

// UpdateActivityStatusHandler меняет статус активности
// PATCH /api/activities/:id/status
func UpdateActivityStatusHandler(c *gin.Context) {
	activityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Находим активность
	var activity models.Activity
	if err := db.DB.First(&activity, activityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
		return
	}

	// Сохраняем старый статус
	oldStatus := activity.Status

	// Обновляем статус
	activity.Status = req.Status
	if err := db.DB.Save(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	// Логируем изменение (опционально)
	log := models.ActivityLog{
		ActivityID: uint(activityID),
		UserID:     activity.OrganizerID,
		Action:     "status_change",
		Details:    "from " + oldStatus + " to " + req.Status,
	}
	db.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"message":     "status updated successfully",
		"activity_id": activityID,
		"old_status":  oldStatus,
		"new_status":  req.Status,
	})
}
