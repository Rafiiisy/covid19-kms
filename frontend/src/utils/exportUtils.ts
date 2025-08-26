import { ETLResult } from '../services/api';

// Export ETL result as JSON
export const exportAsJSON = (data: ETLResult, filename?: string) => {
  const jsonString = JSON.stringify(data, null, 2);
  const blob = new Blob([jsonString], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  
  const link = document.createElement('a');
  link.href = url;
  link.download = filename || `etl-result-${new Date().toISOString().split('T')[0]}.json`;
  
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  
  URL.revokeObjectURL(url);
};

// Export ETL result as CSV
export const exportAsCSV = (data: ETLResult, filename?: string) => {
  // Convert ETL result to CSV format
  const csvRows = [];
  
  // Header row
  csvRows.push(['Metric', 'Value']);
  
  // Basic info
  csvRows.push(['Status', data.status]);
  csvRows.push(['Message', data.message]);
  csvRows.push(['Timestamp', data.timestamp]);
  csvRows.push(['Pipeline Duration', data.pipeline_duration]);
  
  // Extraction data
  if (data.extraction) {
    csvRows.push(['', '']); // Empty row for separation
    csvRows.push(['Extraction', '']);
    csvRows.push(['Query', data.extraction.query]);
    csvRows.push(['Sources Count', Object.keys(data.extraction.sources).length]);
    csvRows.push(['Sources', Object.keys(data.extraction.sources).join(', ')]);
  }
  
  // Transformation data
  if (data.transformation) {
    csvRows.push(['', '']); // Empty row for separation
    csvRows.push(['Transformation', '']);
    csvRows.push(['YouTube Videos', data.transformation.YouTube?.length || 0]);
    csvRows.push(['News Articles', data.transformation.News?.length || 0]);
    csvRows.push(['Transformed At', data.transformation.TransformedAt]);
    if (data.transformation.Summary?.AverageRelevance) {
      csvRows.push(['Average Relevance', data.transformation.Summary.AverageRelevance]);
    }
  }
  
  // Loading data
  if (data.loading) {
    csvRows.push(['', '']); // Empty row for separation
    csvRows.push(['Loading', '']);
    csvRows.push(['Success', data.loading.success]);
    csvRows.push(['Records Count', data.loading.records_count]);
    csvRows.push(['Message', data.loading.message]);
    csvRows.push(['Timestamp', data.loading.timestamp]);
  }
  
  // Summary data
  if (data.summary) {
    csvRows.push(['', '']); // Empty row for separation
    csvRows.push(['Summary', '']);
    Object.entries(data.summary).forEach(([key, value]) => {
      csvRows.push([key, String(value)]);
    });
  }
  
  // Convert to CSV string
  const csvContent = csvRows.map(row => 
    row.map(cell => `"${String(cell).replace(/"/g, '""')}"`).join(',')
  ).join('\n');
  
  // Create and download file
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  
  const link = document.createElement('a');
  link.href = url;
  link.download = filename || `etl-result-${new Date().toISOString().split('T')[0]}.csv`;
  
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  
  URL.revokeObjectURL(url);
};



