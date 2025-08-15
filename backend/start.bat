@echo off

REM Set environment variables
set SERVER_PORT=8000
set SERVER_HOST=localhost
set ENV=development
set API_ENABLE_CORS=true

echo 🚀 Starting COVID-19 KMS Backend Server
echo 📍 Port: %SERVER_PORT%
echo 🌐 Host: %SERVER_HOST%
echo 🔧 Environment: %ENV%
echo 🔓 CORS: %API_ENABLE_CORS%
echo.

REM Run the server
go run cmd/api/main.go

pause
