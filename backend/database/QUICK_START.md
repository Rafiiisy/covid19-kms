# 🚀 Quick Start: Railway PostgreSQL Setup

This guide will get you up and running with PostgreSQL on Railway in under 10 minutes!

## ⚡ **Step 1: Create Railway Account (2 min)**

1. **Go to:** https://railway.app/
2. **Click "Start a New Project"**
3. **Sign in with GitHub**
4. **Create new project**

## 🗄️ **Step 2: Add PostgreSQL (1 min)**

1. **Click "New Service"**
2. **Choose "Database"**
3. **Select "PostgreSQL"**
4. **Click "Add PostgreSQL"**

## 🔑 **Step 3: Get Connection String (1 min)**

1. **Click on your database service**
2. **Go to "Connect" tab**
3. **Copy the connection string**
4. **Note your credentials**

## 📋 **Step 4: Create Tables (1 min)**

1. **Go to "Query" tab**
2. **Run this SQL:**

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

-- Create indexes
CREATE INDEX idx_raw_data_source ON raw_data(source);
CREATE INDEX idx_processed_data_source ON processed_data(source);
CREATE INDEX idx_processed_data_timestamp ON processed_data(processed_at);
```

## 🔧 **Step 5: Update Environment (1 min)**

1. **Copy `env.example` to `.env`**
2. **Update with your Railway credentials:**

```bash
# Railway Database
DATABASE_URL=postgresql://username:password@host:port/database
DATABASE_HOST=your-railway-host
DATABASE_PORT=5432
DATABASE_NAME=your-database-name
DATABASE_USER=your-username
DATABASE_PASSWORD=your-password
```

## 🚀 **Step 6: Test Your Setup (2 min)**

```bash
# Start your server
go run cmd/api/main.go

# Test database connection
curl http://localhost:8000/api/etl/data

# Test ETL pipeline
curl -X POST http://localhost:8000/api/etl/run

# Check data in database
curl http://localhost:8000/api/etl/data/stats
```

## ✅ **You're Done!**

Your COVID-19 KMS now has:
- ✅ **PostgreSQL database** on Railway
- ✅ **Automatic table creation**
- ✅ **Data loading from ETL pipeline**
- ✅ **API endpoints for data retrieval**
- ✅ **Database statistics**

## 🔍 **Monitor Your Database**

1. **Go to Railway dashboard**
2. **Click on PostgreSQL service**
3. **Check "Query" tab for data**
4. **Monitor performance metrics**

## 🚨 **Troubleshooting**

### **Connection Failed?**
- Check DATABASE_URL in .env
- Verify Railway database is running
- Ensure SSL mode is enabled

### **Tables Not Created?**
- Check database user permissions
- Run CREATE TABLE statements manually

### **Data Not Loading?**
- Check ETL pipeline logs
- Verify table schema matches data

---

**Need help?** Check the full setup guide in `railway_setup.md` 📚
