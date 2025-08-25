package etl

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// DataTransformer handles data cleaning, transformation, and enrichment
type DataTransformer struct {
	covidKeywords []string
}

// TransformedData represents the structure of transformed data
type TransformedData struct {
	YouTube       []TransformedVideo   `json:"youtube"`
	News          []TransformedArticle `json:"news"`
	Summary       DataSummary          `json:"summary"`
	TransformedAt string               `json:"transformed_at"`
}

// TransformedVideo represents a transformed YouTube video
type TransformedVideo struct {
	ID                  string                 `json:"id"`
	Title               string                 `json:"title"`
	Description         string                 `json:"description"`
	PublishedAt         string                 `json:"published_at"`
	ChannelTitle        string                 `json:"channel_title"`
	ThumbnailURL        string                 `json:"thumbnail_url"`
	Source              string                 `json:"source"`
	CovidRelevanceScore float64                `json:"covid_relevance_score"`
	Language            string                 `json:"language"`
	WordCount           int                    `json:"word_count"`
	ExtractedAt         string                 `json:"extracted_at"`
	TransformedAt       string                 `json:"transformed_at"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// TransformedArticle represents a transformed news article
type TransformedArticle struct {
	ID                  string  `json:"id"`
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	Content             string  `json:"content"`
	URL                 string  `json:"url"`
	Source              string  `json:"source"`
	CovidRelevanceScore float64 `json:"covid_relevance_score"`
	Language            string  `json:"language"`
	WordCount           int     `json:"word_count"`
	ExtractedAt         string  `json:"extracted_at"`
	TransformedAt       string  `json:"transformed_at"`
}

// DataSummary represents summary statistics
type DataSummary struct {
	TotalVideos         int     `json:"total_videos"`
	TotalArticles       int     `json:"total_articles"`
	AverageRelevance    float64 `json:"average_relevance"`
	ProcessingTimestamp string  `json:"processing_timestamp"`
}

// NewDataTransformer creates a new DataTransformer instance
func NewDataTransformer() *DataTransformer {
	return &DataTransformer{
		covidKeywords: []string{
			"covid", "coronavirus", "pandemic", "vaccine", "vaccination",
			"lockdown", "quarantine", "social distancing", "mask",
			"indonesia", "jakarta", "jawa", "sulawesi", "sumatra",
		},
	}
}

// TransformData transforms all extracted data
func (dt *DataTransformer) TransformData(youtubeData, newsData, instagramData interface{}) *TransformedData {
	log.Println("Starting data transformation...")

	transformedData := &TransformedData{
		TransformedAt: time.Now().Format(time.RFC3339),
	}

	// Transform YouTube data
	if youtubeData != nil {
		transformedData.YouTube = dt.transformYouTubeData(youtubeData)
	}

	// Transform news data (can be single source or slice of sources)
	if newsData != nil {
		switch v := newsData.(type) {
		case []interface{}:
			// Handle multiple news sources
			for _, source := range v {
				if source != nil {
					transformedArticles := dt.transformNewsData(source)
					transformedData.News = append(transformedData.News, transformedArticles...)
				}
			}
		default:
			// Handle single news source
			transformedData.News = dt.transformNewsData(newsData)
		}
	}

	// Transform Instagram data
	if instagramData != nil {
		transformedData.News = append(transformedData.News, dt.transformInstagramData(instagramData)...)
	}

	// Create summary
	transformedData.Summary = dt.createSummary(transformedData.YouTube, transformedData.News)

	log.Println("Data transformation completed")
	return transformedData
}

// transformYouTubeData transforms YouTube data (now comments with video metadata)
func (dt *DataTransformer) transformYouTubeData(data interface{}) []TransformedVideo {
	var transformedVideos []TransformedVideo

	log.Println("Transforming YouTube data (comments)...")

	// Handle different data structures from different sources
	switch v := data.(type) {
	case *YouTubeData:
		// Handle YouTube API response structure - now contains comments with video metadata
		if v.Videos != nil {
			if commentsList, ok := v.Videos.([]interface{}); ok {
				log.Printf("Transforming %d YouTube comments", len(commentsList))
				for _, commentData := range commentsList {
					if commentMap, ok := commentData.(map[string]interface{}); ok {
						// Extract comment and video info
						if comment, exists := commentMap["comment"]; exists {
							if video, exists := commentMap["video"]; exists {
								transformedVideo := dt.transformYouTubeComment(comment, video)
								if transformedVideo != nil {
									transformedVideos = append(transformedVideos, *transformedVideo)
								}
							}
						}
					}
				}
			}
		}
	case map[string]interface{}:
		// Handle other YouTube API response structures
		if videos, ok := v["videos"]; ok {
			if videosList, ok := videos.([]interface{}); ok {
				for _, video := range videosList {
					if videoMap, ok := video.(map[string]interface{}); ok {
						transformedVideo := dt.transformYouTubeVideo(videoMap)
						if transformedVideo != nil {
							transformedVideos = append(transformedVideos, *transformedVideo)
						}
					}
				}
			}
		}
	}

	log.Printf("Transformed %d YouTube comments", len(transformedVideos))
	return transformedVideos
}

// transformYouTubeComment transforms a YouTube comment with video metadata
func (dt *DataTransformer) transformYouTubeComment(comment interface{}, video interface{}) *TransformedVideo {
	if commentMap, ok := comment.(map[string]interface{}); ok {
		if videoMap, ok := video.(map[string]interface{}); ok {
			// Extract comment content
			content := ""
			if commentContent, exists := commentMap["content"]; exists {
				content = fmt.Sprintf("%v", commentContent)
			}

			// Extract video metadata
			title := ""
			if videoTitle, exists := videoMap["title"]; exists {
				title = fmt.Sprintf("%v", videoTitle)
			}

			// Calculate COVID relevance score based on content
			relevanceScore := dt.calculateCOVIDRelevance(content)

			// Create rich metadata
			metadata := map[string]interface{}{
				"video": map[string]interface{}{
					"title":     videoMap["title"],
					"videoId":   videoMap["videoId"],
					"url":       videoMap["url"],
					"views":     videoMap["views"],
					"duration":  videoMap["duration"],
					"author":    videoMap["author"],
					"published": videoMap["published"],
				},
				"comment": map[string]interface{}{
					"author":            commentMap["author"],
					"content":           commentMap["content"],
					"publishedTimeText": commentMap["publishedTimeText"],
					"replies":           commentMap["stats"].(map[string]interface{})["replies"],
					"votes":             commentMap["stats"].(map[string]interface{})["votes"],
					"commentId":         commentMap["commentId"],
				},
			}

			// Create transformed video entry (representing a comment)
			return &TransformedVideo{
				ID:                  fmt.Sprintf("comment_%v", time.Now().UnixNano()),
				Title:               title,
				Description:         content, // Comment content goes in description
				PublishedAt:         time.Now().Format(time.RFC3339),
				ChannelTitle:        "YouTube Comments",
				ThumbnailURL:        "",
				Source:              "YouTube",
				CovidRelevanceScore: relevanceScore,
				Language:            "en",
				WordCount:           len(strings.Split(content, " ")),
				ExtractedAt:         time.Now().Format(time.RFC3339),
				TransformedAt:       time.Now().Format(time.RFC3339),
				Metadata:            metadata,
			}
		}
	}
	return nil
}

// calculateCOVIDRelevance calculates relevance score for COVID-19 content
func (dt *DataTransformer) calculateCOVIDRelevance(content string) float64 {
	contentLower := strings.ToLower(content)
	score := 0.0

	// Check for COVID-related keywords
	for _, keyword := range dt.covidKeywords {
		if strings.Contains(contentLower, strings.ToLower(keyword)) {
			score += 0.2
		}
	}

	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}

	// Minimum relevance for any comment
	if score < 0.1 {
		score = 0.1
	}

	return score
}

// transformInstagramData transforms Instagram data to TransformedArticle format
func (dt *DataTransformer) transformInstagramData(data interface{}) []TransformedArticle {
	var transformedArticles []TransformedArticle

	log.Println("Transforming Instagram data...")

	// Handle different Instagram data structures
	switch v := data.(type) {
	case *InstagramData:
		// Handle Instagram API response structure
		if v.Posts != nil {
			if postsList, ok := v.Posts.([]interface{}); ok {
				log.Printf("Transforming %d Instagram posts", len(postsList))
				for _, post := range postsList {
					if postMap, ok := post.(map[string]interface{}); ok {
						transformedArticle := dt.transformInstagramPost(postMap)
						if transformedArticle != nil {
							transformedArticles = append(transformedArticles, *transformedArticle)
						}
					}
				}
			}
		}
	case map[string]interface{}:
		// Handle other Instagram API response structures
		if posts, ok := v["posts"]; ok {
			if postsList, ok := posts.([]interface{}); ok {
				for _, post := range postsList {
					if postMap, ok := post.(map[string]interface{}); ok {
						transformedArticle := dt.transformInstagramPost(postMap)
						if transformedArticle != nil {
							transformedArticles = append(transformedArticles, *transformedArticle)
						}
					}
				}
			}
		}
	}

	log.Printf("Transformed %d Instagram posts", len(transformedArticles))
	return transformedArticles
}

// transformYouTubeVideo transforms a single YouTube video
func (dt *DataTransformer) transformYouTubeVideo(videoMap map[string]interface{}) *TransformedVideo {
	// Extract title
	title := ""
	if titleVal, ok := videoMap["title"]; ok {
		title = dt.cleanText(fmt.Sprintf("%v", titleVal))
	}

	// Extract description
	description := ""
	if descVal, ok := videoMap["descriptionSnippet"]; ok {
		description = dt.cleanText(fmt.Sprintf("%v", descVal))
	}

	// Extract published date
	publishedAt := ""
	if publishedVal, ok := videoMap["publishedTimeText"]; ok {
		publishedAt = fmt.Sprintf("%v", publishedVal)
	}

	// Extract channel title
	channelTitle := ""
	if authorVal, ok := videoMap["author"].(map[string]interface{}); ok {
		if authorTitleVal, ok := authorVal["title"]; ok {
			channelTitle = fmt.Sprintf("%v", authorTitleVal)
		}
	}

	// Extract thumbnail URL
	thumbnailURL := ""
	if thumbnailsVal, ok := videoMap["thumbnails"].([]interface{}); ok && len(thumbnailsVal) > 0 {
		if firstThumb, ok := thumbnailsVal[0].(map[string]interface{}); ok {
			if urlVal, ok := firstThumb["url"]; ok {
				thumbnailURL = fmt.Sprintf("%v", urlVal)
			}
		}
	}

	// Extract video ID
	_ = ""
	if idVal, ok := videoMap["videoId"]; ok {
		_ = fmt.Sprintf("%v", idVal)
	}

	// Calculate COVID-19 relevance score
	relevanceScore := dt.calculateCovidRelevance(title + " " + description)

	// Detect language
	language := dt.detectLanguage(title + " " + description)

	// Calculate word count
	wordCount := len(strings.Fields(title + " " + description))

	// Generate unique ID
	id := dt.generateVideoID(videoMap)

	// Create transformed video
	transformedVideo := &TransformedVideo{
		ID:                  id,
		Title:               title,
		Description:         description,
		PublishedAt:         publishedAt,
		ChannelTitle:        channelTitle,
		ThumbnailURL:        thumbnailURL,
		Source:              "YouTube",
		CovidRelevanceScore: relevanceScore,
		Language:            language,
		WordCount:           wordCount,
		ExtractedAt:         time.Now().Format(time.RFC3339),
		TransformedAt:       time.Now().Format(time.RFC3339),
	}

	return transformedVideo
}

// transformNewsData transforms news data
func (dt *DataTransformer) transformNewsData(data interface{}) []TransformedArticle {
	var transformedArticles []TransformedArticle

	log.Println("Transforming news data...")
	log.Printf("Debug: News data type: %T", data)
	log.Printf("Debug: News data value: %+v", data)

	// Handle different data structures from different sources
	switch v := data.(type) {
	case *IndonesiaNewsData:
		// Handle Indonesia News API response structure
		if v.Sources != nil {
			if items, ok := v.Sources["items"]; ok {
				if itemsList, ok := items.([]interface{}); ok {
					log.Printf("Transforming %d Indonesia news items", len(itemsList))
					for _, item := range itemsList {
						if articleMap, ok := item.(map[string]interface{}); ok {
							transformedArticle := dt.transformNewsItem(articleMap)
							if transformedArticle != nil {
								transformedArticles = append(transformedArticles, *transformedArticle)
							}
						}
					}
				}
			}
		}
	case *InstagramData:
		// Handle Instagram posts structure
		if v.Posts != nil {
			if postsList, ok := v.Posts.([]interface{}); ok {
				log.Printf("Transforming %d Instagram posts", len(postsList))
				for _, post := range postsList {
					if postMap, ok := post.(map[string]interface{}); ok {
						transformedArticle := dt.transformInstagramPost(postMap)
						if transformedArticle != nil {
							transformedArticles = append(transformedArticles, *transformedArticle)
						}
					}
				}
			}
		}
	case *NewsData:
		// Handle Real-Time News API response structure
		if v.Articles != nil {
			if articlesList, ok := v.Articles.([]interface{}); ok {
				log.Printf("Transforming %d Real-Time news articles", len(articlesList))
				for _, article := range articlesList {
					if articleMap, ok := article.(map[string]interface{}); ok {
						transformedArticle := dt.transformNewsItem(articleMap)
						if transformedArticle != nil {
							transformedArticles = append(transformedArticles, *transformedArticle)
						}
					}
				}
			}
		}
	case map[string]interface{}:
		// Handle other news API response structures
		log.Printf("Debug: Processing map[string]interface{} with keys: %v", getMapKeys(v))

		if items, ok := v["items"]; ok {
			if itemsList, ok := items.([]interface{}); ok {
				log.Printf("Debug: Found 'items' key with %d items", len(itemsList))

				// Check if this looks like Indonesia News data by examining the first item
				isIndonesiaNews := false
				if len(itemsList) > 0 {
					if firstItem, ok := itemsList[0].(map[string]interface{}); ok {
						log.Printf("Debug: First item keys: %v", getMapKeys(firstItem))
						// Check for Indonesia News specific fields
						for key := range firstItem {
							if key == "namakanal" || key == "idberita" || key == "namaparent" || key == "namasubkanal" {
								isIndonesiaNews = true
								log.Printf("Debug: Detected Indonesia News field: %s", key)
								break
							}
						}
					}
				}

				if isIndonesiaNews {
					log.Printf("Transforming %d Indonesia news items (flattened structure)", len(itemsList))
					for _, item := range itemsList {
						if articleMap, ok := item.(map[string]interface{}); ok {
							transformedArticle := dt.transformNewsItem(articleMap)
							if transformedArticle != nil {
								transformedArticles = append(transformedArticles, *transformedArticle)
							}
						}
					}
				} else {
					log.Printf("Debug: Not Indonesia News, processing as generic news")
					// Handle other news sources
					for _, item := range itemsList {
						if articleMap, ok := item.(map[string]interface{}); ok {
							transformedArticle := dt.transformNewsItem(articleMap)
							if transformedArticle != nil {
								transformedArticles = append(transformedArticles, *transformedArticle)
							}
						}
					}
				}
			}
		}
		// Handle Instagram posts structure
		if posts, ok := v["posts"]; ok {
			if postsList, ok := posts.([]interface{}); ok {
				for _, post := range postsList {
					if postMap, ok := post.(map[string]interface{}); ok {
						transformedArticle := dt.transformInstagramPost(postMap)
						if transformedArticle != nil {
							transformedArticles = append(transformedArticles, *transformedArticle)
						}
					}
				}
			}
		}
	}

	log.Printf("Transformed %d news articles", len(transformedArticles))
	return transformedArticles
}

// transformNewsItem transforms a single news item to TransformedArticle
func (dt *DataTransformer) transformNewsItem(articleMap map[string]interface{}) *TransformedArticle {
	// Debug: Log available fields for Indonesia News
	log.Printf("Debug: transformNewsItem called with fields: %v", getMapKeys(articleMap))

	// Extract title
	title := ""
	if titleVal, ok := articleMap["title"]; ok {
		title = dt.cleanText(fmt.Sprintf("%v", titleVal))
	}

	// Extract description/summary
	description := ""
	if descVal, ok := articleMap["summary"]; ok {
		description = dt.cleanText(fmt.Sprintf("%v", descVal))
	} else if descVal, ok := articleMap["description"]; ok {
		description = dt.cleanText(fmt.Sprintf("%v", descVal))
	} else if descVal, ok := articleMap["snippet"]; ok {
		description = dt.cleanText(fmt.Sprintf("%v", descVal))
	}

	// Extract content (use description if no content)
	content := description
	if contentVal, ok := articleMap["content"]; ok {
		content = dt.cleanText(fmt.Sprintf("%v", contentVal))
	}

	// Extract URL
	url := ""
	if urlVal, ok := articleMap["url"]; ok {
		url = fmt.Sprintf("%v", urlVal)
	} else if urlVal, ok := articleMap["link"]; ok {
		url = fmt.Sprintf("%v", urlVal)
	}

	// Extract source
	source := ""

	// For Real-Time News, check first by looking for specific fields
	if _, hasArticleID := articleMap["article_id"]; hasArticleID {
		log.Printf("Debug: Found article_id field, setting source to 'Real-Time News'")
		source = "Real-Time News"
	} else if _, hasSourceName := articleMap["source_name"]; hasSourceName {
		log.Printf("Debug: Found source_name field, setting source to 'Real-Time News'")
		source = "Real-Time News"
	}

	// For Indonesia News, check next
	if source == "" {
		log.Printf("Debug: Source is empty, checking for Indonesia News fields")
		// Check if this is from Indonesia News by looking for specific fields
		if _, hasNamakanal := articleMap["namakanal"]; hasNamakanal {
			log.Printf("Debug: Found namakanal field, setting source to 'Indonesia News'")
			source = "Indonesia News"
		} else if _, hasIdberita := articleMap["idberita"]; hasIdberita {
			log.Printf("Debug: Found idberita field, setting source to 'Indonesia News'")
			source = "Indonesia News"
		}
	}

	// Fallback to other source fields if still empty
	if source == "" {
		if sourceVal, ok := articleMap["source"]; ok {
			source = fmt.Sprintf("%v", sourceVal)
			log.Printf("Debug: Found 'source' field: %s", source)
		} else if sourceVal, ok := articleMap["source_name"]; ok {
			source = fmt.Sprintf("%v", sourceVal)
			log.Printf("Debug: Found 'source_name' field: %s", source)
		}
	}

	// Debug: Log available fields for Indonesia News
	if source == "" {
		// Check if this looks like Indonesia News data
		hasIndonesiaFields := false
		for key := range articleMap {
			if key == "namakanal" || key == "idberita" || key == "namaparent" || key == "namasubkanal" {
				hasIndonesiaFields = true
				break
			}
		}
		if hasIndonesiaFields {
			source = "Indonesia News"
			// Log the available fields for debugging
			log.Printf("Debug: Indonesia News article with fields: %v", getMapKeys(articleMap))
		}
	}

	// Check URL-based detection for Indonesian news sources
	if source == "" && url != "" {
		urlLower := strings.ToLower(url)
		if strings.Contains(urlLower, "detik.com") ||
			strings.Contains(urlLower, "kompas.com") ||
			strings.Contains(urlLower, "cnnindonesia.com") {
			source = "Indonesia News"
			log.Printf("Debug: Detected Indonesian news source from URL: %s", url)
		}
	}

	// Extract published date (for future use if needed)
	_ = ""
	if dateVal, ok := articleMap["published_at"]; ok {
		_ = dt.parseDateTime(fmt.Sprintf("%v", dateVal))
	} else if dateVal, ok := articleMap["date"]; ok {
		if dateMap, ok := dateVal.(map[string]interface{}); ok {
			if publishVal, ok := dateMap["publish"]; ok {
				_ = dt.parseDateTime(fmt.Sprintf("%v", publishVal))
			}
		}
	} else if dateVal, ok := articleMap["published_datetime_utc"]; ok {
		_ = dt.parseDateTime(fmt.Sprintf("%v", dateVal))
	}

	// Extract author (for future use if needed)
	_ = ""
	if authorVal, ok := articleMap["author"]; ok {
		_ = fmt.Sprintf("%v", authorVal)
	} else if authorVal, ok := articleMap["penulis"]; ok {
		_ = fmt.Sprintf("%v", authorVal)
	} else if authorVal, ok := articleMap["editor"]; ok {
		_ = fmt.Sprintf("%v", authorVal)
	}

	// Calculate COVID-19 relevance score
	relevanceScore := dt.calculateCovidRelevance(title + " " + description + " " + content)

	// Detect language
	language := dt.detectLanguage(title + " " + description + " " + content)

	// Calculate word count
	wordCount := len(strings.Fields(title + " " + description + " " + content))

	// Generate unique ID
	id := dt.generateArticleID(articleMap)

	// Create transformed article
	transformedArticle := &TransformedArticle{
		ID:                  id,
		Title:               title,
		Description:         description,
		Content:             content,
		URL:                 url,
		Source:              source,
		CovidRelevanceScore: relevanceScore,
		Language:            language,
		WordCount:           wordCount,
		ExtractedAt:         time.Now().Format(time.RFC3339),
		TransformedAt:       time.Now().Format(time.RFC3339),
	}

	return transformedArticle
}

// transformInstagramPost transforms a single Instagram post to TransformedArticle
func (dt *DataTransformer) transformInstagramPost(postMap map[string]interface{}) *TransformedArticle {
	// Extract caption text
	caption := ""
	if captionVal, ok := postMap["caption_text"]; ok {
		caption = dt.cleanText(fmt.Sprintf("%v", captionVal))
	}

	// Extract post code/ID
	postCode := ""
	if codeVal, ok := postMap["code"]; ok {
		postCode = fmt.Sprintf("%v", codeVal)
	}

	// Extract like count
	likeCount := 0
	if likeVal, ok := postMap["like_count"]; ok {
		if likeInt, ok := likeVal.(float64); ok {
			likeCount = int(likeInt)
		}
	}

	// Extract comment count
	commentCount := 0
	if commentVal, ok := postMap["comment_count"]; ok {
		if commentInt, ok := commentVal.(float64); ok {
			commentCount = int(commentInt)
		}
	}

	// Extract user info
	username := ""
	if userVal, ok := postMap["user"].(map[string]interface{}); ok {
		if usernameVal, ok := userVal["username"]; ok {
			username = fmt.Sprintf("%v", usernameVal)
		}
	}

	// Extract timestamp
	timestamp := ""
	if timeVal, ok := postMap["taken_at"]; ok {
		timestamp = fmt.Sprintf("%v", timeVal)
	}

	// Create a description combining caption and engagement metrics
	description := caption
	if likeCount > 0 || commentCount > 0 {
		description += fmt.Sprintf(" (Likes: %d, Comments: %d)", likeCount, commentCount)
	}

	// Calculate COVID-19 relevance score
	relevanceScore := dt.calculateCovidRelevance(caption)

	// Detect language
	language := dt.detectLanguage(caption)

	// Calculate word count
	wordCount := len(strings.Fields(caption))

	// Generate unique ID
	id := dt.generateInstagramPostID(postMap)

	// Create transformed article
	transformedArticle := &TransformedArticle{
		ID:                  id,
		Title:               fmt.Sprintf("Instagram Post by @%s", username),
		Description:         description,
		Content:             caption,
		URL:                 fmt.Sprintf("https://instagram.com/p/%s", postCode),
		Source:              fmt.Sprintf("Instagram (@%s)", username),
		CovidRelevanceScore: relevanceScore,
		Language:            language,
		WordCount:           wordCount,
		ExtractedAt:         timestamp,
		TransformedAt:       time.Now().Format(time.RFC3339),
	}

	return transformedArticle
}

// cleanText cleans and normalizes text
func (dt *DataTransformer) cleanText(text string) string {
	if text == "" {
		return ""
	}

	// Remove extra whitespace
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Remove special characters (keep basic punctuation)
	text = regexp.MustCompile(`[^\w\s.,!?-]`).ReplaceAllString(text, "")

	return text
}

// calculateCovidRelevance calculates relevance score for COVID-19 content
func (dt *DataTransformer) calculateCovidRelevance(text string) float64 {
	if text == "" {
		return 0.0
	}

	text = strings.ToLower(text)
	score := 0.0

	for _, keyword := range dt.covidKeywords {
		if strings.Contains(text, keyword) {
			score += 1.0
		}
	}

	// Normalize score to 0-1 range
	maxPossibleScore := float64(len(dt.covidKeywords))
	if maxPossibleScore > 0 {
		score = score / maxPossibleScore
	}

	return score
}

// detectLanguage detects the language of the text (simplified)
func (dt *DataTransformer) detectLanguage(text string) string {
	if text == "" {
		return "unknown"
	}

	// Simple language detection based on common words
	text = strings.ToLower(text)

	// Indonesian words
	indonesianWords := []string{"yang", "dan", "atau", "dengan", "untuk", "dari", "ke", "di", "pada"}
	for _, word := range indonesianWords {
		if strings.Contains(text, word) {
			return "id"
		}
	}

	// English words
	englishWords := []string{"the", "and", "or", "with", "for", "from", "to", "in", "on", "at"}
	for _, word := range englishWords {
		if strings.Contains(text, word) {
			return "en"
		}
	}

	return "unknown"
}

// parseDateTime parses datetime strings
func (dt *DataTransformer) parseDateTime(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// Try to parse various date formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed.Format(time.RFC3339)
		}
	}

	// Return original if parsing fails
	return dateStr
}

// generateArticleID generates a unique ID for an article
func (dt *DataTransformer) generateArticleID(article interface{}) string {
	// Generate a hash-based ID from article content to ensure uniqueness
	if articleMap, ok := article.(map[string]interface{}); ok {
		// Use title and URL to create a unique hash
		title := ""
		url := ""
		if titleVal, ok := articleMap["title"]; ok {
			title = fmt.Sprintf("%v", titleVal)
		}
		if urlVal, ok := articleMap["url"]; ok {
			url = fmt.Sprintf("%v", urlVal)
		}

		// Create a simple hash from title + url + timestamp
		content := title + url + fmt.Sprintf("%d", time.Now().UnixNano())
		hash := 0
		for _, char := range content {
			hash = ((hash << 5) - hash + int(char)) & 0xffffffff
		}
		return fmt.Sprintf("article_%x", hash)
	}

	// Fallback to timestamp-based ID
	return fmt.Sprintf("article_%d", time.Now().UnixNano())
}

// generateVideoID generates a unique ID for a YouTube video
func (dt *DataTransformer) generateVideoID(video interface{}) string {
	// Generate a hash-based ID from video content to ensure uniqueness
	if videoMap, ok := video.(map[string]interface{}); ok {
		// Use video ID to create a unique hash
		videoID := ""
		if idVal, ok := videoMap["videoId"]; ok {
			videoID = fmt.Sprintf("%v", idVal)
		}

		// Create a simple hash from video ID + timestamp
		content := videoID + fmt.Sprintf("%d", time.Now().UnixNano())
		hash := 0
		for _, char := range content {
			hash = ((hash << 5) - hash + int(char)) & 0xffffffff
		}
		return fmt.Sprintf("video_%x", hash)
	}

	// Fallback to timestamp-based ID
	return fmt.Sprintf("video_%d", time.Now().UnixNano())
}

// generateInstagramPostID generates a unique ID for an Instagram post
func (dt *DataTransformer) generateInstagramPostID(post interface{}) string {
	// Generate a hash-based ID from post content to ensure uniqueness
	if postMap, ok := post.(map[string]interface{}); ok {
		// Use post code and timestamp to create a unique hash
		postCode := ""
		timestamp := ""
		if codeVal, ok := postMap["code"]; ok {
			postCode = fmt.Sprintf("%v", codeVal)
		}
		if timeVal, ok := postMap["taken_at"]; ok {
			timestamp = fmt.Sprintf("%v", timeVal)
		}

		// Create a simple hash from post code + timestamp + current time
		content := postCode + timestamp + fmt.Sprintf("%d", time.Now().UnixNano())
		hash := 0
		for _, char := range content {
			hash = ((hash << 5) - hash + int(char)) & 0xffffffff
		}
		return fmt.Sprintf("instagram_%x", hash)
	}

	// Fallback to timestamp-based ID
	return fmt.Sprintf("instagram_%d", time.Now().UnixNano())
}

// createSummary creates summary statistics
func (dt *DataTransformer) createSummary(videos []TransformedVideo, articles []TransformedArticle) DataSummary {
	totalVideos := len(videos)
	totalArticles := len(articles)

	// Calculate average relevance
	totalRelevance := 0.0
	count := 0

	for _, video := range videos {
		totalRelevance += video.CovidRelevanceScore
		count++
	}

	for _, article := range articles {
		totalRelevance += article.CovidRelevanceScore
		count++
	}

	averageRelevance := 0.0
	if count > 0 {
		averageRelevance = totalRelevance / float64(count)
	}

	return DataSummary{
		TotalVideos:         totalVideos,
		TotalArticles:       totalArticles,
		AverageRelevance:    averageRelevance,
		ProcessingTimestamp: time.Now().Format(time.RFC3339),
	}
}

// getMapKeys returns the keys of a map as a slice of strings
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
