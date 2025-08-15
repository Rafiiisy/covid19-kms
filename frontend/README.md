# COVID-19 KMS Frontend Dashboard

A React-based dashboard for monitoring and controlling the COVID-19 Knowledge Management System ETL pipeline.

## Features

### 🎯 **Real-time ETL Monitoring**
- Live pipeline status and health monitoring
- Real-time data processing metrics
- Pipeline execution history

### 📊 **Data Visualization**
- Interactive charts showing data source distribution
- Pipeline performance metrics
- Data processing statistics

### 🔄 **Manual ETL Control**
- **Refresh Button**: Manually trigger the complete ETL pipeline
- Individual stage execution (extract, transform, load)
- Real-time progress tracking

### 📈 **Comprehensive Metrics**
- YouTube videos count
- News articles count
- Total records processed
- Pipeline execution duration
- Average data relevance scores
- Extraction source counts

### 💾 **Data Export**
- Export ETL results as JSON
- Export ETL results as CSV
- Automatic filename generation with timestamps

## Architecture

### Components
- **Dashboard**: Main dashboard interface
- **DataCharts**: Chart.js-based data visualization
- **StatusIndicator**: Pipeline status display
- **ErrorDisplay**: Error handling and display
- **LoadingSpinner**: Loading states during ETL execution

### Hooks
- **useETL**: Custom hook for ETL state management and API calls

### Services
- **API Service**: Axios-based HTTP client for backend communication
- **Export Utils**: Data export functionality

## API Integration

The dashboard integrates with the Go backend ETL APIs:

- `POST /api/etl/run` - Execute complete ETL pipeline
- `GET /api/etl/status` - Get pipeline status
- `GET /api/health` - Health check
- `POST /api/etl/extract` - Run extraction stage
- `POST /api/etl/transform` - Run transformation stage
- `POST /api/etl/load` - Run loading stage

## Setup & Installation

### Prerequisites
- Node.js 16+ 
- Go backend running on `localhost:8000`

### Installation
```bash
cd frontend
npm install
```

### Development
```bash
npm run dev
```

### Build
```bash
npm run build
```

## Usage

### 1. **Start the Dashboard**
- Ensure the Go backend is running on port 8000
- Start the frontend with `npm run dev`
- Open `http://localhost:3000` in your browser

### 2. **Monitor Pipeline Status**
- View current pipeline status in the navigation bar
- Check system health status
- Monitor real-time metrics

### 3. **Run ETL Pipeline**
- Click the **"Refresh Data"** button to execute the complete ETL pipeline
- Watch real-time progress with the loading spinner
- View execution results and metrics

### 4. **Analyze Data**
- View interactive charts showing data distribution
- Analyze pipeline performance metrics
- Export results for further analysis

### 5. **Export Data**
- Use **Export JSON** for raw data export
- Use **Export CSV** for spreadsheet analysis
- Files are automatically named with timestamps

## Data Flow

1. **User clicks "Refresh Data"**
2. **Frontend calls** `POST /api/etl/run`
3. **Backend executes** complete ETL pipeline
4. **Results returned** to frontend
5. **Dashboard updates** with new metrics and charts
6. **Data visualization** reflects latest pipeline execution

## Error Handling

- **Network errors**: Automatic retry and user notification
- **Pipeline failures**: Detailed error display with clear messages
- **API timeouts**: 30-second timeout for ETL operations
- **User feedback**: Clear status indicators and error messages

## Responsive Design

- **Mobile-friendly**: Optimized for all screen sizes
- **Tailwind CSS**: Modern, responsive UI components
- **Accessibility**: Screen reader friendly with proper ARIA labels

## Dependencies

- **React 18** - UI framework
- **TypeScript** - Type safety
- **Chart.js** - Data visualization
- **Tailwind CSS** - Styling
- **Axios** - HTTP client
- **Lucide React** - Icons

## Development Notes

- **State Management**: Uses React hooks for local state
- **API Calls**: Centralized in service layer
- **Error Boundaries**: Comprehensive error handling
- **Performance**: Optimized re-renders and API calls
- **Testing**: Ready for unit and integration tests

## Troubleshooting

### Common Issues

1. **Backend Connection Failed**
   - Ensure Go backend is running on port 8000
   - Check CORS configuration in backend

2. **Charts Not Displaying**
   - Verify Chart.js dependencies are installed
   - Check browser console for JavaScript errors

3. **Export Not Working**
   - Ensure ETL pipeline has been executed
   - Check browser download permissions

4. **Slow Performance**
   - Monitor backend ETL execution time
   - Check network latency between frontend and backend

## Contributing

1. Follow TypeScript best practices
2. Use functional components with hooks
3. Maintain responsive design principles
4. Add proper error handling
5. Update documentation for new features
