package etl

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// DataTransformer handles data cleaning, transformation, and enrichment
type DataTransformer struct {
	covidKeywords []string
}

// TransformedData represents the structure of transformed data
type TransformedData struct {
	YouTube      []TransformedVideo `json:"youtube"`
	News         []TransformedArticle `json:"news"`
	Summary      DataSummary         `json:"summary"`
	TransformedAt string             `json:"transformed_at"`
}

// TransformedVideo represents a transformed YouTube video
type TransformedVideo struct {
	ID                  string  `json:"id"`
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	PublishedAt         string  `json:"published_at"`
	ChannelTitle        string  `json:"channel_title"`
	ThumbnailURL        string  `json:"thumbnail_url"`
	Source              string  `json:"source"`
	CovidRelevanceScore float64 `json:"covid_relevance_score"`
	Language            string  `json:"language"`
	WordCount           int     `json:"word_count"`
	ExtractedAt         string  `json:"extracted_at"`
	TransformedAt       string  `json:"transformed_at"`
}

// TransformedArticle represents a transformed news article
type TransformedArticle struct {
	ID                  string  `json:"id"`
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	Content             string  `json:"content"`
	URL                 string  `json:"url"`
	Source              string  `json:"source"`
	CovidRelevanceScore float64 `json:"covid_relevance_score"`
	Language            string  `json:"language"`
	WordCount           int     `json:"word_count"`
	ExtractedAt         string  `json:"extracted_at"`
	TransformedAt       string  `json:"transformed_at"`
}

// DataSummary represents summary statistics
type DataSummary struct {
	TotalVideos         int     `json:"total_videos"`
	TotalArticles       int     `json:"total_articles"`
	AverageRelevance    float64 `json:"average_relevance"`
	ProcessingTimestamp string  `json:"processing_timestamp"`
}

// NewDataTransformer creates a new DataTransformer instance
func NewDataTransformer() *DataTransformer {
	return &DataTransformer{
		covidKeywords: []string{
			"covid", "coronavirus", "pandemic", "vaccine", "vaccination",
			"lockdown", "quarantine", "social distancing", "mask",
			"indonesia", "jakarta", "jawa", "sulawesi", "sumatra",
		},
	}
}

// TransformData transforms all extracted data
func (dt *DataTransformer) TransformData(youtubeData, newsData interface{}) *TransformedData {
	log.Println("Starting data transformation...")

	transformedData := &TransformedData{
		TransformedAt: time.Now().Format(time.RFC3339),
	}

	// Transform YouTube data
	if youtubeData != nil {
		transformedData.YouTube = dt.transformYouTubeData(youtubeData)
	}

	// Transform news data
	if newsData != nil {
		transformedData.News = dt.transformNewsData(newsData)
	}

	// Create summary
	transformedData.Summary = dt.createSummary(transformedData.YouTube, transformedData.News)

	log.Println("Data transformation completed")
	return transformedData
}

// transformYouTubeData transforms YouTube data
func (dt *DataTransformer) transformYouTubeData(data interface{}) []TransformedVideo {
	var transformedVideos []TransformedVideo

	// Type assertion and processing would go here
	// For now, return empty slice as placeholder
	log.Println("Transforming YouTube data...")

	return transformedVideos
}

// transformNewsData transforms news data
func (dt *DataTransformer) transformNewsData(data interface{}) []TransformedArticle {
	var transformedArticles []TransformedArticle

	// Type assertion and processing would go here
	// For now, return empty slice as placeholder
	log.Println("Transforming news data...")

	return transformedArticles
}

// cleanText cleans and normalizes text
func (dt *DataTransformer) cleanText(text string) string {
	if text == "" {
		return ""
	}

	// Remove extra whitespace
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Remove special characters (keep basic punctuation)
	text = regexp.MustCompile(`[^\w\s.,!?-]`).ReplaceAllString(text, "")

	return text
}

// calculateCovidRelevance calculates relevance score for COVID-19 content
func (dt *DataTransformer) calculateCovidRelevance(text string) float64 {
	if text == "" {
		return 0.0
	}

	text = strings.ToLower(text)
	score := 0.0

	for _, keyword := range dt.covidKeywords {
		if strings.Contains(text, keyword) {
			score += 1.0
		}
	}

	// Normalize score to 0-1 range
	maxPossibleScore := float64(len(dt.covidKeywords))
	if maxPossibleScore > 0 {
		score = score / maxPossibleScore
	}

	return score
}

// detectLanguage detects the language of the text (simplified)
func (dt *DataTransformer) detectLanguage(text string) string {
	if text == "" {
		return "unknown"
	}

	// Simple language detection based on common words
	text = strings.ToLower(text)
	
	// Indonesian words
	indonesianWords := []string{"yang", "dan", "atau", "dengan", "untuk", "dari", "ke", "di", "pada"}
	for _, word := range indonesianWords {
		if strings.Contains(text, word) {
			return "id"
		}
	}

	// English words
	englishWords := []string{"the", "and", "or", "with", "for", "from", "to", "in", "on", "at"}
	for _, word := range englishWords {
		if strings.Contains(text, word) {
			return "en"
		}
	}

	return "unknown"
}

// parseDateTime parses datetime strings
func (dt *DataTransformer) parseDateTime(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// Try to parse various date formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed.Format(time.RFC3339)
		}
	}

	// Return original if parsing fails
	return dateStr
}

// generateArticleID generates a unique ID for an article
func (dt *DataTransformer) generateArticleID(article interface{}) string {
	// Implementation would generate a hash or UUID based on article content
	// For now, return a timestamp-based ID
	return fmt.Sprintf("article_%d", time.Now().Unix())
}

// createSummary creates summary statistics
func (dt *DataTransformer) createSummary(videos []TransformedVideo, articles []TransformedArticle) DataSummary {
	totalVideos := len(videos)
	totalArticles := len(articles)

	// Calculate average relevance
	totalRelevance := 0.0
	count := 0

	for _, video := range videos {
		totalRelevance += video.CovidRelevanceScore
		count++
	}

	for _, article := range articles {
		totalRelevance += article.CovidRelevanceScore
		count++
	}

	averageRelevance := 0.0
	if count > 0 {
		averageRelevance = totalRelevance / float64(count)
	}

	return DataSummary{
		TotalVideos:         totalVideos,
		TotalArticles:       totalArticles,
		AverageRelevance:    averageRelevance,
		ProcessingTimestamp: time.Now().Format(time.RFC3339),
	}
}
