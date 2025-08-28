import { useState, useCallback, useEffect } from 'react';
import { etlAPI, ETLResult, PipelineStatus, HealthStatus } from '../services/api';

export interface ETLState {
  isLoading: boolean;
  lastResult: ETLResult | null;
  pipelineStatus: PipelineStatus | null;
  healthStatus: HealthStatus | null;
  error: string | null;
  lastUpdated: string;
  // Database data state
  databaseData: {
    youtube: any[] | null;
    googleNews: any[] | null;
    instagram: any[] | null;
    indonesiaNews: any[] | null;
    summary: any | null;
    sentimentDistribution: any | null;
    wordFrequency: any | null;
  } | null;
}

export const useETL = () => {
  const [state, setState] = useState<ETLState>({
    isLoading: false,
    lastResult: null,
    pipelineStatus: null,
    healthStatus: null,
    error: null,
    lastUpdated: 'Never',
    databaseData: null,
  });

  // Run the complete ETL pipeline
  const runPipeline = useCallback(async () => {
    setState(prev => ({ ...prev, isLoading: true, error: null }));
    
    try {
      const result = await etlAPI.runPipeline();
      setState(prev => ({
        ...prev,
        lastResult: result,
        lastUpdated: new Date().toLocaleString(),
        isLoading: false,
      }));
      return result;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to run ETL pipeline';
      setState(prev => ({
        ...prev,
        error: errorMessage,
        isLoading: false,
      }));
      throw error;
    }
  }, []);

  // Get pipeline status
  const getStatus = useCallback(async () => {
    try {
      const status = await etlAPI.getStatus();
      setState(prev => ({ ...prev, pipelineStatus: status }));
      return status;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to get pipeline status';
      setState(prev => ({ ...prev, error: errorMessage }));
      throw error;
    }
  }, []);

  // Health check
  const checkHealth = useCallback(async () => {
    try {
      const health = await etlAPI.healthCheck();
      setState(prev => ({ ...prev, healthStatus: health }));
      return health;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to check health';
      setState(prev => ({ ...prev, error: errorMessage }));
      throw error;
    }
  }, []);

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }));
  }, []);

  // Fetch all database data
  const fetchAllDatabaseData = useCallback(async () => {
    try {
      console.log('🔄 Fetching database data...');
      
      // Fetch data from all sources concurrently
      const [youtubeData, googleNewsData, instagramData, indonesiaNewsData, summaryData, sentimentDistributionData, wordFrequencyData] = await Promise.all([
        etlAPI.getYouTubeData(),
        etlAPI.getGoogleNewsData(),
        etlAPI.getInstagramData(),
        etlAPI.getIndonesiaNewsData(),
        etlAPI.getDataSummary(),
        etlAPI.getSentimentDistribution(),
        etlAPI.getWordFrequency(),
      ]);

      const databaseData = {
        youtube: youtubeData.data || [],
        googleNews: googleNewsData.data || [],
        instagram: instagramData.data || [],
        indonesiaNews: indonesiaNewsData.data || [],
        summary: summaryData.summary || {},
        sentimentDistribution: sentimentDistributionData.distribution || {},
        wordFrequency: wordFrequencyData.wordFrequency || {},
      };

      setState(prev => ({ ...prev, databaseData }));
      console.log('✅ Database data fetched successfully:', databaseData);
      
      return databaseData;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch database data';
      console.error('❌ Error fetching database data:', error);
      setState(prev => ({ ...prev, error: errorMessage }));
      throw error;
    }
  }, []);

  // Fetch all database data on component mount
  useEffect(() => {
    fetchAllDatabaseData();
  }, [fetchAllDatabaseData]);

  // Get metrics for dashboard display
  const getMetrics = useCallback(() => {
    // Use database data if available, fallback to ETL result
    if (state.databaseData?.summary) {
      const { summary } = state.databaseData;
      return {
        totalRecords: summary.total_records || 0,
        youtubeVideos: state.databaseData.youtube?.length || 0,
        newsArticles: (state.databaseData.googleNews?.length || 0) + (state.databaseData.indonesiaNews?.length || 0),
        extractionSources: 4, // Fixed number of sources
        averageRelevance: summary.average_relevance || 0,
        pipelineStatus: state.lastResult?.status || 'unknown',
        duration: state.lastResult?.pipeline_duration || 'N/A',
        timestamp: summary.latest_update || 'Never',
      };
    }

    // Fallback to ETL result data
    if (state.lastResult) {
      const { extraction, transformation, loading, summary } = state.lastResult;
      return {
        totalRecords: loading?.records_count || 0,
        youtubeVideos: transformation?.YouTube?.length || 0,
        newsArticles: transformation?.News?.length || 0,
        extractionSources: extraction?.sources ? Object.keys(extraction.sources).length : 0,
        averageRelevance: transformation?.Summary?.AverageRelevance || 0,
        pipelineStatus: state.lastResult.status,
        duration: state.lastResult.pipeline_duration,
        timestamp: state.lastResult.timestamp,
      };
    }

    return null;
  }, [state.databaseData, state.lastResult]);

  return {
    ...state,
    runPipeline,
    getStatus,
    checkHealth,
    clearError,
    getMetrics,
    fetchAllDatabaseData,
  };
};

