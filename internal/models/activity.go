package models

import "gorm.io/gorm"

type Activity struct {
	gorm.Model
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	Category       string        `json:"category"`                     // Например, "sports"
	Date           string        `json:"date"`                         // Для простоты пока строка, потом перейдем на time.Time
	Status         string        `json:"status" gorm:"default:active"` // active, cancelled, finished
	MaxPeople      int           `json:"max_people"`
	OrganizerID    uint          `json:"organizer_id"`
	NeedModeration bool          `json:"need_moderation"` // true = ручное подтверждение
	Organizer      User          `json:"organizer,omitempty" gorm:"foreignKey:OrganizerID"`
	Participants   []Participant `json:"participants,omitempty"`
}
