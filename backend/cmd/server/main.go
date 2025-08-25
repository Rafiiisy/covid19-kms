package main

import (
	"log"
	"os"

	"covid19-kms/database"
	"covid19-kms/internal/api"
	"covid19-kms/internal/etl"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"service":   "COVID-19 KMS Backend",
			"language":  "Go",
			"framework": "Gin",
		})
	})

	// Initialize database connection
	if err := database.InitDatabase(); err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		log.Println("⚠️ Some endpoints may not work without database connection")
	} else {
		log.Println("✅ Database connection established")
	}

	// Initialize data handler
	dataHandler := api.NewDataHandler()

	// API routes
	apiGroup := r.Group("/api")
	{
		// ETL endpoints
		apiGroup.POST("/etl/run", func(c *gin.Context) {
			// Create and run ETL pipeline
			orchestrator := etl.NewETLOrchestrator()
			result := orchestrator.RunETLPipeline()

			c.JSON(200, gin.H{
				"status":  "success",
				"message": "ETL pipeline completed successfully",
				"result":  result,
			})
		})

		apiGroup.GET("/etl/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":      "success",
				"timestamp":   "2025-08-25T17:55:28.298311Z",
				"service":     "COVID-19 KMS Backend",
				"version":     "1.0.0",
				"endpoints":   []string{"/api/etl/run", "/api/etl/status", "/api/etl/data/*"},
				"description": "ETL Pipeline API for COVID-19 data processing",
			})
		})

		// Data endpoints
		apiGroup.GET("/etl/data/youtube", func(c *gin.Context) {
			dataHandler.GetYouTubeData(c.Writer, c.Request)
		})

		apiGroup.GET("/etl/data/google-news", func(c *gin.Context) {
			dataHandler.GetGoogleNewsData(c.Writer, c.Request)
		})

		apiGroup.GET("/etl/data/indonesia-news", func(c *gin.Context) {
			dataHandler.GetIndonesiaNewsData(c.Writer, c.Request)
		})

		apiGroup.GET("/etl/data/instagram", func(c *gin.Context) {
			dataHandler.GetInstagramData(c.Writer, c.Request)
		})

		// Summary endpoint
		apiGroup.GET("/etl/data/summary", func(c *gin.Context) {
			dataHandler.GetDataSummary(c.Writer, c.Request)
		})

		// Analytics endpoints
		apiGroup.GET("/analytics/summary", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"total_records": 0,
				"sources":       []string{"youtube", "google_news", "indonesia_news", "covid_news"},
			})
		})
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Note: Database connection will be maintained throughout server lifetime
	// It will be closed when the process terminates

	log.Printf("Starting COVID-19 KMS Backend on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
