package models

import "time"

type ActivityLog struct {
	ID         uint      `gorm:"primaryKey"`
	ActivityID uint      `json:"activity_id"`
	UserID     uint      `json:"user_id"`
	Action     string    `json:"action"`  // join_request, leave, approved, rejected
	Details    string    `json:"details"` // доп. информация
	Timestamp  time.Time `json:"timestamp" gorm:"autoCreateTime"`
}
