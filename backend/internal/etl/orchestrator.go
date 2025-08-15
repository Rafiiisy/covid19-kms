package etl

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// ETLOrchestrator coordinates the entire ETL pipeline
type ETLOrchestrator struct {
	extractor   *DataExtractor
	transformer *DataTransformer
	loader      *DataLoader
}

// ETLResult represents the result of the entire ETL pipeline
type ETLResult struct {
	Status           string                 `json:"status"`
	Message          string                 `json:"message"`
	Timestamp        string                 `json:"timestamp"`
	PipelineDuration string                 `json:"pipeline_duration"`
	Extraction       *ExtractedData        `json:"extraction,omitempty"`
	Transformation   *TransformedData      `json:"transformation,omitempty"`
	Loading          *LoadResult           `json:"loading,omitempty"`
	Summary          map[string]interface{} `json:"summary,omitempty"`
	Error            string                 `json:"error,omitempty"`
}

// NewETLOrchestrator creates a new ETL orchestrator
func NewETLOrchestrator() *ETLOrchestrator {
	return &ETLOrchestrator{
		extractor:   NewDataExtractor(),
		transformer: NewDataTransformer(),
		loader:      NewDataLoader(),
	}
}

// RunETLPipeline executes the complete ETL pipeline
func (eo *ETLOrchestrator) RunETLPipeline() *ETLResult {
	startTime := time.Now()
	log.Println("🚀 Starting ETL pipeline...")

	result := &ETLResult{
		Timestamp: startTime.Format(time.RFC3339),
	}

	// Step 1: Extract data from all sources
	log.Println("📊 Step 1: Data Extraction")
	extractedData, err := eo.extractData()
	if err != nil {
		result.Status = "error"
		result.Message = "ETL pipeline failed during extraction"
		result.Error = err.Error()
		result.PipelineDuration = time.Since(startTime).String()
		return result
	}
	result.Extraction = extractedData

	// Step 2: Transform and clean data
	log.Println("🔄 Step 2: Data Transformation")
	transformedData, err := eo.transformData(extractedData)
	if err != nil {
		result.Status = "error"
		result.Message = "ETL pipeline failed during transformation"
		result.Error = err.Error()
		result.PipelineDuration = time.Since(startTime).String()
		return result
	}
	result.Transformation = transformedData

	// Step 3: Load data to destinations
	log.Println("💾 Step 3: Data Loading")
	loadResult, err := eo.loadData(extractedData, transformedData)
	if err != nil {
		result.Status = "error"
		result.Message = "ETL pipeline failed during loading"
		result.Error = err.Error()
		result.PipelineDuration = time.Since(startTime).String()
		return result
	}
	result.Loading = loadResult

	// Create summary
	result.Summary = eo.createSummary(extractedData, transformedData, loadResult)

	// Calculate pipeline duration
	duration := time.Since(startTime)
	result.PipelineDuration = duration.String()

	// Set final status
	result.Status = "success"
	result.Message = "ETL pipeline completed successfully"

	log.Printf("✅ ETL pipeline completed in %s", duration)
	return result
}

// extractData extracts data from all sources
func (eo *ETLOrchestrator) extractData() (*ExtractedData, error) {
	log.Println("🔄 Starting data extraction...")
	
	extractedData := eo.extractor.ExtractAllSources()
	
	if extractedData == nil {
		return nil, fmt.Errorf("data extraction returned nil")
	}

	log.Printf("✅ Data extraction completed. Sources: %d", len(extractedData.Sources))
	return extractedData, nil
}

// transformData transforms and cleans the extracted data
func (eo *ETLOrchestrator) transformData(extractedData *ExtractedData) (*TransformedData, error) {
	log.Println("🔄 Starting data transformation...")

	// Extract YouTube and news data for transformation
	var youtubeData, newsData interface{}
	
	if source, exists := extractedData.Sources["youtube"]; exists {
		youtubeData = source
	}
	
	if source, exists := extractedData.Sources["google_news"]; exists {
		newsData = source
	}

	transformedData := eo.transformer.TransformData(youtubeData, newsData)
	
	if transformedData == nil {
		return nil, fmt.Errorf("data transformation returned nil")
	}

	log.Printf("✅ Data transformation completed. Videos: %d, Articles: %d", 
		len(transformedData.YouTube), len(transformedData.News))
	
	return transformedData, nil
}

// loadData loads data to local storage
func (eo *ETLOrchestrator) loadData(extractedData *ExtractedData, transformedData *TransformedData) (*LoadResult, error) {
	log.Println("🔄 Starting data loading...")

	// Load raw data to local storage
	rawLoadResult := eo.loader.LoadRawData(extractedData)
	if !rawLoadResult.Success {
		log.Printf("⚠️ Raw data loading failed: %s", rawLoadResult.Error)
	}

	// Load transformed data to local storage
	processedLoadResult := eo.loader.LoadData(transformedData)
	if !processedLoadResult.Success {
		log.Printf("⚠️ Processed data loading failed: %s", processedLoadResult.Error)
	}

	// Return the processed data load result as primary
	return processedLoadResult, nil
}

// createSummary creates a comprehensive summary of the ETL pipeline
func (eo *ETLOrchestrator) createSummary(extractedData *ExtractedData, transformedData *TransformedData, loadResult *LoadResult) map[string]interface{} {
	summary := map[string]interface{}{
		"pipeline_status": "completed",
		"extraction": map[string]interface{}{
			"timestamp": extractedData.Timestamp,
			"query":     extractedData.Query,
			"sources":   len(extractedData.Sources),
		},
		"transformation": map[string]interface{}{
			"timestamp":        transformedData.TransformedAt,
			"videos_count":     len(transformedData.YouTube),
			"articles_count":   len(transformedData.News),
			"average_relevance": transformedData.Summary.AverageRelevance,
		},
		"loading": map[string]interface{}{
			"success":       loadResult.Success,
			"message":       loadResult.Message,
			"records_count": loadResult.RecordsCount,
			"timestamp":     loadResult.Timestamp,
		},
		"load_report": eo.loader.GetLoadReport(),
	}

	return summary
}

// ToJSON converts the ETL result to JSON
func (er *ETLResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(er, "", "  ")
}

// GetPipelineMetrics returns key metrics from the ETL pipeline
func (er *ETLResult) GetPipelineMetrics() map[string]interface{} {
	metrics := map[string]interface{}{
		"status":            er.Status,
		"duration":          er.PipelineDuration,
		"timestamp":         er.Timestamp,
		"extraction_sources": 0,
		"transformed_videos": 0,
		"transformed_articles": 0,
		"loaded_records":    0,
	}

	if er.Extraction != nil {
		metrics["extraction_sources"] = len(er.Extraction.Sources)
	}

	if er.Transformation != nil {
		metrics["transformed_videos"] = len(er.Transformation.YouTube)
		metrics["transformed_articles"] = len(er.Transformation.News)
	}

	if er.Loading != nil {
		metrics["loaded_records"] = er.Loading.RecordsCount
	}

	return metrics
}
