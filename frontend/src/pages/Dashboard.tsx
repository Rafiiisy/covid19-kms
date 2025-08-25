import React, { useEffect, useState } from 'react';
import { RefreshCw, Download, BarChart3, TrendingUp, Globe, Video, Database, Activity, Clock } from 'lucide-react';
import { useETL } from '../hooks/useETL';
import { DataCharts } from '../components/DataCharts';
import { StatusIndicator } from '../components/StatusIndicator';
import { ErrorDisplay } from '../components/ErrorDisplay';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { RecordsPopup } from '../components/RecordsPopup';
import { exportAsJSON, exportAsCSV } from '../utils/exportUtils';

const Dashboard: React.FC = () => {
  const [showRecordsPopup, setShowRecordsPopup] = useState(false);
  
  const {
    isLoading,
    lastResult,
    pipelineStatus,
    healthStatus,
    error,
    lastUpdated,
    databaseData,
    runPipeline,
    getStatus,
    checkHealth,
    clearError,
    getMetrics,
    fetchAllDatabaseData,
  } = useETL();

  // Load initial status on component mount
  useEffect(() => {
    getStatus();
    checkHealth();
  }, [getStatus, checkHealth]);

  const handleRefreshData = async () => {
    try {
      console.log('🚀 Starting ETL pipeline...');
      
      // 1. Run ETL Pipeline
      await runPipeline();
      console.log('✅ ETL pipeline completed');
      
      // 2. Refresh status after pipeline completion
      await getStatus();
      await checkHealth();
      
      // 3. Fetch fresh data from database
      console.log('📊 Fetching updated database data...');
      await fetchAllDatabaseData();
      console.log('✅ Database data updated in UI');
      
    } catch (error) {
      console.error('❌ Failed to refresh data:', error);
    }
  };

  const metrics = getMetrics();

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Navigation */}
      <nav className="bg-blue-600 text-white shadow-lg">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <h1 className="text-xl font-bold">COVID-19 KMS</h1>
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-sm">Last Updated: {lastUpdated}</span>
              {pipelineStatus && (
                <StatusIndicator 
                  status={pipelineStatus.status} 
                  timestamp={pipelineStatus.timestamp}
                />
              )}
            </div>
          </div>
        </div>
        </nav>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h2 className="text-3xl font-bold text-gray-900">COVID-19 Knowledge Management Dashboard</h2>
          <p className="mt-2 text-gray-600">Real-time data analysis from multiple sources</p>
        </div>

        {/* Error Display */}
        <ErrorDisplay error={error} onClear={clearError} />

        {/* Control Panel */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-8">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Control Panel</h3>
          <div className="flex flex-wrap gap-4 items-center">
            <button
              onClick={handleRefreshData}
              disabled={isLoading}
              className="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white px-6 py-3 rounded-lg font-medium flex items-center space-x-2 transition-colors"
            >
              <RefreshCw className={`w-5 h-5 ${isLoading ? 'animate-spin' : ''}`} />
              <span>{isLoading ? 'Processing...' : 'Refresh Data'}</span>
            </button>
            
            <div className="flex space-x-2">
              <button 
                onClick={() => lastResult && exportAsJSON(lastResult)}
                disabled={!lastResult}
                className="bg-gray-100 hover:bg-gray-200 disabled:bg-gray-50 disabled:text-gray-400 text-gray-700 px-4 py-3 rounded-lg font-medium flex items-center space-x-2 transition-colors disabled:cursor-not-allowed"
              >
                <Download className="w-5 h-5" />
                <span>Export JSON</span>
              </button>
              <button 
                onClick={() => lastResult && exportAsCSV(lastResult)}
                disabled={!lastResult}
                className="bg-gray-100 hover:bg-gray-200 disabled:bg-gray-50 disabled:text-gray-400 text-gray-700 px-4 py-3 rounded-lg font-medium flex items-center space-x-2 transition-colors disabled:cursor-not-allowed"
              >
                <Download className="w-5 h-5" />
                <span>Export CSV</span>
              </button>
            </div>
          </div>
          
          {/* Pipeline Status */}
          {lastResult && (
            <div className="mt-4 pt-4 border-t border-gray-200">
              <div className="flex items-center space-x-4">
                <span className="text-sm font-medium text-gray-700">Pipeline Status:</span>
                <StatusIndicator 
                  status={lastResult.status} 
                  message={lastResult.message}
                  timestamp={lastResult.timestamp}
                />
                {lastResult.pipeline_duration && (
                  <span className="text-sm text-gray-600">
                    Duration: {lastResult.pipeline_duration}
                  </span>
                )}
              </div>
            </div>
          )}
        </div>

        {/* Statistics Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <div className="bg-gradient-to-br from-purple-500 to-purple-600 text-white rounded-lg p-6 cursor-pointer hover:from-purple-600 hover:to-purple-700 transition-all duration-200" onClick={() => setShowRecordsPopup(true)}>
            <div className="flex items-center">
              <Database className="w-8 h-8 mr-3" />
              <div>
                <p className="text-purple-100 text-sm">Total Records</p>
                <p className="text-2xl font-bold">{metrics?.totalRecords || 0}</p>
                <p className="text-purple-200 text-xs mt-1">Click to view details</p>
              </div>
            </div>
          </div>
          
          <div className="bg-gradient-to-br from-orange-500 to-orange-600 text-white rounded-lg p-6">
            <div className="flex items-center">
              <Activity className="w-8 h-8 mr-3" />
              <div>
                <p className="text-orange-100 text-sm">Pipeline Status</p>
                <p className="text-lg font-bold">
                  {metrics?.pipelineStatus ? metrics.pipelineStatus.charAt(0).toUpperCase() + metrics.pipelineStatus.slice(1) : '-'}
                </p>
              </div>
            </div>
          </div>
        </div>



        {/* Data Visualization */}
        {isLoading ? (
          <div className="bg-white rounded-lg shadow-md p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Processing ETL Pipeline</h3>
            <LoadingSpinner message="Running ETL Pipeline..." size="lg" />
          </div>
        ) : (
          <DataCharts etlResult={lastResult} databaseData={databaseData} />
        )}
      </div>

      {/* Records Popup Modal */}
      {showRecordsPopup && (
        <RecordsPopup 
          isOpen={showRecordsPopup}
          onClose={() => setShowRecordsPopup(false)}
          databaseData={databaseData}
        />
      )}
    </div>
  );
};

export default Dashboard;
