package models

import (
	"time"

	"gorm.io/gorm"
)

type Participant struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `json:"user_id"`
	ActivityID uint           `json:"activity_id"`
	Status     string         `json:"status"` // pending, joined, waitlist
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	User       User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Activity   Activity       `json:"activity,omitempty" gorm:"foreignKey:ActivityID"`
}
