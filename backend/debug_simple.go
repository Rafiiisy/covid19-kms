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

	fmt.Println("🔍 Simple Debug: Checking Database State")
	fmt.Println("========================================")

	// Initialize database
	if err := database.InitDatabase(); err != nil {
		fmt.Printf("❌ Failed to initialize database: %v\n", err)
		return
	}
	defer database.CloseDatabase()

	// Get current counts
	counts, err := database.GetDataCount()
	if err != nil {
		fmt.Printf("❌ Failed to get data count: %v\n", err)
		return
	}
	fmt.Printf("📊 Current database counts: %+v\n", counts)

	// Get processed data
	fmt.Println("\n🔍 Fetching processed data...")
	allData, err := database.GetLatestProcessedData(200)
	if err != nil {
		fmt.Printf("❌ Failed to get data: %v\n", err)
		return
	}

	fmt.Printf("📊 Retrieved %d records from database\n", len(allData))

	// Check source distribution
	fmt.Println("\n📊 Source Distribution Analysis:")
	fmt.Println("================================")

	sourceCounts := make(map[string]int)
	for _, data := range allData {
		sourceCounts[data.Source]++
	}

	for source, count := range sourceCounts {
		fmt.Printf("  %s: %d records\n", source, count)
	}

	// Check for any data that might be Real-Time News
	fmt.Println("\n🔍 Looking for potential Real-Time News data...")
	var potentialRealTimeNews []database.ProcessedData

	for _, data := range allData {
		processedDataLower := strings.ToLower(data.ProcessedData)
		if strings.Contains(processedDataLower, "article_id") ||
			strings.Contains(processedDataLower, "source_name") ||
			strings.Contains(processedDataLower, "published_datetime_utc") {
			potentialRealTimeNews = append(potentialRealTimeNews, data)
		}
	}

	if len(potentialRealTimeNews) > 0 {
		fmt.Printf("Found %d potential Real-Time News records\n", len(potentialRealTimeNews))
		fmt.Printf("First record source: %s\n", potentialRealTimeNews[0].Source)
		fmt.Printf("First record processed data: %s\n", truncateString(potentialRealTimeNews[0].ProcessedData, 200))
	} else {
		fmt.Println("❌ No Real-Time News data found in any form")
	}

	fmt.Println("\n🎯 Next Steps:")
	fmt.Println("===============")
	fmt.Println("1. Real-Time News data is missing from database")
	fmt.Println("2. Need to run ETL again to see what's happening")
	fmt.Println("3. Check if Real-Time News extraction is working")
	fmt.Println("4. Verify transformer is processing Real-Time News correctly")

	fmt.Println("\n🎉 Simple debug completed!")
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
