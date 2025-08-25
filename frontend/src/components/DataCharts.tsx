import React from 'react';
import { Chart as ChartJS, ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, Title } from 'chart.js';
import { Pie, Bar } from 'react-chartjs-2';
import { ETLResult } from '../services/api';

ChartJS.register(ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, Title);

interface DataChartsProps {
  etlResult: ETLResult | null;
  databaseData?: {
    youtube: any[] | null;
    googleNews: any[] | null;
    instagram: any[] | null;
    indonesiaNews: any[] | null;
    summary: any | null;
  } | null;
}

export const DataCharts: React.FC<DataChartsProps> = ({ etlResult, databaseData }) => {
  // Use database data if available, otherwise show ETL result message
  if (!etlResult && !databaseData) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Data Visualization</h3>
        <div className="h-64 bg-gray-100 rounded-lg flex items-center justify-center">
          <p className="text-gray-500">Run ETL pipeline to see data visualization</p>
        </div>
      </div>
    );
  }

  // Use database data for charts if available
  const useDatabaseData = databaseData && databaseData.summary;

  // Prepare data for pie chart (data source distribution)
  const sourceData = {
    labels: ['YouTube Videos', 'Google News', 'Instagram', 'Indonesia News'],
    datasets: [
      {
        data: useDatabaseData ? [
          databaseData.youtube?.length || 0,
          databaseData.googleNews?.length || 0,
          databaseData.instagram?.length || 0,
          databaseData.indonesiaNews?.length || 0,
        ] : [
          etlResult?.transformation?.YouTube?.length || 0,
          etlResult?.transformation?.News?.length || 0,
        ],
        backgroundColor: [
          'rgba(255, 99, 132, 0.8)',
          'rgba(54, 162, 235, 0.8)',
          'rgba(153, 102, 255, 0.8)',
          'rgba(255, 206, 86, 0.8)',
        ],
        borderColor: [
          'rgba(255, 99, 132, 1)',
          'rgba(54, 162, 235, 1)',
          'rgba(153, 102, 255, 1)',
          'rgba(255, 206, 86, 1)',
        ],
        borderWidth: 1,
      },
    ],
  };

  // Prepare data for bar chart (pipeline stages performance)
  const pipelineData = {
    labels: ['YouTube', 'Google News', 'Instagram', 'Indonesia News'],
    datasets: [
      {
        label: 'Records Processed',
        data: useDatabaseData ? [
          databaseData.youtube?.length || 0,
          databaseData.googleNews?.length || 0,
          databaseData.instagram?.length || 0,
          databaseData.indonesiaNews?.length || 0,
        ] : [
          etlResult?.extraction?.sources ? Object.keys(etlResult.extraction.sources).length : 0,
          (etlResult?.transformation?.YouTube?.length || 0) + (etlResult?.transformation?.News?.length || 0),
          etlResult?.loading?.records_count || 0,
        ],
        backgroundColor: [
          'rgba(255, 99, 132, 0.8)',
          'rgba(255, 99, 132, 0.8)',
          'rgba(54, 162, 235, 0.8)',
          'rgba(153, 102, 255, 0.8)',
        ],
        borderColor: [
          'rgba(255, 99, 132, 1)',
          'rgba(255, 99, 132, 1)',
          'rgba(54, 162, 235, 1)',
          'rgba(153, 102, 255, 1)',
        ],
        borderWidth: 1,
      },
    ],
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'bottom' as const,
      },
    },
  };

  const barChartOptions = {
    ...chartOptions,
    plugins: {
      ...chartOptions.plugins,
      title: {
        display: true,
        text: 'Data Source Distribution',
      },
    },
    scales: {
      y: {
        beginAtZero: true,
      },
    },
  };

  return (
    <div className="space-y-6">
      {/* Data Source Distribution */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Data Source Distribution</h3>
        <div className="h-64">
          <Pie data={sourceData} options={chartOptions} />
        </div>
      </div>

      {/* Pipeline Performance */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Pipeline Stages Performance</h3>
        <div className="h-64">
          <Bar data={pipelineData} options={barChartOptions} />
        </div>
      </div>

      {/* Pipeline Summary */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Pipeline Summary</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="bg-gray-50 p-4 rounded-lg">
            <h4 className="font-medium text-gray-700 mb-2">Execution Details</h4>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600">Status:</span>
                <span className={`font-medium ${
                  etlResult?.status === 'success' ? 'text-green-600' : 'text-red-600'
                }`}>
                  {etlResult?.status || 'Unknown'}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Duration:</span>
                <span className="font-medium">{etlResult?.pipeline_duration || 'N/A'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Timestamp:</span>
                <span className="font-medium">{etlResult?.timestamp ? new Date(etlResult.timestamp).toLocaleString() : 'N/A'}</span>
              </div>
            </div>
          </div>
          
          <div className="bg-gray-50 p-4 rounded-lg">
            <h4 className="font-medium text-gray-700 mb-2">Data Summary</h4>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600">YouTube Videos:</span>
                <span className="font-medium">{etlResult?.transformation?.YouTube?.length || 0}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">News Articles:</span>
                <span className="font-medium">{etlResult?.transformation?.News?.length || 0}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Total Records:</span>
                <span className="font-medium">{etlResult?.loading?.records_count || 0}</span>
              </div>
              {etlResult?.transformation?.Summary?.AverageRelevance && (
                <div className="flex justify-between">
                  <span className="text-gray-600">Avg Relevance:</span>
                  <span className="font-medium">
                    {etlResult.transformation.Summary.AverageRelevance.toFixed(2)}
                  </span>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

