package etl

import (
	"testing"
	"time"
)

func TestNewDataExtractor(t *testing.T) {
	extractor := NewDataExtractor()
	if extractor == nil {
		t.Error("NewDataExtractor returned nil")
	}
}

func TestNewDataTransformer(t *testing.T) {
	transformer := NewDataTransformer()
	if transformer == nil {
		t.Error("NewDataTransformer returned nil")
	}
}

func TestNewDataLoader(t *testing.T) {
	loader := NewDataLoader()
	if loader == nil {
		t.Error("NewDataLoader returned nil")
	}
}

func TestNewETLOrchestrator(t *testing.T) {
	orchestrator := NewETLOrchestrator()
	if orchestrator == nil {
		t.Error("NewETLOrchestrator returned nil")
	}
}

func TestDataTransformerCleanText(t *testing.T) {
	transformer := NewDataTransformer()
	
	// Test clean text functionality
	cleanText := transformer.cleanText("  Hello   World!  ")
	expected := "Hello World!"
	if cleanText != expected {
		t.Errorf("cleanText failed: expected '%s', got '%s'", expected, cleanText)
	}
}

func TestDataTransformerCalculateCovidRelevance(t *testing.T) {
	transformer := NewDataTransformer()
	
	// Test COVID relevance calculation
	score := transformer.calculateCovidRelevance("covid vaccine indonesia")
	if score <= 0 {
		t.Error("COVID relevance score should be greater than 0 for COVID-related text")
	}
	
	// Test non-COVID text
	score = transformer.calculateCovidRelevance("cooking recipe food")
	if score != 0 {
		t.Error("COVID relevance score should be 0 for non-COVID text")
	}
}

func TestDataTransformerDetectLanguage(t *testing.T) {
	transformer := NewDataTransformer()
	
	// Test Indonesian detection
	lang := transformer.detectLanguage("yang dan atau dengan untuk dari ke di pada")
	if lang != "id" {
		t.Errorf("Language detection failed for Indonesian: expected 'id', got '%s'", lang)
	}
	
	// Test English detection
	lang = transformer.detectLanguage("the and or with for from to in on at")
	if lang != "en" {
		t.Errorf("Language detection failed for English: expected 'en', got '%s'", lang)
	}
}

func TestDataTransformerParseDateTime(t *testing.T) {
	transformer := NewDataTransformer()
	
	// Test RFC3339 format
	dateStr := "2023-12-25T10:30:00Z"
	parsed := transformer.parseDateTime(dateStr)
	if parsed == "" {
		t.Error("DateTime parsing failed for RFC3339 format")
	}
	
	// Test invalid format
	invalidDate := "invalid-date"
	parsed = transformer.parseDateTime(invalidDate)
	if parsed != invalidDate {
		t.Error("DateTime parsing should return original string for invalid format")
	}
}

func TestDataTransformerCreateSummary(t *testing.T) {
	transformer := NewDataTransformer()
	
	// Create test data
	testVideos := []TransformedVideo{
		{
			ID:                  "test1",
			Title:               "COVID Vaccine Update",
			CovidRelevanceScore: 0.8,
			Language:            "en",
			WordCount:           100,
		},
		{
			ID:                  "test2",
			Title:               "Indonesia COVID Cases",
			CovidRelevanceScore: 0.9,
			Language:            "id",
			WordCount:           150,
		},
	}
	
	testArticles := []TransformedArticle{
		{
			ID:                  "article1",
			Title:               "Pandemic Response",
			CovidRelevanceScore: 0.7,
			Language:            "en",
			WordCount:           200,
		},
	}
	
	summary := transformer.createSummary(testVideos, testArticles)
	
	if summary.TotalVideos != 2 {
		t.Errorf("Expected 2 videos, got %d", summary.TotalVideos)
	}
	
	if summary.TotalArticles != 1 {
		t.Errorf("Expected 1 article, got %d", summary.TotalArticles)
	}
	
	if summary.AverageRelevance <= 0 {
		t.Error("Average relevance should be greater than 0")
	}
}

func TestDataLoaderSaveLocally(t *testing.T) {
	loader := NewDataLoader()
	
	// Create test data
	testData := &TransformedData{
		YouTube: []TransformedVideo{
			{
				ID:    "test1",
				Title: "Test Video",
			},
		},
		News: []TransformedArticle{
			{
				ID:    "article1",
				Title: "Test Article",
			},
		},
	}
	
	result := loader.saveLocally(testData)
	
	if !result.Success {
		t.Error("Save locally should succeed")
	}
	
	if result.RecordsCount != 2 {
		t.Errorf("Expected 2 records, got %d", result.RecordsCount)
	}
}

func TestETLOrchestratorRunPipeline(t *testing.T) {
	orchestrator := NewETLOrchestrator()
	
	// Test that the orchestrator can be created and has the required components
	if orchestrator.extractor == nil {
		t.Error("Orchestrator should have an extractor")
	}
	
	if orchestrator.transformer == nil {
		t.Error("Orchestrator should have a transformer")
	}
	
	if orchestrator.loader == nil {
		t.Error("Orchestrator should have a loader")
	}
}

func TestExtractedDataStructure(t *testing.T) {
	data := &ExtractedData{
		Timestamp: time.Now().Format(time.RFC3339),
		Query:     "covid19",
		Sources:   make(map[string]interface{}),
	}
	
	if data.Timestamp == "" {
		t.Error("Timestamp should not be empty")
	}
	
	if data.Query != "covid19" {
		t.Error("Query should be 'covid19'")
	}
	
	if data.Sources == nil {
		t.Error("Sources should be initialized")
	}
}

func TestTransformedDataStructure(t *testing.T) {
	data := &TransformedData{
		TransformedAt: time.Now().Format(time.RFC3339),
		YouTube:       []TransformedVideo{},
		News:          []TransformedArticle{},
	}
	
	if data.TransformedAt == "" {
		t.Error("TransformedAt should not be empty")
	}
	
	if data.YouTube == nil {
		t.Error("YouTube should be initialized")
	}
	
	if data.News == nil {
		t.Error("News should be initialized")
	}
}
