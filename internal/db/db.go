package db

import (
	"github.com/Erzhan/weekend-warrior-backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// Создаст файл warriors.db в корне проекта
	DB, err = gorm.Open(sqlite.Open("warriors.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Автоматическое создание таблицы на основе структуры User
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Activity{})
	DB.AutoMigrate(&models.Participant{})
	DB.AutoMigrate(&models.ActivityLog{})
	DB.AutoMigrate(&models.Message{})
}
