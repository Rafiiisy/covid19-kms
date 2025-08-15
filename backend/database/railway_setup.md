# Railway PostgreSQL Setup Guide

This guide will help you set up PostgreSQL on Railway and integrate it with your COVID-19 KMS backend.

## 🚀 **Step 1: Create Railway Account**

1. **Go to:** https://railway.app/
2. **Click "Start a New Project"**
3. **Sign in with GitHub** (recommended)
4. **Create new project**

## 🗄️ **Step 2: Add PostgreSQL Database**

1. **Click "New Service"**
2. **Choose "Database"**
3. **Select "PostgreSQL"**
4. **Click "Add PostgreSQL"**
5. **Wait for database to be created**

## 🔑 **Step 3: Get Connection Details**

1. **Click on your database service**
2. **Go to "Connect" tab**
3. **Copy the connection string**
4. **Note your database credentials**

## 📋 **Step 4: Database Schema Setup**

### **Create Tables**
Go to **"Query" tab** and run:

```sql
-- Raw data table
CREATE TABLE raw_data (
  id SERIAL PRIMARY KEY,
  source VARCHAR(50) NOT NULL,
  extracted_at TIMESTAMP DEFAULT NOW(),
  raw_data JSONB NOT NULL,
  query VARCHAR(255)
);

-- Processed data table
CREATE TABLE processed_data (
  id SERIAL PRIMARY KEY,
  source VARCHAR(50) NOT NULL,
  processed_at TIMESTAMP DEFAULT NOW(),
  title TEXT,
  content TEXT,
  relevance_score DECIMAL(3,2),
  sentiment VARCHAR(20),
  processed_data JSONB NOT NULL
);

-- Create indexes for better performance
CREATE INDEX idx_raw_data_source ON raw_data(source);
CREATE INDEX idx_processed_data_source ON processed_data(source);
CREATE INDEX idx_processed_data_timestamp ON processed_data(processed_at);
```

## 🔧 **Step 5: Update Backend Configuration**

### **Install PostgreSQL Driver**
```bash
cd backend
go get github.com/lib/pq
```

### **Update Environment Variables**
Add to your `backend/.env`:
```bash
# Railway Database
DATABASE_URL=postgresql://username:password@host:port/database
DATABASE_HOST=your-railway-host
DATABASE_PORT=5432
DATABASE_NAME=your-database-name
DATABASE_USER=your-username
DATABASE_PASSWORD=your-password
```

## 📁 **Step 6: Database Connection Files**

### **Database Connection Manager**
Create `backend/database/connection.go`:
```go
package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    _ "github.com/lib/pq"
)

var DB *sql.DB

// InitDatabase initializes the database connection
func InitDatabase() error {
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        // Fallback to individual environment variables
        connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
            os.Getenv("DATABASE_HOST"),
            os.Getenv("DATABASE_PORT"),
            os.Getenv("DATABASE_USER"),
            os.Getenv("DATABASE_PASSWORD"),
            os.Getenv("DATABASE_NAME"),
        )
    }

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("failed to connect to database: %v", err)
    }

    // Test the connection
    if err = DB.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %v", err)
    }

    log.Println("✅ Database connection established")
    return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
    if DB != nil {
        return DB.Close()
    }
    return nil
}
```

