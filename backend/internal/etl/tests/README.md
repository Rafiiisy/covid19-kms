# Go ETL Test Suite

This directory contains comprehensive tests for the Go ETL pipeline, mirroring the structure and functionality of the Python tests in `_python/tests/`.

## 🏗️ Structure

```
tests/
├── main.go                    # Main ETL pipeline test runner
├── test_apis/                # Individual API tests
│   ├── test_youtube.go       # YouTube API tests
│   ├── test_google_news.go   # Google News API tests
│   ├── test_instagram.go     # Instagram API tests
│   └── test_indo_news.go     # Indonesia News API tests
├── fixtures/                  # Test data fixtures
├── output/                    # Test output files
└── README.md                  # This file
```

## 🚀 Quick Start

### Run Complete ETL Pipeline Test
```bash
cd backend/internal/etl/tests
go run main.go
```

### Run Individual API Tests
```bash
# YouTube API tests
go run test_apis/test_youtube.go

# Google News API tests
go run test_apis/test_google_news.go

# Instagram API tests
go run test_apis/test_instagram.go

# Indonesia News API tests
go run test_apis/test_indo_news.go
```

## 📊 Test Output

All tests generate timestamped output files in the `tests/output/` directory:

- **Pipeline Results**: `pipeline_results_YYYYMMDD_HHMMSS.json`
- **Extracted Data**: `extracted_data_YYYYMMDD_HHMMSS.json`
- **Transformed Data**: `transformed_data_YYYYMMDD_HHMMSS.json`
- **Load Reports**: `load_report_YYYYMMDD_HHMMSS.json`
- **API Test Results**: `{api}_test_{n}_YYYYMMDD_HHMMSS.json`
- **Latest Files**: `*_latest.json` (symlinks to most recent results)

## 🔧 Test Components

### 1. Main Pipeline Test (`main.go`)
- **Extraction Stage**: Simulates data extraction from all sources
- **Transformation Stage**: Transforms extracted data into normalized format
- **Loading Stage**: Simulates database loading and CSV export
- **Comprehensive Logging**: Emoji-based progress indicators and detailed output

### 2. API Tests (`test_apis/`)
Each API test file includes:
- **Individual Function Tests**: Tests specific API functionality
- **Mock Data Generation**: Realistic test data for validation
- **Error Handling**: Tests both success and failure scenarios
- **Performance Metrics**: Timing and record count validation

### 3. Output Management
- **Timestamped Files**: All outputs include timestamps for tracking
- **Latest Symlinks**: Easy access to most recent results
- **JSON Format**: Structured output for analysis and debugging
- **CSV Export**: Dimension table exports for database loading

## 📈 Test Coverage

### YouTube API Tests
- ✅ Hashtag video extraction
- ✅ Video comments retrieval
- ✅ Full workflow validation
- ✅ Data structure validation

### Google News API Tests
- ✅ COVID-19 news search
- ✅ Indonesian language support
- ✅ Content analysis and metrics
- ✅ Source distribution analysis

### Instagram API Tests
- ✅ Hashtag media extraction
- ✅ Post engagement metrics
- ✅ Content sentiment analysis
- ✅ User interaction data

### Indonesia News API Tests
- ✅ Multi-source news extraction
- ✅ Local news coverage
- ✅ Content relevance scoring
- ✅ Language detection

## 🎯 Key Features

1. **Comprehensive Testing**: Covers all ETL pipeline stages
2. **Mock Data**: Realistic test data without external API calls
3. **Error Handling**: Tests failure scenarios and recovery
4. **Performance Metrics**: Timing and throughput measurements
5. **Structured Output**: JSON format for easy analysis
6. **Visual Feedback**: Emoji-based progress indicators
7. **Timestamped Results**: Track test execution over time

## 🔍 Debugging

### View Test Output
```bash
# Check latest pipeline results
cat tests/output/pipeline_results_latest.json

# View specific API test results
cat tests/output/youtube_test_1_YYYYMMDD_HHMMSS.json
```

### Common Issues
1. **Permission Errors**: Ensure write access to `tests/output/` directory
2. **Missing Dependencies**: Run `go mod tidy` in the backend directory
3. **Path Issues**: Run tests from the correct directory (`backend/internal/etl/tests`)

## 📝 Adding New Tests

### Create New API Test
1. Copy existing test file structure
2. Modify test functions for your API
3. Update mock data and validation logic
4. Add to test runner if needed

### Extend Pipeline Test
1. Add new stage to `main.go`
2. Implement stage logic and validation
3. Update result structures and output
4. Test with `go run main.go`

## 🚀 Next Steps

- [ ] Add Instagram API tests
- [ ] Add Indonesia News API tests
- [ ] Implement CSV export functionality
- [ ] Add database integration tests
- [ ] Create test fixtures for consistent data
- [ ] Add performance benchmarking
- [ ] Implement test result comparison

## 📚 Related Files

- **Python Tests**: `_python/tests/` - Reference implementation
- **Go ETL**: `../` - Main ETL package
- **Output Data**: `tests/output/` - Generated test results
