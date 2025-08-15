# 🗄️ Database Package - PostgreSQL Integration

This package provides PostgreSQL database integration for the COVID-19 KMS backend, replacing local file storage with a robust, scalable database solution.

## 🏗️ **Architecture Overview**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   ETL Pipeline  │───►│  Database Layer │───►│  PostgreSQL DB  │
│                 │    │                 │    │   (Railway)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
   Extract Data          Insert/Query           Store Data
   Transform Data        Manage Schema          Index Data
   Load to DB           Handle Errors          Backup Data
```

## 📁 **Package Structure**

```
database/
├── connection.go      # Database connection management
├── models.go          # Data structures and table creation
├── operations.go      # CRUD operations
├── database.go        # Package entry point
├── README.md          # This file
├── QUICK_START.md     # Quick setup guide
└── railway_setup.md   # Detailed Railway setup
```

## 🔧 **Core Components**

### **1. Connection Management (`connection.go`)**
- **`InitDatabase()`**: Establishes PostgreSQL connection
- **`CloseDatabase()`**: Gracefully closes connection
- **Environment variable support**: DATABASE_URL or individual credentials
- **SSL support**: Required for Railway deployment

### **2. Data Models (`models.go`)**
- **`RawData`**: Raw extracted data structure
- **`ProcessedData`**: Transformed data structure
- **`CreateTables()`**: Automatic table creation with indexes

### **3. Database Operations (`operations.go`)**
- **`InsertRawData()`**: Store raw extracted data
- **`InsertProcessedData()`**: Store processed data
- **`GetLatestProcessedData()`**: Retrieve latest data
- **`GetDataBySource()`**: Filter data by source
- **`GetDataCount()`**: Get record statistics

## 🚀 **Quick Start**

### **1. Install Dependencies**
```bash
go get github.com/lib/pq
```

### **2. Set Environment Variables**
```bash
# Railway Database
DATABASE_URL=postgresql://username:password@host:port/database
DATABASE_HOST=your-railway-host
DATABASE_PORT=5432
DATABASE_NAME=your-database-name
DATABASE_USER=your-username
DATABASE_PASSWORD=your-password
```

### **3. Initialize Database**
```go
import "covid19-kms/database"

func main() {
    // Initialize database connection
    if err := database.InitDatabase(); err != nil {
        log.Fatal(err)
    }
    defer database.CloseDatabase()

    // Create tables
    if err := database.CreateTables(); err != nil {
        log.Fatal(err)
    }
}
```

## 📊 **Database Schema**

### **Raw Data Table**
```sql
CREATE TABLE raw_data (
  id SERIAL PRIMARY KEY,
  source VARCHAR(50) NOT NULL,
  extracted_at TIMESTAMP DEFAULT NOW(),
  raw_data JSONB NOT NULL,
  query VARCHAR(255)
);
```

### **Processed Data Table**
```sql
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
```

### **Indexes**
```sql
CREATE INDEX idx_raw_data_source ON raw_data(source);
CREATE INDEX idx_processed_data_source ON processed_data(source);
CREATE INDEX idx_processed_data_timestamp ON processed_data(processed_at);
```

## 🔄 **Integration with ETL Pipeline**

### **Data Loading**
```go
// In your ETL loader
func (dl *DataLoader) LoadData(data *TransformedData) *LoadResult {
    for _, video := range data.YouTube {
        processedData := &database.ProcessedData{
            Source:         "youtube",
            Title:          video.Title,
            Content:        video.Description,
            RelevanceScore: video.CovidRelevanceScore,
            Sentiment:      "neutral",
            ProcessedData:  video.ToJSON(),
        }
        
        if err := database.InsertProcessedData(processedData); err != nil {
            log.Printf("Failed to insert video data: %v", err)
        }
    }
    
    return &LoadResult{...}
}
```

### **Data Retrieval**
```go
// In your API handler
func (h *DataHandler) GetLatestData(w http.ResponseWriter, r *http.Request) {
    data, err := database.GetLatestProcessedData(100)
    if err != nil {
        http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
        return
    }
    
    // Process and return data
}
```

## 📈 **API Endpoints**

### **Data Retrieval**
- **`GET /api/etl/data`**: Get latest processed data
- **`GET /api/etl/data/source?source=youtube`**: Get data by source
- **`GET /api/etl/data/stats`**: Get database statistics

### **Example Responses**

#### **Latest Data**
```json
{
  "status": "success",
  "timestamp": "2025-08-15T15:30:00+07:00",
  "data": [
    {
      "source": "youtube",
      "title": "COVID-19 Update",
      "content": "Latest information...",
      "relevance_score": 0.95,
      "sentiment": "neutral",
      "processed_at": "2025-08-15T15:25:00+07:00"
    }
  ],
  "total_count": 1
}
```

#### **Database Stats**
```json
{
  "status": "success",
  "timestamp": "2025-08-15T15:30:00+07:00",
  "stats": {
    "raw_data": 25,
    "processed_data": 18
  }
}
```

## 🚨 **Error Handling**

### **Common Errors**
1. **Connection Failed**: Check DATABASE_URL and Railway status
2. **Table Creation Failed**: Verify user permissions
3. **Data Insertion Failed**: Check data format and schema

### **Debug Commands**
```bash
# Check environment variables
echo $DATABASE_URL

