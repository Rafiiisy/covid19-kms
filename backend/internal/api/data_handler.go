package api

import (
	"encoding/json"
	"net/http"
	"time"

	"covid19-kms/database"
)

// DataHandler handles data retrieval from PostgreSQL database
type DataHandler struct {
	// No need for dataDir since we're using database
}

// NewDataHandler creates a new data handler
func NewDataHandler() *DataHandler {
	return &DataHandler{}
}

// GetLatestData retrieves the latest data from PostgreSQL database
func (h *DataHandler) GetLatestData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Get data from database
	data, err := h.retrieveLatestData()
	if err != nil {
		http.Error(w, "Failed to retrieve data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]interface{}{
		"status":     "success",
		"timestamp":  time.Now().Format(time.RFC3339),
		"data":       data,
		"total_count": len(data),
	}

	json.NewEncoder(w).Encode(response)
}

// retrieveLatestData fetches latest data from PostgreSQL database
func (h *DataHandler) retrieveLatestData() ([]map[string]interface{}, error) {
	// Get latest processed data from database
	processedData, err := database.GetLatestProcessedData(100)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, data := range processedData {
		// Convert database model to response format
		result := map[string]interface{}{
			"source":         data.Source,
			"title":          data.Title,
			"content":        data.Content,
			"relevance_score": data.RelevanceScore,
			"sentiment":      data.Sentiment,
			"processed_at":   data.ProcessedAt.Format(time.RFC3339),
			"processed_data": data.ProcessedData,
		}
		results = append(results, result)
	}

	return results, nil
}

// GetDataBySource retrieves data filtered by source
func (h *DataHandler) GetDataBySource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Get source from query parameter
	source := r.URL.Query().Get("source")
	if source == "" {
		http.Error(w, "Source parameter is required", http.StatusBadRequest)
		return
	}

	// Get data by source from database
	data, err := database.GetDataBySource(source, 100)
	if err != nil {
		http.Error(w, "Failed to retrieve data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var results []map[string]interface{}
	for _, item := range data {
		result := map[string]interface{}{
			"source":         item.Source,
			"title":          item.Title,
			"content":        item.Content,
			"relevance_score": item.RelevanceScore,
			"sentiment":      item.Sentiment,
			"processed_at":   item.ProcessedAt.Format(time.RFC3339),
			"processed_data": item.ProcessedData,
		}
		results = append(results, result)
	}

	// Return response
	response := map[string]interface{}{
		"status":     "success",
		"timestamp":  time.Now().Format(time.RFC3339),
		"source":     source,
		"data":       results,
		"total_count": len(results),
	}

	json.NewEncoder(w).Encode(response)
}

// GetDataStats retrieves database statistics
func (h *DataHandler) GetDataStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set CORS headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Get data counts from database
	counts, err := database.GetDataCount()
	if err != nil {
		http.Error(w, "Failed to retrieve stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]interface{}{
		"status":     "success",
		"timestamp":  time.Now().Format(time.RFC3339),
		"stats":      counts,
	}

	json.NewEncoder(w).Encode(response)
}
