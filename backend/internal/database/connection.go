package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"days/internal/db"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Database struct {
	DB      *sql.DB
	Queries *db.Queries
}

// NewConfig creates database config from environment variables
func NewConfig() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "days"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// Connect establishes a connection to the database
func Connect(config *Config) (*Database, error) {
	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	// Open database connection
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)

	// Create queries instance
	queries := db.New(sqlDB)

	log.Printf("Connected to database: %s:%s/%s", config.Host, config.Port, config.DBName)

	return &Database{
		DB:      sqlDB,
		Queries: queries,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}

// getEnv gets environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
