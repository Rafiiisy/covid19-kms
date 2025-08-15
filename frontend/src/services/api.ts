import axios from 'axios';

const API_BASE_URL = 'http://localhost:8000';

// API client with base configuration
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000, // 30 seconds timeout for ETL operations
  headers: {
    'Content-Type': 'application/json',
  },
});

// Types based on your Go backend ETL structures
export interface ETLResult {
  status: string;
  message: string;
  timestamp: string;
  pipeline_duration: string;
  extraction?: ExtractedData;
  transformation?: TransformedData;
  loading?: LoadResult;
  summary?: Record<string, any>;
  error?: string;
}

export interface ExtractedData {
  timestamp: string;
  query: string;
  sources: Record<string, any>;
}

export interface TransformedData {
  YouTube: any[];
  News: any[];
  TransformedAt: string;
  Summary?: {
    AverageRelevance: number;
  };
}

export interface LoadResult {
  success: boolean;
  message: string;
  records_count: number;
  timestamp: string;
}

export interface PipelineStatus {
  status: string;
  timestamp: string;
  service: string;
  version: string;
  endpoints: string[];
  description: string;
}

export interface HealthStatus {
  status: string;
  timestamp: string;
  service: string;
  uptime: string;
}

// ETL API functions
export const etlAPI = {
  // Run the complete ETL pipeline
  runPipeline: async (): Promise<ETLResult> => {
    const response = await apiClient.post('/api/etl/run');
    return response.data;
  },

  // Get current pipeline status
  getStatus: async (): Promise<PipelineStatus> => {
    const response = await apiClient.get('/api/etl/status');
    return response.data;
  },

  // Health check
  healthCheck: async (): Promise<HealthStatus> => {
    const response = await apiClient.get('/api/health');
    return response.data;
  },

  // Run individual ETL stages (if needed)
  extractData: async (): Promise<any> => {
    const response = await apiClient.post('/api/etl/extract');
    return response.data;
  },

  transformData: async (): Promise<any> => {
    const response = await apiClient.post('/api/etl/transform');
    return response.data;
  },

  loadData: async (): Promise<any> => {
    const response = await apiClient.post('/api/etl/load');
    return response.data;
  },
};

export default apiClient;
