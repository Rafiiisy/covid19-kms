package database

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"
)

// InsertRawData inserts raw data into the database
func InsertRawData(source, query string, rawData interface{}) error {
	jsonData, err := json.Marshal(rawData)
	if err != nil {
		return fmt.Errorf("failed to marshal raw data: %v", err)
	}

	sqlQuery := `
		INSERT INTO raw_data (source, query, raw_data)
		VALUES ($1, $2, $3)
	`

	_, err = DB.Exec(sqlQuery, source, query, string(jsonData))
	if err != nil {
		return fmt.Errorf("failed to insert raw data: %v", err)
	}

	return nil
}

// InsertProcessedData inserts processed data into the database
func InsertProcessedData(data *ProcessedData) error {
	sqlQuery := `
		INSERT INTO processed_data (source, title, content, relevance_score, sentiment, sentiment_score, sentiment_confidence, processed_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := DB.Exec(sqlQuery,
		data.Source,
		data.Title,
		data.Content,
		data.RelevanceScore,
		data.Sentiment,
		data.SentimentScore,
		data.SentimentConfidence,
		data.ProcessedData,
	)
	if err != nil {
		return fmt.Errorf("failed to insert processed data: %v", err)
	}

	return nil
}

// GetLatestProcessedData retrieves the latest processed data
func GetLatestProcessedData(limit int) ([]ProcessedData, error) {
	// Check if database is connected and ensure connection is alive
	if err := EnsureConnection(); err != nil {
		return []ProcessedData{}, fmt.Errorf("database connection issue: %v", err)
	}

	sqlQuery := `
		SELECT id, source, processed_at, title, content, relevance_score, sentiment, sentiment_score, sentiment_confidence, processed_data
		FROM processed_data 
		ORDER BY processed_at DESC 
		LIMIT $1
	`

	rows, err := DB.Query(sqlQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query processed data: %v", err)
	}
	defer rows.Close()

	var results []ProcessedData
	for rows.Next() {
		var data ProcessedData
		err := rows.Scan(
			&data.ID,
			&data.Source,
			&data.ProcessedAt,
			&data.Title,
			&data.Content,
			&data.RelevanceScore,
			&data.Sentiment,
			&data.SentimentScore,
			&data.SentimentConfidence,
			&data.ProcessedData,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		results = append(results, data)
	}

	return results, nil
}

// GetDataBySource retrieves data by source
func GetDataBySource(source string, limit int) ([]ProcessedData, error) {
	// Check if database is connected and ensure connection is alive
	if err := EnsureConnection(); err != nil {
		return []ProcessedData{}, fmt.Errorf("database connection issue: %v", err)
	}

	var sqlQuery string
	var args []interface{}

	if limit > 0 {
		// If limit specified, use it
		sqlQuery = `
			SELECT id, source, processed_at, title, content, relevance_score, sentiment, sentiment_score, sentiment_confidence, processed_data
			FROM processed_data 
			WHERE source = $1
			ORDER BY processed_at DESC 
			LIMIT $2
		`
		args = []interface{}{source, limit}
	} else {
		// If no limit (or limit = 0), get ALL data
		sqlQuery = `
			SELECT id, source, processed_at, title, content, relevance_score, sentiment, sentiment_score, sentiment_confidence, processed_data
			FROM processed_data 
			WHERE source = $1
			ORDER BY processed_at DESC
		`
		args = []interface{}{source}
	}

	rows, err := DB.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query data by source: %v", err)
	}
	defer rows.Close()

	var results []ProcessedData
	for rows.Next() {
		var data ProcessedData
		err := rows.Scan(
			&data.ID,
			&data.Source,
			&data.ProcessedAt,
			&data.Title,
			&data.Content,
			&data.RelevanceScore,
			&data.Sentiment,
			&data.SentimentScore,
			&data.SentimentConfidence,
			&data.ProcessedData,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		results = append(results, data)
	}

	return results, nil
}

// GetDataCount returns the total count of records
func GetDataCount() (map[string]int, error) {
	// Check if database is connected and ensure connection is alive
	if err := EnsureConnection(); err != nil {
		return map[string]int{"raw_data": 0, "processed_data": 0}, fmt.Errorf("database connection issue: %v", err)
	}

	counts := make(map[string]int)

	// Count raw data
	var rawCount int
	err := DB.QueryRow("SELECT COUNT(*) FROM raw_data").Scan(&rawCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count raw data: %v", err)
	}
	counts["raw_data"] = rawCount

	// Count processed data
	var processedCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM processed_data").Scan(&processedCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count processed data: %v", err)
	}
	counts["processed_data"] = processedCount

	return counts, nil
}

// GetDataSummary returns a comprehensive summary of all data
func GetDataSummary() (map[string]interface{}, error) {
	// Check if database is connected and ensure connection is alive
	if err := EnsureConnection(); err != nil {
		return map[string]interface{}{
			"error":             "Database connection issue",
			"source_counts":     map[string]int{"youtube": 0, "google_news": 0, "instagram": 0, "indonesia_news": 0},
			"average_relevance": 0.0,
			"total_records":     0,
			"latest_update":     "Never",
		}, fmt.Errorf("database connection issue: %v", err)
	}

	summary := make(map[string]interface{})

	// Get counts by source
	sources := []string{"youtube", "google_news", "instagram", "indonesia_news"}
	sourceCounts := make(map[string]int)

	for _, source := range sources {
		var count int
		err := DB.QueryRow("SELECT COUNT(*) FROM processed_data WHERE source = $1", source).Scan(&count)
		if err != nil {
			// Log error but continue with other sources
			fmt.Printf("Warning: failed to count %s data: %v\n", source, err)
			sourceCounts[source] = 0
		} else {
			sourceCounts[source] = count
		}
	}

	// Get average relevance score
	var avgRelevance float64
	err := DB.QueryRow("SELECT AVG(relevance_score) FROM processed_data WHERE relevance_score IS NOT NULL").Scan(&avgRelevance)
	if err != nil {
		avgRelevance = 0.0
	}

	// Get total records
	var totalRecords int
	err = DB.QueryRow("SELECT COUNT(*) FROM processed_data").Scan(&totalRecords)
	if err != nil {
		totalRecords = 0
	}

	// Get latest update timestamp
	var latestUpdate string
	err = DB.QueryRow("SELECT MAX(processed_at) FROM processed_data").Scan(&latestUpdate)
	if err != nil {
		latestUpdate = "Never"
	}

	summary["source_counts"] = sourceCounts
	summary["average_relevance"] = avgRelevance
	summary["total_records"] = totalRecords
	summary["latest_update"] = latestUpdate

	return summary, nil
}

// GetSentimentDistribution returns sentiment distribution across all sources
func GetSentimentDistribution() (map[string]interface{}, error) {
	// Check if database is connected and ensure connection is alive
	if err := EnsureConnection(); err != nil {
		return map[string]interface{}{
			"error": "Database connection issue",
			"sources": map[string]interface{}{
				"youtube":        map[string]interface{}{"positive": 0, "negative": 0, "neutral": 0},
				"google_news":    map[string]interface{}{"positive": 0, "negative": 0, "neutral": 0},
				"instagram":      map[string]interface{}{"positive": 0, "negative": 0, "neutral": 0},
				"indonesia_news": map[string]interface{}{"positive": 0, "negative": 0, "neutral": 0},
			},
		}, fmt.Errorf("database connection issue: %v", err)
	}

	distribution := make(map[string]interface{})
	sources := []string{"youtube", "google_news", "instagram", "indonesia_news"}
	sentiments := []string{"positive", "negative", "neutral"}

	// Initialize distribution structure
	sourceDistribution := make(map[string]interface{})
	for _, source := range sources {
		sourceDistribution[source] = make(map[string]interface{})
		for _, sentiment := range sentiments {
			sourceDistribution[source].(map[string]interface{})[sentiment] = 0
		}
	}

	// Query sentiment distribution for each source
	for _, source := range sources {
		for _, sentiment := range sentiments {
			var count int
			query := "SELECT COUNT(*) FROM processed_data WHERE source = $1 AND sentiment = $2"
			err := DB.QueryRow(query, source, sentiment).Scan(&count)
			if err != nil {
				// Log error but continue
				fmt.Printf("Warning: failed to count %s %s data: %v\n", source, sentiment, err)
				sourceDistribution[source].(map[string]interface{})[sentiment] = 0
			} else {
				sourceDistribution[source].(map[string]interface{})[sentiment] = count
			}
		}
	}

	distribution["sources"] = sourceDistribution

	// Calculate totals
	totalPositive := 0
	totalNegative := 0
	totalNeutral := 0

	for _, source := range sources {
		sourceData := sourceDistribution[source].(map[string]interface{})
		totalPositive += sourceData["positive"].(int)
		totalNegative += sourceData["negative"].(int)
		totalNeutral += sourceData["neutral"].(int)
	}

	distribution["totals"] = map[string]interface{}{
		"positive": totalPositive,
		"negative": totalNegative,
		"neutral":  totalNeutral,
		"total":    totalPositive + totalNegative + totalNeutral,
	}

	return distribution, nil
}

// GetWordFrequency returns word frequency analysis across all sources
func GetWordFrequency() (map[string]interface{}, error) {
	// Check if database is connected and ensure connection is alive
	if err := EnsureConnection(); err != nil {
		return map[string]interface{}{
			"error": "Database connection issue",
			"words": []map[string]interface{}{},
		}, fmt.Errorf("database connection issue: %v", err)
	}

	// Query to get all titles and content for word analysis
	query := `
		SELECT 
			source,
			title,
			content,
			sentiment,
			sentiment_score
		FROM processed_data 
		WHERE title IS NOT NULL OR content IS NOT NULL
		ORDER BY processed_at DESC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query word frequency data: %v", err)
	}
	defer rows.Close()

	// Process text and count words
	wordCounts := make(map[string]map[string]interface{})
	stopWords := getStopWords()

	for rows.Next() {
		var source, title, content, sentiment string
		var sentimentScore *float64

		err := rows.Scan(&source, &title, &content, &sentiment, &sentimentScore)
		if err != nil {
			continue
		}

		// Combine title and content for analysis
		combinedText := title + " " + content
		words := tokenizeText(combinedText)

		for _, word := range words {
			wordLower := strings.ToLower(strings.TrimSpace(word))

			// Skip stop words, short words, and non-alphabetic
			if len(wordLower) < 3 || contains(stopWords, wordLower) || !isAlphabetic(wordLower) {
				continue
			}

			// Initialize word entry if not exists
			if _, exists := wordCounts[wordLower]; !exists {
				wordCounts[wordLower] = map[string]interface{}{
					"word":           wordLower,
					"count":          0,
					"positive_count": 0,
					"negative_count": 0,
					"neutral_count":  0,
					"sources":        make(map[string]int),
					"avg_sentiment":  0.0,
				}
			}

			// Update counts
			wordCounts[wordLower]["count"] = wordCounts[wordLower]["count"].(int) + 1
			wordCounts[wordLower]["sources"].(map[string]int)[source] = wordCounts[wordLower]["sources"].(map[string]int)[source] + 1

			// Update sentiment counts
			switch sentiment {
			case "positive":
				wordCounts[wordLower]["positive_count"] = wordCounts[wordLower]["positive_count"].(int) + 1
			case "negative":
				wordCounts[wordLower]["negative_count"] = wordCounts[wordLower]["negative_count"].(int) + 1
			case "neutral":
				wordCounts[wordLower]["neutral_count"] = wordCounts[wordLower]["neutral_count"].(int) + 1
			}

			// Update average sentiment score
			if sentimentScore != nil {
				currentAvg := wordCounts[wordLower]["avg_sentiment"].(float64)
				currentCount := wordCounts[wordLower]["count"].(int)
				newAvg := (currentAvg*float64(currentCount-1) + *sentimentScore) / float64(currentCount)
				wordCounts[wordLower]["avg_sentiment"] = newAvg
			}
		}
	}

	// Convert to sorted list and limit to top words
	var wordList []map[string]interface{}
	for _, wordData := range wordCounts {
		wordList = append(wordList, wordData)
	}

	// Sort by frequency (descending) and take top 100
	sort.Slice(wordList, func(i, j int) bool {
		return wordList[i]["count"].(int) > wordList[j]["count"].(int)
	})

	if len(wordList) > 100 {
		wordList = wordList[:100]
	}

	return map[string]interface{}{
		"words":              wordList,
		"total_words":        len(wordList),
		"analysis_timestamp": time.Now().Format(time.RFC3339),
	}, nil
}

