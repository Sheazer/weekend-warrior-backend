package models

import "gorm.io/gorm"

type Activity struct {
	gorm.Model
	Title          string `json:"title"`
	Description    string `json:"description"`
	Category       string `json:"category"` // Например, "sports"
	Date           string `json:"date"`     // Для простоты пока строка, потом перейдем на time.Time
	Status         string `json:"status"`   // "active", "cancelled", "finished"
	MaxPeople      int    `json:"max_people"`
	OrganizerID    uint   `json:"organizer_id"`
	NeedModeration bool   `json:"need_moderation"` // true = ручное подтверждение
}

type Message struct {
	gorm.Model
	SenderID	   int    `json:"sender_id"` 
	ActivityID     int    `json:"activity_id"`
	Content        string `json:"content"`
	SentAt         string `json:"sent_at"`
}
