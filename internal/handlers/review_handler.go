package handlers

import (
	"net/http"
	"strconv"

	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type CreateReviewRequest struct {
	ReviewerID uint   `json:"reviewer_id" binding:"required"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Comment    string `json:"comment" binding:"required,max=500"`
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

	// 3. Проверяем, что reviewer УЧАСТВОВАЛ в этой активности (статус joined)
	var participant models.Participant
	if err := db.DB.Where("activity_id = ? AND user_id = ? AND status = ?",
		activityID, req.ReviewerID, "joined").First(&participant).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you must be a participant to leave a review",
		})
		return
	}

	// 4. Проверяем, что reviewer НЕ организатор (нельзя оставить отзыв самому себе)
	if req.ReviewerID == activity.OrganizerID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "organizer cannot review themselves",
		})
		return
	}

	// 5. Проверяем, что отзыв ещё не оставлен (один отзыв на активность от пользователя)
	var existingReview models.Review
	if err := db.DB.Where("reviewer_id = ? AND activity_id = ?",
		req.ReviewerID, activityID).First(&existingReview).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "you have already left a review for this activity",
		})
		return
	}

	// 6. Создаём отзыв
	review := models.Review{
		ReviewerID: req.ReviewerID,
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
		"review":  review,
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