// Helper functions for word frequency analysis
func getStopWords() map[string]bool {
	stopWords := map[string]bool{
		// English stop words
		"the": true, "and": true, "or": true, "but": true, "in": true, "on": true, "at": true,
		"to": true, "for": true, "of": true, "with": true, "by": true, "from": true, "up": true,
		"about": true, "into": true, "through": true, "during": true, "before": true, "after": true,
		"above": true, "below": true, "between": true, "among": true, "within": true, "without": true,
		"is": true, "are": true, "was": true, "were": true, "be": true, "been": true, "being": true,
		"have": true, "has": true, "had": true, "do": true, "does": true, "did": true, "will": true,
		"would": true, "could": true, "should": true, "may": true, "might": true, "can": true,
		"this": true, "that": true, "these": true, "those": true, "i": true, "you": true, "he": true,
		"she": true, "it": true, "we": true, "they": true, "me": true, "him": true, "her": true,
		"us": true, "them": true, "my": true, "your": true, "his": true, "its": true,
		"our": true, "their": true, "mine": true, "yours": true, "hers": true, "ours": true, "theirs": true,

		// Indonesian stop words
		"yang": true, "dan": true, "atau": true, "tetapi": true, "di": true, "ke": true, "dari": true,
		"untuk": true, "dengan": true, "oleh": true, "tentang": true, "antara": true, "dalam": true,
		"adalah": true, "akan": true, "sudah": true, "belum": true, "tidak": true, "bukan": true,
		"ini": true, "itu": true, "saya": true, "anda": true, "dia": true, "kami": true, "mereka": true,
		"kita": true,

		// Common words to filter out
		"covid": true, "coronavirus": true, "virus": true, "pandemic": true, "epidemic": true,
		"case": true, "cases": true, "death": true, "deaths": true, "recovery": true, "recoveries": true,
		"vaccine": true, "vaccination": true, "lockdown": true, "quarantine": true, "isolation": true,
		"test": true, "testing": true, "positive": true, "negative": true, "confirmed": true,
		"report": true, "reported": true, "announced": true, "announcement": true, "update": true,
		"news": true, "article": true, "post": true, "comment": true, "video": true, "media": true,
	}
	return stopWords
}

func tokenizeText(text string) []string {
	// Split by whitespace and punctuation
	words := strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})

	// Clean words (remove empty strings, normalize)
	var cleanedWords []string
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word != "" && len(word) > 1 {
			cleanedWords = append(cleanedWords, word)
		}
	}

	return cleanedWords
}

func contains(slice map[string]bool, item string) bool {
	_, exists := slice[item]
	return exists
}

func isAlphabetic(word string) bool {
	for _, char := range word {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}
