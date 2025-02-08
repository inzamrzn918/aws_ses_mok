package models

import (
	"time"

	"gorm.io/gorm"
)

// EmailLog stores metadata about sent emails
type EmailLog struct {
	ID        uint      `gorm:"primaryKey"`
	From      string    `gorm:"not null"`
	To        string    `gorm:"not null"`
	Status    string    `gorm:"not null"` // e.g., "sent", "blocked", "failed"
	CreatedAt time.Time // Timestamp of email attempt
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
