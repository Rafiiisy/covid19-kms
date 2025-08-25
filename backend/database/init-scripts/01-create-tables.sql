-- =============================================================================
-- COVID-19 KMS Database Initialization Script
-- Creates all necessary tables for the Knowledge Management System
-- =============================================================================

-- Raw data table for storing extracted data
CREATE TABLE IF NOT EXISTS raw_data (
    id SERIAL PRIMARY KEY,
    source VARCHAR(50) NOT NULL,
    extracted_at TIMESTAMP DEFAULT NOW(),
    raw_data JSONB NOT NULL,
    query VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Processed data table for storing transformed data
CREATE TABLE IF NOT EXISTS processed_data (
    id SERIAL PRIMARY KEY,
    source VARCHAR(50) NOT NULL,
    processed_at TIMESTAMP DEFAULT NOW(),
    title TEXT,
    content TEXT,
    relevance_score DECIMAL(3,2),
    sentiment VARCHAR(20),
    processed_data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ETL job logs table for tracking pipeline runs
CREATE TABLE IF NOT EXISTS etl_logs (
    id SERIAL PRIMARY KEY,
    job_id VARCHAR(100) NOT NULL,
    job_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    start_time TIMESTAMP DEFAULT NOW(),
    end_time TIMESTAMP,
    records_processed INTEGER DEFAULT 0,
    error_message TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Performance indexes
CREATE INDEX IF NOT EXISTS idx_raw_data_source ON raw_data(source);
CREATE INDEX IF NOT EXISTS idx_raw_data_extracted_at ON raw_data(extracted_at);
CREATE INDEX IF NOT EXISTS idx_processed_data_source ON processed_data(source);
CREATE INDEX IF NOT EXISTS idx_processed_data_processed_at ON processed_data(processed_at);
CREATE INDEX IF NOT EXISTS idx_processed_data_relevance ON processed_data(relevance_score);
CREATE INDEX IF NOT EXISTS idx_etl_logs_job_id ON etl_logs(job_id);
CREATE INDEX IF NOT EXISTS idx_etl_logs_status ON etl_logs(status);

-- Comments for documentation
COMMENT ON TABLE raw_data IS 'Stores raw data extracted from various sources (YouTube, News, etc.)';
COMMENT ON TABLE processed_data IS 'Stores processed and transformed data ready for analysis';
COMMENT ON TABLE etl_logs IS 'Tracks ETL pipeline execution logs and statistics';

-- Print completion message
DO $$
BEGIN
    RAISE NOTICE 'COVID-19 KMS database tables created successfully!';
END $$;

