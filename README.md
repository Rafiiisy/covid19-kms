# COVID-19 Knowledge Management System

A comprehensive knowledge management system for COVID-19 data analysis and visualization, built with Go backend, React frontend, and Google Cloud Platform. This system integrates data from multiple RapidAPI sources to provide insights and analytics for COVID-19 research.

## 🎯 Project Overview

This system is designed for big data research on COVID-19, providing:
- **Multi-source data integration** from various RapidAPI endpoints
- **High-performance ETL pipeline** built with Go for concurrent data processing
- **Modern React dashboard** for real-time data visualization and analysis
- **Google Cloud Platform** deployment with scalable infrastructure

## 📊 Data Sources

The system integrates data from the following RapidAPI sources:

### Primary Data Sources
1. **YouTube API** - COVID-19 related video content and statistics
   - **Endpoint**: YouTube Data API v3 via RapidAPI
   - **Data**: Video metadata, views, duration, comments, likes
   - **Query**: COVID-19 related hashtags and keywords

2. **Google News API** - Latest COVID-19 news articles and updates
   - **Endpoint**: Google News Search via RapidAPI
   - **Data**: Article headlines, content, publication dates, sources
   - **Query**: COVID-19 news from global sources

3. **Instagram API** - COVID-19 related social media posts
   - **Endpoint**: Instagram Hashtag Media via RapidAPI
   - **Data**: Post captions, engagement metrics, hashtags
   - **Query**: COVID-19 related hashtags and content

4. **Indonesia News API** - Indonesia-specific COVID-19 news
   - **Endpoint**: Indonesian News Sources via RapidAPI
   - **Data**: News from CNN Indonesia, Detik, Kompas, Tempo
   - **Query**: COVID-19 news from Indonesian media outlets

### Data Processing
- **Real-time extraction** from all sources via RapidAPI
- **Automatic data transformation** and cleaning
- **Sentiment analysis** for news and social media content
- **Relevance scoring** based on COVID-19 keywords
- **Metadata enrichment** with source-specific information

## 🏗️ System Architecture

```
RepoCloud/
├── README.md                    # Project documentation
├── backend/                     # Go backend application
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # Go application entry point
│   ├── internal/
│   │   ├── api/                # HTTP handlers and routes
│   │   ├── etl/                # ETL pipeline (extract, transform, load)
│   │   ├── models/             # Data structures and schemas
│   │   ├── services/           # Business logic layer
│   │   └── config/             # Configuration management
│   ├── go.mod                  # Go module dependencies
│   ├── go.sum                  # Go module checksums
│   └── Dockerfile              # Container configuration
├── frontend/                    # React frontend application
│   ├── src/
│   │   ├── components/         # Reusable React components
│   │   ├── pages/              # Page-level components
│   │   ├── services/           # API integration services
│   │   ├── hooks/              # Custom React hooks
│   │   ├── types/              # TypeScript type definitions
│   │   ├── utils/              # Utility functions
│   │   ├── App.tsx             # Main application component
│   │   └── index.tsx           # Application entry point
│   ├── public/                 # Static assets
│   ├── package.json            # Node.js dependencies
│   ├── tsconfig.json           # TypeScript configuration
│   ├── tailwind.config.js      # Tailwind CSS configuration
│   └── vite.config.ts          # Vite build configuration
├── infrastructure/              # GCP deployment configuration
│   ├── terraform/              # Infrastructure as Code
│   ├── cloudbuild/             # CI/CD pipeline configuration
│   └── kubernetes/             # K8s manifests (optional)
├── shared/                      # Shared types and utilities
│   └── types/                  # Common type definitions
└── data/                        # Local data storage
    ├── raw/                     # Raw data
    └── processed/               # Processed data
```

## 🚀 Features

### ETL Pipeline (Go Backend)
- **High-performance concurrent processing** using Go goroutines
- **On-demand execution** via API endpoint (no scheduling required)
- **Multi-source data extraction** from RapidAPI endpoints
- **Data transformation and cleaning** with efficient Go processing
- **BigQuery integration** for data warehousing
- **Cloud Storage backup** for raw data

### Dashboard (React Frontend)
- **Real-time data visualization** with interactive charts
- **COVID-19 metrics tracking** and trend analysis
- **News sentiment analysis** and categorization
- **Geographic data visualization** for Indonesia
- **Export capabilities** for research purposes
- **Responsive design** for desktop and mobile

### API Integration
- **RESTful API endpoints** built with Go for high performance
- **WebSocket support** for real-time data updates
- **Authentication and authorization** for secure access
- **Rate limiting** to comply with API quotas
- **Error handling and retry mechanisms**

## 🛠️ Technology Stack

### Backend
- **Go 1.21+** - High-performance programming language
- **Gin/Echo** - Fast HTTP web framework
- **Goroutines** - Concurrent ETL processing
- **GORM** - Database ORM for data models
- **JWT** - Authentication and authorization

### Frontend
- **React 18+** - Modern UI library
- **TypeScript** - Type-safe development
- **Vite** - Fast build tool and dev server
- **Tailwind CSS** - Utility-first CSS framework
- **Chart.js/Recharts** - Data visualization
- **React Query** - Server state management

### Cloud Infrastructure
- **Google Cloud Platform (GCP)**
  - **Cloud Run** - Serverless Go backend hosting
  - **BigQuery** - Data warehouse and analytics
  - **Cloud Storage** - Object storage for raw data
  - **Cloud Build** - CI/CD pipeline automation
  - **Cloud Monitoring** - Performance monitoring

## 📋 Prerequisites

Before setting up the project, ensure you have:

