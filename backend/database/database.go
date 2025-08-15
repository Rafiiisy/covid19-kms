package database

// This file serves as the main entry point for the database package
// It re-exports all the main components for easy importing

// Main Database Components
// - Connection: Database connection management
// - Models: Data structures and table creation
// - Operations: CRUD operations for data

// Database Operations
// - InitDatabase: Initialize PostgreSQL connection
// - CreateTables: Create database schema
// - InsertRawData: Store raw extracted data
// - InsertProcessedData: Store processed data
// - GetLatestProcessedData: Retrieve latest data
// - GetDataBySource: Filter data by source
// - GetDataCount: Get record counts

// Usage Example:
// ```
// // Initialize database
// if err := database.InitDatabase(); err != nil {
//     log.Fatal(err)
// }
// defer database.CloseDatabase()
//
// // Create tables
// if err := database.CreateTables(); err != nil {
//     log.Fatal(err)
// }
//
// // Insert data
// data := &database.ProcessedData{
//     Source: "youtube",
//     Title: "COVID-19 Update",
//     Content: "Latest information...",
//     RelevanceScore: 0.95,
//     Sentiment: "neutral",
//     ProcessedData: "{}",
// }
// database.InsertProcessedData(data)
//
// // Retrieve data
// results, err := database.GetLatestProcessedData(10)
// if err != nil {
//     log.Printf("Error: %v", err)
// }
// ```
