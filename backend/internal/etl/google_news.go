package etl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// RealTimeNewsAPI represents the Real-Time News Data API client for RapidAPI
type RealTimeNewsAPI struct {
	APIKey string
	Host   string
	Client *http.Client
}

// RealTimeNewsResponse represents the API response structure
type RealTimeNewsResponse struct {
	Status    string      `json:"status"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"` // Can be string or object
	Query     string      `json:"query,omitempty"`
	Country   string      `json:"country,omitempty"`
	Lang      string      `json:"lang,omitempty"`
	Limit     int         `json:"limit,omitempty"`
}

// NewsData represents the extracted news data
type NewsData struct {
	Timestamp string      `json:"timestamp"`
	Articles  interface{} `json:"articles"`
}

// NewRealTimeNewsAPI creates a new Real-Time News Data API client
func NewRealTimeNewsAPI() *RealTimeNewsAPI {
	apiKey := os.Getenv("RAPIDAPI_KEY")
	if apiKey == "" {
		apiKey = "your_rapidapi_key_here"
	}

	return &RealTimeNewsAPI{
		APIKey: apiKey,
		Host:   "real-time-news-data.p.rapidapi.com",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchNews searches for news articles with the given parameters
func (rt *RealTimeNewsAPI) SearchNews(query, country, lang string, limit int, timePublished string) (*RealTimeNewsResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("query", query)
	params.Set("limit", fmt.Sprintf("%d", limit))

	if timePublished != "" {
		params.Set("time_published", timePublished)
	} else {
		params.Set("time_published", "anytime") // Default to anytime
	}

	if country != "" {
		params.Set("country", country)
	}

	if lang != "" {
		params.Set("lang", lang)
	}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/search?%s", rt.Host, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("x-rapidapi-key", rt.APIKey)
	req.Header.Set("x-rapidapi-host", rt.Host)

	// Make request
	resp, err := rt.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result RealTimeNewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Set additional fields
	result.Query = query
	result.Country = country
	result.Lang = lang
	result.Limit = limit

	// Check HTTP status and set response status
	if resp.StatusCode == http.StatusOK {
		if result.Status == "" {
			result.Status = "success"
		}
	} else {
		result.Status = "error"
		if result.Error == nil { // Check if Error is nil, indicating no error object
			result.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}
	}

	return &result, nil
}
