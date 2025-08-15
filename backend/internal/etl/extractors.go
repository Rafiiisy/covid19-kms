package etl

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// DataExtractor orchestrates data extraction from all API sources
type DataExtractor struct {
	youtubeAPI      *YouTubeAPI
	googleNewsAPI   *GoogleNewsAPI
	instagramAPI    *InstagramAPI
	indonesiaNewsAPI *IndonesiaNewsAPI
}

// ExtractedData represents the structure of extracted data from all sources
type ExtractedData struct {
	Timestamp string                 `json:"timestamp"`
	Query     string                 `json:"query"`
	Sources   map[string]interface{} `json:"sources"`
}

// NewDataExtractor creates a new DataExtractor instance
func NewDataExtractor() *DataExtractor {
	return &DataExtractor{
		youtubeAPI:      NewYouTubeAPI(),
		googleNewsAPI:   NewGoogleNewsAPI(),
		instagramAPI:    NewInstagramAPI(),
		indonesiaNewsAPI: NewIndonesiaNewsAPI(),
	}
}

// ExtractAllSources extracts data from all sources concurrently using goroutines
func (de *DataExtractor) ExtractAllSources() *ExtractedData {
	log.Println("🔄 Starting data extraction from all sources...")

	extractedData := &ExtractedData{
		Timestamp: time.Now().Format(time.RFC3339),
		Query:     "covid19",
		Sources:   make(map[string]interface{}),
	}

	// Create channels for concurrent extraction
	youtubeChan := make(chan interface{})
	googleNewsChan := make(chan interface{})
	instagramChan := make(chan interface{})
	indonesiaNewsChan := make(chan interface{})

	// Extract YouTube data concurrently
	go func() {
		log.Println("📺 Extracting YouTube data...")
		data, err := de.extractYouTubeData()
		if err != nil {
			log.Printf("❌ YouTube extraction failed: %v", err)
			youtubeChan <- map[string]string{"error": err.Error()}
		} else {
			// Check if videos data exists and get length
			if data.Videos != nil {
				if videos, ok := data.Videos.([]interface{}); ok {
					log.Printf("✅ YouTube: %d videos extracted", len(videos))
				} else {
					log.Printf("✅ YouTube: data extracted (type: %T)", data.Videos)
				}
			} else {
				log.Printf("✅ YouTube: data extracted")
			}
			youtubeChan <- data
		}
	}()

	// Extract Google News data concurrently
	go func() {
		log.Println("📰 Extracting Google News data...")
		data, err := de.extractGoogleNewsData()
		if err != nil {
			log.Printf("❌ Google News extraction failed: %v", err)
			googleNewsChan <- map[string]string{"error": err.Error()}
		} else {
			// Check if articles data exists and get length
			if data.Articles != nil {
				if articles, ok := data.Articles.([]interface{}); ok {
					log.Printf("✅ Google News: %d articles extracted", len(articles))
				} else {
					log.Printf("✅ Google News: data extracted (type: %T)", data.Articles)
				}
			} else {
				log.Printf("✅ Google News: data extracted")
			}
			googleNewsChan <- data
		}
	}()

	// Extract Instagram data concurrently
	go func() {
		log.Println("📱 Extracting Instagram data...")
		data, err := de.extractInstagramData()
		if err != nil {
			log.Printf("❌ Instagram extraction failed: %v", err)
			instagramChan <- map[string]string{"error": err.Error()}
		} else {
			// Check if posts data exists and get length
			if data.Posts != nil {
				if posts, ok := data.Posts.([]interface{}); ok {
					log.Printf("✅ Instagram: %d posts extracted", len(posts))
				} else {
					log.Printf("✅ Instagram: data extracted (type: %T)", data.Posts)
				}
			} else {
				log.Printf("✅ Instagram: data extracted")
			}
			instagramChan <- data
		}
	}()

	// Extract Indonesia News data concurrently
	go func() {
		log.Println("🇮🇩 Extracting Indonesia News data...")
		data, err := de.extractIndonesiaNewsData()
		if err != nil {
			log.Printf("❌ Indonesia News extraction failed: %v", err)
			indonesiaNewsChan <- map[string]string{"error": err.Error()}
		} else {
			totalArticles := 0
			for _, source := range data.Sources {
				if sourceData, ok := source.(map[string]interface{}); ok {
					if articles, exists := sourceData["articles"]; exists {
						if articleList, ok := articles.([]interface{}); ok {
							totalArticles += len(articleList)
						}
					}
				}
			}
			log.Printf("✅ Indonesia News: %d articles extracted", totalArticles)
			indonesiaNewsChan <- data
		}
	}()

	// Collect results from all channels
	extractedData.Sources["youtube"] = <-youtubeChan
	extractedData.Sources["google_news"] = <-googleNewsChan
	extractedData.Sources["instagram"] = <-instagramChan
	extractedData.Sources["indonesia_news"] = <-indonesiaNewsChan

	log.Println("🎉 Data extraction completed!")
	return extractedData
}

