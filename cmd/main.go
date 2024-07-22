package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"media-service/infrastructure/router"
	"os"
)

func main() {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the router
	r := router.Route()

	// Start the server
	startServer(r)
}

func startServer(router *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
