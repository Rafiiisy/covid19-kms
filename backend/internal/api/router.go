package api

import (
	"encoding/json"
	"net/http"
)

// Router handles HTTP routing for the ETL API
type Router struct {
	etlHandler  *ETLHandler
	dataHandler *DataHandler
}

// NewRouter creates a new router instance
func NewRouter() *Router {
	return &Router{
		etlHandler:  NewETLHandler(),
		dataHandler: NewDataHandler(),
	}
}

// SetupRoutes configures all API routes
func (r *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Wrap all routes with CORS middleware
	mux.HandleFunc("/", r.corsMiddleware(r.handleRoot))
	mux.HandleFunc("/api", r.corsMiddleware(r.handleAPIInfo))
	mux.HandleFunc("/api/etl/run", r.corsMiddleware(r.etlHandler.RunETLPipeline))
	mux.HandleFunc("/api/etl/status", r.corsMiddleware(r.etlHandler.GetPipelineStatus))
	mux.HandleFunc("/api/etl/extract", r.corsMiddleware(r.etlHandler.ExtractData))
	mux.HandleFunc("/api/etl/transform", r.corsMiddleware(r.etlHandler.TransformData))
	mux.HandleFunc("/api/etl/load", r.corsMiddleware(r.etlHandler.LoadData))
	mux.HandleFunc("/api/etl/data", r.corsMiddleware(r.dataHandler.GetLatestData))
	mux.HandleFunc("/api/etl/data/source", r.corsMiddleware(r.dataHandler.GetDataBySource))
	mux.HandleFunc("/api/etl/data/stats", r.corsMiddleware(r.dataHandler.GetDataStats))

	// New database query endpoints for individual sources
	mux.HandleFunc("/api/etl/data/youtube", r.corsMiddleware(r.dataHandler.GetYouTubeData))
	mux.HandleFunc("/api/etl/data/google-news", r.corsMiddleware(r.dataHandler.GetGoogleNewsData))
	mux.HandleFunc("/api/etl/data/instagram", r.corsMiddleware(r.dataHandler.GetInstagramData))
	mux.HandleFunc("/api/etl/data/indonesia-news", r.corsMiddleware(r.dataHandler.GetIndonesiaNewsData))
	mux.HandleFunc("/api/etl/data/summary", r.corsMiddleware(r.dataHandler.GetDataSummary))

	mux.HandleFunc("/health", r.corsMiddleware(r.etlHandler.HealthCheck))
	mux.HandleFunc("/api/health", r.corsMiddleware(r.etlHandler.HealthCheck))

	return mux
}

// handleRoot handles the root endpoint
func (r *Router) handleRoot(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"service":     "COVID-19 Knowledge Management System",
		"version":     "1.0.0",
		"description": "ETL Pipeline API for COVID-19 data processing",
		"endpoints": map[string]interface{}{
			"root":     "/",
			"api_info": "/api",
			"etl": map[string]string{
				"run_pipeline":   "/api/etl/run",
				"status":         "/api/etl/status",
				"extract":        "/api/etl/extract",
				"transform":      "/api/etl/transform",
				"load":           "/api/etl/load",
				"data":           "/api/etl/data",
				"data_by_source": "/api/etl/data/source?source=youtube",
				"data_stats":     "/api/etl/data/stats",
			},
			"health": "/api/health",
		},
		"documentation": "API documentation available at /api",
	}

	// Convert to JSON and send response
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// handleAPIInfo handles the /api endpoint
func (r *Router) handleAPIInfo(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/api" {
		http.NotFound(w, req)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"api_name":    "ETL Pipeline API",
		"version":     "1.0.0",
		"description": "RESTful API for COVID-19 ETL pipeline operations",
		"base_url":    "/api",
		"endpoints": map[string]interface{}{
			"etl": map[string]interface{}{
				"run_pipeline": map[string]interface{}{
					"method":      "POST",
					"url":         "/api/etl/run",
					"description": "Run the complete ETL pipeline",
					"body":        "none",
					"response":    "ETLResult with pipeline execution details",
				},
				"status": map[string]interface{}{
					"method":      "GET",
					"url":         "/api/etl/status",
					"description": "Get current pipeline status and API information",
					"body":        "none",
					"response":    "API status and endpoint information",
				},
				"extract": map[string]interface{}{
					"method":      "POST",
					"url":         "/api/etl/extract",
					"description": "Run only the data extraction stage",
					"body":        "none",
					"response":    "ExtractedData with raw data from all sources",
				},
				"transform": map[string]interface{}{
					"method":      "POST",
					"url":         "/api/etl/transform",
					"description": "Run only the data transformation stage",
					"body":        "none",
					"response":    "TransformedData with cleaned and enriched data",
				},
				"load": map[string]interface{}{
					"method":      "POST",
					"url":         "/api/etl/load",
					"description": "Run only the data loading stage",
					"body":        "none",
					"response":    "LoadResult with loading operation details",
				},
			},
			"health": map[string]interface{}{
				"method":      "GET",
				"url":         "/api/health",
				"description": "Health check endpoint for monitoring",
				"body":        "none",
				"response":    "Health status information",
			},
		},
		"data_sources": []string{
			"YouTube (videos and comments)",
			"Google News (articles)",
			"Instagram (posts and media)",
			"Indonesia News (local news sources)",
		},
		"supported_formats": []string{
			"JSON (request/response)",
			"CSV (data export)",
			"SQLite (local storage)",
		},
		"authentication": "Currently none (development mode)",
		"rate_limiting":  "Not implemented (development mode)",
	}

	// Convert to JSON and send response
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CORS middleware for handling cross-origin requests
func (r *Router) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Set comprehensive CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin, Cache-Control, X-File-Name")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Content-Disposition")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, req)
	}
}
