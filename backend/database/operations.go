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
	sqlQuery := `
		SELECT id, source, processed_at, title, content, relevance_score, sentiment, processed_data
		FROM processed_data 
		WHERE source = $1
		ORDER BY processed_at DESC 
		LIMIT $2
	`

	rows, err := DB.Query(sqlQuery, source, limit)
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
