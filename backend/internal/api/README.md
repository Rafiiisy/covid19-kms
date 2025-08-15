# ETL API Layer

This directory contains the HTTP API layer for the COVID-19 Knowledge Management System ETL pipeline.

## 🚀 Quick Start

### Running the API Server

```bash
# From backend directory
go run cmd/api/main.go

# Or build and run
go build -o bin/api cmd/api/main.go
./bin/api
```

### Default Configuration

- **Server**: `localhost:8080`
- **Environment**: Development (API keys optional)
- **Database**: SQLite (default)

## 📡 API Endpoints

### Root Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | API information and available endpoints |
| `GET` | `/api` | Detailed API documentation |

### ETL Pipeline Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/etl/run` | Run complete ETL pipeline |
| `GET` | `/api/etl/status` | Get pipeline status and API info |
| `POST` | `/api/etl/extract` | Run only data extraction stage |
| `POST` | `/api/etl/transform` | Run only data transformation stage |
| `POST` | `/api/etl/load` | Run only data loading stage |

### Health & Monitoring

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/health` | Health check endpoint |

## 🔧 Configuration

### Environment Variables

Copy `config/env.example` to `.env` and configure:

```bash
# Server
SERVER_PORT=8080
SERVER_HOST=localhost

# Environment
ENV=development

# API Keys (required for production)
YOUTUBE_API_KEY=your_key_here
GOOGLE_NEWS_API_KEY=your_key_here
INSTAGRAM_API_KEY=your_key_here
INDONESIA_NEWS_API_KEY=your_key_here
```

### Configuration Structure

```go
type Config struct {
    Server       ServerConfig       // HTTP server settings
    ETL          ETLConfig          // Pipeline configuration
    API          APIConfig          // API behavior settings
    Database     DatabaseConfig     // Database connection
    ExternalAPIs ExternalAPIsConfig // External API settings
    Logging      LoggingConfig      // Logging configuration
}
```

## 📊 Usage Examples

### 1. Run Complete ETL Pipeline

```bash
curl -X POST http://localhost:8080/api/etl/run
```

**Response:**
```json
{
  "status": "completed",
  "timestamp": "2025-08-15T12:00:00Z",
  "stages": {
    "extraction": "success",
    "transformation": "success",
    "loading": "success"
  },
  "summary": {
    "total_videos": 25,
    "total_articles": 50,
    "total_posts": 15
  }
}
```

### 2. Check API Status

```bash
curl http://localhost:8080/api/etl/status
```

**Response:**
```json
{
  "status": "ready",
  "timestamp": "2025-08-15T12:00:00Z",
  "service": "ETL Pipeline API",
  "version": "1.0.0",
  "endpoints": [
    "/api/etl/run",
    "/api/etl/status",
    "/api/etl/extract",
    "/api/etl/transform",
    "/api/etl/load"
  ]
}
```

### 3. Run Individual Stages

```bash
# Extract data only
curl -X POST http://localhost:8080/api/etl/extract

# Transform data only
curl -X POST http://localhost:8080/api/etl/transform

# Load data only
curl -X POST http://localhost:8080/api/etl/load
```

## 🏗️ Architecture

### Components

1. **Server** (`server.go`)
   - HTTP server management
   - Graceful shutdown
   - Signal handling

2. **Router** (`router.go`)
   - Route configuration
   - Middleware setup
   - CORS handling

3. **ETL Handler** (`etl_handler.go`)
   - HTTP request handling
   - ETL pipeline orchestration
   - Response formatting

4. **Configuration** (`config/`)
   - Environment variable loading
   - Configuration validation
   - Default values

### Middleware

- **CORS**: Cross-origin request handling
- **Logging**: Request/response logging
- **Validation**: Configuration validation

## 🔒 Security

### Development Mode
- API keys optional
- CORS enabled for all origins
- No authentication required

### Production Mode
- API keys required
- CORS restricted
- Authentication recommended (not implemented)

## 🧪 Testing

### Run API Tests

```bash
# From backend directory
go test -v ./internal/api

# From api directory
go test -v .
```

### Test Individual Components

```bash
# Test configuration
go test -v ./internal/config

# Test ETL handler
go test -v ./internal/api -run TestETLHandler
```

## 📝 Logging

### Log Levels

- **DEBUG**: Detailed debugging information
- **INFO**: General information messages
- **WARN**: Warning messages
- **ERROR**: Error messages

### Log Output

- **Development**: Console output
- **Production**: File output with rotation

## 🚨 Error Handling

### HTTP Status Codes

- `200 OK`: Success
- `400 Bad Request`: Invalid request
- `405 Method Not Allowed`: Unsupported HTTP method
- `500 Internal Server Error`: Server error

### Error Response Format

```json
{
  "error": "Error message",
  "timestamp": "2025-08-15T12:00:00Z",
  "details": "Additional error details"
}
```

## 🔄 Development Workflow

1. **Modify API handlers** in `etl_handler.go`
2. **Update routes** in `router.go`
3. **Configure settings** in `config/`
4. **Test locally** with `go run cmd/api/main.go`
5. **Build and deploy** with `go build`

## 📚 Dependencies

- **Standard Library**: `net/http`, `encoding/json`, `log`
- **Internal**: `covid19-kms/internal/etl`, `covid19-kms/internal/config`

## 🎯 Next Steps

- [ ] Add authentication middleware
- [ ] Implement rate limiting
- [ ] Add request/response validation
- [ ] Create OpenAPI/Swagger documentation
- [ ] Add metrics and monitoring
- [ ] Implement caching layer
