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
}

// YouTubeData represents the extracted YouTube data
type YouTubeData struct {
	Timestamp string      `json:"timestamp"`
	Videos    interface{} `json:"videos"`
}

// NewYouTubeAPI creates a new YouTube API client
func NewYouTubeAPI() *YouTubeAPI {
	apiKey := os.Getenv("RAPIDAPI_KEY")
	if apiKey == "" {
		apiKey = "your_rapidapi_key_here"
	}

	return &YouTubeAPI{
		APIKey: apiKey,
		Host:   "yt-api.p.rapidapi.com",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetHashtagVideos retrieves videos for a specific hashtag
func (yt *YouTubeAPI) GetHashtagVideos(tag, geo string) (*YouTubeResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("tag", tag)
	if geo != "" {
		params.Set("geo", geo)
	}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/hashtag?%s", yt.Host, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("x-rapidapi-key", yt.APIKey)
	req.Header.Set("x-rapidapi-host", yt.Host)

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
	result.Tag = tag
	result.Geo = geo

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
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
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/comments?%s", yt.Host, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("x-rapidapi-key", yt.APIKey)
	req.Header.Set("x-rapidapi-host", yt.Host)

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

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		result.Status = "error"
		if result.Error == "" {
			result.Error = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}
	}

	return &result, nil
}
