package etl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// YouTubeAPI represents the YouTube API client for RapidAPI
type YouTubeAPI struct {
	APIKey string
	Host   string
	Client *http.Client
}

// YouTubeResponse represents the API response structure
type YouTubeResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Tag     string      `json:"tag,omitempty"`
	Geo     string      `json:"geo,omitempty"`
	VideoID string      `json:"video_id,omitempty"`

	// Direct API response fields
	Contents         []interface{} `json:"contents,omitempty"`
	CursorNext       string        `json:"cursorNext,omitempty"`
	EstimatedResults int64         `json:"estimatedResults,omitempty"`
	FilterGroups     interface{}   `json:"filterGroups,omitempty"`
	Refinements      interface{}   `json:"refinements,omitempty"`

	// Comments API response fields
	Comments           []interface{} `json:"comments,omitempty"`
	TotalCommentsCount int64         `json:"totalCommentsCount,omitempty"`
	Filters            interface{}   `json:"filters,omitempty"`
}

// YouTubeData represents the extracted YouTube data
type YouTubeData struct {
	Timestamp string      `json:"timestamp"`
	Videos    interface{} `json:"videos"`
}

// NewYouTubeAPI creates a new YouTube API client
func NewYouTubeAPI(apiKey string) *YouTubeAPI {
	fmt.Printf("🔧 Creating YouTube API client with key: %s...\n", apiKey[:10])

	if apiKey == "" {
		apiKey = "your_rapidapi_key_here"
		fmt.Printf("⚠️ Warning: Using default API key: %s\n", apiKey)
	}

	// Debug: Print the API key being used (first 10 chars)
	fmt.Printf("YouTube API Key: %s...\n", apiKey[:10])

	// Get host from environment variable or use default
	host := os.Getenv("YOUTUBE_HOST")
	if host == "" {
		host = "yt-api.p.rapidapi.com" // Default fallback
	}

	client := &YouTubeAPI{
		APIKey: apiKey,
		Host:   host,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	fmt.Printf("✅ YouTube API client created successfully\n")
	fmt.Printf("✅ Host: %s\n", client.Host)
	fmt.Printf("✅ Timeout: %v\n", client.Client.Timeout)

	return client
}

// SearchVideos searches for videos using the correct YouTube API endpoint
func (yt *YouTubeAPI) SearchVideos(query, lang, geo string) (*YouTubeResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("q", query)
	if lang != "" {
		params.Set("hl", lang)
	}
	if geo != "" {
		params.Set("gl", geo)
	}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/search/?%s", yt.Host, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-RapidAPI-Key", yt.APIKey)
	req.Header.Set("X-RapidAPI-Host", yt.Host)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// Debug: Print request details
	fmt.Printf("YouTube Search Request URL: %s\n", req.URL.String())
	fmt.Printf("YouTube Search Headers: %v\n", req.Header)

	// Make request
	resp, err := yt.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result YouTubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Set additional fields
	result.Tag = query
	result.Geo = geo

	// Check HTTP status and set response status
	if resp.StatusCode == http.StatusOK {
		result.Status = "success"
	} else {
		result.Status = "error"
		if result.Error == "" {
			result.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}
	}

	return &result, nil
}

// GetVideoComments retrieves comments for a specific video
func (yt *YouTubeAPI) GetVideoComments(videoID string) (*YouTubeResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("id", videoID)

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/video/comments/?%s", yt.Host, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-RapidAPI-Key", yt.APIKey)
	req.Header.Set("X-RapidAPI-Host", yt.Host)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// Make request
	resp, err := yt.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result YouTubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Set additional fields
	result.VideoID = videoID

	// Check HTTP status and set response status
	if resp.StatusCode == http.StatusOK {
		result.Status = "success"
	} else {
		result.Status = "error"
		if result.Error == "" {
			result.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}
	}

	return &result, nil
}
