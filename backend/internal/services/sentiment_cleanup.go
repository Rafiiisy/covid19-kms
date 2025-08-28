package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// SentimentCleanupService handles cleaning up sentiment data in the database
type SentimentCleanupService struct {
	db                *sql.DB
	sentimentAnalyzer *SentimentAnalyzer
}

// CleanupResult represents the result of a sentiment cleanup operation
type CleanupResult struct {
	TotalRecords     int           `json:"total_records"`
	ProcessedRecords int           `json:"processed_records"`
	UpdatedRecords   int           `json:"updated_records"`
	ErrorRecords     int           `json:"error_records"`
	ProcessingTime   time.Duration `json:"processing_time"`
	Errors           []string      `json:"errors,omitempty"`
	Status           string        `json:"status"`
}

// NewSentimentCleanupService creates a new sentiment cleanup service
func NewSentimentCleanupService(db *sql.DB) *SentimentCleanupService {
	return &SentimentCleanupService{
		db:                db,
		sentimentAnalyzer: NewSentimentAnalyzer(),
	}
}

// CleanAllSentiments cleans sentiment data for all records in the database
func (scs *SentimentCleanupService) CleanAllSentiments() *CleanupResult {
	log.Println("🧹 Starting sentiment cleanup for all records...")

	startTime := time.Now()
	result := &CleanupResult{
		Status: "processing",
	}

	// Get total count of records
	totalCount, err := scs.getTotalRecordCount()
	if err != nil {
		result.Status = "error"
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to get record count: %v", err))
		return result
	}
	result.TotalRecords = totalCount

	// Process records in batches
	batchSize := 100
	offset := 0

	for offset < totalCount {
		// Get batch of records
		records, err := scs.getRecordsBatch(offset, batchSize)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to get batch at offset %d: %v", offset, err))
			offset += batchSize
			continue
		}

		// Process batch
		batchResult := scs.processBatch(records)
		result.ProcessedRecords += batchResult.ProcessedRecords
		result.UpdatedRecords += batchResult.UpdatedRecords
		result.ErrorRecords += batchResult.ErrorRecords
		result.Errors = append(result.Errors, batchResult.Errors...)

		// Log progress
		log.Printf("📊 Processed batch: %d/%d records (%.1f%%)",
			result.ProcessedRecords, totalCount,
			float64(result.ProcessedRecords)/float64(totalCount)*100)

		offset += batchSize
	}

	result.ProcessingTime = time.Since(startTime)

	if len(result.Errors) == 0 {
		result.Status = "completed"
		log.Printf("✅ Sentiment cleanup completed successfully in %v", result.ProcessingTime)
	} else {
		result.Status = "completed_with_errors"
		log.Printf("⚠️  Sentiment cleanup completed with %d errors in %v", len(result.Errors), result.ProcessingTime)
	}

	return result
}

// CleanSentimentBySource cleans sentiment data for a specific source
func (scs *SentimentCleanupService) CleanSentimentBySource(source string) *CleanupResult {
	log.Printf("🧹 Starting sentiment cleanup for source: %s", source)

	startTime := time.Now()
	result := &CleanupResult{
		Status: "processing",
	}

	// Get total count of records for this source
	totalCount, err := scs.getRecordCountBySource(source)
	if err != nil {
		result.Status = "error"
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to get record count for source %s: %v", source, err))
		return result
	}
	result.TotalRecords = totalCount

	// Process records in batches
	batchSize := 100
	offset := 0

	for offset < totalCount {
		// Get batch of records for this source
		records, err := scs.getRecordsBySourceBatch(source, offset, batchSize)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to get batch for source %s at offset %d: %v", source, offset, err))
			offset += batchSize
			continue
		}

		// Process batch
		batchResult := scs.processBatch(records)
		result.ProcessedRecords += batchResult.ProcessedRecords
		result.UpdatedRecords += batchResult.UpdatedRecords
		result.ErrorRecords += batchResult.ErrorRecords
		result.Errors = append(result.Errors, batchResult.Errors...)

		// Log progress
		log.Printf("📊 Processed batch for %s: %d/%d records (%.1f%%)",
			source, result.ProcessedRecords, totalCount,
			float64(result.ProcessedRecords)/float64(totalCount)*100)

		offset += batchSize
	}

	result.ProcessingTime = time.Since(startTime)

	if len(result.Errors) == 0 {
		result.Status = "completed"
		log.Printf("✅ Sentiment cleanup for %s completed successfully in %v", source, result.ProcessingTime)
	} else {
		result.Status = "completed_with_errors"
		log.Printf("⚠️  Sentiment cleanup for %s completed with %d errors in %v", source, len(result.Errors), result.ProcessingTime)
	}

	return result
}

