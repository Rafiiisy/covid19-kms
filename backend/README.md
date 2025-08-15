# COVID-19 Knowledge Management System - Backend

This directory contains the Go-based backend for the COVID-19 Knowledge Management System, implementing a comprehensive ETL (Extract, Transform, Load) pipeline for processing COVID-19 related data from multiple sources.

## 🏗️ **Backend Architecture**

### **High-Level Architecture**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend API   │    │   ETL Pipeline  │
│   (React/TS)    │◄──►│   (HTTP/REST)   │◄──►│   (Go)         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Configuration │    │   Data Sources  │
                       │   Management    │    │   (APIs)        │
                       └─────────────────┘    └─────────────────┘
```

### **Component Architecture**
```
backend/
├── cmd/                    # Executable entry points
│   └── api/              # Main API server executable
├── internal/              # Private application code
│   ├── api/              # HTTP API layer (handlers, routing, server)
│   ├── config/           # Configuration management
│   └── etl/              # ETL pipeline core
│       ├── tests/        # ETL test suite
│       └── data/         # Data storage
└── go.mod                 # Go module definition
```

## 🔄 **Data Flow Architecture**

### **1. ETL Pipeline Flow**
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Extract   │───►│ Transform  │───►│   Load      │───►│   Store     │
│   (APIs)    │    │ (Clean/    │    │ (Database)  │    │ (SQLite/    │
│             │    │  Enrich)   │    │             │    │  CSV)       │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

### **2. API Request Flow**
```
Client Request → Router → Handler → ETL Pipeline → Response
     ↓              ↓        ↓         ↓           ↓
  HTTP/JSON    Route Match  Business  Data       JSON
  (curl/      (endpoint)    Logic     Processing Response
   browser)                                    (Status/Data)
```

### **3. Configuration Flow**
```
Environment Variables → .env File → Config Loader → Application
     ↓                    ↓           ↓            ↓
  System/User         Local File   Validation   Runtime Use
  Settings            Override     & Defaults   (Server/ETL)
