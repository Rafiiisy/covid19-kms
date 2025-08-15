package etl

// This file serves as the main entry point for the ETL package
// It re-exports all the main components for easy importing

// Main ETL Components
// - DataExtractor: Orchestrates data extraction from all API sources
// - DataTransformer: Handles data cleaning, transformation, and enrichment
// - DataLoader: Manages loading data to local storage
// - ETLOrchestrator: Coordinates the entire ETL pipeline

// API Extractors
// - YouTubeAPI: Extracts YouTube videos and comments data
// - GoogleNewsAPI: Extracts Google News articles
// - InstagramAPI: Extracts Instagram posts and media
// - IndonesiaNewsAPI: Extracts Indonesian news from multiple sources

// Data Structures
// - ExtractedData: Raw data from all sources
// - TransformedData: Cleaned and enriched data
// - ETLResult: Complete pipeline execution result
// - LoadResult: Data loading operation result

// Usage Example:
// ```
// orchestrator := etl.NewETLOrchestrator()
// result := orchestrator.RunETLPipeline()
// 
// if result.Status == "success" {
//     fmt.Printf("Pipeline completed in %s\n", result.PipelineDuration)
//     fmt.Printf("Extracted from %d sources\n", len(result.Extraction.Sources))
//     fmt.Printf("Transformed %d videos and %d articles\n", 
//         len(result.Transformation.YouTube), 
//         len(result.Transformation.News))
// }
// ```
