package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"covid19-kms/internal/config"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	router *Router
	server *http.Server
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) *Server {
	router := NewRouter()
	
	return &Server{
		config: cfg,
		router: router,
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
			Handler:      router.SetupRoutes(),
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Create a channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("🚀 Starting ETL API server on %s:%s", s.config.Server.Host, s.config.Server.Port)
		log.Printf("📊 Environment: %s", s.getEnvironment())
		log.Printf("🔗 API Documentation: http://%s:%s/api", s.config.Server.Host, s.config.Server.Port)
		log.Printf("🏥 Health Check: http://%s:%s/api/health", s.config.Server.Host, s.config.Server.Port)
		
		serverErrors <- s.server.ListenAndServe()
	}()

	// Create a channel to listen for an interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking select waiting for either a server error or a shutdown signal
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("🛑 Start shutdown... Signal: %v", sig)
		
		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Gracefully shutdown the server
		if err := s.server.Shutdown(ctx); err != nil {
			log.Printf("❌ Could not stop server gracefully: %v", err)
			if err := s.server.Close(); err != nil {
				return fmt.Errorf("could not force close server: %w", err)
			}
		}
	}

	return nil
}

// Stop gracefully stops the server
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// getEnvironment returns the current environment
func (s *Server) getEnvironment() string {
	if s.config.IsProduction() {
		return "production"
	}
	return "development"
}

// RunServer is a convenience function to start the server with default configuration
func RunServer() error {
	// Load environment variables
	if err := config.LoadDefaultEnv(); err != nil {
		log.Printf("⚠️ Warning: Could not load .env file: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration based on environment
	if err := validateConfiguration(cfg); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Create and start server
	server := NewServer(cfg)
	return server.Start()
}

// validateConfiguration validates the configuration based on the environment
func validateConfiguration(cfg *config.Config) error {
	if cfg.IsProduction() {
		// In production, API keys are required
		requiredKeys := config.GetRequiredEnvsForProduction()
		if err := config.ValidateRequiredEnvs(requiredKeys); err != nil {
			return fmt.Errorf("production environment requires API keys: %w", err)
		}
		log.Println("✅ Production configuration validated")
	} else {
		// In development/test, API keys are optional
		log.Println("✅ Development configuration loaded (API keys optional)")
	}

	return nil
}

// StartServerWithConfig starts the server with a specific configuration
func StartServerWithConfig(cfg *config.Config) error {
	server := NewServer(cfg)
	return server.Start()
}
