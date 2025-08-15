package etl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// InstagramAPI represents the Instagram API client for RapidAPI
type InstagramAPI struct {
	APIKey string
	Host   string
	Client *http.Client
}

// InstagramResponse represents the API response structure
type InstagramResponse struct {
	Status   string      `json:"status"`
	Data     interface{} `json:"data,omitempty"`
	Error    string      `json:"error,omitempty"`
	Hashtag  string      `json:"hashtag,omitempty"`
	MaxID    string      `json:"max_id,omitempty"`
	MediaID  string      `json:"media_id,omitempty"`
	Amount   int         `json:"amount,omitempty"`
}

// InstagramData represents the extracted Instagram data
type InstagramData struct {
	Timestamp string      `json:"timestamp"`
	Posts     interface{} `json:"posts"`
}

// NewInstagramAPI creates a new Instagram API client
func NewInstagramAPI() *InstagramAPI {
	apiKey := os.Getenv("RAPIDAPI_KEY")
	if apiKey == "" {
		apiKey = "your_rapidapi_key_here"
	}

	return &InstagramAPI{
		APIKey: apiKey,
		Host:   "instagram-premium-api-2023.p.rapidapi.com",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetHashtagMedia retrieves media for a specific hashtag
func (ig *InstagramAPI) GetHashtagMedia(name, maxID string) (*InstagramResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("name", name)
	if maxID != "" {
		params.Set("max_id", maxID)
	}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/v1/hashtag/medias/top/recent/chunk?%s", ig.Host, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("x-rapidapi-key", ig.APIKey)
	req.Header.Set("x-rapidapi-host", ig.Host)

	// Make request
	resp, err := ig.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result InstagramResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Set additional fields
	result.Hashtag = name
	result.MaxID = maxID

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		result.Status = "error"
		if result.Error == "" {
			result.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}
	}

	return &result, nil
}

// GetMediaComments retrieves comments for a specific media post
func (ig *InstagramAPI) GetMediaComments(mediaID string, amount int) (*InstagramResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("amount", fmt.Sprintf("%d", amount))
	params.Set("id", mediaID)

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/v1/media/comments?%s", ig.Host, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("x-rapidapi-key", ig.APIKey)
	req.Header.Set("x-rapidapi-host", ig.Host)

	// Make request
	resp, err := ig.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result InstagramResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Set additional fields
	result.MediaID = mediaID
	result.Amount = amount

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		result.Status = "error"
		if result.Error == "" {
			result.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}
	}

	return &result, nil
}
