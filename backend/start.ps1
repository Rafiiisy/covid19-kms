# Set environment variables
$env:SERVER_PORT = "8000"
$env:SERVER_HOST = "localhost"
$env:ENV = "development"
$env:API_ENABLE_CORS = "true"

Write-Host "🚀 Starting COVID-19 KMS Backend Server" -ForegroundColor Green
Write-Host "📍 Port: $env:SERVER_PORT" -ForegroundColor Cyan
Write-Host "Host: $env:SERVER_HOST" -ForegroundColor Cyan
Write-Host "🔧 Environment: $env:ENV" -ForegroundColor Cyan
Write-Host "🔓 CORS: $env:API_ENABLE_CORS" -ForegroundColor Cyan
Write-Host ""

# Run the server
go run cmd/api/main.go
