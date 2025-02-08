package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/inzamrzn918/aws-ses-mock/routes"
)

func main() {
	r := gin.Default()

	// Load API routes
	routes.SetupRoutes(r)

	// Start server
	port := "8080"
	fmt.Println("Server is running on port:", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
