package etl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// IndonesiaNewsAPI represents the Indonesia News API client for RapidAPI
type IndonesiaNewsAPI struct {
	APIKey string
	Host   string
	Client *http.Client
}

// IndonesiaNewsResponse represents the actual API response structure from RapidAPI
type IndonesiaNewsResponse struct {
	Items    []interface{}          `json:"items,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Status   string                 `json:"status,omitempty"`
	Error    string                 `json:"error,omitempty"`
	Source   string                 `json:"source,omitempty"`
	Query    string                 `json:"query,omitempty"`
	Params   interface{}            `json:"params,omitempty"`
}

// IndonesiaNewsData represents the extracted Indonesia news data
type IndonesiaNewsData struct {
	Timestamp string                 `json:"timestamp"`
	Sources   map[string]interface{} `json:"sources"`
}

// NewIndonesiaNewsAPI creates a new Indonesia News API client
func NewIndonesiaNewsAPI() *IndonesiaNewsAPI {
	apiKey := os.Getenv("RAPIDAPI_KEY")
	if apiKey == "" {
		apiKey = "your_rapidapi_key_here"
	}

	return &IndonesiaNewsAPI{
		APIKey: apiKey,
		Host:   "indonesia-news.p.rapidapi.com",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchNews searches for news from different Indonesian sources
func (in *IndonesiaNewsAPI) SearchNews(source, query string, params map[string]interface{}) (*IndonesiaNewsResponse, error) {
	var endpoint string

	// Build endpoint based on source
	switch source {
	case "cnn":
		page := getIntParam(params, "page", 1)
		limit := getIntParam(params, "limit", 100)
		endpoint = fmt.Sprintf("/search/cnn?query=%s&page=%d&limit=%d", url.QueryEscape(query), page, limit)
	case "detik":
		limit := getIntParam(params, "limit", 10)
		page := getIntParam(params, "page", 1)
		endpoint = fmt.Sprintf("/search/detik?keyword=%s&limit=%d&page=%d", url.QueryEscape(query), limit, page)
	case "kompas":
		page := getIntParam(params, "page", 1)
		limit := getIntParam(params, "limit", 10)
		endpoint = fmt.Sprintf("/search/kompas?command=%s&page=%d&limit=%d", url.QueryEscape(query), page, limit)
	default:
		return &IndonesiaNewsResponse{
			Status: "error",
			Error:  fmt.Sprintf("Unsupported source: %s. Supported sources: cnn, detik, kompas", source),
			Source: source,
			Query:  query,
		}, nil
	}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s%s", in.Host, endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("x-rapidapi-key", in.APIKey)
	req.Header.Set("x-rapidapi-host", in.Host)

	// Make request
	resp, err := in.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status first
	if resp.StatusCode != http.StatusOK {
		return &IndonesiaNewsResponse{
			Status: "error",
			Error:  fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
			Source: source,
			Query:  query,
		}, nil
	}

	// Parse response into the actual API structure
	var apiResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Create our response structure
	result := &IndonesiaNewsResponse{
		Status: "success",
		Source: source,
		Query:  query,
		Params: params,
	}

	// Extract items based on source-specific response structure
	switch source {
	case "cnn":
		// CNN uses "items" field directly
		if items, ok := apiResponse["items"]; ok {
			if itemsArray, ok := items.([]interface{}); ok {
				result.Items = itemsArray
			}
		}
	case "detik":
		// DETIK uses "item" field (singular)
		if items, ok := apiResponse["item"]; ok {
			if itemsArray, ok := items.([]interface{}); ok {
				result.Items = itemsArray
			}
		}
	case "kompas":
		// KOMPAS uses nested structure: xml.pencarian.item
		if xmlData, ok := apiResponse["xml"].(map[string]interface{}); ok {
			if pencarian, ok := xmlData["pencarian"].(map[string]interface{}); ok {
				if items, ok := pencarian["item"]; ok {
					if itemsArray, ok := items.([]interface{}); ok {
						result.Items = itemsArray
					}
				}
			}
		}
	}

	// Extract metadata if it exists
	if metadata, ok := apiResponse["metadata"]; ok {
		if metadataMap, ok := metadata.(map[string]interface{}); ok {
			result.Metadata = metadataMap
		}
	}

	// If no items found, check if there's an error
	if len(result.Items) == 0 {
		if errorMsg, ok := apiResponse["error"]; ok {
			result.Status = "error"
			result.Error = fmt.Sprintf("%v", errorMsg)
		} else {
			// No items but no error - this might be normal for some sources
			result.Status = "success"
			result.Error = "" // Clear any error message
		}
	}

	return result, nil
}

// GetNewsDetail retrieves detailed news article
func (in *IndonesiaNewsAPI) GetNewsDetail(source, identifier string) (*IndonesiaNewsResponse, error) {
	var endpoint string

	// Build endpoint based on source
	switch source {
	case "cnn":
		endpoint = fmt.Sprintf("/detail/cnn?url=%s", url.QueryEscape(identifier))
	case "detik":
		endpoint = fmt.Sprintf("/detail/detik?url=%s", url.QueryEscape(identifier))
	case "kompas":
		endpoint = fmt.Sprintf("/detail/kompas?guid=%s", url.QueryEscape(identifier))
	default:
		return &IndonesiaNewsResponse{
			Status: "error",
			Error:  fmt.Sprintf("Unsupported source: %s. Supported sources: cnn, detik, kompas", source),
			Source: source,
		}, nil
	}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s%s", in.Host, endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("x-rapidapi-key", in.APIKey)
	req.Header.Set("x-rapidapi-host", in.Host)

	// Make request
	resp, err := in.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status first
	if resp.StatusCode != http.StatusOK {
		return &IndonesiaNewsResponse{
			Status: "error",
			Error:  fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
			Source: source,
		}, nil
	}

	// Parse response
	var apiResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Create our response structure
	result := &IndonesiaNewsResponse{
		Status: "success",
		Source: source,
	}

	// Extract items based on source-specific response structure
	switch source {
	case "cnn":
		if items, ok := apiResponse["items"]; ok {
			if itemsArray, ok := items.([]interface{}); ok {
				result.Items = itemsArray
			}
		}
	case "detik":
		if items, ok := apiResponse["item"]; ok {
			if itemsArray, ok := items.([]interface{}); ok {
				result.Items = itemsArray
			}
		}
	case "kompas":
		if xmlData, ok := apiResponse["xml"].(map[string]interface{}); ok {
			if pencarian, ok := xmlData["pencarian"].(map[string]interface{}); ok {
				if items, ok := pencarian["item"]; ok {
					if itemsArray, ok := items.([]interface{}); ok {
						result.Items = itemsArray
					}
				}
			}
		}
	}

	// Extract metadata if it exists
	if metadata, ok := apiResponse["metadata"]; ok {
		if metadataMap, ok := metadata.(map[string]interface{}); ok {
			result.Metadata = metadataMap
		}
	}

	// If no items found, check if there's an error
	if len(result.Items) == 0 {
		if errorMsg, ok := apiResponse["error"]; ok {
			result.Status = "error"
			result.Error = fmt.Sprintf("%v", errorMsg)
		} else {
			// No items but no error - this might be normal for some sources
			result.Status = "success"
			result.Error = "" // Clear any error message
		}
	}

	return result, nil
}

// getIntParam safely extracts an integer parameter from the params map
func getIntParam(params map[string]interface{}, key string, defaultValue int) int {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if parsed, err := strconv.Atoi(v); err == nil {
				return parsed
			}
		}
	}
	return defaultValue
}
