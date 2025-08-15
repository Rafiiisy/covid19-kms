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
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fallback to individual environment variables
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_PASSWORD"),
			os.Getenv("DATABASE_NAME"),
		)
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("✅ Database connection established")
	return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
