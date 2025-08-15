import { useState, useCallback } from 'react';
import { etlAPI, ETLResult, PipelineStatus, HealthStatus } from '../services/api';

export interface ETLState {
  isLoading: boolean;
  lastResult: ETLResult | null;
  pipelineStatus: PipelineStatus | null;
  healthStatus: HealthStatus | null;
  error: string | null;
  lastUpdated: string;
}

export const useETL = () => {
  const [state, setState] = useState<ETLState>({
    isLoading: false,
    lastResult: null,
    pipelineStatus: null,
    healthStatus: null,
    error: null,
    lastUpdated: 'Never',
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

  // Get metrics for dashboard display
  const getMetrics = useCallback(() => {
    if (!state.lastResult) return null;

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
  }, [state.lastResult]);

  return {
    ...state,
    runPipeline,
    getStatus,
    checkHealth,
    clearError,
    getMetrics,
  };
};
