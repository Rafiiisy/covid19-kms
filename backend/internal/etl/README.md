# Go ETL Package for COVID-19 Knowledge Management System

This package provides a complete ETL (Extract, Transform, Load) pipeline implementation in Go, ported from the original Python implementation.

## 🏗️ **Package Structure**

```
backend/internal/etl/
├── etl.go              # Main package documentation and exports
├── extractors.go       # Main data extraction orchestrator
├── youtube.go          # YouTube API client
├── google_news.go      # Google News API client
├── instagram.go        # Instagram API client
├── indo_news.go        # Indonesia News API client
├── transformers.go     # Data transformation and cleaning
├── loaders.go          # Data loading to local storage
├── orchestrator.go     # Main ETL pipeline coordinator
├── etl_test.go         # Unit tests
└── README.md           # This file
```

## 🚀 **Key Features**

### **1. Concurrent Data Extraction**
- **YouTube API**: Extract videos and comments using RapidAPI
- **Google News API**: Search for COVID-19 related news articles
- **Instagram API**: Extract posts and media with hashtag filtering
- **Indonesia News API**: Multi-source Indonesian news extraction
- **Goroutines**: All extractions run concurrently for optimal performance

### **2. Data Transformation**
- **Text Cleaning**: Remove special characters and normalize whitespace
- **COVID-19 Relevance Scoring**: Calculate relevance based on keywords
- **Language Detection**: Simple Indonesian/English detection
- **Data Enrichment**: Add metadata and processing timestamps

### **3. Data Loading**
- **Local Storage**: Load transformed data to local file system
- **JSON Format**: Store data in structured JSON files
- **Timestamped Files**: Organize data with timestamps
- **Error Handling**: Graceful fallbacks and comprehensive error reporting

### **4. Pipeline Orchestration**
- **End-to-End Pipeline**: Complete ETL workflow coordination
- **Progress Tracking**: Real-time pipeline status and metrics
- **Error Recovery**: Robust error handling and reporting
- **Performance Metrics**: Pipeline duration and record counts

## 📊 **Data Flow**

```
API Sources → Extractors → Transformer → Loader → Destinations
     ↓           ↓           ↓          ↓         ↓
  YouTube    Concurrent   Clean &    Local      Local
  News       Processing   Enrich     Storage   Storage
  Instagram
  Indonesia
```

## 🔧 **Usage Examples**

### **Basic ETL Pipeline Execution**

```go
package main

import (
    "log"
    "your-project/internal/etl"
)

func main() {
    // Create ETL orchestrator
    orchestrator := etl.NewETLOrchestrator()
    
    // Run the complete pipeline
    result := orchestrator.RunETLPipeline()
    
    if result.Status == "success" {
        log.Printf("✅ Pipeline completed in %s", result.PipelineDuration)
        log.Printf("📊 Extracted from %d sources", len(result.Extraction.Sources))
        log.Printf("🔄 Transformed %d videos and %d articles", 
            len(result.Transformation.YouTube), 
            len(result.Transformation.News))
    } else {
        log.Printf("❌ Pipeline failed: %s", result.Error)
    }
}
```

### **Individual Component Usage**

```go
// Data Extraction
extractor := etl.NewDataExtractor()
extractedData := extractor.ExtractAllSources()

// Data Transformation
transformer := etl.NewDataTransformer()
transformedData := transformer.TransformData(youtubeData, newsData)

// Data Loading
loader := etl.NewDataLoader()
loadResult := loader.LoadData(transformedData)
```

## 🌐 **API Integration**

### **YouTube API**
```go
youtubeAPI := etl.NewYouTubeAPI()
videos, err := youtubeAPI.GetHashtagVideos("covid19", "")
comments, err := youtubeAPI.GetVideoComments("video_id")
```

### **Google News API**
```go
newsAPI := etl.NewGoogleNewsAPI()
articles, err := newsAPI.SearchNews("COVID-19", "id", "id-ID")
```

### **Instagram API**
```go
instagramAPI := etl.NewInstagramAPI()
posts, err := instagramAPI.GetHashtagMedia("covid19", "")
```

### **Indonesia News API**
```go
indoNewsAPI := etl.NewIndonesiaNewsAPI()
tempoNews, err := indoNewsAPI.SearchNews("tempo", "COVID-19", nil)
kompasNews, err := indoNewsAPI.SearchNews("kompas", "COVID-19", map[string]interface{}{
    "page": 1,
    "limit": 10,
})
```

## ⚙️ **Configuration**

### **Environment Variables**

```bash
# RapidAPI Configuration
RAPIDAPI_KEY=your_rapidapi_key_here

# Local Storage Configuration
OUTPUT_DIR=data
```

### **Default Values**
- **Output Directory**: `data`
- **Raw Data**: `data/raw/`
- **Processed Data**: `data/processed/`

## 🧪 **Testing**

Run the test suite:

```bash
cd backend/internal/etl
go test -v
```

Test coverage includes:
- Component initialization
- Text cleaning and transformation
- COVID-19 relevance scoring
- Language detection
- DateTime parsing
- Summary generation
- Load reporting
- Pipeline metrics

## 📈 **Performance Benefits**

### **Go vs Python Implementation**

1. **Concurrency**: Native goroutines for parallel API calls
2. **Memory Efficiency**: Lower memory footprint and better garbage collection
3. **Compilation**: Single binary deployment, no runtime dependencies
4. **Type Safety**: Compile-time error checking with Go's type system
5. **HTTP Client**: Optimized HTTP client with connection pooling

### **Expected Improvements**
- **Extraction Speed**: 3-5x faster due to concurrent processing
- **Memory Usage**: 2-3x lower memory consumption
- **Deployment**: Single binary, easier containerization
- **Maintenance**: Strong typing reduces runtime errors

## 🔄 **Migration Status**

### **✅ Completed**
- [x] Complete ETL package structure
- [x] All API extractors (YouTube, Google News, Instagram, Indonesia News)
- [x] Data transformation and cleaning
- [x] Data loading with local storage
- [x] Pipeline orchestration
- [x] Comprehensive error handling
- [x] Unit tests
- [x] Documentation

### **🚧 Next Steps**
1. **Integration**: Connect ETL package to main Go server
2. **Real API Testing**: Test with actual RapidAPI credentials
3. **Local Storage Setup**: Configure local file storage directories
4. **Frontend Integration**: Connect React frontend to Go backend
5. **Performance Tuning**: Optimize based on real-world usage

## 📚 **Dependencies**

### **Go Modules**
- `net/http`: HTTP client for API calls
- `encoding/json`: JSON marshaling/unmarshaling
- `time`: Time handling and formatting
- `os`: Environment variable access
- `log`: Logging functionality

### **External Dependencies**
- **RapidAPI**: All data source APIs
- **Local File System**: JSON file storage

## 🎯 **Architecture Principles**

1. **Separation of Concerns**: Each component has a single responsibility
2. **Error Handling**: Comprehensive error handling with graceful fallbacks
3. **Concurrency**: Leverage Go's goroutines for parallel processing
4. **Configuration**: Environment-based configuration with sensible defaults
5. **Testing**: Comprehensive unit tests for all components
6. **Documentation**: Clear documentation and usage examples

---

**Note**: This package is ready for integration with your main Go server. The next step is to import it into your main application and replace the placeholder ETL endpoint with actual pipeline execution.
