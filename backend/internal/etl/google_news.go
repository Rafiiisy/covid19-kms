package etl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// GoogleNewsAPI represents the Google News API client for RapidAPI
type GoogleNewsAPI struct {
	APIKey string
	Host   string
	Client *http.Client
}

// GoogleNewsResponse represents the API response structure
type GoogleNewsResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Keyword string      `json:"keyword,omitempty"`
	Lang    string      `json:"lang,omitempty"`
	LR      string      `json:"lr,omitempty"`
}

// NewsData represents the extracted news data
type NewsData struct {
	Timestamp string      `json:"timestamp"`
	Articles  interface{} `json:"articles"`
}

// NewGoogleNewsAPI creates a new Google News API client
func NewGoogleNewsAPI() *GoogleNewsAPI {
	apiKey := os.Getenv("RAPIDAPI_KEY")
	if apiKey == "" {
		apiKey = "your_rapidapi_key_here"
	}

	return &GoogleNewsAPI{
		APIKey: apiKey,
		Host:   "google-news13.p.rapidapi.com",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchNews searches for news articles with the given keyword
func (gn *GoogleNewsAPI) SearchNews(keyword, lang, lr string) (*GoogleNewsResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("keyword", keyword)
	params.Set("lr", lr)
	if lang != "" {
		params.Set("lang", lang)
	}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/search?%s", gn.Host, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("x-rapidapi-key", gn.APIKey)
	req.Header.Set("x-rapidapi-host", gn.Host)

	// Make request
	resp, err := gn.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result GoogleNewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Set additional fields
	result.Keyword = keyword
	result.Lang = lang
	result.LR = lr

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		result.Status = "error"
		if result.Error == "" {
			result.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}
	}

	return &result, nil
}
