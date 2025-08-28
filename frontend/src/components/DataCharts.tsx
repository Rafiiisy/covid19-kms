import React from 'react';
import { Chart as ChartJS, ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, Title } from 'chart.js';
import { Pie, Bar } from 'react-chartjs-2';
import { ETLResult } from '../services/api';
import { WordCloud } from './WordCloud';

ChartJS.register(ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, Title);

interface DataChartsProps {
  etlResult: ETLResult | null;
  databaseData?: {
    youtube: any[] | null;
    googleNews: any[] | null;
    instagram: any[] | null;
    indonesiaNews: any[] | null;
    summary: any | null;
    sentimentDistribution: any | null;
    wordFrequency: any | null;
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

  // Prepare data for stacked bar chart (sentiment distribution)
  const sentimentData = {
    labels: ['YouTube', 'Google News', 'Instagram', 'Indonesia News'],
    datasets: [
      {
        label: 'Positive',
        data: useDatabaseData && databaseData.sentimentDistribution ? [
          databaseData.sentimentDistribution.sources?.youtube?.positive || 0,
          databaseData.sentimentDistribution.sources?.google_news?.positive || 0,
          databaseData.sentimentDistribution.sources?.instagram?.positive || 0,
          databaseData.sentimentDistribution.sources?.indonesia_news?.positive || 0,
        ] : [0, 0, 0, 0],
        backgroundColor: 'rgba(34, 197, 94, 0.8)',
        borderColor: 'rgba(34, 197, 94, 1)',
        borderWidth: 1,
      },
      {
        label: 'Negative',
        data: useDatabaseData && databaseData.sentimentDistribution ? [
          databaseData.sentimentDistribution.sources?.youtube?.negative || 0,
          databaseData.sentimentDistribution.sources?.google_news?.negative || 0,
          databaseData.sentimentDistribution.sources?.instagram?.negative || 0,
          databaseData.sentimentDistribution.sources?.indonesia_news?.negative || 0,
        ] : [0, 0, 0, 0],
        backgroundColor: 'rgba(239, 68, 68, 0.8)',
        borderColor: 'rgba(239, 68, 68, 1)',
        borderWidth: 1,
      },
      {
        label: 'Neutral',
        data: useDatabaseData && databaseData.sentimentDistribution ? [
          databaseData.sentimentDistribution.sources?.youtube?.neutral || 0,
          databaseData.sentimentDistribution.sources?.google_news?.neutral || 0,
          databaseData.sentimentDistribution.sources?.instagram?.neutral || 0,
          databaseData.sentimentDistribution.sources?.indonesia_news?.neutral || 0,
        ] : [0, 0, 0, 0],
        backgroundColor: 'rgba(156, 163, 175, 0.8)',
        borderColor: 'rgba(156, 163, 175, 1)',
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

  const sentimentChartOptions = {
    ...chartOptions,
    plugins: {
      ...chartOptions.plugins,
      title: {
        display: true,
        text: 'Sentiment Distribution by Source',
      },
    },
    scales: {
      x: {
        stacked: true,
      },
      y: {
        stacked: true,
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

      {/* Sentiment Distribution */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Sentiment Distribution by Source</h3>
        <div className="h-64">
          <Bar data={sentimentData} options={sentimentChartOptions} />
        </div>
        {useDatabaseData && databaseData.sentimentDistribution && (
          <div className="mt-4 grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div className="text-center">
              <div className="text-green-600 font-semibold">
                {databaseData.sentimentDistribution.totals?.positive || 0}
              </div>
              <div className="text-gray-600">Total Positive</div>
            </div>
            <div className="text-center">
              <div className="text-red-600 font-semibold">
                {databaseData.sentimentDistribution.totals?.negative || 0}
              </div>
              <div className="text-gray-600">Total Negative</div>
            </div>
            <div className="text-center">
              <div className="text-gray-600 font-semibold">
                {databaseData.sentimentDistribution.totals?.neutral || 0}
              </div>
              <div className="text-gray-600">Total Neutral</div>
            </div>
            <div className="text-center">
              <div className="text-blue-600 font-semibold">
                {databaseData.sentimentDistribution.totals?.total || 0}
              </div>
              <div className="text-gray-600">Total Records</div>
            </div>
          </div>
        )}
      </div>



      {/* Word Cloud */}
      <WordCloud wordFrequency={databaseData?.wordFrequency} />
    </div>
  );
};