// CleanSentimentByDateRange cleans sentiment data for records within a date range
func (scs *SentimentCleanupService) CleanSentimentByDateRange(startDate, endDate time.Time) *CleanupResult {
	log.Printf("🧹 Starting sentiment cleanup for date range: %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	startTime := time.Now()
	result := &CleanupResult{
		Status: "processing",
	}

	// Get total count of records in date range
	totalCount, err := scs.getRecordCountByDateRange(startDate, endDate)
	if err != nil {
		result.Status = "error"
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to get record count for date range: %v", err))
		return result
	}
	result.TotalRecords = totalCount

	// Process records in batches
	batchSize := 100
	offset := 0

	for offset < totalCount {
		// Get batch of records in date range
		records, err := scs.getRecordsByDateRangeBatch(startDate, endDate, offset, batchSize)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to get batch for date range at offset %d: %v", offset, err))
			offset += batchSize
			continue
		}

		// Process batch
		batchResult := scs.processBatch(records)
		result.ProcessedRecords += batchResult.ProcessedRecords
		result.UpdatedRecords += batchResult.UpdatedRecords
		result.ErrorRecords += batchResult.ErrorRecords
		result.Errors = append(result.Errors, batchResult.Errors...)

		// Log progress
		log.Printf("📊 Processed batch for date range: %d/%d records (%.1f%%)",
			result.ProcessedRecords, totalCount,
			float64(result.ProcessedRecords)/float64(totalCount)*100)

		offset += batchSize
	}

	result.ProcessingTime = time.Since(startTime)

	if len(result.Errors) == 0 {
		result.Status = "completed"
		log.Printf("✅ Sentiment cleanup for date range completed successfully in %v", result.ProcessingTime)
	} else {
		result.Status = "completed_with_errors"
		log.Printf("⚠️  Sentiment cleanup for date range completed with %d errors in %v", len(result.Errors), result.ProcessingTime)
	}

	return result
}

