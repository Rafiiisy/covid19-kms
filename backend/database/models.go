package database

import (
	"fmt"
	"log"
	"time"
)

// RawData represents raw extracted data
type RawData struct {
	ID          int       `json:"id"`
	Source      string    `json:"source"`
	ExtractedAt time.Time `json:"extracted_at"`
	RawData     string    `json:"raw_data"` // JSON string
	Query       string    `json:"query"`
}

// ProcessedData represents processed data
type ProcessedData struct {
	ID                  int       `json:"id"`
	Source              string    `json:"source"`
	ProcessedAt         time.Time `json:"processed_at"`
	Title               string    `json:"title"`
	Content             string    `json:"content"`
	RelevanceScore      float64   `json:"relevance_score"`
	Sentiment           string    `json:"sentiment"`
	SentimentScore      *float64  `json:"sentiment_score,omitempty"`
	SentimentConfidence *float64  `json:"sentiment_confidence,omitempty"`
	ProcessedData       string    `json:"processed_data"` // JSON string
}

// CreateTables creates all necessary tables
func CreateTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS raw_data (
			id SERIAL PRIMARY KEY,
			source VARCHAR(50) NOT NULL,
			extracted_at TIMESTAMP DEFAULT NOW(),
			raw_data JSONB NOT NULL,
			query VARCHAR(255)
		)`,
		`CREATE TABLE IF NOT EXISTS processed_data (
			id SERIAL PRIMARY KEY,
			source VARCHAR(50) NOT NULL,
			processed_at TIMESTAMP DEFAULT NOW(),
			title TEXT,
			content TEXT,
			relevance_score DECIMAL(3,2),
			sentiment VARCHAR(20),
			sentiment_score DECIMAL(3,2),
			sentiment_confidence DECIMAL(3,2),
			processed_data JSONB NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_raw_data_source ON raw_data(source)`,
		`CREATE INDEX IF NOT EXISTS idx_processed_data_source ON processed_data(source)`,
		`CREATE INDEX IF NOT EXISTS idx_processed_data_timestamp ON processed_data(processed_at)`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	log.Println("✅ Database tables created successfully")
	return nil
}
