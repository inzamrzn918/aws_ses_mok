package controllers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	emailCount int
	mu         sync.Mutex
)

// SendEmail - It will send email to the given email address
// It will return the message id and status of the email
func SendEmail(c *gin.Context) {
	var request struct {
		To      string `json:"to" binding:"required"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	mu.Lock()
	emailCount++
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"MessageId": "mock-message-id-12345",
		"Status":    "Success",
	})
}

// GetStats - Returns number of emails "sent"
func GetStats(c *gin.Context) {
	mu.Lock()
	count := emailCount
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"emails_sent": count,
	})
}
