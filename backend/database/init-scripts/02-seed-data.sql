-- =============================================================================
-- COVID-19 KMS Database Seed Data
-- Inserts sample data for development and testing
-- =============================================================================

-- Insert sample raw data
INSERT INTO raw_data (source, query, raw_data) VALUES 
('youtube', 'covid19 vaccine', '{"video_id": "sample123", "title": "COVID-19 Vaccine Update", "description": "Latest vaccine information"}'),
('news', 'covid19 indonesia', '{"article_id": "news456", "title": "COVID-19 Cases in Indonesia", "content": "Current situation report"}'),
('instagram', 'covid19 health', '{"post_id": "insta789", "caption": "Health guidelines for COVID-19", "hashtags": ["#covid19", "#health"]}');

-- Insert sample processed data
INSERT INTO processed_data (source, title, content, relevance_score, sentiment, processed_data) VALUES 
(
    'youtube',
    'COVID-19 Vaccine Update',
    'Latest information about COVID-19 vaccines and their effectiveness',
    0.95,
    'positive',
    '{"video_id": "sample123", "duration": 300, "views": 10000, "language": "en"}'
),
(
    'news',
    'COVID-19 Cases in Indonesia',
    'Current COVID-19 situation in Indonesia with latest statistics',
    0.88,
    'neutral',
    '{"article_id": "news456", "publish_date": "2025-08-20", "source_url": "https://example.com", "language": "id"}'
),
(
    'instagram',
    'COVID-19 Health Guidelines',
    'Important health guidelines and safety measures for COVID-19',
    0.75,
    'positive',
    '{"post_id": "insta789", "likes": 500, "comments": 25, "hashtags": ["#covid19", "#health"]}'
);

-- Insert sample ETL job logs
INSERT INTO etl_logs (job_id, job_type, status, records_processed, metadata) VALUES 
('job_001', 'extract', 'completed', 10, '{"sources": ["youtube", "news"], "duration_seconds": 45}'),
('job_002', 'transform', 'completed', 10, '{"avg_relevance": 0.86, "duration_seconds": 30}'),
('job_003', 'load', 'completed', 10, '{"target": "postgresql", "duration_seconds": 15}');

-- Print completion message
DO $$
BEGIN
    RAISE NOTICE 'COVID-19 KMS sample data inserted successfully!';
    RAISE NOTICE 'Raw data records: 3';
    RAISE NOTICE 'Processed data records: 3'; 
    RAISE NOTICE 'ETL log records: 3';
END $$;

