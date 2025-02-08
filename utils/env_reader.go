package utils

import (
	"log"
	"os"
	"strconv"
)

// Get email limit per hour from .env
func GetEmailLimitPerHour() int {
	limit, err := strconv.Atoi(os.Getenv("EMAIL_LIMIT_PER_HOUR"))
	if err != nil {
		return 5 // Default fallback
	}
	return limit
}

// Log email sending (for debugging only, not storing content)
func LogEmailSent(from string, to []string) {
	log.Printf("Email sent from %s to %v", from, to)
}
