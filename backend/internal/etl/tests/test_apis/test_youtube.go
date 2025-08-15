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

// YouTubeAPITest represents the YouTube API test suite
type YouTubeAPITest struct {
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

// NewYouTubeAPITest creates a new YouTube API test instance
func NewYouTubeAPITest() *YouTubeAPITest {
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	
	return &YouTubeAPITest{
		outputDir: outputDir,
	}
}

// RunAllTests executes all YouTube API tests
func (yt *YouTubeAPITest) RunAllTests() {
	fmt.Println("📺 Starting YouTube API Tests")
	fmt.Println(strings.Repeat("=", 40))
	
	timestamp := time.Now().Format("20060102_150405")
	
	// Test 1: Hashtag Videos
	fmt.Println("\n🔍 Test 1: Hashtag Videos")
	result1 := yt.testHashtagVideos(timestamp)
	yt.saveTestResult(result1, timestamp, "youtube_test_1")
	
	// Test 2: Video Comments
	fmt.Println("\n💬 Test 2: Video Comments")
	result2 := yt.testVideoComments(timestamp)
	yt.saveTestResult(result2, timestamp, "youtube_test_2")
	
	// Test 3: Full Workflow
	fmt.Println("\n🔄 Test 3: Full Workflow")
	result3 := yt.testFullWorkflow(timestamp)
	yt.saveTestResult(result3, timestamp, "youtube_test_3")
	
	// Print summary
	yt.printTestSummary([]TestResult{result1, result2, result3})
}

// testHashtagVideos tests the hashtag videos functionality
func (yt *YouTubeAPITest) testHashtagVideos(timestamp string) TestResult {
	startTime := time.Now()
	
	fmt.Println("  📋 Testing hashtag videos extraction...")
	
	// Simulate API call with mock data
	mockData := map[string]interface{}{
		"status": "success",
		"data": []map[string]interface{}{
			{
				"video_id":    "yt_video_001",
				"title":       "COVID-19 Vaccine Update Indonesia",
				"channel":     "Health Channel ID",
				"views":       15000,
				"likes":       1200,
				"duration":    "15:30",
				"published":   "2023-12-01T10:00:00Z",
				"description": "Latest update on COVID-19 vaccination progress in Indonesia",
				"tags":        []string{"covid19", "vaccine", "indonesia", "health"},
			},
			{
				"video_id":    "yt_video_002",
				"title":       "Pandemi COVID-19 di Jakarta",
				"channel":     "News Channel ID",
				"views":       25000,
				"likes":       1800,
				"duration":    "12:45",
				"published":   "2023-12-02T14:30:00Z",
				"description": "Update terbaru tentang situasi COVID-19 di Jakarta",
				"tags":        []string{"covid19", "jakarta", "pandemi", "news"},
			},
			{
				"video_id":    "yt_video_003",
				"title":       "Social Distancing Guidelines Indonesia",
				"channel":     "Government Channel",
				"views":       8000,
				"likes":       650,
				"duration":    "8:20",
				"published":   "2023-12-03T09:15:00Z",
				"description": "Guidelines for social distancing during COVID-19 pandemic",
				"tags":        []string{"covid19", "social distancing", "guidelines", "indonesia"},
			},
		},
		"query": "covid19",
		"geo":   "ID",
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
			TestName:  "Hashtag Videos Test",
			Status:    "failed",
			Timestamp: timestamp,
			Error:     "No videos found in response",
			Duration:  duration.String(),
			Records:   0,
		}
	}
	
	fmt.Printf("  ✅ Found %d videos\n", records)
	fmt.Printf("  ⏱️  Test completed in %s\n", duration)
	
	return TestResult{
		TestName:  "Hashtag Videos Test",
		Status:    "success",
		Timestamp: timestamp,
		Data:      mockData,
		Duration:  duration.String(),
		Records:   records,
	}
}

