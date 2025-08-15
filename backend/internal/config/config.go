package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the ETL system
type Config struct {
	// Server configuration
	Server ServerConfig `json:"server"`
	
	// ETL Pipeline configuration
	ETL ETLConfig `json:"etl"`
	
	// API configuration
	API APIConfig `json:"api"`
	
	// Database configuration
	Database DatabaseConfig `json:"database"`
	
	// External APIs configuration
	ExternalAPIs ExternalAPIsConfig `json:"external_apis"`
	
	// Logging configuration
	Logging LoggingConfig `json:"logging"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string        `json:"port"`
	Host         string        `json:"host"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// ETLConfig holds ETL pipeline configuration
type ETLConfig struct {
	MaxConcurrentExtractions int           `json:"max_concurrent_extractions"`
	ExtractionTimeout        time.Duration `json:"extraction_timeout"`
	TransformationTimeout    time.Duration `json:"transformation_timeout"`
	LoadingTimeout           time.Duration `json:"loading_timeout"`
	BatchSize                int           `json:"batch_size"`
	RetryAttempts            int           `json:"retry_attempts"`
	RetryDelay               time.Duration `json:"retry_delay"`
}

// APIConfig holds API-related configuration
type APIConfig struct {
	EnableCORS        bool   `json:"enable_cors"`
	EnableLogging     bool   `json:"enable_logging"`
	EnableMetrics     bool   `json:"enable_metrics"`
	RateLimitRequests int    `json:"rate_limit_requests"`
	RateLimitWindow   string `json:"rate_limit_window"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type      string `json:"type"`       // "sqlite", "postgres", "mysql"
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	SSLMode   string `json:"ssl_mode"`
	MaxConns  int    `json:"max_connections"`
	IdleConns int    `json:"idle_connections"`
}

// ExternalAPIsConfig holds external API configuration
type ExternalAPIsConfig struct {
	YouTube      YouTubeConfig      `json:"youtube"`
	GoogleNews   GoogleNewsConfig   `json:"google_news"`
	Instagram    InstagramConfig    `json:"instagram"`
	IndonesiaNews IndonesiaNewsConfig `json:"indonesia_news"`
}

// YouTubeConfig holds YouTube API configuration
type YouTubeConfig struct {
	APIKey     string `json:"api_key"`
	Host       string `json:"host"`
	MaxResults int    `json:"max_results"`
	Timeout    int    `json:"timeout"`
}

// GoogleNewsConfig holds Google News API configuration
type GoogleNewsConfig struct {
	APIKey     string `json:"api_key"`
	Host       string `json:"host"`
	MaxResults int    `json:"max_results"`
	Language   string `json:"language"`
	Country    string `json:"country"`
}

// InstagramConfig holds Instagram API configuration
type InstagramConfig struct {
	APIKey     string `json:"api_key"`
	Host       string `json:"host"`
	MaxResults int    `json:"max_results"`
	Timeout    int    `json:"timeout"`
}

// IndonesiaNewsConfig holds Indonesia News API configuration
type IndonesiaNewsConfig struct {
	APIKey     string `json:"api_key"`
	Host       string `json:"host"`
	MaxResults int    `json:"max_results"`
	Sources    string `json:"sources"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `json:"level"`       // "debug", "info", "warn", "error"
	Format     string `json:"format"`      // "json", "text"
	Output     string `json:"output"`      // "stdout", "file"
	FilePath   string `json:"file_path"`
	MaxSize    int    `json:"max_size"`    // MB
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`     // days
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8000"),
			Host:         getEnv("SERVER_HOST", "localhost"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		ETL: ETLConfig{
			MaxConcurrentExtractions: getIntEnv("ETL_MAX_CONCURRENT_EXTRACTIONS", 5),
			ExtractionTimeout:        getDurationEnv("ETL_EXTRACTION_TIMEOUT", 5*time.Minute),
			TransformationTimeout:    getDurationEnv("ETL_TRANSFORMATION_TIMEOUT", 2*time.Minute),
			LoadingTimeout:           getDurationEnv("ETL_LOADING_TIMEOUT", 3*time.Minute),
			BatchSize:                getIntEnv("ETL_BATCH_SIZE", 100),
			RetryAttempts:            getIntEnv("ETL_RETRY_ATTEMPTS", 3),
			RetryDelay:               getDurationEnv("ETL_RETRY_DELAY", 5*time.Second),
		},
		API: APIConfig{
			EnableCORS:        getBoolEnv("API_ENABLE_CORS", true),
			EnableLogging:     getBoolEnv("API_ENABLE_LOGGING", true),
			EnableMetrics:     getBoolEnv("API_ENABLE_METRICS", true),
			RateLimitRequests: getIntEnv("API_RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:   getEnv("API_RATE_LIMIT_WINDOW", "1m"),
		},
		Database: DatabaseConfig{
			Type:      getEnv("DB_TYPE", "sqlite"),
			Host:      getEnv("DB_HOST", "localhost"),
			Port:      getIntEnv("DB_PORT", 5432),
			Username:  getEnv("DB_USERNAME", ""),
			Password:  getEnv("DB_PASSWORD", ""),
			Database:  getEnv("DB_DATABASE", "covid19_kms"),
			SSLMode:   getEnv("DB_SSL_MODE", "disable"),
			MaxConns:  getIntEnv("DB_MAX_CONNECTIONS", 10),
			IdleConns: getIntEnv("DB_IDLE_CONNECTIONS", 5),
		},
		ExternalAPIs: ExternalAPIsConfig{
			YouTube: YouTubeConfig{
				APIKey:     getEnv("YOUTUBE_API_KEY", ""),
				Host:       getEnv("YOUTUBE_HOST", "yt-api.p.rapidapi.com"),
				MaxResults: getIntEnv("YOUTUBE_MAX_RESULTS", 50),
				Timeout:    getIntEnv("YOUTUBE_TIMEOUT", 30),
			},
			GoogleNews: GoogleNewsConfig{
				APIKey:     getEnv("GOOGLE_NEWS_API_KEY", ""),
				Host:       getEnv("GOOGLE_NEWS_HOST", "google-news.p.rapidapi.com"),
				MaxResults: getIntEnv("GOOGLE_NEWS_MAX_RESULTS", 100),
				Language:   getEnv("GOOGLE_NEWS_LANGUAGE", "id"),
				Country:    getEnv("GOOGLE_NEWS_COUNTRY", "ID"),
			},
			Instagram: InstagramConfig{
				APIKey:     getEnv("INSTAGRAM_API_KEY", ""),
				Host:       getEnv("INSTAGRAM_HOST", "instagram-bulk-profile-scrapper.p.rapidapi.com"),
				MaxResults: getIntEnv("INSTAGRAM_MAX_RESULTS", 50),
				Timeout:    getIntEnv("INSTAGRAM_TIMEOUT", 30),
			},
			IndonesiaNews: IndonesiaNewsConfig{
				APIKey:     getEnv("INDONESIA_NEWS_API_KEY", ""),
				Host:       getEnv("INDONESIA_NEWS_HOST", "indonesia-news.p.rapidapi.com"),
				MaxResults: getIntEnv("INDONESIA_NEWS_MAX_RESULTS", 100),
				Sources:    getEnv("INDONESIA_NEWS_SOURCES", "tempo,kompas,detik"),
			},
		},
		Logging: LoggingConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "text"),
			Output:     getEnv("LOG_OUTPUT", "stdout"),
			FilePath:   getEnv("LOG_FILE_PATH", "logs/etl.log"),
			MaxSize:    getIntEnv("LOG_MAX_SIZE", 100),
			MaxBackups: getIntEnv("LOG_MAX_BACKUPS", 3),
			MaxAge:     getIntEnv("LOG_MAX_AGE", 7),
		},
	}

	return config, nil
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	switch c.Database.Type {
	case "sqlite":
		return fmt.Sprintf("%s.db", c.Database.Database)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Database.Host, c.Database.Port, c.Database.Username, c.Database.Password,
			c.Database.Database, c.Database.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			c.Database.Username, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Database)
	default:
		return fmt.Sprintf("%s.db", c.Database.Database)
	}
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return getEnv("ENV", "development") == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return getEnv("ENV", "development") == "production"
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
