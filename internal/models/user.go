package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // Авто-поля: ID, дата создания и т.д.
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"` // Почта должна быть уникальной
}