```

## 📁 **Directory Structure & Functions**

### **`cmd/api/` - Main Executable**
```
cmd/api/
└── main.go                 # Application entry point
```

**Functions:**
- **`main()`**: Application startup and configuration validation
- **`validateConfiguration()`**: Ensures required settings are present
- **`getEnvironmentString()`**: Returns human-readable environment info

### **`internal/api/` - HTTP API Layer**
```
internal/api/
├── etl_handler.go         # ETL operation HTTP handlers
├── router.go              # API routing and middleware
├── server.go              # HTTP server management
└── README.md              # API documentation
```

**Functions:**

#### **`etl_handler.go`**
- **`NewETLHandler()`**: Creates new ETL handler instance
- **`RunETLPipeline()`**: Handles complete ETL pipeline execution
- **`GetPipelineStatus()`**: Returns API status and endpoint information
- **`ExtractData()`**: Runs data extraction stage only
- **`TransformData()`**: Runs data transformation stage only
- **`LoadData()`**: Runs data loading stage only
- **`HealthCheck()`**: Health monitoring endpoint

#### **`router.go`**
- **`NewRouter()`**: Creates new router instance
- **`SetupRoutes()`**: Configures all API endpoints
- **`handleRoot()`**: Root endpoint with API information
- **`handleAPIInfo()`**: Detailed API documentation
- **`enableCORS()`**: Cross-origin request handling
- **`loggingMiddleware()`**: Request/response logging

#### **`server.go`**
- **`NewServer()`**: Creates new HTTP server instance
- **`Start()`**: Starts the HTTP server with graceful shutdown
- **`Stop()`**: Gracefully stops the server
- **`RunServer()`**: Convenience function for default configuration
- **`StartServerWithConfig()`**: Starts server with specific configuration

### **`internal/config/` - Configuration Management**
```
internal/config/
├── config.go              # Configuration structures and loading
├── env.go                 # Environment variable management
├── env.example            # Environment template
└── README.md              # Configuration documentation
```

**Functions:**

#### **`config.go`**
- **`LoadConfig()`**: Loads configuration from environment variables
- **`GetDatabaseDSN()`**: Returns database connection string
- **`IsDevelopment()`**: Checks if running in development mode
- **`IsProduction()`**: Checks if running in production mode

#### **`env.go`**
- **`LoadEnvFile()`**: Loads environment variables from .env file
- **`LoadDefaultEnv()`**: Loads from common .env locations
- **`ValidateRequiredEnvs()`**: Validates required environment variables
- **`GetRequiredEnv()`**: Gets required environment variable

### **`internal/etl/` - ETL Pipeline Core**
```
internal/etl/
├── etl.go                 # Main ETL package entry point
├── extractors.go          # Data extraction logic
├── transformers.go        # Data transformation logic
├── loaders.go             # Data loading logic
├── orchestrator.go        # ETL pipeline orchestration
├── etl_test.go            # ETL unit tests
├── tests/                 # ETL test suite
│   ├── main.go           # Main ETL pipeline test
│   ├── test_apis/        # API-specific tests
│   │   ├── test_youtube.go      # YouTube API tests
│   │   └── test_google_news.go  # Google News API tests
│   └── README.md         # Test suite documentation
├── data/                  # Data storage
│   └── processed/        # Processed data files
└── [api-specific files]  # YouTube, Google News, Instagram, Indonesia News
```

**Functions:**

#### **`etl.go`**
- **Package documentation**: Main ETL package overview
- **Re-exports**: Makes ETL components available

#### **`extractors.go`**
- **`NewDataExtractor()`**: Creates new data extractor
- **`ExtractAllSources()`**: Extracts data from all configured sources
- **`extractFromYouTube()`**: YouTube data extraction
- **`extractFromGoogleNews()`**: Google News data extraction
- **`extractFromInstagram()`**: Instagram data extraction
- **`extractFromIndonesiaNews()`**: Indonesia News data extraction

#### **`transformers.go`**
- **`NewDataTransformer()`**: Creates new data transformer
- **`TransformData()`**: Transforms extracted data
- **`cleanText()`**: Cleans and normalizes text data
- **`calculateCovidRelevance()`**: Calculates COVID-19 relevance scores
- **`detectLanguage()`**: Detects text language
- **`parseDateTime()`**: Parses and standardizes dates
- **`createSummary()`**: Creates data summaries

#### **`loaders.go`**
- **`NewDataLoader()`**: Creates new data loader
- **`LoadData()`**: Loads data to local storage
- **`LoadRawData()`**: Loads raw data to local storage
- **`saveLocally()`**: Saves data to local files
- **`GetLoadReport()`**: Returns loading operation report

#### **`orchestrator.go`**
- **`NewETLOrchestrator()`**: Creates new ETL orchestrator
- **`RunETLPipeline()`**: Runs complete ETL pipeline
- **`extractData()`**: Orchestrates data extraction
- **`transformData()`**: Orchestrates data transformation
- **`loadData()`**: Orchestrates data loading
- **`GetPipelineMetrics()`**: Returns pipeline execution metrics

## 🚀 **Setup & Installation**

### **Prerequisites**
- **Go 1.21+**: [Download Go](https://golang.org/dl/)
- **Git**: [Download Git](https://git-scm.com/)
- **API Keys**: RapidAPI keys for external services (optional for development)

### **1. Clone and Navigate**
```bash
# Clone the repository
git clone <repository-url>
cd RepoCloud/backend

# Verify Go installation
go version
```

### **2. Install Dependencies**
```bash
# Download Go modules
go mod download

# Verify dependencies
go mod tidy
```

### **3. Environment Configuration**
```bash
# Copy environment template
cp internal/config/env.example .env

# Edit environment variables
nano .env  # or use your preferred editor
```

**Required Environment Variables:**
```bash
# Environment
ENV=development

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# API Keys (required for production)
YOUTUBE_API_KEY=your_youtube_api_key_here
GOOGLE_NEWS_API_KEY=your_google_news_api_key_here
INSTAGRAM_API_KEY=your_instagram_api_key_here
INDONESIA_NEWS_API_KEY=your_indonesia_news_api_key_here
```

### **4. Build the Application**
```bash
# Build the API server
go build -o api cmd/api/main.go

# Verify the binary was created
ls -la api
```

## 🏃 **Running the Backend**

### **1. Run the API Server**
```bash
# Method 1: Direct execution
go run cmd/api/main.go

# Method 2: Build and run
go build -o api cmd/api/main.go
./api

# Method 3: Run from specific directory
cd cmd/api
go run main.go
```

### **2. Verify the Server**
```bash
# Check if server is running
curl http://localhost:8080/api/health

# Get API documentation
curl http://localhost:8080/api

# Check pipeline status
curl http://localhost:8080/api/etl/status
```

### **3. Test ETL Pipeline**
```bash
# Run complete ETL pipeline
curl -X POST http://localhost:8080/api/etl/run

# Run individual stages
curl -X POST http://localhost:8080/api/etl/extract
curl -X POST http://localhost:8080/api/etl/transform
curl -X POST http://localhost:8080/api/etl/load
```

## 🧪 **Testing**

### **Run All Tests**
```bash
# Test entire backend
go test -v ./...

# Test specific packages
go test -v ./internal/api
go test -v ./internal/config
go test -v ./internal/etl
```

### **Run ETL Test Suite**
```bash
# Navigate to ETL tests
cd internal/etl/tests

