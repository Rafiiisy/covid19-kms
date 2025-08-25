@echo off
REM =============================================================================
REM COVID-19 KMS ETL Containers Startup Script (Windows)
REM =============================================================================

echo 🚀 Starting COVID-19 KMS ETL Containers

REM Check if Docker is running
docker info >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker is not running. Please start Docker first.
    pause
    exit /b 1
)

REM Load environment variables from env file
if exist env (
    for /f "tokens=1,* delims==" %%a in (env) do (
        if not "%%a"=="" if not "%%a:~0,1%"=="#" (
            set %%a=%%b
        )
    )
    echo 📋 Environment variables loaded from env file
) else (
    echo ⚠️ Warning: env file not found. Using default values.
    echo 💡 Copy env.example to env and customize for your environment.
)

REM Set default values if not in env file
if not defined HOST set HOST=0.0.0.0
if not defined PORT set PORT=8000
if not defined ENV set ENV=development

echo 📍 Port: %PORT%
echo 🌐 Host: %HOST%
echo 🔧 Environment: %ENV%
echo.

REM Create necessary directories
if not exist data\raw mkdir data\raw
if not exist data\processed mkdir data\processed
if not exist logs mkdir logs

REM Stop and remove existing ETL containers
echo 🔄 Stopping existing ETL containers...
docker-compose down

REM Build and start ETL containers
echo 🏗️ Building and starting ETL containers...
docker-compose up -d

REM Wait for containers to be ready
echo ⏳ Waiting for ETL containers to be ready...
timeout /t 15 >nul

REM Check container status
echo 🔍 Checking ETL container status...
docker-compose ps

echo.
echo ✅ ETL containers are running!
echo.
echo 📊 Container URLs:
echo 🌐 ETL API: http://localhost:8000
echo.
echo 🎯 To stop ETL containers: docker-compose down
echo 🔧 To view logs: docker-compose logs -f backend
echo 📋 To run ETL pipeline: curl -X POST http://localhost:8000/api/etl/run
echo.
echo 💡 Note: Make sure database is running separately at /database
echo.

pause
