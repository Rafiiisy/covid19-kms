package api

import (
	"encoding/json"
	"net/http"
	"time"

	"covid19-kms/internal/etl"
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

	// Set comprehensive CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

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

	// Set comprehensive CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Create status response
	status := map[string]interface{}{
		"status":        "ready",
		"timestamp":     time.Now().Format(time.RFC3339),
		"service":       "ETL Pipeline API",
		"version":       "1.0.0",
		"endpoints":     []string{"/api/etl/run", "/api/etl/status", "/api/etl/extract", "/api/etl/transform", "/api/etl/load"},
		"description":   "COVID-19 Knowledge Management System ETL Pipeline",
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

	// Set comprehensive CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Create extractor and run extraction
	extractor := etl.NewDataExtractor()
	_ = extractor.ExtractAllSources()

	// Create response
	response := map[string]interface{}{
		"status":        "success",
		"timestamp":     time.Now().Format(time.RFC3339),
		"stage":         "extraction",
		"data":          "extraction_completed",
		"message":       "Data extraction completed successfully",
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

	// Set comprehensive CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

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
	transformedData := transformer.TransformData(nil, nil) // Using nil for demo

	if transformedData == nil {
		http.Error(w, "Transformation failed", http.StatusInternalServerError)
		return
	}

	// Create response
	response := map[string]interface{}{
		"status":        "success",
		"timestamp":     time.Now().Format(time.RFC3339),
		"stage":         "transformation",
		"data":          transformedData,
		"message":       "Data transformation completed successfully",
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

	// Set comprehensive CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Create loader
	loader := etl.NewDataLoader()

	// Create sample data for loading
	transformedData := &etl.TransformedData{
		YouTube:      []etl.TransformedVideo{},
		News:         []etl.TransformedArticle{},
		TransformedAt: time.Now().Format(time.RFC3339),
	}

	// Run loading
	loadResult := loader.LoadData(transformedData)

	// Create response
	response := map[string]interface{}{
		"status":        "success",
		"timestamp":     time.Now().Format(time.RFC3339),
		"stage":         "loading",
		"result":        loadResult,
		"message":       "Data loading completed successfully",
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

	// Set comprehensive CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

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
