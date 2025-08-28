package etl

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"covid19-kms/database"
)

// DataLoader handles loading data to PostgreSQL database
type DataLoader struct {
	// No outputDir needed for database
}

// LoadResult represents the result of a data loading operation
type LoadResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Timestamp    string `json:"timestamp"`
	RecordsCount int    `json:"records_count"`
	Error        string `json:"error,omitempty"`
}

// NewDataLoader creates a new DataLoader instance
func NewDataLoader() *DataLoader {
	return &DataLoader{}
}

// LoadData loads transformed data to PostgreSQL database
func (dl *DataLoader) LoadData(data *TransformedData) *LoadResult {
	log.Println("Loading data to PostgreSQL database...")

	// Count total records
	totalRecords := len(data.YouTube) + len(data.News)

	// Save to database
	for _, video := range data.YouTube {
		// Convert video to JSON string
		videoJSON, err := json.Marshal(video)
		if err != nil {
			log.Printf("Failed to marshal video data: %v", err)
			continue
		}

		processedData := &database.ProcessedData{
			Source:              "youtube",
			Title:               video.Title,
			Content:             video.Description,
			RelevanceScore:      video.CovidRelevanceScore,
			Sentiment:           video.Sentiment,
			SentimentScore:      &video.SentimentScore,
			SentimentConfidence: &video.SentimentConfidence,
			ProcessedData:       string(videoJSON),
		}

		if err := database.InsertProcessedData(processedData); err != nil {
			log.Printf("Failed to insert video data: %v", err)
		}
	}

	for _, article := range data.News {
		// Convert article to JSON string
		articleJSON, err := json.Marshal(article)
		if err != nil {
			log.Printf("Failed to marshal article data: %v", err)
			continue
		}

		// Determine the specific source based on the article source field
		sourceName := "news" // default
		if article.Source != "" {
			switch article.Source {
			case "CNN", "DETIK", "KOMPAS", "Indonesia News":
				sourceName = "indonesia_news"
			case "Real-Time News":
				sourceName = "google_news" // Store as google_news for backward compatibility
			case "Instagram":
				sourceName = "instagram"
			default:
				// Check if it contains Instagram-related keywords
				if strings.Contains(strings.ToLower(article.Source), "instagram") {
					sourceName = "instagram"
				} else if strings.Contains(strings.ToLower(article.Source), "indonesia") {
					sourceName = "indonesia_news"
				} else {
					sourceName = "news"
				}
			}
		}

		processedData := &database.ProcessedData{
			Source:              sourceName,
			Title:               article.Title,
			Content:             article.Content,
			RelevanceScore:      article.CovidRelevanceScore,
			Sentiment:           article.Sentiment,
			SentimentScore:      &article.SentimentScore,
			SentimentConfidence: &article.SentimentConfidence,
			ProcessedData:       string(articleJSON),
		}

		if err := database.InsertProcessedData(processedData); err != nil {
			log.Printf("Failed to insert article data: %v", err)
		}
	}

	return &LoadResult{
		Success:      true,
		Message:      "Data successfully loaded to PostgreSQL database",
		Timestamp:    time.Now().Format(time.RFC3339),
		RecordsCount: totalRecords,
	}
}

// LoadRawData loads raw extracted data to PostgreSQL database
func (dl *DataLoader) LoadRawData(data *ExtractedData) *LoadResult {
	log.Println("Loading raw data to PostgreSQL database...")

	// Save raw data to database
	for sourceName, sourceData := range data.Sources {
		if err := database.InsertRawData(sourceName, data.Query, sourceData); err != nil {
			log.Printf("Failed to insert raw data for source %s: %v", sourceName, err)
		}
	}

	return &LoadResult{
		Success:      true,
		Message:      "Raw data successfully loaded to PostgreSQL database",
		Timestamp:    time.Now().Format(time.RFC3339),
		RecordsCount: len(data.Sources),
	}
}

// GetLoadReport generates a load report
func (dl *DataLoader) GetLoadReport() map[string]interface{} {
	return map[string]interface{}{
		"storage_type": "postgresql",
		"timestamp":    time.Now().Format(time.RFC3339),
	}
}
