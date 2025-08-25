package main

import (
	"covid19-kms/database"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Load environment variables
	loadEnv()

	fmt.Println("🔍 Debugging Generic 'news' Source Articles")
	fmt.Println("============================================")

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		fmt.Printf("❌ Failed to initialize database: %v\n", err)
		return
	}
	defer database.CloseDatabase()

	// Get all data
	fmt.Println("\n🔍 Fetching all data from database...")
	allData, err := database.GetLatestProcessedData(200)
	if err != nil {
		fmt.Printf("❌ Failed to get data: %v\n", err)
		return
	}

	fmt.Printf("📊 Retrieved %d records from database\n", len(allData))

	// Find articles with generic "news" source
	var newsSourceArticles []database.ProcessedData
	for _, data := range allData {
		if data.Source == "news" {
			newsSourceArticles = append(newsSourceArticles, data)
		}
	}

	fmt.Printf("\n📰 Found %d articles with generic 'news' source\n", len(newsSourceArticles))

	if len(newsSourceArticles) == 0 {
		fmt.Println("✅ No generic 'news' source articles found!")
		return
	}

	// Show sample of these articles
	fmt.Println("\n📋 Sample of Generic 'news' Source Articles:")
	fmt.Println("=============================================")
	
	for i, article := range newsSourceArticles {
		if i >= 5 { // Show only first 5
			break
		}
		fmt.Printf("\n%d. Title: %s\n", i+1, article.Title)
		fmt.Printf("   Content: %s\n", truncateString(article.Content, 100))
		fmt.Printf("   Processed Data: %s\n", truncateString(article.ProcessedData, 150))
		fmt.Printf("   ---")
	}

	// Analyze the processed data to understand the structure
	fmt.Println("\n🔍 Analyzing Processed Data Structure:")
	fmt.Println("======================================")
	
	if len(newsSourceArticles) > 0 {
		firstArticle := newsSourceArticles[0]
		fmt.Printf("First article processed data (first 200 chars):\n%s\n", 
			truncateString(firstArticle.ProcessedData, 200))
	}

	fmt.Println("\n🎯 Recommendations:")
	fmt.Println("===================")
	fmt.Println("1. Check if these articles have the expected fields (article_id, source_name)")
	fmt.Println("2. Verify if they're from Real-Time News API or another source")
	fmt.Println("3. Update the transformer logic to properly categorize them")
	fmt.Println("4. Consider adding more field detection patterns")

	fmt.Println("\n🎉 Debug analysis completed!")
}

func loadEnv() {
	envFile := "env"
	if _, err := os.Stat(envFile); err == nil {
		content, err := os.ReadFile(envFile)
		if err != nil {
			fmt.Printf("Warning: Could not read env file: %v\n", err)
			return
		}

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					os.Setenv(parts[0], parts[1])
				}
			}
		}
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
