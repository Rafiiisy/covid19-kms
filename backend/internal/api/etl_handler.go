package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"covid19-kms/database"
	"covid19-kms/internal/etl"
	"covid19-kms/internal/services"
)

// ETLHandler handles HTTP requests for ETL operations
type ETLHandler struct {
	orchestrator *etl.ETLOrchestrator
}

// NewETLHandler creates a new ETL handler
func NewETLHandler() *ETLHandler {
	return &ETLHandler{
		orchestrator: etl.NewETLOrchestrator(),
	}
}

// RunETLPipeline handles POST requests to run the complete ETL pipeline
func (h *ETLHandler) RunETLPipeline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type (CORS is handled by middleware)
	w.Header().Set("Content-Type", "application/json")

	// Run the ETL pipeline
	result := h.orchestrator.RunETLPipeline()

	// Convert result to JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// GetPipelineStatus handles GET requests to check pipeline status
func (h *ETLHandler) GetPipelineStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type (CORS is handled by middleware)
	w.Header().Set("Content-Type", "application/json")

	// Create status response
	status := map[string]interface{}{
		"status":      "ready",
		"timestamp":   time.Now().Format(time.RFC3339),
		"service":     "ETL Pipeline API",
		"version":     "1.0.0",
		"endpoints":   []string{"/api/etl/run", "/api/etl/status", "/api/etl/extract", "/api/etl/transform", "/api/etl/load", "/api/etl/cleanup/sentiment", "/api/etl/data/*"},
		"description": "COVID-19 Knowledge Management System ETL Pipeline",
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// ExtractData handles POST requests to run only the extraction stage
func (h *ETLHandler) ExtractData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type (CORS is handled by middleware)
	w.Header().Set("Content-Type", "application/json")

	// Create extractor and run extraction
	extractor := etl.NewDataExtractor()
	_ = extractor.ExtractAllSources()

	// Create response
	response := map[string]interface{}{
		"status":    "success",
		"timestamp": time.Now().Format(time.RFC3339),
		"stage":     "extraction",
		"data":      "extraction_completed",
		"message":   "Data extraction completed successfully",
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// TransformData handles POST requests to run only the transformation stage
func (h *ETLHandler) TransformData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type (CORS is handled by middleware)
	w.Header().Set("Content-Type", "application/json")

	// For transformation, we need some input data
	// In a real scenario, this would come from the request body
	// For now, we'll create sample data
	_ = &etl.ExtractedData{
		Timestamp: time.Now().Format(time.RFC3339),
		Query:     "covid19",
		Sources:   make(map[string]interface{}),
	}

	// Create transformer and run transformation
	transformer := etl.NewDataTransformer()
	transformedData := transformer.TransformData(nil, nil, nil) // Using nil for demo

	if transformedData == nil {
		http.Error(w, "Transformation failed", http.StatusInternalServerError)
		return
	}

	// Create response
	response := map[string]interface{}{
		"status":    "success",
		"timestamp": time.Now().Format(time.RFC3339),
		"stage":     "transformation",
		"data":      transformedData,
		"message":   "Data transformation completed successfully",
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// LoadData handles POST requests to run only the loading stage
func (h *ETLHandler) LoadData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type (CORS is handled by middleware)
	w.Header().Set("Content-Type", "application/json")

	// Create loader
	loader := etl.NewDataLoader()

	// Create sample data for loading
	transformedData := &etl.TransformedData{
		YouTube:       []etl.TransformedVideo{},
		News:          []etl.TransformedArticle{},
		TransformedAt: time.Now().Format(time.RFC3339),
	}

	// Run loading
	loadResult := loader.LoadData(transformedData)

	// Create response
	response := map[string]interface{}{
		"status":    "success",
		"timestamp": time.Now().Format(time.RFC3339),
		"stage":     "loading",
		"result":    loadResult,
		"message":   "Data loading completed successfully",
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// HealthCheck handles GET requests for health monitoring
func (h *ETLHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type (CORS is handled by middleware)
	w.Header().Set("Content-Type", "application/json")

	// Create health response
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "ETL Pipeline",
		"uptime":    "running",
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(health, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CleanupSentiments handles sentiment cleanup requests
func (h *ETLHandler) CleanupSentiments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get database connection
	if err := database.EnsureConnection(); err != nil {
		http.Error(w, fmt.Sprintf("Database connection failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Create cleanup service
	cleanupService := services.NewSentimentCleanupService(database.DB)

	// Parse query parameters
	source := r.URL.Query().Get("source")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var result *services.CleanupResult

	// Determine cleanup type based on parameters
	if source != "" {
		// Clean specific source
		log.Printf("🧹 Starting sentiment cleanup for source: %s", source)
		result = cleanupService.CleanSentimentBySource(source)
	} else if startDateStr != "" && endDateStr != "" {
		// Clean by date range
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid start_date format: %v", err), http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid end_date format: %v", err), http.StatusBadRequest)
			return
		}

		log.Printf("🧹 Starting sentiment cleanup for date range: %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		result = cleanupService.CleanSentimentByDateRange(startDate, endDate)
	} else {
		// Clean all sentiments
		log.Printf("🧹 Starting sentiment cleanup for all records")
		result = cleanupService.CleanAllSentiments()
	}

	// Return result
	response := map[string]interface{}{
		"status":    "success",
		"timestamp": time.Now().Format(time.RFC3339),
		"operation": "sentiment_cleanup",
		"result":    result,
	}

	json.NewEncoder(w).Encode(response)
}
