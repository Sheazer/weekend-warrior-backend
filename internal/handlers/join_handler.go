package handlers

import (
	"net/http"
	"strconv"

	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JoinRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// JoinActivityHandler обрабатывает POST /api/activities/:id/join
func JoinActivityHandler(c *gin.Context) {
    activityID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
        return
    }

    // 🔥 БЕРЕМ ИЗ КОНТЕКСТА: Кто авторизован — тот и присоединяется
    userIDFromContext, exists := c.Get("user_id") 
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authorized"})
        return
    }

    // Безопасное приведение типа интерфейса к uint
    var currentUserID uint
    switch v := userIDFromContext.(type) {
    case uint:
        currentUserID = v
    case int:
        currentUserID = uint(v)
    case float64:
        currentUserID = uint(v)
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id format in token"})
        return
    }

    // НАЧАЛО ТРАНЗАКЦИИ (гарантирует конкурентность)
    err = db.DB.Transaction(func(tx *gorm.DB) error {
        // 1. Получаем активность с блокировкой строки (FOR UPDATE)
        var activity models.Activity
        if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&activity, activityID).Error; err != nil {
            return err
        }

        // 2. Проверяем, не участвует ли уже пользователь (используем currentUserID)
        var existing models.Participant
        if err := tx.Where("user_id = ? AND activity_id = ?", currentUserID, activityID).
            First(&existing).Error; err == nil {
            return &UserAlreadyJoinedError{Status: existing.Status}
        }

        // 3. Считаем текущее количество подтверждённых участников
        var joinedCount int64
        tx.Model(&models.Participant{}).
            Where("activity_id = ? AND status = ?", activityID, "joined").
            Count(&joinedCount)

        // 4. Определяем статус нового участника
        var newStatus string
        var message string

        if joinedCount < int64(activity.MaxPeople) {
            if activity.NeedModeration {
                newStatus = "pending"
                message = "Your request is pending approval by the organizer"
            } else {
                newStatus = "joined"
                message = "Successfully joined the activity"
            }
        } else {
            // Мест нет → лист ожидания
            newStatus = "waitlist"
            message = "Activity is full, you are on the waitlist"
        }

        // 5. Создаём запись участника (используем currentUserID)
        participant := models.Participant{
            UserID:     currentUserID,
            ActivityID: uint(activityID),
            Status:     newStatus,
        }
        if err := tx.Create(&participant).Error; err != nil {
            return err
        }

        // 6. Логируем действие (для отладки конкурентности)
        log := models.ActivityLog{
            ActivityID: uint(activityID),
            UserID:     currentUserID,
            Action:     "join_request",
            Details:    "status=" + newStatus,
        }
        tx.Create(&log)

        c.JSON(http.StatusOK, gin.H{
            "status":  newStatus,
            "message": message,
            "activity": gin.H{
                "id":           activity.ID,
                "title":        activity.Title,
                "joined_count": joinedCount + 1,
                "max_people":   activity.MaxPeople,
            },
        })
        return nil
    })

    if err != nil {
        switch e := err.(type) {
        case *UserAlreadyJoinedError:
            c.JSON(http.StatusConflict, gin.H{
                "error":  "user already participated",
                "status": e.Status,
            })
        default:
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "failed to join activity",
            })
        }
    }
}

// UserAlreadyJoinedError кастомная ошибка
type UserAlreadyJoinedError struct {
	Status string
}

func (e *UserAlreadyJoinedError) Error() string {
	return "user already joined with status: " + e.Status
}
