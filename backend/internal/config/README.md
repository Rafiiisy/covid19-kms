# Configuration Layer

This directory contains the configuration management system for the COVID-19 Knowledge Management System ETL pipeline.

## 🚀 Quick Start

### 1. Copy Environment Template

```bash
# Copy the example environment file
cp env.example .env

# Edit with your actual values
nano .env
```

### 2. Set Required Variables

```bash
# Environment
ENV=development

# API Keys (for production)
YOUTUBE_API_KEY=your_actual_key_here
GOOGLE_NEWS_API_KEY=your_actual_key_here
INSTAGRAM_API_KEY=your_actual_key_here
INDONESIA_NEWS_API_KEY=your_actual_key_here
```

### 3. Load Configuration in Code

```go
import "covid19-kms/internal/config"

// Load configuration
cfg, err := config.LoadConfig()
if err != nil {
    log.Fatal(err)
}

// Use configuration
fmt.Printf("Server will run on %s:%s\n", cfg.Server.Host, cfg.Server.Port)
```

## 🔧 Configuration Structure

### Server Configuration

```go
type ServerConfig struct {
    Port         string        // Server port (e.g., "8080")
    Host         string        // Server host (e.g., "localhost")
    ReadTimeout  time.Duration // Request read timeout
    WriteTimeout time.Duration // Response write timeout
    IdleTimeout  time.Duration // Connection idle timeout
}
```

**Default Values:**
- Port: `8080`
- Host: `localhost`
- ReadTimeout: `30s`
- WriteTimeout: `30s`
- IdleTimeout: `60s`

### ETL Pipeline Configuration

```go
type ETLConfig struct {
    MaxConcurrentExtractions int           // Max concurrent API calls
    ExtractionTimeout        time.Duration // Extraction stage timeout
    TransformationTimeout    time.Duration // Transformation stage timeout
    LoadingTimeout           time.Duration // Loading stage timeout
    BatchSize                int           // Records per batch
    RetryAttempts            int           // Number of retry attempts
    RetryDelay               time.Duration // Delay between retries
}
```

**Default Values:**
- MaxConcurrentExtractions: `5`
- ExtractionTimeout: `5m`
- TransformationTimeout: `2m`
- LoadingTimeout: `3m`
- BatchSize: `100`
- RetryAttempts: `3`
- RetryDelay: `5s`

### API Configuration

```go
type APIConfig struct {
    EnableCORS        bool   // Enable CORS headers
    EnableLogging     bool   // Enable request logging
    EnableMetrics     bool   // Enable metrics collection
    RateLimitRequests int    // Requests per time window
    RateLimitWindow   string // Rate limit time window
}
```

**Default Values:**
- EnableCORS: `true`
- EnableLogging: `true`
- EnableMetrics: `true`
- RateLimitRequests: `100`
- RateLimitWindow: `1m`

### Database Configuration

```go
type DatabaseConfig struct {
    Type      string // Database type: "sqlite", "postgres", "mysql"
    Host      string // Database host
    Port      int    // Database port
    Username  string // Database username
    Password  string // Database password
    Database  string // Database name
    SSLMode   string // SSL mode (for postgres)
    MaxConns  int    // Maximum connections
    IdleConns int    // Idle connections
}
```

**Default Values:**
- Type: `sqlite`
- Host: `localhost`
- Port: `5432`
- Database: `covid19_kms`
- SSLMode: `disable`
- MaxConns: `10`
- IdleConns: `5`

### External APIs Configuration

#### YouTube API

```go
type YouTubeConfig struct {
    APIKey     string // RapidAPI key for YouTube
    Host       string // API host
    MaxResults int    // Maximum results per query
    Timeout    int    // Request timeout in seconds
}
```

**Default Values:**
- Host: `yt-api.p.rapidapi.com`
- MaxResults: `50`
- Timeout: `30`

#### Google News API

```go
type GoogleNewsConfig struct {
    APIKey     string // RapidAPI key for Google News
    Host       string // API host
    MaxResults int    // Maximum results per query
    Language   string // News language
    Country    string // News country
}
```

**Default Values:**
- Host: `google-news.p.rapidapi.com`
- MaxResults: `100`
- Language: `id` (Indonesian)
- Country: `ID` (Indonesia)

#### Instagram API

```go
type InstagramConfig struct {
    APIKey     string // RapidAPI key for Instagram
    Host       string // API host
    MaxResults int    // Maximum results per query
    Timeout    int    // Request timeout in seconds
}
```

**Default Values:**
- Host: `instagram-bulk-profile-scrapper.p.rapidapi.com`
- MaxResults: `50`
- Timeout: `30`

#### Indonesia News API

```go
type IndonesiaNewsConfig struct {
    APIKey     string // RapidAPI key for Indonesia News
    Host       string // API host
    MaxResults int    // Maximum results per query
    Sources    string // News sources (comma-separated)
}
```

**Default Values:**
- Host: `indonesia-news.p.rapidapi.com`
- MaxResults: `100`
- Sources: `tempo,kompas,detik`

### Logging Configuration

```go
type LoggingConfig struct {
    Level      string // Log level: "debug", "info", "warn", "error"
    Format     string // Log format: "json", "text"
    Output     string // Output destination: "stdout", "file"
    FilePath   string // Log file path
    MaxSize    int    // Maximum file size in MB
    MaxBackups int    // Maximum number of backup files
    MaxAge     int    // Maximum age of backup files in days
}
```

**Default Values:**
- Level: `info`
- Format: `text`
- Output: `stdout`
- FilePath: `logs/etl.log`
- MaxSize: `100`
- MaxBackups: `3`
- MaxAge: `7`

