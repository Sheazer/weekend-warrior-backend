package handlers

import (
	"net/http"
	"strconv"

	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ApproveParticipantHandler подтверждает участника (pending -> joined)
func ApproveParticipantHandler(c *gin.Context) {
	activityID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := strconv.Atoi(c.Param("user_id"))

	// Кто подтверждает? (нужно получить из JWT или заголовка)
	organizerID := getOrganizerID(c) // пока заглушка

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Проверяем, что текущий пользователь - организатор
		var activity models.Activity
		if err := tx.First(&activity, activityID).Error; err != nil {
			return err
		}
		if activity.OrganizerID != uint(organizerID) {
			return &NotOrganizerError{}
		}

		// 2. Находим заявку со статусом pending
		var participant models.Participant
		if err := tx.Where("activity_id = ? AND user_id = ? AND status = ?",
			activityID, userID, "pending").
			First(&participant).Error; err != nil {
			return err
		}

		// 3. Проверяем, есть ли ещё свободные места
		var joinedCount int64
		tx.Model(&models.Participant{}).
			Where("activity_id = ? AND status = ?", activityID, "joined").
			Count(&joinedCount)

		if joinedCount >= int64(activity.MaxPeople) {
			return &NoFreeSlotsError{}
		}

		// 4. Подтверждаем
		participant.Status = "joined"
		tx.Save(&participant)

		// 5. Логируем
		log := models.ActivityLog{
			ActivityID: uint(activityID),
			UserID:     uint(userID),
			Action:     "approved",
			Details:    "approved by organizer",
		}
		tx.Create(&log)

		c.JSON(http.StatusOK, gin.H{
			"status":  "approved",
			"message": "User has been approved to join",
		})
		return nil
	})

	if err != nil {
		switch err.(type) {
		case *NotOrganizerError:
			c.JSON(http.StatusForbidden, gin.H{"error": "only organizer can approve"})
		case *NoFreeSlotsError:
			c.JSON(http.StatusConflict, gin.H{"error": "no free slots left"})
		default:
			c.JSON(http.StatusNotFound, gin.H{"error": "pending request not found"})
		}
	}
}

// RejectParticipantHandler отклоняет заявку
func RejectParticipantHandler(c *gin.Context) {
	activityID, _ := strconv.Atoi(c.Param("id"))
	userID, _ := strconv.Atoi(c.Param("user_id"))

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Проверяем, что текущий пользователь - организатор
		var activity models.Activity
		tx.First(&activity, activityID)
		if activity.OrganizerID != uint(getOrganizerID(c)) {
			return &NotOrganizerError{}
		}

		// Удаляем заявку (или меняем статус на rejected)
		if err := tx.Where("activity_id = ? AND user_id = ? AND status = ?",
			activityID, userID, "pending").
			Delete(&models.Participant{}).Error; err != nil {
			return err
		}

		// Логируем
		log := models.ActivityLog{
			ActivityID: uint(activityID),
			UserID:     uint(userID),
			Action:     "rejected",
			Details:    "rejected by organizer",
		}
		tx.Create(&log)

		c.JSON(http.StatusOK, gin.H{
			"status":  "rejected",
			"message": "User request has been rejected",
		})
		return nil
	})

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "only organizer can reject"})
	}
}

// Вспомогательные функции и ошибки
func getOrganizerID(c *gin.Context) int {
	// TODO: заменить на реальное получение user_id из JWT токена
	// Пока заглушка
	return 1
}

type NotOrganizerError struct{}

func (e *NotOrganizerError) Error() string { return "not organizer" }

type NoFreeSlotsError struct{}

func (e *NoFreeSlotsError) Error() string { return "no free slots" }
