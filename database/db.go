package database

import (
	"log"
	"messaging-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "postgresql://postgres:postgres@localhost:5432/messaging_service?sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	// Auto-migrate the models
	DB.AutoMigrate(&models.User{}, &models.Message{})
	log.Println("Database connection established and tables migrated")
}