### **Database Models**
Create `backend/database/models.go`:
```go
package database

import (
    "time"
)

// RawData represents raw extracted data
type RawData struct {
    ID          int       `json:"id"`
    Source      string    `json:"source"`
    ExtractedAt time.Time `json:"extracted_at"`
    RawData     string    `json:"raw_data"` // JSON string
    Query       string    `json:"query"`
}

// ProcessedData represents processed data
type ProcessedData struct {
    ID             int       `json:"id"`
    Source         string    `json:"source"`
    ProcessedAt    time.Time `json:"processed_at"`
    Title          string    `json:"title"`
    Content        string    `json:"content"`
    RelevanceScore float64   `json:"relevance_score"`
    Sentiment      string    `json:"sentiment"`
    ProcessedData  string    `json:"processed_data"` // JSON string
}

// CreateTables creates all necessary tables
func CreateTables() error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS raw_data (
            id SERIAL PRIMARY KEY,
            source VARCHAR(50) NOT NULL,
            extracted_at TIMESTAMP DEFAULT NOW(),
            raw_data JSONB NOT NULL,
            query VARCHAR(255)
        )`,
        `CREATE TABLE IF NOT EXISTS processed_data (
            id SERIAL PRIMARY KEY,
            source VARCHAR(50) NOT NULL,
            processed_at TIMESTAMP DEFAULT NOW(),
            title TEXT,
            content TEXT,
            relevance_score DECIMAL(3,2),
            sentiment VARCHAR(20),
            processed_data JSONB NOT NULL
        )`,
        `CREATE INDEX IF NOT EXISTS idx_raw_data_source ON raw_data(source)`,
        `CREATE INDEX IF NOT EXISTS idx_processed_data_source ON processed_data(source)`,
        `CREATE INDEX IF NOT EXISTS idx_processed_data_timestamp ON processed_data(processed_at)`,
    }

    for _, query := range queries {
        if _, err := DB.Exec(query); err != nil {
            return fmt.Errorf("failed to execute query: %v", err)
        }
    }

    log.Println("✅ Database tables created successfully")
    return nil
}
```

### **Database Operations**
Create `backend/database/operations.go`:
```go
package database

import (
    "encoding/json"
    "fmt"
    "time"
)

// InsertRawData inserts raw data into the database
func InsertRawData(source, query string, rawData interface{}) error {
    jsonData, err := json.Marshal(rawData)
    if err != nil {
        return fmt.Errorf("failed to marshal raw data: %v", err)
    }

    query := `
        INSERT INTO raw_data (source, query, raw_data)
        VALUES ($1, $2, $3)
    `

    _, err = DB.Exec(query, source, query, string(jsonData))
    if err != nil {
        return fmt.Errorf("failed to insert raw data: %v", err)
    }

    return nil
}

