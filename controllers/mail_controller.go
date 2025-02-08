package controllers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inzamrzn918/aws-ses-mock/database"
	"github.com/inzamrzn918/aws-ses-mock/models"
	"github.com/inzamrzn918/aws-ses-mock/utils"
)

var emailCount = make(map[string]int)
var emailCountLock = sync.Mutex{}

func SendEmail(c *gin.Context) {
	var request struct {
		From    string   `json:"from" binding:"required,email"`
		To      []string `json:"to" binding:"required"`
		Subject string   `json:"subject" binding:"required"`
		Body    string   `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Check if sender is blocked
	if isBlocked, reason := isBlockedEmail(request.From); isBlocked {
		database.DB.Create(&models.EmailLog{
			From:   request.From,
			To:     "N/A",
			Status: "blocked",
		})
		c.JSON(http.StatusForbidden, gin.H{"error": "Sender is blocked", "reason": reason})
		return
	}

	// Check if any recipient is blocked
	for _, recipient := range request.To {
		if blocked, _ := isBlockedEmail(recipient); blocked {
			database.DB.Create(&models.EmailLog{
				From:   request.From,
				To:     recipient,
				Status: "blocked",
			})
			c.JSON(http.StatusForbidden, gin.H{"error": "One or more recipients are blocked"})
			return
		}
	}

	// Enforce email sending limits
	emailCountLock.Lock()
	currentHour := time.Now().Format("2006-01-02 15")
	senderKey := request.From + "_" + currentHour

	if emailCount[senderKey] >= utils.GetEmailLimitPerHour() {
		database.DB.Create(&models.EmailLog{
			From:   request.From,
			To:     "N/A",
			Status: "rate_limited",
		})
		emailCountLock.Unlock()
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Email sending limit reached"})
		return
	}
	emailCount[senderKey]++
	emailCountLock.Unlock()

	// Log email attempt in database (without content)
	for _, recipient := range request.To {
		database.DB.Create(&models.EmailLog{
			From:   request.From,
			To:     recipient,
			Status: "sent",
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email processed successfully"})
}

func GetStats(c *gin.Context) {
	var totalEmailsSent int64
	var totalBlocked int64
	var totalRateLimited int64

	database.DB.Model(&models.EmailLog{}).Where("status = ?", "sent").Count(&totalEmailsSent)
	database.DB.Model(&models.EmailLog{}).Where("status = ?", "blocked").Count(&totalBlocked)
	database.DB.Model(&models.EmailLog{}).Where("status = ?", "rate_limited").Count(&totalRateLimited)

	c.JSON(http.StatusOK, gin.H{
		"total_emails_sent":    totalEmailsSent,
		"total_emails_blocked": totalBlocked,
		"total_rate_limited":   totalRateLimited,
	})
}

// Block an email
func BlockEmail(c *gin.Context) {
	var request struct {
		Email  string `json:"email" binding:"required,email"`
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Check if email already exists
	var email models.Email
	result := database.DB.Where("email_address = ?", request.Email).First(&email)

	if result.RowsAffected > 0 {
		email.IsBlocked = true
		email.Reason = request.Reason
		database.DB.Save(&email)
	} else {
		// Create new entry if email was never stored
		newEmail := models.Email{
			EmailAddress: request.Email,
			IsBlocked:    true,
			Reason:       request.Reason,
		}
		database.DB.Create(&newEmail)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email blocked successfully"})
}

// Unblock an email
func UnblockEmail(c *gin.Context) {
	email := c.Param("email")

	// Check if email exists
	var emailEntry models.Email
	result := database.DB.Where("email_address = ?", email).First(&emailEntry)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found in blocked list"})
		return
	}

	// Unblock email
	emailEntry.IsBlocked = false
	emailEntry.Reason = ""
	database.DB.Save(&emailEntry)

	c.JSON(http.StatusOK, gin.H{"message": "Email unblocked successfully"})
}

// List all blocked emails
func ListBlockedEmails(c *gin.Context) {
	var blockedEmails []models.Email
	database.DB.Select("email_address, is_blocked, reason, self_cooldown_date, self_cooldown_days").Where("is_blocked = ?", true).Find(&blockedEmails)

	c.JSON(http.StatusOK, blockedEmails)
}

// Check if an email is blocked
func isBlockedEmail(email string) (bool, string) {
	var emailEntry models.Email
	if database.DB.Where("email_address = ?", email).First(&emailEntry).RowsAffected > 0 {
		if emailEntry.IsBlocked {
			return true, "Email is blocked: " + emailEntry.Reason
		}
	}
	return false, ""
}

// Self Block an Email (Cooldown System)
func SelfBlockEmail(c *gin.Context) {
	var request struct {
		Email            string `json:"email" binding:"required,email"`
		SelfCooldownDays int    `json:"self_cooldown_days" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Check if email already exists
	var email models.Email
	result := database.DB.Where("email_address = ?", request.Email).First(&email)

	if result.RowsAffected > 0 {
		email.IsBlocked = true
		email.Reason = "Self cooldown"
		email.SelfCooldownDate = &[]time.Time{time.Now().AddDate(0, 0, request.SelfCooldownDays)}[0]
		email.SelfCooldownDays = request.SelfCooldownDays
		database.DB.Save(&email)
	} else {
		// Create new entry if email was never stored
		newEmail := models.Email{
			EmailAddress:     request.Email,
			IsBlocked:        true,
			Reason:           "Self cooldown",
			SelfCooldownDate: &[]time.Time{time.Now().AddDate(0, 0, request.SelfCooldownDays)}[0],
			SelfCooldownDays: request.SelfCooldownDays,
		}
		database.DB.Create(&newEmail)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email is now in self-cooldown"})
}