// processBatch processes a batch of records and updates their sentiment
func (scs *SentimentCleanupService) processBatch(records []ProcessedDataRecord) *CleanupResult {
	result := &CleanupResult{}

	log.Printf("🔄 Processing batch of %d records", len(records))

	for _, record := range records {
		result.ProcessedRecords++

		// Analyze sentiment for the record
		combinedText := record.Title + " " + record.Content
		log.Printf("📝 Analyzing record %d: '%s'", record.ID, combinedText[:min(len(combinedText), 50)])

		sentimentResult := scs.sentimentAnalyzer.AnalyzeSentiment(combinedText)
		log.Printf("🎯 Record %d sentiment result: %s (%.3f, %.3f)",
			record.ID, sentimentResult.Category, sentimentResult.Score, sentimentResult.Confidence)

		// Update the record in database
		err := scs.updateRecordSentiment(record.ID, sentimentResult)
		if err != nil {
			result.ErrorRecords++
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to update record %d: %v", record.ID, err))
			log.Printf("❌ Failed to update record %d: %v", record.ID, err)
		} else {
			result.UpdatedRecords++
			log.Printf("✅ Successfully updated record %d", record.ID)
		}
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// updateRecordSentiment updates the sentiment fields for a single record
func (scs *SentimentCleanupService) updateRecordSentiment(recordID int, sentimentResult *SentimentResult) error {
	query := `
		UPDATE processed_data 
		SET sentiment = $1, 
		    sentiment_score = $2, 
		    sentiment_confidence = $3,
		    processed_at = $4
		WHERE id = $5
	`

	log.Printf("🔧 Updating record %d: sentiment='%s', score=%.3f, confidence=%.3f",
		recordID, sentimentResult.Category, sentimentResult.Score, sentimentResult.Confidence)

	result, err := scs.db.Exec(query,
		sentimentResult.Category,
		sentimentResult.Score,
		sentimentResult.Confidence,
		time.Now(),
		recordID,
	)

	if err != nil {
		log.Printf("❌ Failed to update record %d: %v", recordID, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Updated record %d: %d rows affected", recordID, rowsAffected)

	return nil
}

// Database helper functions
func (scs *SentimentCleanupService) getTotalRecordCount() (int, error) {
	var count int
	err := scs.db.QueryRow("SELECT COUNT(*) FROM processed_data").Scan(&count)
	return count, err
}

func (scs *SentimentCleanupService) getRecordCountBySource(source string) (int, error) {
	var count int
	err := scs.db.QueryRow("SELECT COUNT(*) FROM processed_data WHERE source = $1", source).Scan(&count)
	return count, err
}

func (scs *SentimentCleanupService) getRecordCountByDateRange(startDate, endDate time.Time) (int, error) {
	var count int
	err := scs.db.QueryRow("SELECT COUNT(*) FROM processed_data WHERE processed_at BETWEEN $1 AND $2", startDate, endDate).Scan(&count)
	return count, err
}

func (scs *SentimentCleanupService) getRecordsBatch(offset, limit int) ([]ProcessedDataRecord, error) {
	query := `
		SELECT id, source, title, content, relevance_score, sentiment, processed_at, processed_data
		FROM processed_data 
		ORDER BY id 
		LIMIT $1 OFFSET $2
	`

	rows, err := scs.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ProcessedDataRecord
	for rows.Next() {
		var record ProcessedDataRecord
		err := rows.Scan(
			&record.ID,
			&record.Source,
			&record.Title,
			&record.Content,
			&record.RelevanceScore,
			&record.Sentiment,
			&record.ProcessedAt,
			&record.ProcessedData,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (scs *SentimentCleanupService) getRecordsBySourceBatch(source string, offset, limit int) ([]ProcessedDataRecord, error) {
	query := `
		SELECT id, source, title, content, relevance_score, sentiment, processed_at, processed_data
		FROM processed_data 
		WHERE source = $1
		ORDER BY id 
		LIMIT $2 OFFSET $3
	`

	rows, err := scs.db.Query(query, source, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ProcessedDataRecord
	for rows.Next() {
		var record ProcessedDataRecord
		err := rows.Scan(
			&record.ID,
			&record.Source,
			&record.Title,
			&record.Content,
			&record.RelevanceScore,
			&record.Sentiment,
			&record.ProcessedAt,
			&record.ProcessedData,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (scs *SentimentCleanupService) getRecordsByDateRangeBatch(startDate, endDate time.Time, offset, limit int) ([]ProcessedDataRecord, error) {
	query := `
		SELECT id, source, title, content, relevance_score, sentiment, processed_at, processed_data
		FROM processed_data 
		WHERE processed_at BETWEEN $1 AND $2
		ORDER BY id 
		LIMIT $3 OFFSET $4
	`

	rows, err := scs.db.Query(query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ProcessedDataRecord
	for rows.Next() {
		var record ProcessedDataRecord
		err := rows.Scan(
			&record.ID,
			&record.Source,
			&record.Title,
			&record.Content,
			&record.RelevanceScore,
			&record.Sentiment,
			&record.ProcessedAt,
			&record.ProcessedData,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// ProcessedDataRecord represents a record from the processed_data table
type ProcessedDataRecord struct {
	ID             int       `json:"id"`
	Source         string    `json:"source"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	RelevanceScore float64   `json:"relevance_score"`
	Sentiment      string    `json:"sentiment"`
	ProcessedAt    time.Time `json:"processed_at"`
	ProcessedData  string    `json:"processed_data"`
}
