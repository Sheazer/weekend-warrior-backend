package models

import "gorm.io/gorm"

type Activity struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Date        string `json:"date"`
	Status      string `json:"status" gorm:"default:active"`
	MaxPeople   int    `json:"max_people"`

	// координаты
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	OrganizerID    uint `json:"organizer_id"`
	NeedModeration bool `json:"need_moderation"`

	Organizer    User          `json:"organizer,omitempty" gorm:"foreignKey:OrganizerID"`
	Participants []Participant `json:"participants,omitempty"`
}

type Message struct {
	gorm.Model
	SenderID   int    `json:"sender_id"`
	ActivityID int    `json:"activity_id"`
	Content    string `json:"content"`
}