1. **Go 1.21+** installed on your system
2. **Node.js 18+** and npm installed
3. **Google Cloud Platform Account** with billing enabled
4. **RapidAPI Account** with access to required APIs:
   - **YouTube Data API v3** - For video content and statistics
   - **Google News Search API** - For global news articles
   - **Instagram Hashtag Media API** - For social media posts
   - **Indonesia News API** - For Indonesian media sources

## 🔧 Installation & Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd RepoCloud
```

### 2. Set Up Backend (Go)
```bash
cd backend

# Install Go dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env file with your API keys and GCP settings

# Run the backend
go run cmd/api/main.go
```

### 3. Set Up Frontend (React)
```bash
cd frontend

# Install Node.js dependencies
npm install

# Set up environment variables
cp .env.example .env
# Edit .env file with your backend API URL

# Run the frontend
npm run dev
```

### 4. Set Up Google Cloud
```bash
# Authenticate with Google Cloud
gcloud auth login
gcloud config set project YOUR_PROJECT_ID

# Enable required APIs
gcloud services enable run.googleapis.com
gcloud services enable bigquery.googleapis.com
gcloud services enable storage.googleapis.com
gcloud services enable cloudbuild.googleapis.com
```

### 5. Run the Complete System
```bash
# Terminal 1: Backend
cd backend
go run cmd/api/main.go

# Terminal 2: Frontend
cd frontend
npm run dev
```

## 🎮 Usage

### Starting the ETL Process
1. Access the dashboard at `http://localhost:5173` (React dev server)
2. Click the "Refresh Data" button to trigger the ETL pipeline
3. Monitor the progress in the dashboard
4. View processed data in the analytics section

### API Endpoints
- `GET /api/health` - Health check
- `POST /api/etl/run` - Trigger ETL pipeline
- `GET /api/etl/status` - Get pipeline status
- `GET /api/etl/data/youtube` - YouTube data with metadata
- `GET /api/etl/data/google-news` - Google News data
- `GET /api/etl/data/instagram` - Instagram data with engagement metrics
- `GET /api/etl/data/indonesia-news` - Indonesia News data
- `GET /api/etl/data/summary` - Overall data summary
- `GET /api/etl/data/source` - Data by source
- `GET /api/etl/data/stats` - Data statistics
- `GET /api/etl/data/sentiment-distribution` - Sentiment analysis across all sources
- `GET /api/etl/data/word-frequency` - Word frequency analysis and trending topics
- `POST /api/etl/cleanup/sentiment` - Cleanup and recalculate sentiment scores

## 📊 Dashboard Features

### Data Visualization
- **Multi-Source Data Dashboard** - Unified view of all data sources
- **Source Distribution Charts** - Pie charts showing data from YouTube, Google News, Instagram, Indonesia News
- **Sentiment Distribution Charts** - Stacked bar charts showing sentiment analysis across all sources
- **Word Cloud Visualization** - Interactive word frequency analysis with sentiment color coding
- **Content Metadata Display** - Rich metadata for videos, posts, and articles
- **Real-time Data Refresh** - On-demand ETL pipeline execution
- **Interactive Records Popup** - Detailed view of all processed data with tabbed interface

### Analytics
- **Data Source Analytics** - Comprehensive analysis across all sources
- **Content Performance Metrics** - Engagement data from YouTube and Instagram
- **News Sentiment Analysis** - Sentiment scoring for articles and posts
- **Relevance Scoring** - AI-powered relevance assessment for COVID-19 content
- **Source Comparison** - Cross-platform data analysis and insights
- **Real-time Metrics** - Live dashboard with current data counts and status

## 🔒 Security

- **API Key Management** - Secure storage of RapidAPI keys
- **JWT Authentication** - Secure user authentication
- **CORS Configuration** - Cross-origin resource sharing setup
- **Rate Limiting** - API usage throttling
- **Data Encryption** - At-rest and in-transit encryption

## 🧪 Testing

### Backend Testing
```bash
cd backend
# Run all tests
go test ./...

# Run specific tests
go test ./internal/api/...
go test ./internal/etl/...
```

### Frontend Testing
```bash
cd frontend
# Run all tests
npm test

# Run specific tests
npm test -- --testPathPattern=components
```

## 📈 Monitoring & Logging

- **Application Logs** - Structured logging with Cloud Logging
- **Performance Monitoring** - Cloud Monitoring integration
- **Error Tracking** - Automatic error reporting
- **Health Checks** - Application health monitoring

## 🚀 Deployment

### Google Cloud Run (Backend)
```bash
cd backend
# Deploy to Cloud Run
gcloud run deploy covid-kms-backend \
  --source . \
  --platform managed \
  --region asia-southeast1 \
  --allow-unauthenticated
```

### Google Cloud Storage (Frontend)
```bash
cd frontend
# Build the frontend
npm run build

# Deploy to Cloud Storage
gsutil -m cp -r dist/* gs://your-bucket-name/
```

### Local Development
```bash
# Backend
cd backend
go run cmd/api/main.go

# Frontend
cd frontend
npm run dev
```

## 📚 Documentation

This README contains all the essential information to get started. For additional details:
- Check the code comments for API usage
- Review the ETL pipeline in the `backend/internal/etl/` folder
- Examine the React components in `frontend/src/components/`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Team

- **Research Lead** - [Your Name]
- **Backend Developer** - [Team Member]
- **Frontend Developer** - [Team Member]
- **DevOps Engineer** - [Team Member]

## 📞 Support

For support and questions:
- Create an issue in the GitHub repository
- Contact the development team
- Check the documentation in the `docs/` folder

---

**Note**: This is a research project for big data analysis of COVID-19 data. Ensure compliance with data privacy regulations and API usage terms when using this system. 