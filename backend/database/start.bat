@echo off
REM =============================================================================
REM COVID-19 KMS Database Startup Script (Windows)
REM =============================================================================

echo 🚀 Starting COVID-19 KMS PostgreSQL Database

REM Check if Docker is running
docker info >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker is not running. Please start Docker first.
    pause
    exit /b 1
)

REM Create logs directory if it doesn't exist
if not exist logs mkdir logs

REM Stop and remove existing containers
echo 🔄 Stopping existing containers...
docker-compose down

REM Build and start the database
echo 🏗️ Building and starting database...
docker-compose up -d postgres

REM Wait for database to be ready
echo ⏳ Waiting for database to be ready...
timeout /t 10 >nul

REM Check if database is ready
docker-compose exec postgres pg_isready -U postgres -d covid19_kms >nul 2>&1
if errorlevel 0 (
    echo ✅ Database is ready!
    echo 📊 Database URL: postgresql://postgres:password@localhost:5432/covid19_kms
    echo 🔗 Connection: Host=localhost, Port=5432, Database=covid19_kms, User=postgres
    
    REM Show database info
    echo.
    echo 📋 Database Information:
    docker-compose exec postgres psql -U postgres -d covid19_kms -c "\dt"
    
    echo.
    echo 🎯 To stop the database: docker-compose down
    echo 🔧 To start pgAdmin: docker-compose --profile admin up -d pgadmin
    echo 🌐 pgAdmin URL: http://localhost:8080 (admin@covid19kms.com / admin123)
    echo.
    echo 📋 Container logs:
    docker-compose logs --tail=5 postgres
) else (
    echo ❌ Database failed to start properly
    echo 📋 Container logs:
    docker-compose logs postgres
    pause
    exit /b 1
)

pause
