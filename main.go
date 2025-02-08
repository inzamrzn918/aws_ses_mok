package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/inzamrzn918/aws-ses-mock/database"
	"github.com/inzamrzn918/aws-ses-mock/routes"
	"github.com/joho/godotenv"
)

func init() {
	setupLogger()
	database.InitDB()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func setupLogger() {
	file, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}

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
