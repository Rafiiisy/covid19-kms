package database

import (
	"encoding/json"
	"fmt"
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
		INSERT INTO processed_data (source, title, content, relevance_score, sentiment, processed_data)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := DB.Exec(sqlQuery,
		data.Source,
		data.Title,
		data.Content,
		data.RelevanceScore,
		data.Sentiment,
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
		SELECT id, source, processed_at, title, content, relevance_score, sentiment, processed_data
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
			SELECT id, source, processed_at, title, content, relevance_score, sentiment, processed_data
			FROM processed_data 
			WHERE source = $1
			ORDER BY processed_at DESC 
			LIMIT $2
		`
		args = []interface{}{source, limit}
	} else {
		// If no limit (or limit = 0), get ALL data
		sqlQuery = `
			SELECT id, source, processed_at, title, content, relevance_score, sentiment, processed_data
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
