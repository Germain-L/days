package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"days/internal/database"
	"days/internal/handlers"
	"days/internal/services"

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

	// Initialize services
	userService := services.NewUserService(db.Queries)
	calendarService := services.NewCalendarService(db.Queries)
	// colorMeaningService := services.NewColorMeaningService(db.Queries, calendarService)
	// dayEntryService := services.NewDayEntryService(db.Queries, calendarService, colorMeaningService)

	// Initialize server with handlers
	server := handlers.NewServer(userService, calendarService)

	// Setup routes
	mux := server.SetupRoutes()

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("API endpoints:")
	log.Printf("  POST   /api/users          - Create user")
	log.Printf("  POST   /api/auth/login     - Login")
	log.Printf("  GET    /api/users/{id}     - Get user")
	log.Printf("  GET    /api/calendars      - Get user calendars")
	log.Printf("  POST   /api/calendars      - Create calendar")
	log.Printf("  GET    /api/calendars/{id} - Get calendar")
	log.Printf("  PUT    /api/calendars/{id} - Update calendar")
	log.Printf("  DELETE /api/calendars/{id} - Delete calendar")
	log.Printf("  GET    /health             - Health check")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