## 🌍 Environment Variables

### Server Variables

```bash
SERVER_PORT=8080
SERVER_HOST=localhost
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=60s
```

### ETL Variables

```bash
ETL_MAX_CONCURRENT_EXTRACTIONS=5
ETL_EXTRACTION_TIMEOUT=5m
ETL_TRANSFORMATION_TIMEOUT=2m
ETL_LOADING_TIMEOUT=3m
ETL_BATCH_SIZE=100
ETL_RETRY_ATTEMPTS=3
ETL_RETRY_DELAY=5s
```

### API Variables

```bash
API_ENABLE_CORS=true
API_ENABLE_LOGGING=true
API_ENABLE_METRICS=true
API_RATE_LIMIT_REQUESTS=100
API_RATE_LIMIT_WINDOW=1m
```

### Database Variables

```bash
DB_TYPE=sqlite
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=
DB_PASSWORD=
DB_DATABASE=covid19_kms
DB_SSL_MODE=disable
DB_MAX_CONNECTIONS=10
DB_IDLE_CONNECTIONS=5
```

### External API Variables

```bash
# YouTube
YOUTUBE_API_KEY=your_key_here
YOUTUBE_HOST=yt-api.p.rapidapi.com
YOUTUBE_MAX_RESULTS=50
YOUTUBE_TIMEOUT=30

# Google News
GOOGLE_NEWS_API_KEY=your_key_here
GOOGLE_NEWS_HOST=google-news.p.rapidapi.com
GOOGLE_NEWS_MAX_RESULTS=100
GOOGLE_NEWS_LANGUAGE=id
GOOGLE_NEWS_COUNTRY=ID

# Instagram
INSTAGRAM_API_KEY=your_key_here
INSTAGRAM_HOST=instagram-bulk-profile-scrapper.p.rapidapi.com
INSTAGRAM_MAX_RESULTS=50
INSTAGRAM_TIMEOUT=30

# Indonesia News
INDONESIA_NEWS_API_KEY=your_key_here
INDONESIA_NEWS_HOST=indonesia-news.p.rapidapi.com
INDONESIA_NEWS_MAX_RESULTS=100
INDONESIA_NEWS_SOURCES=tempo,kompas,detik
```

### Logging Variables

```bash
LOG_LEVEL=info
LOG_FORMAT=text
LOG_OUTPUT=stdout
LOG_FILE_PATH=logs/etl.log
LOG_MAX_SIZE=100
LOG_MAX_BACKUPS=3
LOG_MAX_AGE=7
```

## 🔄 Environment Loading

### Automatic Loading

The configuration system automatically tries to load environment variables from:

1. `.env` (current directory)
2. `config/.env` (config subdirectory)
3. `../.env` (parent directory)

### Manual Loading

```go
import "covid19-kms/internal/config"

// Load from specific file
err := config.LoadEnvFile(".env")

// Load if file exists (no error if missing)
err := config.LoadEnvFileIfExists(".env")

// Load from default locations
err := config.LoadDefaultEnv()
```

## ✅ Configuration Validation

### Development Mode

- API keys are optional
- Uses default values for missing configurations
- CORS enabled for all origins

### Production Mode

- API keys are required
- Validates all required environment variables
- Stricter security settings

### Validation Functions

```go
// Check if environment variable is set
value, err := config.GetRequiredEnv("YOUTUBE_API_KEY")
if err != nil {
    log.Fatal("YouTube API key is required")
}

// Validate multiple required variables
requiredKeys := []string{"YOUTUBE_API_KEY", "GOOGLE_NEWS_API_KEY"}
err := config.ValidateRequiredEnvs(requiredKeys)

// Check environment
if config.IsProduction() {
    // Production-specific logic
}
```

## 🧪 Testing

### Run Configuration Tests

```bash
# From backend directory
go test -v ./internal/config

# From config directory
go test -v .
```

### Test Environment Loading

```bash
# Test with sample .env file
cp env.example .env.test
go test -v . -run TestLoadEnvFile
```

## 📁 File Structure

```
backend/internal/config/
├── README.md           # This documentation
├── config.go           # Main configuration structures
├── env.go              # Environment variable loading
└── env.example         # Environment template
```

## 🎯 Best Practices

### 1. Environment-Specific Files

```bash
# Development
.env.development

# Production
.env.production

# Testing
.env.test
```

### 2. Sensitive Information

- Never commit `.env` files to version control
- Use `.env.example` for templates
- Store production secrets in secure systems

### 3. Configuration Hierarchy

1. Environment variables (highest priority)
2. `.env` file
3. Default values (lowest priority)

### 4. Validation

- Always validate configuration on startup
- Provide clear error messages for missing values
- Use appropriate defaults for development

## 🚨 Troubleshooting

### Common Issues

1. **Configuration not loaded**
   - Check if `.env` file exists
   - Verify file permissions
   - Check for syntax errors in `.env`

2. **Missing environment variables**
   - Copy from `env.example`
   - Set required variables
   - Check environment mode

3. **Database connection issues**
   - Verify database type
   - Check connection parameters
   - Ensure database is running

### Debug Mode

```bash
# Enable debug logging
LOG_LEVEL=debug

# Check loaded configuration
curl http://localhost:8080/api/etl/status
```

## 🔮 Future Enhancements

- [ ] Configuration hot-reloading
- [ ] Configuration encryption
- [ ] Configuration validation schemas
- [ ] Configuration backup/restore
- [ ] Configuration monitoring
- [ ] Configuration templates for different environments