# Run main ETL pipeline test
go run main.go

# Run individual API tests
go run test_apis/test_youtube.go
go run test_apis/test_google_news.go
```

### **Run Unit Tests**
```bash
# Run ETL unit tests
go test -v ./internal/etl

# Run with coverage
go test -v -cover ./internal/etl
```

## 🔧 **Development Commands**

### **Code Quality**
```bash
# Format code
go fmt ./...

# Run linter (if installed)
golangci-lint run

# Check for race conditions
go test -race ./...
```

### **Dependency Management**
```bash
# Add new dependency
go get github.com/example/package

# Update dependencies
go get -u ./...

# Clean module cache
go clean -modcache
```

### **Build Variants**
```bash
# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o api-linux cmd/api/main.go
GOOS=windows GOARCH=amd64 go build -o api-windows.exe cmd/api/main.go
GOOS=darwin GOARCH=amd64 go build -o api-macos cmd/api/main.go

# Build with debug info
go build -gcflags="-N -l" -o api-debug cmd/api/main.go
```

## 📊 **API Endpoints**

### **Core Endpoints**
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | API information and available endpoints |
| `GET` | `/api` | Detailed API documentation |
| `GET` | `/api/health` | Health check endpoint |

### **ETL Pipeline Endpoints**
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/etl/run` | Run complete ETL pipeline |
| `GET` | `/api/etl/status` | Get pipeline status and API info |
| `POST` | `/api/etl/extract` | Run only data extraction stage |
| `POST` | `/api/etl/transform` | Run only data transformation stage |
| `POST` | `/api/etl/load` | Run only data loading stage |

## 🗄️ **Data Storage**

### **Local Storage**
- **SQLite Database**: `covid_knowledge_warehouse.db`
- **JSON Files**: Processed data in `internal/etl/data/processed/`
- **CSV Exports**: Star schema exports for analysis

### **Data Flow**
```
Raw APIs → Extract → Transform → Load → Storage
   ↓         ↓         ↓        ↓       ↓
YouTube   YouTube   Cleaned   SQLite   CSV
News      News      Enriched  JSON     Analysis
Social    Social    Structured
```

## 🔒 **Security & Configuration**

### **Development Mode**
- API keys optional
- CORS enabled for all origins
- No authentication required
- Debug logging enabled

### **Production Mode**
- API keys required
- CORS restricted
- Authentication recommended
- Structured logging

## 🚨 **Troubleshooting**

### **Common Issues**

#### **1. Port Already in Use**
```bash
# Check what's using port 8080
netstat -ano | findstr :8080  # Windows
lsof -i :8080                 # Linux/Mac

# Kill the process or change port in .env
SERVER_PORT=8081
```

#### **2. Configuration Not Loaded**
```bash
# Check if .env file exists
ls -la .env

# Verify environment variables
echo $SERVER_PORT
echo $ENV
```

#### **3. Go Module Issues**
```bash
# Clean and reinstall modules
go clean -modcache
go mod download
go mod tidy
```

#### **4. Build Errors**
```bash
# Check Go version
go version

# Verify module
go mod verify

# Clean build cache
go clean -cache
```

### **Debug Mode**
```bash
# Enable debug logging
LOG_LEVEL=debug

# Check server logs
tail -f logs/etl.log
```

## 📚 **Additional Resources**

### **Documentation**
- [API Documentation](./internal/api/README.md)
- [Configuration Guide](./internal/config/README.md)
- [ETL Test Suite](./internal/etl/tests/README.md)

### **External Dependencies**
- [Go Documentation](https://golang.org/doc/)
- [RapidAPI Documentation](https://rapidapi.com/docs)
- [SQLite Documentation](https://www.sqlite.org/docs.html)

## 🎯 **Next Steps**

### **Immediate Improvements**
- [ ] Add authentication middleware
- [ ] Implement rate limiting
- [ ] Add request/response validation
- [ ] Create OpenAPI/Swagger documentation

### **Future Enhancements**
- [ ] Add metrics and monitoring
- [ ] Implement caching layer
- [ ] Add database migrations
- [ ] Create deployment scripts
- [ ] Add CI/CD pipeline

## 🤝 **Contributing**

### **Development Workflow**
1. Create feature branch
2. Make changes
3. Run tests: `go test -v ./...`
4. Build: `go build ./cmd/api`
5. Submit pull request

### **Code Standards**
- Follow Go formatting: `go fmt`
- Add tests for new functionality
- Update documentation
- Use meaningful commit messages

---

**Backend Status**: ✅ **Ready for Development**
**Last Updated**: August 2025
**Go Version**: 1.21+
**Architecture**: ETL Pipeline with REST API
