#!/bin/bash

# Set environment variables
export SERVER_PORT=8000
export SERVER_HOST=localhost
export ENV=development
export API_ENABLE_CORS=true

echo "🚀 Starting COVID-19 KMS Backend Server"
echo "📍 Port: $SERVER_PORT"
echo "🌐 Host: $SERVER_HOST"
echo "🔧 Environment: $ENV"
echo "🔓 CORS: $API_ENABLE_CORS"
echo ""

# Run the server
go run cmd/api/main.go
