package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inzamrzn918/aws-ses-mock/database"
	"github.com/inzamrzn918/aws-ses-mock/models"
)

// Subscribe to Newsletter
func SubscribeNewsletter(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Category string `json:"category" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Check if email already exists
	var email models.Email
	result := database.DB.Where("email_address = ?", request.Email).First(&email)

	if result.RowsAffected > 0 {
		email.IsSubscribed = true
		email.Category = request.Category
		database.DB.Save(&email)
	} else {
		// Create new entry if email was never stored
		newEmail := models.Email{
			EmailAddress: request.Email,
			IsSubscribed: true,
			Category:     request.Category,
		}
		database.DB.Create(&newEmail)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscribed successfully"})
}

// Unsubscribe from Newsletter
func UnsubscribeNewsletter(c *gin.Context) {
	email := c.Param("email")

	// Check if email exists
	var emailEntry models.Email
	result := database.DB.Where("email_address = ?", email).First(&emailEntry)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found in subscriber list"})
		return
	}

	// Unsubscribe email
	emailEntry.IsSubscribed = false
	emailEntry.Category = ""
	database.DB.Save(&emailEntry)

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}

// List all newsletter subscribers
func ListNewsletterSubscribers(c *gin.Context) {
	var subscribers []models.Email
	database.DB.Select("email_address, category").Where("is_subscribed = ?", true).Find(&subscribers)

	c.JSON(http.StatusOK, subscribers)
}
