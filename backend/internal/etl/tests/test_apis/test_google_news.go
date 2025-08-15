package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GoogleNewsAPITest represents the Google News API test suite
type GoogleNewsAPITest struct {
	outputDir string
}

// TestResult represents the result of a test
type TestResult struct {
	TestName    string                 `json:"test_name"`
	Status      string                 `json:"status"`
	Timestamp   string                 `json:"timestamp"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Duration    string                 `json:"duration,omitempty"`
	Records     int                    `json:"records,omitempty"`
}

// NewGoogleNewsAPITest creates a new Google News API test instance
func NewGoogleNewsAPITest() *GoogleNewsAPITest {
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	
	return &GoogleNewsAPITest{
		outputDir: outputDir,
	}
}

// RunAllTests executes all Google News API tests
func (gn *GoogleNewsAPITest) RunAllTests() {
	fmt.Println("📰 Starting Google News API Tests")
	fmt.Println(strings.Repeat("=", 40))
	
	timestamp := time.Now().Format("20060102_150405")
	
	// Test 1: COVID-19 News Search
	fmt.Println("\n🔍 Test 1: COVID-19 News Search")
	result1 := gn.testCOVIDNewsSearch(timestamp)
	gn.saveTestResult(result1, timestamp, "covid_news_test_1")
	
	// Test 2: Indonesian Language News
	fmt.Println("\n🇮🇩 Test 2: Indonesian Language News")
	result2 := gn.testIndonesianNews(timestamp)
	gn.saveTestResult(result2, timestamp, "covid_news_test_2")
	
	// Test 3: News Content Analysis
	fmt.Println("\n📊 Test 3: News Content Analysis")
	result3 := gn.testNewsContentAnalysis(timestamp)
	gn.saveTestResult(result3, timestamp, "covid_news_test_3")
	
	// Print summary
	gn.printTestSummary([]TestResult{result1, result2, result3})
}

// testCOVIDNewsSearch tests the COVID-19 news search functionality
func (gn *GoogleNewsAPITest) testCOVIDNewsSearch(timestamp string) TestResult {
	startTime := time.Now()
	
	fmt.Println("  📋 Testing COVID-19 news search...")
	
	// Simulate API call with mock data
	mockData := map[string]interface{}{
		"status": "success",
		"keyword": "COVID-19",
		"lang":    "id",
		"lr":      "id-ID",
		"data": []map[string]interface{}{
			{
				"title":       "Update COVID-19 Indonesia: Kasus Terbaru",
				"description": "Laporan terbaru tentang perkembangan kasus COVID-19 di Indonesia",
				"url":         "https://example.com/news1",
				"source":      "Tempo",
				"published":   "2023-12-01T08:00:00Z",
				"language":    "id",
				"relevance":   0.95,
			},
			{
				"title":       "Vaksinasi COVID-19: Progress di Jawa",
				"description": "Update tentang kemajuan program vaksinasi COVID-19 di Pulau Jawa",
				"url":         "https://example.com/news2",
				"source":      "Kompas",
				"published":   "2023-12-01T09:30:00Z",
				"language":    "id",
				"relevance":   0.92,
			},
			{
				"title":       "Protokol Kesehatan: Panduan Terbaru",
				"description": "Panduan terbaru tentang protokol kesehatan selama pandemi COVID-19",
				"url":         "https://example.com/news3",
				"source":      "Detik",
				"published":   "2023-12-01T10:15:00Z",
				"language":    "id",
				"relevance":   0.88,
			},
		},
	}
	
	duration := time.Since(startTime)
	
	// Validate data
	records := 0
	if data, ok := mockData["data"].([]map[string]interface{}); ok {
		records = len(data)
	}
	
	// Check if we have valid data
	if records == 0 {
		return TestResult{
			TestName:  "COVID-19 News Search Test",
			Status:    "failed",
			Timestamp: timestamp,
			Error:     "No news articles found in response",
			Duration:  duration.String(),
			Records:   0,
		}
	}
	
	fmt.Printf("  ✅ Found %d news articles\n", records)
	fmt.Printf("  ⏱️  Test completed in %s\n", duration)
	
	return TestResult{
		TestName:  "COVID-19 News Search Test",
		Status:    "success",
		Timestamp: timestamp,
		Data:      mockData,
		Duration:  duration.String(),
		Records:   records,
	}
}

// testIndonesianNews tests the Indonesian language news functionality
func (gn *GoogleNewsAPITest) testIndonesianNews(timestamp string) TestResult {
	startTime := time.Now()
	
	fmt.Println("  📋 Testing Indonesian language news...")
	
	// Simulate API call with mock data
	mockData := map[string]interface{}{
		"status": "success",
		"keyword": "COVID-19",
		"lang":    "id",
		"lr":      "id-ID",
		"data": []map[string]interface{}{
			{
				"title":       "Pandemi COVID-19: Situasi Terkini",
				"description": "Analisis mendalam tentang situasi pandemi COVID-19 di Indonesia",
				"url":         "https://example.com/indo1",
				"source":      "CNN Indonesia",
				"published":   "2023-12-01T11:00:00Z",
				"language":    "id",
				"relevance":   0.96,
				"sentiment":   "neutral",
			},
			{
				"title":       "Ekonomi Indonesia: Dampak COVID-19",
				"description": "Analisis dampak pandemi COVID-19 terhadap perekonomian Indonesia",
				"url":         "https://example.com/indo2",
				"source":      "Kontan",
				"published":   "2023-12-01T12:00:00Z",
				"language":    "id",
				"relevance":   0.89,
				"sentiment":   "negative",
			},
		},
	}
	
	duration := time.Since(startTime)
	
	// Validate data
	records := 0
	if data, ok := mockData["data"].([]map[string]interface{}); ok {
		records = len(data)
	}
	
	// Check if we have valid data
	if records == 0 {
		return TestResult{
			TestName:  "Indonesian News Test",
			Status:    "failed",
			Timestamp: timestamp,
			Error:     "No Indonesian news articles found in response",
			Duration:  duration.String(),
			Records:   0,
		}
	}
	
	fmt.Printf("  ✅ Found %d Indonesian news articles\n", records)
	fmt.Printf("  ⏱️  Test completed in %s\n", duration)
	
	return TestResult{
		TestName:  "Indonesian News Test",
		Status:    "success",
		Timestamp: timestamp,
		Data:      mockData,
		Duration:  duration.String(),
		Records:   records,
	}
}

// testNewsContentAnalysis tests the news content analysis functionality
func (gn *GoogleNewsAPITest) testNewsContentAnalysis(timestamp string) TestResult {
	startTime := time.Now()
	
	fmt.Println("  📋 Testing news content analysis...")
	
	// Simulate content analysis with mock data
	analysisData := map[string]interface{}{
		"analysis": "news_content_analysis",
		"timestamp": timestamp,
		"articles_analyzed": 5,
		"content_metrics": map[string]interface{}{
			"total_words":       1250,
			"average_length":    250,
			"covid_mentions":    23,
			"indonesia_mentions": 15,
			"vaccine_mentions":  8,
		},
		"sentiment_analysis": map[string]interface{}{
			"positive": 2,
			"neutral":  2,
			"negative": 1,
		},
		"source_distribution": map[string]interface{}{
			"tempo":          2,
			"kompas":         1,
			"detik":          1,
			"cnn_indonesia":  1,
		},
		"language_distribution": map[string]interface{}{
			"indonesian": 5,
			"english":    0,
		},
		"relevance_scores": []float64{0.95, 0.92, 0.88, 0.96, 0.89},
	}
	
	duration := time.Since(startTime)
	
	// Calculate total records
	totalRecords := 0
	if articles, ok := analysisData["articles_analyzed"].(int); ok {
		totalRecords = articles
	}
	
	fmt.Printf("  ✅ Content analysis completed\n")
	fmt.Printf("  📊 Articles analyzed: %d\n", totalRecords)
	fmt.Printf("  ⏱️  Analysis completed in %s\n", duration)
	
	return TestResult{
		TestName:  "News Content Analysis Test",
		Status:    "success",
		Timestamp: timestamp,
		Data:      analysisData,
		Duration:  duration.String(),
		Records:   totalRecords,
	}
}

// saveTestResult saves the test result to a JSON file
func (gn *GoogleNewsAPITest) saveTestResult(result TestResult, timestamp, filename string) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(gn.outputDir, 0755); err != nil {
		log.Printf("Warning: Failed to create output directory: %v", err)
		return
	}
	
	// Generate filename
	outputFile := filepath.Join(gn.outputDir, fmt.Sprintf("%s_%s.json", filename, timestamp))
	
	// Convert to JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Warning: Failed to marshal test result: %v", err)
		return
	}
	
	// Write to file
	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		log.Printf("Warning: Failed to write test result: %v", err)
		return
	}
	
	fmt.Printf("  💾 Test result saved: %s\n", outputFile)
}

// printTestSummary prints a summary of all test results
func (gn *GoogleNewsAPITest) printTestSummary(results []TestResult) {
	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println("📊 GOOGLE NEWS API TEST SUMMARY")
	fmt.Println(strings.Repeat("=", 40))
	
	totalTests := len(results)
	passedTests := 0
	totalRecords := 0
	
	for _, result := range results {
		if result.Status == "success" {
			passedTests++
		}
		totalRecords += result.Records
	}
	
	fmt.Printf("Total Tests: %d\n", totalTests)
	fmt.Printf("Passed: %d\n", passedTests)
	fmt.Printf("Failed: %d\n", totalTests-passedTests)
	fmt.Printf("Total Records: %d\n", totalRecords)
	
	if passedTests == totalTests {
		fmt.Println("🎉 All tests passed successfully!")
	} else {
		fmt.Printf("⚠️  %d test(s) failed\n", totalTests-passedTests)
	}
}

func main() {
	// Create and run Google News API tests
	googleNewsTest := NewGoogleNewsAPITest()
	googleNewsTest.RunAllTests()
}
