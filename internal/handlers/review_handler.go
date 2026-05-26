package handlers

import (
	"net/http"
	"strconv"

	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type CreateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required"`
	Comment string `json:"comment" binding:"required"`
}

// CreateReviewHandler создаёт новый отзыв об организаторе
// POST /api/activities/:id/review
func CreateReviewHandler(c *gin.Context) {
	activityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Проверяем, что активность существует
	var activity models.Activity
	if err := db.DB.First(&activity, activityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
		return
	}

	// 2. Проверяем, что активность ЗАВЕРШЕНА (только тогда можно оставить отзыв)
	if activity.Status != "finished" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "can only leave a review for finished activities",
			"current_status": activity.Status,
		})
		return
	}

	// 🔒 Извлекаем ID авторизованного пользователя из контекста JWT
	userIDFromContext, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authorized"})
		return
	}

	var currentReviewerID uint
	switch v := userIDFromContext.(type) {
	case uint: currentReviewerID = v
	case int: currentReviewerID = uint(v)
	case float64: currentReviewerID = uint(v)
	}

	// 3. Проверяем, что reviewer УЧАСТВОВАЛ в этой активности
	var participant models.Participant
	if err := db.DB.Where("activity_id = ? AND user_id = ? AND status = ?", 
		activityID, currentReviewerID, "joined").First(&participant).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "you must be a participant to leave a review"})
		return
	}

	// 4. Проверяем, что reviewer НЕ организатор
	if currentReviewerID == activity.OrganizerID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organizer cannot review themselves"})
		return
	}

	// 5. Проверяем, что отзыв ещё не оставлен
	var existingReview models.Review
	if err := db.DB.Where("reviewer_id = ? AND activity_id = ?", 
		currentReviewerID, activityID).First(&existingReview).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "you have already left a review for this activity"})
		return
	}

	// 6. Создаём отзыв (используем безопасный currentReviewerID)
	review := models.Review{
		ReviewerID: currentReviewerID,
		RevieweeID: activity.OrganizerID,
		ActivityID: uint(activityID),
		Rating:     req.Rating,
		Comment:    req.Comment,
	}

	if err := db.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
			"message": "review created successfully",
			"review": gin.H{
				"id":          review.ID, // или review.id, в зависимости от вашей базовой модели gorm.Model
				"reviewer_id": review.ReviewerID,
				"reviewee_id": review.RevieweeID,
				"activity_id": review.ActivityID,
				"rating":      review.Rating,
				"comment":     review.Comment,
			},
		})
}

// GetOrganizerFeedbackHandler получает отзывы об организаторе
// GET /api/users/:id/feedback
func GetOrganizerFeedbackHandler(c *gin.Context) {
	organizerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Проверяем, что пользователь существует
	var user models.User
	if err := db.DB.First(&user, organizerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Получаем все отзывы об этом организаторе
	var reviews []models.Review
	db.DB.Where("reviewee_id = ?", organizerID).
		Preload("Reviewer"). // подгружаем данные того, кто оставил отзыв
		Preload("Activity"). // подгружаем данные активности
		Order("created_at DESC").
		Find(&reviews)

	// Считаем средний рейтинг
	var avgRating float64
	if len(reviews) > 0 {
		var sum int
		for _, r := range reviews {
			sum += r.Rating
		}
		avgRating = float64(sum) / float64(len(reviews))
	}

	c.JSON(http.StatusOK, gin.H{
		"organizer_id":   organizerID,
		"organizer_name": user.Name,
		"total_reviews":  len(reviews),
		"average_rating": avgRating,
		"reviews":        reviews,
	})
}

// GetActivityReviewsHandler получает все отзывы об активности (для организатора)
// GET /api/activities/:id/reviews
func GetActivityReviewsHandler(c *gin.Context) {
	activityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	var reviews []models.Review
	db.DB.Where("activity_id = ?", activityID).
		Preload("Reviewer").
		Order("created_at DESC").
		Find(&reviews)

	c.JSON(http.StatusOK, gin.H{
		"activity_id":   activityID,
		"total_reviews": len(reviews),
		"reviews":       reviews,
	})
}
