package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"covid19-kms/database"
	"covid19-kms/internal/etl"
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

	// Set content type (CORS is handled by middleware)
	w.Header().Set("Content-Type", "application/json")

	// Get data from database
	data, err := h.retrieveLatestData()
	if err != nil {
		http.Error(w, "Failed to retrieve data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]interface{}{
		"status":      "success",
		"timestamp":   time.Now().Format(time.RFC3339),
		"data":        data,
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
			"source":          data.Source,
			"title":           data.Title,
			"content":         data.Content,
			"relevance_score": data.RelevanceScore,
			"sentiment":       data.Sentiment,
			"processed_at":    data.ProcessedAt.Format(time.RFC3339),
			"processed_data":  data.ProcessedData,
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

	// Set content type (CORS is handled by middleware)
	w.Header().Set("Content-Type", "application/json")

	// Get source from query parameter
	source := r.URL.Query().Get("source")
	if source == "" {
		http.Error(w, "Source parameter is required", http.StatusBadRequest)
		return
	}

	// Get data by source from database (get ALL data by passing limit = 0)
	data, err := database.GetDataBySource(source, 0)
	if err != nil {
		http.Error(w, "Failed to retrieve data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var results []map[string]interface{}
	for _, item := range data {
		result := map[string]interface{}{
			"source":          item.Source,
			"title":           item.Title,
			"content":         item.Content,
			"relevance_score": item.RelevanceScore,
			"sentiment":       item.Sentiment,
			"processed_at":    item.ProcessedAt.Format(time.RFC3339),
			"processed_data":  item.ProcessedData,
		}
		results = append(results, result)
	}

	// Return response
	response := map[string]interface{}{
		"status":      "success",
		"timestamp":   time.Now().Format(time.RFC3339),
		"source":      source,
		"data":        results,
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

	// Set content type (CORS is handled by middleware)
	w.Header().Set("Content-Type", "application/json")

	// Get data counts from database
	counts, err := database.GetDataCount()
	if err != nil {
		http.Error(w, "Failed to retrieve stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]interface{}{
		"status":    "success",
		"timestamp": time.Now().Format(time.RFC3339),
		"stats":     counts,
	}

	json.NewEncoder(w).Encode(response)
}

// GetYouTubeData retrieves YouTube data from database
func (h *DataHandler) GetYouTubeData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get YouTube data from database (get ALL data by passing limit = 0)
	data, err := database.GetDataBySource("youtube", 0)
	if err != nil {
		http.Error(w, "Failed to retrieve YouTube data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Process and enrich the data
	var enrichedData []map[string]interface{}
	for _, item := range data {
		// Create the enriched item with the expected structure
		enrichedItem := map[string]interface{}{
			"id":                    item.ID,
			"title":                 item.Title,
			"description":           item.Content,
			"covid_relevance_score": item.RelevanceScore,
		}

		// Parse processed_data JSON to extract metadata
		var metadata map[string]interface{}
		if item.ProcessedData != "" {
			if err := json.Unmarshal([]byte(item.ProcessedData), &metadata); err == nil {
				// Create the metadata structure as expected
				metadataStructure := map[string]interface{}{
					"video": map[string]interface{}{
						"id":       metadata["video_id"],
						"title":    item.Title,
						"views":    metadata["views"],
						"duration": metadata["duration"],
						"likes":    metadata["likes"],
					},
					"comment": map[string]interface{}{
						"id":       item.ID,
						"content":  item.Content,
						"language": metadata["language"],
					},
				}
				enrichedItem["metadata"] = metadataStructure
			}
		}

		enrichedData = append(enrichedData, enrichedItem)
	}

	response := map[string]interface{}{
		"data": enrichedData,
	}

	json.NewEncoder(w).Encode(response)
}

// GetGoogleNewsData retrieves Google News data from database
func (h *DataHandler) GetGoogleNewsData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get Google News data from database (get ALL data by passing limit = 0)
	data, err := database.GetDataBySource("google_news", 0)
	if err != nil {
		http.Error(w, "Failed to retrieve Google News data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Process and enrich the data
	var enrichedData []map[string]interface{}
	for _, item := range data {
		enrichedItem := map[string]interface{}{
			"id":              item.ID,
			"source":          item.Source,
			"processed_at":    item.ProcessedAt,
			"title":           item.Title,
			"content":         item.Content,
			"relevance_score": item.RelevanceScore,
			"sentiment":       item.Sentiment,
		}

		// Parse processed_data JSON to extract metadata
		if item.ProcessedData != "" {
			var metadata map[string]interface{}
			if err := json.Unmarshal([]byte(item.ProcessedData), &metadata); err == nil {
				// Extract news-specific metadata
				if url, ok := metadata["url"]; ok {
					enrichedItem["url"] = url
				}
				if author, ok := metadata["author"]; ok {
					enrichedItem["author"] = author
				}
				if publishedAt, ok := metadata["published_at"]; ok {
					enrichedItem["published_at"] = publishedAt
				}
				if source, ok := metadata["source"]; ok {
					enrichedItem["news_source"] = source
				}
				if language, ok := metadata["language"]; ok {
					enrichedItem["language"] = language
				}
				if category, ok := metadata["category"]; ok {
					enrichedItem["category"] = category
				}
			}
		}

		enrichedData = append(enrichedData, enrichedItem)
	}

	response := map[string]interface{}{
		"status":      "success",
		"timestamp":   time.Now().Format(time.RFC3339),
		"source":      "google_news",
		"data":        enrichedData,
		"total_count": len(enrichedData),
	}

	json.NewEncoder(w).Encode(response)
}

// GetInstagramData retrieves Instagram data from database
func (h *DataHandler) GetInstagramData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get Instagram data from database (get ALL data by passing limit = 0)
	data, err := database.GetDataBySource("instagram", 0)
	if err != nil {
		http.Error(w, "Failed to retrieve Instagram data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Process and enrich the data
	var enrichedData []map[string]interface{}
	for _, item := range data {
		enrichedItem := map[string]interface{}{
			"id":              item.ID,
			"source":          item.Source,
			"processed_at":    item.ProcessedAt,
			"title":           item.Title,
			"content":         item.Content,
			"relevance_score": item.RelevanceScore,
			"sentiment":       item.Sentiment,
		}

		// Parse processed_data JSON to extract metadata
		if item.ProcessedData != "" {
			var metadata map[string]interface{}
			if err := json.Unmarshal([]byte(item.ProcessedData), &metadata); err == nil {
				// Extract Instagram-specific metadata
				if likes, ok := metadata["likes"]; ok {
					enrichedItem["likes"] = likes
				}
				if comments, ok := metadata["comments"]; ok {
					enrichedItem["comments"] = comments
				}
				if postID, ok := metadata["post_id"]; ok {
					enrichedItem["post_id"] = postID
				}
				if hashtags, ok := metadata["hashtags"]; ok {
					enrichedItem["hashtags"] = hashtags
				}
				if followers, ok := metadata["followers"]; ok {
					enrichedItem["followers"] = followers
				}
				if mediaType, ok := metadata["media_type"]; ok {
					enrichedItem["media_type"] = mediaType
				}
			}
		}

		enrichedData = append(enrichedData, enrichedItem)
	}

	response := map[string]interface{}{
		"status":      "success",
		"timestamp":   time.Now().Format(time.RFC3339),
		"source":      "instagram",
		"data":        enrichedData,
		"total_count": len(enrichedData),
	}

	json.NewEncoder(w).Encode(response)
}

// GetIndonesiaNewsData retrieves Indonesia News data from database or fresh from API
func (h *DataHandler) GetIndonesiaNewsData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Check if user wants fresh data
	refresh := r.URL.Query().Get("refresh")

	if refresh == "true" {
		// Fetch fresh data from the scraper
		h.getFreshIndonesiaNewsData(w, r)
		return
	}

	// Get Indonesia News data from database (get ALL data by passing limit = 0)
	data, err := database.GetDataBySource("indonesia_news", 0)
	if err != nil {
		http.Error(w, "Failed to retrieve Indonesia News data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Process and enrich the data
	var enrichedData []map[string]interface{}
	for _, item := range data {
		enrichedItem := map[string]interface{}{
			"id":              item.ID,
			"source":          item.Source,
			"processed_at":    item.ProcessedAt,
			"title":           item.Title,
			"content":         item.Content,
			"relevance_score": item.RelevanceScore,
			"sentiment":       item.Sentiment,
		}

		// Parse processed_data JSON to extract metadata
		if item.ProcessedData != "" {
			var metadata map[string]interface{}
			if err := json.Unmarshal([]byte(item.ProcessedData), &metadata); err == nil {
				// Extract Indonesia news-specific metadata
				if url, ok := metadata["url"]; ok {
					enrichedItem["url"] = url
				}
				if author, ok := metadata["author"]; ok {
					enrichedItem["author"] = author
				}
				if publishedAt, ok := metadata["published_at"]; ok {
					enrichedItem["published_at"] = publishedAt
				}
				if newsSource, ok := metadata["news_source"]; ok {
					enrichedItem["news_source"] = newsSource
				}
				if language, ok := metadata["language"]; ok {
					enrichedItem["language"] = language
				}
				if category, ok := metadata["category"]; ok {
					enrichedItem["category"] = category
				}
				if region, ok := metadata["region"]; ok {
					enrichedItem["region"] = region
				}
			}
		}

		enrichedData = append(enrichedData, enrichedItem)
	}

	response := map[string]interface{}{
		"status":      "success",
		"timestamp":   time.Now().Format(time.RFC3339),
		"source":      "indonesia_news",
		"data":        enrichedData,
		"total_count": len(enrichedData),
	}

	json.NewEncoder(w).Encode(response)
}

// getFreshIndonesiaNewsData fetches fresh data directly from the Indonesia news scraper
func (h *DataHandler) getFreshIndonesiaNewsData(w http.ResponseWriter, r *http.Request) {
	// Create ETL extractor to get fresh data
	extractor := etl.NewDataExtractor()

	// Extract all sources data (we'll filter for Indonesia news)
	extractedData := extractor.ExtractAllSources()

	// Debug logging
	fmt.Printf("DEBUG: Extracted data sources: %v\n", len(extractedData.Sources))
	for key := range extractedData.Sources {
		fmt.Printf("DEBUG: Source key: %s\n", key)
	}

	// Get Indonesia news data from the extracted sources
	indonesiaData, ok := extractedData.Sources["indonesia_news"]
	if !ok {
		http.Error(w, "Failed to extract Indonesia News data", http.StatusInternalServerError)
		return
	}

	// For debugging, let's return the raw data structure
	response := map[string]interface{}{
		"status":    "success",
		"timestamp": time.Now().Format(time.RFC3339),
		"source":    "indonesia_news",
		"raw_data":  indonesiaData,
		"note":      "Raw ETL data for debugging",
	}

	json.NewEncoder(w).Encode(response)
}

// GetDataSummary retrieves overall data summary from database
func (h *DataHandler) GetDataSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get summary data from database
	summary, err := database.GetDataSummary()
	if err != nil {
		http.Error(w, "Failed to retrieve data summary: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":    "success",
		"timestamp": time.Now().Format(time.RFC3339),
		"summary":   summary,
	}

	json.NewEncoder(w).Encode(response)
}
