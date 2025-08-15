package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"covid19-kms/internal/api"
	"covid19-kms/database"
)

func main() {
	log.Println("🚀 Starting COVID-19 KMS ETL API Server")

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	// Create tables if they don't exist
	if err := database.CreateTables(); err != nil {
		log.Fatalf("❌ Failed to create database tables: %v", err)
	}

	// Create router
	router := api.NewRouter()

	// Create server
	server := &http.Server{
		Addr:    ":8000",
		Handler: router.SetupRoutes(),
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🌐 Server starting on port 8000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🔄 Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
