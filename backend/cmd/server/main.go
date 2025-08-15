package main

import (
	"log"
	"os"

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

	// API routes
	api := r.Group("/api")
	{
		// ETL endpoints
		api.POST("/etl/refresh", func(c *gin.Context) {
			// Import the ETL package and run the pipeline
			// This is a placeholder - in a real implementation, you would:
			// 1. Import the etl package
			// 2. Create an orchestrator
			// 3. Run the pipeline
			// 4. Return the result
			
			c.JSON(200, gin.H{
				"status":  "success",
				"message": "ETL pipeline triggered (placeholder - ETL package ready)",
				"note":    "ETL package has been created and is ready for integration",
			})
		})

		// Data endpoints
		api.GET("/data/youtube", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "success",
				"data":   []string{},
				"source": "youtube",
			})
		})

		api.GET("/data/news", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "success",
				"data":   []string{},
				"source": "news",
			})
		})

		// Analytics endpoints
		api.GET("/analytics/summary", func(c *gin.Context) {
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

	log.Printf("Starting COVID-19 KMS Backend on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