# Test database connection
go run -c "database.InitDatabase()"

# Check database status
curl http://localhost:8000/api/etl/data/stats
```

## 🔒 **Security Features**

- **SSL Required**: Railway enforces SSL connections
- **Parameterized Queries**: Prevents SQL injection
- **Connection Pooling**: Efficient resource management
- **Graceful Shutdown**: Proper connection cleanup

## 📊 **Performance Optimizations**

- **Indexed Queries**: Fast data retrieval by source and timestamp
- **JSONB Storage**: Efficient JSON data storage
- **Connection Pooling**: Reuse database connections
- **Batch Operations**: Efficient bulk data insertion

## 🚀 **Deployment on Railway**

### **Benefits**
✅ **Managed Service**: No server maintenance  
✅ **Automatic Backups**: Data safety  
✅ **Scalable**: Grows with your needs  
✅ **Easy Integration**: Simple connection strings  
✅ **Monitoring**: Built-in performance insights  

### **Setup Steps**
1. **Create Railway account**
2. **Add PostgreSQL service**
3. **Get connection details**
4. **Create database schema**
5. **Update environment variables**
6. **Test connection**

## 🧪 **Testing**

### **Unit Tests**
```bash
# Test database package
go test -v ./database

# Test with coverage
go test -v -cover ./database
```

### **Integration Tests**
```bash
# Start server
go run cmd/api/main.go

# Test endpoints
curl http://localhost:8000/api/etl/data
curl http://localhost:8000/api/etl/data/stats
```

## 🔄 **Migration from Local Storage**

### **Before (Local Files)**
```go
// Old local storage approach
func (dl *DataLoader) LoadData(data *TransformedData) *LoadResult {
    // Save to local JSON files
    filename := fmt.Sprintf("data/processed_%s.json", timestamp)
    return dl.saveLocally(data, filename)
}
```

### **After (PostgreSQL)**
```go
// New database approach
func (dl *DataLoader) LoadData(data *TransformedData) *LoadResult {
    // Save to PostgreSQL database
    for _, item := range data.Items {
        processedData := &database.ProcessedData{...}
        database.InsertProcessedData(processedData)
    }
    return &LoadResult{...}
}
```

## 🎯 **Next Steps**

1. **Deploy to Railway**: Move from local to cloud
2. **Add Data Validation**: Ensure data quality
3. **Implement Caching**: Improve performance
4. **Add Analytics**: Query data for insights
5. **Set up Monitoring**: Track database performance

## 📚 **Additional Resources**

- [Railway PostgreSQL Setup](./railway_setup.md)
- [Quick Start Guide](./QUICK_START.md)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Railway Documentation](https://docs.railway.app/)

---

**Your COVID-19 KMS now has enterprise-grade PostgreSQL database integration!** 🎉
