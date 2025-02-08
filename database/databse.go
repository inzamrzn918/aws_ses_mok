package database

import (
	"fmt"
	"log"

	"github.com/inzamrzn918/aws-ses-mock/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate tables
	err = DB.AutoMigrate(&models.Email{}, &models.EmailLog{})

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database connection established and migrations completed.")
}
