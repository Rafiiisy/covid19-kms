# COVID-19 KMS - New Project Structure

This document explains the new refactored project structure using Go backend + React frontend + GCP.

## 🏗️ **New Project Structure**

```
RepoCloud/
├── _python/                     # 🐍 OLD Python implementation (preserved)
│   ├── app.py                   # FastAPI application
│   ├── etl/                     # ETL pipeline
│   ├── models/                  # Data models
│   ├── utils/                   # Utilities
│   ├── dashboard/               # HTML templates
│   ├── tests/                   # Test files
│   └── requirements.txt         # Python dependencies
├── backend/                     # 🚀 NEW Go backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # Go application entry point
│   ├── internal/
│   │   ├── api/                # HTTP handlers & routes
│   │   ├── etl/                # ETL pipeline
│   │   ├── models/             # Data structures
│   │   ├── services/           # Business logic
│   │   └── config/             # Configuration
│   ├── go.mod                  # Go module dependencies
│   └── go.sum                  # Go module checksums
├── frontend/                    # ⚛️ NEW React frontend
│   ├── src/
│   │   ├── components/         # React components
│   │   ├── pages/              # Page views
│   │   ├── services/           # API services
│   │   ├── hooks/              # Custom React hooks
│   │   ├── types/              # TypeScript types
│   │   ├── utils/              # Utility functions
│   │   ├── App.tsx             # Main application component
│   │   └── index.tsx           # Application entry point
│   ├── public/                 # Static assets
│   ├── package.json            # Node.js dependencies
│   ├── tsconfig.json           # TypeScript configuration
│   ├── tailwind.config.js      # Tailwind CSS configuration
│   ├── postcss.config.js       # PostCSS configuration
│   └── vite.config.ts          # Vite build configuration
├── infrastructure/              # ☁️ GCP deployment configuration
│   ├── terraform/              # Infrastructure as Code
│   ├── cloudbuild/             # CI/CD pipeline configuration
│   └── kubernetes/             # K8s manifests (optional)
├── shared/                      # 🔗 Shared types and utilities
│   └── types/                  # Common type definitions
├── data/                        # 📊 Local data storage
│   ├── raw/                     # Raw data
│   └── processed/               # Processed data
├── README.md                    # Main project documentation
├── README_NEW_STRUCTURE.md      # This file
└── .gitignore                   # Git ignore rules
```

## 🔄 **Migration Status**

### ✅ **Completed:**
- **Directory Structure**: Created new Go + React structure
- **Python Preservation**: Moved all Python files to `_python/` folder
- **Go Backend**: Basic Go server with Gin framework
- **React Frontend**: Modern React app with TypeScript + Tailwind CSS
- **Configuration**: Vite, Tailwind, PostCSS setup

### 🚧 **In Progress:**
- **ETL Pipeline**: Porting from Python to Go
- **API Integration**: Connecting React frontend to Go backend
- **Data Models**: Converting Python models to Go structs

### 📋 **Next Steps:**
1. **Port ETL Logic**: Convert Python ETL to Go
2. **API Development**: Build comprehensive Go API endpoints
3. **Frontend Features**: Add charts, data tables, real-time updates
4. **GCP Integration**: Set up BigQuery, Cloud Storage, Cloud Run
5. **Testing**: Unit tests for both backend and frontend

## 🚀 **How to Run the New Stack**

### **Backend (Go):**
```bash
cd backend
go mod download
go run cmd/server/main.go
# Server runs on http://localhost:8080
```

### **Frontend (React):**
```bash
cd frontend
npm install
npm run dev
# App runs on http://localhost:5173
```

## 🔍 **What's in Each Directory**

### **`_python/` - Old Implementation**
- **Preserved for reference** during migration
- **Can be deleted** once Go implementation is complete
- **Useful for understanding** business logic and data flow

### **`backend/` - Go Backend**
- **High-performance** HTTP server with Gin
- **Concurrent ETL** processing with goroutines
- **GCP integration** for BigQuery and Cloud Storage
- **RESTful API** endpoints for frontend consumption

### **`frontend/` - React Frontend**
- **Modern UI** with TypeScript and Tailwind CSS
- **Real-time updates** and interactive charts
- **Responsive design** for desktop and mobile
- **API integration** with Go backend

### **`infrastructure/` - GCP Deployment**
- **Terraform** for infrastructure as code
- **Cloud Build** for CI/CD automation
- **Kubernetes** manifests (optional)

### **`shared/` - Common Utilities**
- **Type definitions** shared between frontend and backend
- **Utility functions** used across the stack

## 🎯 **Benefits of New Structure**

1. **Performance**: Go's concurrent processing for ETL
2. **Scalability**: GCP's serverless architecture
3. **Developer Experience**: Modern React + TypeScript
4. **Maintainability**: Clear separation of concerns
5. **Future-proof**: Industry-standard technologies

## 📚 **Documentation**

- **Main README**: Updated with new tech stack
- **This File**: Explains the new structure
- **Code Comments**: Inline documentation in Go and React
- **API Docs**: Generated from Go code

---

**Note**: This is a transitional structure. The `_python/` folder will be removed once the Go + React implementation is complete and tested.