// InsertProcessedData inserts processed data into the database
func InsertProcessedData(data *ProcessedData) error {
    query := `
        INSERT INTO processed_data (source, title, content, relevance_score, sentiment, processed_data)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

    _, err := DB.Exec(query, 
        data.Source, 
        data.Title, 
        data.Content, 
        data.RelevanceScore, 
        data.Sentiment, 
        data.ProcessedData,
    )
    if err != nil {
        return fmt.Errorf("failed to insert processed data: %v", err)
    }

    return nil
}

// GetLatestProcessedData retrieves the latest processed data
func GetLatestProcessedData(limit int) ([]ProcessedData, error) {
    query := `
        SELECT id, source, processed_at, title, content, relevance_score, sentiment, processed_data
        FROM processed_data 
        ORDER BY processed_at DESC 
        LIMIT $1
    `

    rows, err := DB.Query(query, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to query processed data: %v", err)
    }
    defer rows.Close()

    var results []ProcessedData
    for rows.Next() {
        var data ProcessedData
        err := rows.Scan(
            &data.ID,
            &data.Source,
            &data.ProcessedAt,
            &data.Title,
            &data.Content,
            &data.RelevanceScore,
            &data.Sentiment,
            &data.ProcessedData,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %v", err)
        }
        results = append(results, data)
    }

    return results, nil
}

// GetDataBySource retrieves data by source
func GetDataBySource(source string, limit int) ([]ProcessedData, error) {
    query := `
        SELECT id, source, processed_at, title, content, relevance_score, sentiment, processed_data
        FROM processed_data 
        WHERE source = $1
        ORDER BY processed_at DESC 
        LIMIT $2
    `

    rows, err := DB.Query(query, source, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to query data by source: %v", err)
    }
    defer rows.Close()

    var results []ProcessedData
    for rows.Next() {
        var data ProcessedData
        err := rows.Scan(
            &data.ID,
            &data.Source,
            &data.ProcessedAt,
            &data.Title,
            &data.Content,
            &data.RelevanceScore,
            &data.Sentiment,
            &data.ProcessedData,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %v", err)
        }
        results = append(results, data)
    }

    return results, nil
}
```

## 🔄 **Step 7: Update ETL Loader**

### **Update loaders.go**
Replace the local storage methods with database methods:

```go
// LoadData loads transformed data to PostgreSQL
func (dl *DataLoader) LoadData(data *TransformedData) *LoadResult {
    log.Println("Loading data to PostgreSQL...")

    // Count total records
    totalRecords := len(data.YouTube) + len(data.News)

    // Save to database
    for _, video := range data.YouTube {
        processedData := &database.ProcessedData{
            Source:         "youtube",
            Title:          video.Title,
            Content:        video.Description,
            RelevanceScore: video.RelevanceScore,
            Sentiment:      video.Sentiment,
            ProcessedData:  video.ToJSON(),
        }

        if err := database.InsertProcessedData(processedData); err != nil {
            log.Printf("Failed to insert video data: %v", err)
        }
    }

    for _, article := range data.News {
        processedData := &database.ProcessedData{
            Source:         "news",
            Title:          article.Title,
            Content:        article.Content,
            RelevanceScore: article.RelevanceScore,
            Sentiment:      article.Sentiment,
            ProcessedData:  article.ToJSON(),
        }

        if err := database.InsertProcessedData(processedData); err != nil {
            log.Printf("Failed to insert article data: %v", err)
        }
    }

    return &LoadResult{
        Success:      true,
        Message:      "Data successfully loaded to PostgreSQL",
        Timestamp:    time.Now().Format(time.RFC3339),
        RecordsCount: totalRecords,
    }
}
```

## 🚀 **Step 8: Update Main Application**

### **Update main.go**
Add database initialization:

```go
func main() {
    log.Println("🚀 Starting COVID-19 KMS ETL API Server")

    // Initialize database
    if err := database.InitDatabase(); err != nil {
        log.Fatalf("❌ Failed to initialize database: %v", err)
    }
    defer database.CloseDatabase()

    // Create tables if they don't exist
    if err := database.CreateTables(); err != nil {
        log.Fatalf("❌ Failed to create database tables: %v", err)
    }

    // ... rest of your main function
}
```

## 🧪 **Step 9: Test the Setup**

### **Test Database Connection**
```bash
# Start your server
go run cmd/api/main.go

# Test the data endpoint
curl http://localhost:8000/api/etl/data
```

### **Test ETL Pipeline**
```bash
# Run ETL pipeline
curl -X POST http://localhost:8000/api/etl/run

# Check if data was loaded to database
curl http://localhost:8000/api/etl/data
```

## 🔍 **Step 10: Monitor Database**

### **Railway Dashboard**
1. **Go to your Railway project**
2. **Click on PostgreSQL service**
3. **Go to "Query" tab**
4. **Run queries to check data:**

```sql
-- Check raw data
SELECT COUNT(*) FROM raw_data;

-- Check processed data
SELECT COUNT(*) FROM processed_data;

-- Check latest data
SELECT * FROM processed_data ORDER BY processed_at DESC LIMIT 5;
```

## 🚨 **Troubleshooting**

### **Common Issues**

1. **Connection Failed**
   - Check DATABASE_URL in environment variables
   - Verify Railway database is running
   - Check SSL mode (Railway requires SSL)

2. **Table Creation Failed**
   - Ensure database user has CREATE permissions
   - Check if tables already exist

3. **Data Insertion Failed**
   - Verify table schema matches your data
   - Check JSON data format

### **Debug Commands**
```bash
# Check environment variables
echo $DATABASE_URL

# Test database connection
go run -c "database.InitDatabase()"
```

## 📊 **Benefits of Railway PostgreSQL**

✅ **Managed Service** - No server maintenance  
✅ **Automatic Backups** - Data safety  
✅ **Scalable** - Grows with your needs  
✅ **Easy Integration** - Simple connection strings  
✅ **Monitoring** - Built-in performance insights  

## 🎯 **Next Steps**

1. **Deploy to Railway** - Move from local to cloud
2. **Add Data Validation** - Ensure data quality
3. **Implement Caching** - Improve performance
4. **Add Analytics** - Query data for insights
5. **Set up Monitoring** - Track database performance

---

**Your COVID-19 KMS now has a robust PostgreSQL database on Railway!** 🎉