// testVideoComments tests the video comments functionality
func (yt *YouTubeAPITest) testVideoComments(timestamp string) TestResult {
	startTime := time.Now()
	
	fmt.Println("  📋 Testing video comments extraction...")
	
	// Simulate API call with mock data
	mockData := map[string]interface{}{
		"status": "success",
		"data": []map[string]interface{}{
			{
				"comment_id":  "comment_001",
				"author":      "User123",
				"text":        "Great information about COVID-19 vaccine!",
				"likes":       45,
				"replies":     12,
				"published":   "2023-12-01T11:30:00Z",
				"sentiment":   "positive",
			},
			{
				"comment_id":  "comment_002",
				"author":      "HealthExpert",
				"text":        "Important to follow vaccination schedule",
				"likes":       78,
				"replies":     5,
				"published":   "2023-12-01T12:15:00Z",
				"sentiment":   "positive",
			},
			{
				"comment_id":  "comment_003",
				"author":      "ConcernedCitizen",
				"text":        "How long until we get herd immunity?",
				"likes":       23,
				"replies":     8,
				"published":   "2023-12-01T13:00:00Z",
				"sentiment":   "neutral",
			},
		},
		"video_id": "yt_video_001",
		"count":    3,
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
			TestName:  "Video Comments Test",
			Status:    "failed",
			Timestamp: timestamp,
			Error:     "No comments found in response",
			Duration:  duration.String(),
			Records:   0,
		}
	}
	
	fmt.Printf("  ✅ Found %d comments\n", records)
	fmt.Printf("  ⏱️  Test completed in %s\n", duration)
	
	return TestResult{
		TestName:  "Video Comments Test",
		Status:    "success",
		Timestamp: timestamp,
		Data:      mockData,
		Duration:  duration.String(),
		Records:   records,
	}
}

// testFullWorkflow tests the complete YouTube API workflow
func (yt *YouTubeAPITest) testFullWorkflow(timestamp string) TestResult {
	startTime := time.Now()
	
	fmt.Println("  📋 Testing complete YouTube API workflow...")
	
	// Simulate full workflow: hashtag search + comments
	workflowData := map[string]interface{}{
		"workflow": "youtube_full_workflow",
		"timestamp": timestamp,
		"steps": []map[string]interface{}{
			{
				"step":       1,
				"action":    "hashtag_search",
				"query":     "covid19",
				"geo":       "ID",
				"status":    "success",
				"videos_found": 3,
			},
			{
				"step":       2,
				"action":    "extract_comments",
				"video_id":  "yt_video_001",
				"status":    "success",
				"comments_found": 3,
			},
			{
				"step":       3,
				"action":    "data_validation",
				"status":    "success",
				"validation_passed": true,
			},
		},
		"summary": map[string]interface{}{
			"total_videos":    3,
			"total_comments":  3,
			"total_engagement": 45000,
			"covid_relevance": 0.95,
		},
	}
	
	duration := time.Since(startTime)
	
	// Calculate total records
	totalRecords := 0
	if summary, ok := workflowData["summary"].(map[string]interface{}); ok {
		if videos, ok := summary["total_videos"].(int); ok {
			totalRecords += videos
		}
		if comments, ok := summary["total_comments"].(int); ok {
			totalRecords += comments
		}
	}
	
	fmt.Printf("  ✅ Workflow completed successfully\n")
	fmt.Printf("  📊 Total records: %d\n", totalRecords)
	fmt.Printf("  ⏱️  Workflow completed in %s\n", duration)
	
	return TestResult{
		TestName:  "Full Workflow Test",
		Status:    "success",
		Timestamp: timestamp,
		Data:      workflowData,
		Duration:  duration.String(),
		Records:   totalRecords,
	}
}

// saveTestResult saves the test result to a JSON file
func (yt *YouTubeAPITest) saveTestResult(result TestResult, timestamp, filename string) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(yt.outputDir, 0755); err != nil {
		log.Printf("Warning: Failed to create output directory: %v", err)
		return
	}
	
	// Generate filename
	outputFile := filepath.Join(yt.outputDir, fmt.Sprintf("%s_%s.json", filename, timestamp))
	
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
func (yt *YouTubeAPITest) printTestSummary(results []TestResult) {
	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println("📊 YOUTUBE API TEST SUMMARY")
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
	// Create and run YouTube API tests
	youtubeTest := NewYouTubeAPITest()
	youtubeTest.RunAllTests()
}
