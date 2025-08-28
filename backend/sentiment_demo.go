package main

import (
	"fmt"

	"covid19-kms/internal/services"
)

func main() {
	fmt.Println("🧪 Testing Sentiment Analysis...")

	// Create sentiment analyzer
	analyzer := services.NewSentimentAnalyzer()

	// Test cases
	testCases := []string{
		"COVID-19 vaccine rollout is successful and effective",
		"Cases are increasing rapidly, this is terrible news",
		"Daily update on coronavirus statistics",
		"Recovery rates are improving, great progress",
		"Lockdown measures are causing economic problems",
		"Vaksin COVID-19 berhasil dan efektif",       // Indonesian positive
		"Kasus meningkat dengan cepat, berita buruk", // Indonesian negative
		"Update harian statistik coronavirus",        // Indonesian neutral
	}

	for i, text := range testCases {
		result := analyzer.AnalyzeSentiment(text)
		fmt.Printf("\n%d. Text: %s\n", i+1, text)
		fmt.Printf("   Sentiment: %s (Score: %.2f, Confidence: %.2f)\n",
			result.Category, result.Score, result.Confidence)
		fmt.Printf("   Keywords: %v\n", result.Keywords)
	}

	fmt.Println("\n✅ Sentiment analysis test completed!")
}
