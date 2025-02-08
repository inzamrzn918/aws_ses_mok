package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/inzamrzn918/aws-ses-mock/controllers"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/send-email", controllers.SendEmail)
	r.GET("/stats", controllers.GetStats)

	r.POST("/block-email", controllers.BlockEmail)
	r.DELETE("/unblock-email/:email", controllers.UnblockEmail)
	r.GET("/blocked-emails", controllers.ListBlockedEmails)

}
