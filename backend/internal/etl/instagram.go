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
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"` // For backward compatibility
	Error   string      `json:"error,omitempty"`
	Hashtag string      `json:"hashtag,omitempty"`
	MaxID   string      `json:"max_id,omitempty"`
	MediaID string      `json:"media_id,omitempty"`
	Amount  int         `json:"amount,omitempty"`

	// Direct API response fields for array structure
	Posts  []interface{} `json:"posts,omitempty"`  // Posts data from first array element
	Cursor string        `json:"cursor,omitempty"` // Cursor token from second array element
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

	// First, try to decode as array to handle the actual API response structure
	var rawResponse []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response as array: %w", err)
	}

	// Create result
	result := &InstagramResponse{
		Hashtag: name,
		MaxID:   maxID,
		Status:  "success",
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		result.Status = "error"
		result.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		return result, nil
	}

	// Handle array response structure
	if len(rawResponse) >= 1 {
		// First element contains the posts data
		if postsData, ok := rawResponse[0].([]interface{}); ok {
			result.Posts = postsData
			result.Data = postsData // Keep backward compatibility
		} else {
			result.Error = "First array element is not a posts array"
		}
	}

	// Second element contains cursor/pagination info
	if len(rawResponse) >= 2 {
		if cursorData, ok := rawResponse[1].(string); ok {
			result.Cursor = cursorData
		}
	}

	return result, nil
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

	// Parse response - comments might also return array structure
	var rawResponse []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, fmt.Errorf("failed to decode comments response: %w", err)
	}

	// Create result
	result := &InstagramResponse{
		MediaID: mediaID,
		Amount:  amount,
		Status:  "success",
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		result.Status = "error"
		result.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		return result, nil
	}

	// Handle array response structure for comments
	if len(rawResponse) >= 1 {
		// First element contains the comments data
		if commentsData, ok := rawResponse[0].([]interface{}); ok {
			result.Posts = commentsData // Reuse Posts field for comments
			result.Data = commentsData  // Keep backward compatibility
		} else {
			result.Error = "First array element is not a comments array"
		}
	}

	return result, nil
}
