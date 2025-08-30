-- Migration: Add sentiment fields to processed_data table
-- Date: 2024-12-19
-- Description: Add sentiment_score and sentiment_confidence fields for better sentiment analysis

-- Add new sentiment fields if they don't exist
DO $$ 
BEGIN
    -- Add sentiment_score field
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'processed_data' 
        AND column_name = 'sentiment_score'
    ) THEN
        ALTER TABLE processed_data ADD COLUMN sentiment_score DECIMAL(3,2);
        RAISE NOTICE 'Added sentiment_score column';
    ELSE
        RAISE NOTICE 'sentiment_score column already exists';
    END IF;

    -- Add sentiment_confidence field
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'processed_data' 
        AND column_name = 'sentiment_confidence'
    ) THEN
        ALTER TABLE processed_data ADD COLUMN sentiment_confidence DECIMAL(3,2);
        RAISE NOTICE 'Added sentiment_confidence column';
    ELSE
        RAISE NOTICE 'sentiment_confidence column already exists';
    END IF;

END $$;

-- Create index on sentiment fields for better query performance
CREATE INDEX IF NOT EXISTS idx_processed_data_sentiment ON processed_data(sentiment, sentiment_score);
CREATE INDEX IF NOT EXISTS idx_processed_data_sentiment_confidence ON processed_data(sentiment_confidence);

-- Update existing records to have default sentiment values
UPDATE processed_data 
SET sentiment_score = 0.0, 
    sentiment_confidence = 0.0 
WHERE sentiment_score IS NULL OR sentiment_confidence IS NULL;

-- Add comments to the new columns
COMMENT ON COLUMN processed_data.sentiment_score IS 'Sentiment score from -1.0 (negative) to +1.0 (positive)';
COMMENT ON COLUMN processed_data.sentiment_confidence IS 'Confidence level of the sentiment analysis (0.0 to 1.0)';
