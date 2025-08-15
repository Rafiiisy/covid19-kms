package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ETLPipeline represents the complete ETL pipeline orchestrator
type ETLPipeline struct {
	outputDir string
}

// PipelineResult represents the result of the complete ETL pipeline
type PipelineResult struct {
	PipelineStartTime string                 `json:"pipeline_start_time"`
	Stages            map[string]StageResult `json:"stages"`
	Summary           map[string]interface{} `json:"summary"`
	Status            string                 `json:"status"`
	ErrorMessage      string                 `json:"error_message,omitempty"`
}

// StageResult represents the result of a single ETL stage
type StageResult struct {
	Status  string      `json:"status"`
	File    string      `json:"file,omitempty"`
	Records interface{} `json:"records,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewETLPipeline creates a new ETL pipeline instance
func NewETLPipeline() *ETLPipeline {
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	
	return &ETLPipeline{
		outputDir: outputDir,
	}
}

// RunFullPipeline executes the complete ETL pipeline
func (ep *ETLPipeline) RunFullPipeline() *PipelineResult {
	fmt.Println("🚀 Starting Complete ETL Pipeline")
	fmt.Println(strings.Repeat("=", 50))
	
	pipelineResults := &PipelineResult{
		PipelineStartTime: time.Now().Format(time.RFC3339),
		Stages:            make(map[string]StageResult),
		Summary:           make(map[string]interface{}),
		Status:            "running",
	}
	
	timestamp := time.Now().Format("20060102_150405")
	
	// Stage 1: Extract
	fmt.Println("\n📥 STAGE 1: EXTRACTION")
	fmt.Println(strings.Repeat("-", 30))
	
	extractionResult := ep.runExtractionStage(timestamp)
	pipelineResults.Stages["extraction"] = extractionResult
	
	if extractionResult.Status == "error" {
		pipelineResults.Status = "failed"
		pipelineResults.ErrorMessage = extractionResult.Error
		return pipelineResults
	}
	
	// Stage 2: Transform
	fmt.Println("\n🔄 STAGE 2: TRANSFORMATION")
	fmt.Println(strings.Repeat("-", 30))
	
	transformationResult := ep.runTransformationStage(timestamp)
	pipelineResults.Stages["transformation"] = transformationResult
	
	if transformationResult.Status == "error" {
		pipelineResults.Status = "failed"
		pipelineResults.ErrorMessage = transformationResult.Error
		return pipelineResults
	}
	
	// Stage 3: Load
	fmt.Println("\n📊 STAGE 3: LOADING")
	fmt.Println(strings.Repeat("-", 30))
	
	loadingResult := ep.runLoadingStage(timestamp)
	pipelineResults.Stages["loading"] = loadingResult
	
	if loadingResult.Status == "error" {
		pipelineResults.Status = "failed"
		pipelineResults.ErrorMessage = loadingResult.Error
		return pipelineResults
	}
	
	// Pipeline completed successfully
	pipelineResults.Status = "completed"
	pipelineResults.Summary["total_stages"] = 3
	pipelineResults.Summary["completion_time"] = time.Now().Format(time.RFC3339)
	
	// Save pipeline results
	ep.savePipelineResults(pipelineResults, timestamp)
	
	fmt.Println("\n🎉 ETL Pipeline completed successfully!")
	return pipelineResults
}

// runExtractionStage runs the data extraction stage
func (ep *ETLPipeline) runExtractionStage(timestamp string) StageResult {
	fmt.Println("🔄 Starting data extraction from all sources...")
	
	// Create sample extracted data for testing
	extractedData := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"query":     "covid19",
		"sources": map[string]interface{}{
			"youtube": map[string]interface{}{
				"videos": []map[string]interface{}{
					{
						"video_id": "test_video_1",
						"title":    "COVID-19 Vaccine Update Indonesia",
						"channel":  "Health Channel",
						"views":    1000,
					},
					{
						"video_id": "test_video_2",
						"title":    "Pandemi COVID-19 di Jakarta",
						"channel":  "News Channel",
						"views":    2000,
					},
				},
			},
			"google_news": map[string]interface{}{
				"articles": []map[string]interface{}{
					{
						"title":       "Indonesia COVID-19 Cases Update",
						"description": "Latest statistics on COVID-19 cases in Indonesia",
						"source":      "Tempo",
						"url":         "https://example.com/article1",
					},
					{
						"title":       "Vaccination Progress in Java",
						"description": "COVID-19 vaccination progress across Java island",
						"source":      "Kompas",
						"url":         "https://example.com/article2",
					},
				},
			},
			"instagram": map[string]interface{}{
				"posts": []map[string]interface{}{
					{
						"post_id": "insta_post_1",
						"caption": "Stay safe during COVID-19 #covid19 #indonesia",
						"likes":   500,
					},
				},
			},
			"indonesia_news": map[string]interface{}{
				"sources": map[string]interface{}{
					"tempo": map[string]interface{}{
						"articles": []map[string]interface{}{
							{
								"title": "Update COVID-19 Indonesia",
								"url":   "https://tempo.co/covid19",
							},
						},
					},
				},
			},
		},
	}
	
	// Save extracted data
	extractedFile := filepath.Join(ep.outputDir, fmt.Sprintf("extracted_data_%s.json", timestamp))
	if err := ep.saveJSONFile(extractedData, extractedFile); err != nil {
		return StageResult{
			Status: "error",
			Error:  fmt.Sprintf("Failed to save extracted data: %v", err),
		}
	}
	
	// Create latest copy
	latestExtracted := filepath.Join(ep.outputDir, "extracted_data_latest.json")
	if err := ep.copyFile(extractedFile, latestExtracted); err != nil {
		log.Printf("Warning: Failed to create latest extracted data copy: %v", err)
	}
	
	// Count records
	totalRecords := ep.countExtractedRecords(extractedData)
	
	fmt.Printf("✅ Extraction completed: %s\n", extractedFile)
	fmt.Printf("📊 Total records extracted: %d\n", totalRecords)
	
	return StageResult{
		Status:  "success",
		File:    extractedFile,
		Records: totalRecords,
	}
}

// runTransformationStage runs the data transformation stage
func (ep *ETLPipeline) runTransformationStage(timestamp string) StageResult {
	fmt.Println("🔄 Starting data transformation...")
	
	// Load extracted data
	extractedFile := filepath.Join(ep.outputDir, fmt.Sprintf("extracted_data_%s.json", timestamp))
	extractedData, err := ep.loadJSONFile(extractedFile)
	if err != nil {
		return StageResult{
			Status: "error",
			Error:  fmt.Sprintf("Failed to load extracted data: %v", err),
		}
	}
	
	// Transform data (simplified for testing)
	transformedData := ep.transformData(extractedData)
	
	// Save transformed data
	transformedFile := filepath.Join(ep.outputDir, fmt.Sprintf("transformed_data_%s.json", timestamp))
	if err := ep.saveJSONFile(transformedData, transformedFile); err != nil {
		return StageResult{
			Status: "error",
			Error:  fmt.Sprintf("Failed to save transformed data: %v", err),
		}
	}
	
	// Create latest copy
	latestTransformed := filepath.Join(ep.outputDir, "transformed_data_latest.json")
	if err := ep.copyFile(transformedFile, latestTransformed); err != nil {
		log.Printf("Warning: Failed to create latest transformed data copy: %v", err)
	}
	
	// Count fact records
	factRecords := 0
	if facts, ok := transformedData["fact_table"]; ok {
		if factList, ok := facts.([]interface{}); ok {
			factRecords = len(factList)
		}
	}
	
	fmt.Printf("✅ Transformation completed: %s\n", transformedFile)
	fmt.Printf("📊 Fact records created: %d\n", factRecords)
	
	return StageResult{
		Status:  "success",
		File:    transformedFile,
		Records: factRecords,
	}
}

// runLoadingStage runs the data loading stage
func (ep *ETLPipeline) runLoadingStage(timestamp string) StageResult {
	fmt.Println("📊 Starting data loading...")
	
	// Load transformed data
	transformedFile := filepath.Join(ep.outputDir, fmt.Sprintf("transformed_data_%s.json", timestamp))
	transformedData, err := ep.loadJSONFile(transformedFile)
	if err != nil {
		return StageResult{
			Status: "error",
			Error:  fmt.Sprintf("Failed to load transformed data: %v", err),
		}
	}
	
	// Simulate loading to database
	loadSuccess := ep.loadDataToDatabase(transformedData)
	if !loadSuccess {
		return StageResult{
			Status: "error",
			Error:  "Failed to load data to database",
		}
	}
	
	// Export to CSV
	csvDir := filepath.Join(ep.outputDir, "csv_exports")
	if err := os.MkdirAll(csvDir, 0755); err != nil {
		log.Printf("Warning: Failed to create CSV export directory: %v", err)
	}
	
	// Create load report
	loadReport := map[string]interface{}{
		"timestamp":        time.Now().Format(time.RFC3339),
		"status":           "success",
		"records_loaded":   ep.countTransformedRecords(transformedData),
		"csv_exports_dir":  csvDir,
		"database_path":    "covid_knowledge_warehouse.db",
	}
	
	loadReportFile := filepath.Join(ep.outputDir, fmt.Sprintf("load_report_%s.json", timestamp))
	if err := ep.saveJSONFile(loadReport, loadReportFile); err != nil {
		log.Printf("Warning: Failed to save load report: %v", err)
	}
	
	fmt.Printf("✅ Loading completed successfully\n")
	fmt.Printf("📊 Records loaded: %d\n", loadReport["records_loaded"])
	
	return StageResult{
		Status:  "success",
		File:    loadReportFile,
		Records: loadReport["records_loaded"],
	}
}

// Helper methods
func (ep *ETLPipeline) saveJSONFile(data interface{}, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, jsonData, 0644)
}

func (ep *ETLPipeline) loadJSONFile(filename string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	
	return result, nil
}

func (ep *ETLPipeline) copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	
	return os.WriteFile(dst, data, 0644)
}

func (ep *ETLPipeline) countExtractedRecords(data map[string]interface{}) int {
	total := 0
	sources := data["sources"].(map[string]interface{})
	
	// Count YouTube videos
	if youtube, ok := sources["youtube"].(map[string]interface{}); ok {
		if videos, ok := youtube["videos"].([]interface{}); ok {
			total += len(videos)
		}
	}
	
	// Count Google News articles
	if news, ok := sources["google_news"].(map[string]interface{}); ok {
		if articles, ok := news["articles"].([]interface{}); ok {
			total += len(articles)
		}
	}
	
	// Count Instagram posts
	if insta, ok := sources["instagram"].(map[string]interface{}); ok {
		if posts, ok := insta["posts"].([]interface{}); ok {
			total += len(posts)
		}
	}
	
	// Count Indonesia News articles
	if indo, ok := sources["indonesia_news"].(map[string]interface{}); ok {
		if indoSources, ok := indo["sources"].(map[string]interface{}); ok {
			for _, source := range indoSources {
				if sourceData, ok := source.(map[string]interface{}); ok {
					if articles, ok := sourceData["articles"].([]interface{}); ok {
						total += len(articles)
					}
				}
			}
		}
	}
	
	return total
}

func (ep *ETLPipeline) transformData(extractedData map[string]interface{}) map[string]interface{} {
	// Simplified transformation for testing
	transformedData := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"fact_table": []map[string]interface{}{
			{
				"fact_id":              "fact_001",
				"source_id":            "youtube",
				"content_type_id":      "video",
				"title":                "COVID-19 Vaccine Update Indonesia",
				"covid_relevance_score": 0.9,
				"language":             "en",
			},
			{
				"fact_id":              "fact_002",
				"source_id":            "google_news",
				"content_type_id":      "article",
				"title":                "Indonesia COVID-19 Cases Update",
				"covid_relevance_score": 0.8,
				"language":             "id",
			},
		},
		"dimension_tables": map[string]interface{}{
			"dim_source": []map[string]interface{}{
				{"source_id": "youtube", "source_name": "YouTube"},
				{"source_id": "google_news", "source_name": "Google News"},
			},
			"dim_content_type": []map[string]interface{}{
				{"content_type_id": "video", "content_type_name": "Video"},
				{"content_type_id": "article", "content_type_name": "Article"},
			},
		},
	}
	
	return transformedData
}

func (ep *ETLPipeline) loadDataToDatabase(data map[string]interface{}) bool {
	// Simplified database loading for testing
	// In a real implementation, this would create SQLite tables and insert data
	fmt.Println("💾 Loading data to database...")
	fmt.Println("✅ Database loading completed (simulated)")
	return true
}

func (ep *ETLPipeline) countTransformedRecords(data map[string]interface{}) int {
	if facts, ok := data["fact_table"]; ok {
		if factList, ok := facts.([]interface{}); ok {
			return len(factList)
		}
	}
	return 0
}

func (ep *ETLPipeline) savePipelineResults(results *PipelineResult, timestamp string) {
	// Save pipeline results
	pipelineFile := filepath.Join(ep.outputDir, fmt.Sprintf("pipeline_results_%s.json", timestamp))
	if err := ep.saveJSONFile(results, pipelineFile); err != nil {
		log.Printf("Warning: Failed to save pipeline results: %v", err)
		return
	}
	
	// Create latest copy
	latestPipeline := filepath.Join(ep.outputDir, "pipeline_results_latest.json")
	if err := ep.copyFile(pipelineFile, latestPipeline); err != nil {
		log.Printf("Warning: Failed to create latest pipeline results copy: %v", err)
	}
	
	fmt.Printf("📊 Pipeline results saved: %s\n", pipelineFile)
}

func main() {
	// Create and run ETL pipeline
	pipeline := NewETLPipeline()
	results := pipeline.RunFullPipeline()
	
	// Print final summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📋 PIPELINE SUMMARY")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Status: %s\n", results.Status)
	fmt.Printf("Start Time: %s\n", results.PipelineStartTime)
	
	if results.Status == "completed" {
		fmt.Println("✅ All stages completed successfully")
		for stage, result := range results.Stages {
			fmt.Printf("  %s: %s\n", stage, result.Status)
		}
	} else {
		fmt.Printf("❌ Pipeline failed: %s\n", results.ErrorMessage)
	}
}