// extractYouTubeData extracts YouTube data
func (de *DataExtractor) extractYouTubeData() (*YouTubeData, error) {
	hashtagResult, err := de.youtubeAPI.GetHashtagVideos("covid19", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get hashtag videos: %w", err)
	}

	if hashtagResult.Status != "success" {
		return nil, fmt.Errorf("YouTube API returned error: %s", hashtagResult.Error)
	}

	return &YouTubeData{
		Timestamp: time.Now().Format(time.RFC3339),
		Videos:    hashtagResult.Data,
	}, nil
}

// extractGoogleNewsData extracts Google News data
func (de *DataExtractor) extractGoogleNewsData() (*NewsData, error) {
	searchResult, err := de.googleNewsAPI.SearchNews("COVID-19", "id", "id-ID")
	if err != nil {
		return nil, fmt.Errorf("failed to search news: %w", err)
	}

	if searchResult.Status != "success" {
		return nil, fmt.Errorf("Google News API returned error: %s", searchResult.Error)
	}

	return &NewsData{
		Timestamp: time.Now().Format(time.RFC3339),
		Articles:  searchResult.Data,
	}, nil
}

// extractInstagramData extracts Instagram data
func (de *DataExtractor) extractInstagramData() (*InstagramData, error) {
	hashtagResult, err := de.instagramAPI.GetHashtagMedia("covid19", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get hashtag media: %w", err)
	}

	if hashtagResult.Status != "success" {
		return nil, fmt.Errorf("Instagram API returned error: %s", hashtagResult.Error)
	}

	return &InstagramData{
		Timestamp: time.Now().Format(time.RFC3339),
		Posts:     hashtagResult.Data,
	}, nil
}

// extractIndonesiaNewsData extracts Indonesia News data
func (de *DataExtractor) extractIndonesiaNewsData() (*IndonesiaNewsData, error) {
	sources := []string{"tempo", "kompas", "detik", "cnn"}
	sourceData := make(map[string]interface{})

	for _, source := range sources {
		searchResult, err := de.indonesiaNewsAPI.SearchNews(source, "COVID-19", nil)
		if err != nil {
			log.Printf("Warning: Failed to extract %s news: %v", source, err)
			sourceData[source] = map[string]string{"error": err.Error()}
			continue
		}

		if searchResult.Status == "success" {
			sourceData[source] = searchResult.Data
		} else {
			sourceData[source] = map[string]string{"error": searchResult.Error}
		}
	}

	return &IndonesiaNewsData{
		Timestamp: time.Now().Format(time.RFC3339),
		Sources:   sourceData,
	}, nil
}

// ToJSON converts the extracted data to JSON
func (ed *ExtractedData) ToJSON() ([]byte, error) {
	return json.MarshalIndent(ed, "", "  ")
}
