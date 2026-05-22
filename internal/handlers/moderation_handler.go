package handlers

import (
	"fmt"
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


	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Проверяем, что текущий пользователь - организатор
		var activity models.Activity
		if err := tx.First(&activity, activityID).Error; err != nil {
			return err
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
        // // 1. Проверяем существование активности и права
        // var activity models.Activity
        // if err := tx.First(&activity, activityID).Error; err != nil {
        //     return err // Активность не найдена
        // }
        
        // if activity.OrganizerID != uint(getOrganizerID(c)) {
        //     return &NotOrganizerError{}
        // }

        // 2. МЕНЯЕМ СТАТУС ВМЕСТО УДАЛЕНИЯ
        // Используем .Model(), чтобы GORM знал, куда писать
        result := tx.Model(&models.Participant{}).
            Where("activity_id = ? AND user_id = ? AND status = ?", activityID, userID, "pending").
            Update("status", "rejected")

        if result.Error != nil {
            return result.Error
        }

        // ПРОВЕРКА: а была ли вообще такая запись? 
        // Если RowsAffected == 0, значит юзера нет в pending (уже одобрен или не подавал)
        if result.RowsAffected == 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "No pending request found for this user"})
            return fmt.Errorf("no records updated")
        }

        // 3. Логируем
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
        // Важно: не всегда ошибка значит "not organizer". Может быть "not found"
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    }
}

type NotOrganizerError struct{}

func (e *NotOrganizerError) Error() string { return "not organizer" }

type NoFreeSlotsError struct{}

func (e *NoFreeSlotsError) Error() string { return "no free slots" }
