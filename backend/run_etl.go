package main

import (
	"covid19-kms/internal/etl"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Load environment variables
	loadEnv()

	fmt.Println("🚀 Running ETL Pipeline to Debug Categorization")
	fmt.Println("================================================")

	// Create ETL orchestrator
	orchestrator := etl.NewETLOrchestrator()

	// Run the ETL pipeline
	fmt.Println("🔄 Starting ETL pipeline...")
	result := orchestrator.RunETLPipeline()

	// Display results
	fmt.Println("\n📊 ETL Pipeline Results:")
	fmt.Println("==========================")
	fmt.Printf("Status: %s\n", result.Status)
	fmt.Printf("Message: %s\n", result.Message)
	fmt.Printf("Duration: %s\n", result.PipelineDuration)
	fmt.Printf("Timestamp: %s\n", result.Timestamp)

	if result.Error != "" {
		fmt.Printf("❌ Error: %s\n", result.Error)
	} else {
		fmt.Println("✅ ETL pipeline completed successfully!")
		
		// Show extraction summary
		if result.Extraction != nil {
			fmt.Println("\n📊 Extraction Summary:")
			fmt.Printf("  Sources: %d\n", len(result.Extraction.Sources))
			for source, data := range result.Extraction.Sources {
				fmt.Printf("  - %s: %T\n", source, data)
			}
		}

		// Show transformation summary
		if result.Transformation != nil {
			fmt.Println("\n🔄 Transformation Summary:")
			fmt.Printf("  YouTube videos: %d\n", len(result.Transformation.YouTube))
			fmt.Printf("  News articles: %d\n", len(result.Transformation.News))
		}

		// Show loading summary
		if result.Loading != nil {
			fmt.Println("\n💾 Loading Summary:")
			fmt.Printf("  Success: %t\n", result.Loading.Success)
			fmt.Printf("  Message: %s\n", result.Loading.Message)
			fmt.Printf("  Records count: %d\n", result.Loading.RecordsCount)
		}
	}

	fmt.Println("\n🎉 ETL execution completed!")
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
