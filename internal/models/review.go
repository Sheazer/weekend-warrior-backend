package models

import (
	"time"
)

type Review struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ReviewerID uint      `json:"reviewer_id"`                                     // кто оставляет отзыв
	RevieweeID uint      `json:"reviewee_id"`                                     // о ком отзыв (organizer)
	ActivityID uint      `json:"activity_id"`                                     // в рамках какой активности
	Rating     int       `json:"rating" gorm:"check:rating >= 1 AND rating <= 5"` // 1-5
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	//DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Reviewer User     `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID"`
	Activity Activity `json:"activity,omitempty" gorm:"foreignKey:ActivityID"`
}
