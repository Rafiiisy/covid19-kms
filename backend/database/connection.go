package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDatabase initializes the database connection
func InitDatabase() error {
	// Check if database should be skipped
	if os.Getenv("SKIP_DATABASE") == "true" {
		log.Println("⚠️ Database initialization skipped (SKIP_DATABASE=true)")
		return nil
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fallback to individual environment variables
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_PASSWORD"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"),
		)
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Configure connection pooling
	DB.SetMaxOpenConns(25)   // Maximum number of open connections
	DB.SetMaxIdleConns(5)    // Maximum number of idle connections
	DB.SetConnMaxLifetime(0) // Connections don't expire

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("✅ Database connection established")
	return nil
}

// EnsureConnection ensures the database connection is alive
func EnsureConnection() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// Ping the database to check if connection is alive
	if err := DB.Ping(); err != nil {
		// Try to reconnect
		log.Println("⚠️ Database connection lost, attempting to reconnect...")
		if err := InitDatabase(); err != nil {
			return fmt.Errorf("failed to reconnect to database: %v", err)
		}
		log.Println("✅ Database reconnection successful")
	}

	return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
