package models

import (
	"time"

	"gorm.io/gorm"
)

// Email model to handle both blocked emails and newsletter subscriptions
type Email struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	EmailAddress     string         `gorm:"uniqueIndex;not null" json:"email_address"`
	IsBlocked        bool           `json:"is_blocked"`
	Reason           string         `json:"reason,omitempty"` // Optional: Why the email is blocked
	SelfCooldownDate *time.Time     `json:"self_cooldown_date,omitempty"`
	SelfCooldownDays int            `json:"self_cooldown_days,omitempty"`
	IsSubscribed     bool           `json:"is_subscribed"`      // Newsletter subscription
	Category         string         `json:"category,omitempty"` // Example: "Tech", "Business", etc.
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete support
}
