package main

import (
	"log"

	"days/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Days Backend Server...")

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load database configuration
	config := database.NewConfig()

	// Log configuration (without password)
	log.Printf("Attempting to connect to database: %s:%s/%s", config.Host, config.Port, config.DBName)

	// Connect to database
	db, err := database.Connect(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("Database connection successful!")
	log.Println("Server setup complete. Press Ctrl+C to exit.")

	// Keep the program running for testing
	select {}
}
