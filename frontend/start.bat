@echo off
REM =============================================================================
REM COVID-19 KMS Frontend Startup Script (Windows)
REM =============================================================================

echo 🚀 Starting COVID-19 KMS Frontend
echo.

REM Check if Node.js is installed
echo 📋 Checking Node.js...
node --version
if errorlevel 1 (
    echo ❌ Node.js is not installed or not in PATH
    echo 📥 Download from: https://nodejs.org/
    pause
    exit /b 1
)

REM Check if npm is installed
echo 📋 Checking npm...
npm --version
if errorlevel 1 (
    echo ❌ npm is not installed or not in PATH
    echo 📥 Install Node.js to get npm
    pause
    exit /b 1
)

echo ✅ Node.js and npm are available
echo.

REM Check if package.json exists
if not exist package.json (
    echo ❌ package.json not found
    echo 💡 Make sure you're in the frontend directory
    pause
    exit /b 1
)

echo ✅ package.json found
echo.

REM Check if node_modules exists
if not exist node_modules (
    echo 📦 Installing dependencies...
    npm install
    if errorlevel 1 (
        echo ❌ Failed to install dependencies
        pause
        exit /b 1
    )
    echo ✅ Dependencies installed successfully
) else (
    echo ✅ Dependencies already installed
)

echo.
echo 🚀 Starting development server...
echo 🌐 Frontend will be available at: http://localhost:5173
echo 🔗 Backend API should be at: http://localhost:8000
echo.
echo 💡 Press Ctrl+C to stop the server
echo.

REM Start the development server
npm run dev

echo.
echo ⚠️ Server stopped or encountered an error
pause
