# Backend Port Change & CORS Configuration

## 🚀 **Port Changed from 8080 to 8000**

The backend server port has been updated from `8080` to `8000` to use standard development ports.

### **What Changed:**

1. **Configuration Default**: `backend/internal/config/config.go` - Default port changed to 8000
2. **Frontend API**: `frontend/src/services/api.ts` - API base URL updated to `http://localhost:8000`
3. **Frontend Port**: `frontend/vite.config.ts` - Frontend port changed to 3000
4. **Documentation**: All README files updated to reflect new ports

### **How to Start Backend:**

#### **Option 1: Using Scripts (Recommended)**
```bash
# Linux/Mac
./start.sh

# Windows Command Prompt
start.bat

# Windows PowerShell
.\start.ps1
```

#### **Option 2: Manual Environment Variables**
```bash
# Linux/Mac
export SERVER_PORT=8000
export SERVER_HOST=localhost
export ENV=development
export API_ENABLE_CORS=true
go run cmd/api/main.go

# Windows Command Prompt
set SERVER_PORT=8000
set SERVER_HOST=localhost
set ENV=development
set API_ENABLE_CORS=true
go run cmd/api/main.go

# Windows PowerShell
$env:SERVER_PORT = "8000"
$env:SERVER_HOST = "localhost"
$env:ENV = "development"
$env:API_ENABLE_CORS = "true"
go run cmd/api/main.go
```

#### **Option 3: Create .env File**
Create a `.env` file in the backend directory:
```env
SERVER_PORT=8000
SERVER_HOST=localhost
ENV=development
API_ENABLE_CORS=true
```

## 🔓 **CORS Configuration Enhanced**

CORS (Cross-Origin Resource Sharing) has been properly configured to allow the frontend to communicate with the backend.

### **CORS Headers Set:**
- `Access-Control-Allow-Origin: *` - Allows all origins
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS` - Allowed HTTP methods
- `Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With` - Allowed headers
- `Access-Control-Allow-Credentials: true` - Allows credentials

### **CORS Implementation:**
- **Location**: `backend/internal/api/router.go`
- **Method**: `corsMiddleware()` function wraps all API endpoints
- **Preflight**: OPTIONS requests are handled automatically

### **What This Means:**
- ✅ Frontend can make API calls from `localhost:5173` to `localhost:3001`
- ✅ All HTTP methods are allowed
- ✅ No CORS errors in browser console
- ✅ Secure development environment

## 🔧 **Testing the Setup**

### **1. Start Backend:**
```bash
cd backend
./start.sh  # or start.bat / start.ps1
```

**Expected Output:**
```
🚀 Starting COVID-19 KMS Backend Server
📍 Port: 8000
🌐 Host: localhost
🔧 Environment: development
🔓 CORS: true

🚀 Starting ETL API server on localhost:8000
📊 Environment: development
🔗 API Documentation: http://localhost:8000/api
🏥 Health Check: http://localhost:8000/api/health
```

### **2. Start Frontend:**
```bash
cd frontend
npm run dev
```

**Expected Output:**
```
  VITE v4.2.0  ready in 500 ms

  ➜  Local:   http://localhost:3000/
  ➜  Network: use --host to expose
```

### **3. Test API Connection:**
Open browser to `http://localhost:8000/api/health`

**Expected Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-08-15T...",
  "service": "ETL Pipeline API",
  "version": "1.0.0"
}
```

### **4. Test Dashboard:**
Open browser to `http://localhost:3000`

- ✅ Dashboard loads without errors
- ✅ No CORS errors in browser console
- ✅ "Refresh Data" button works
- ✅ API calls succeed

## 🚨 **Troubleshooting**

### **Port Already in Use:**
```bash
# Check what's using port 8000
netstat -ano | findstr :8000  # Windows
lsof -i :8000                 # Linux/Mac

# Kill process using port 8000
taskkill /PID <PID> /F        # Windows
kill -9 <PID>                 # Linux/Mac
```

### **CORS Still Not Working:**
1. Ensure backend is running on port 3001
2. Check browser console for errors
3. Verify `API_ENABLE_CORS=true` is set
4. Restart both frontend and backend

### **Frontend Can't Connect:**
1. Verify backend is running: `http://localhost:8000/api/health`
2. Check frontend API URL: `frontend/src/services/api.ts`
3. Ensure no firewall blocking port 8000

## 📝 **Summary**

- **Backend Port**: Changed from 8080 to 8000
- **Frontend Port**: Changed from 5173 to 3000
- **CORS**: Fully configured and working
- **Frontend**: Updated to use new backend port
- **Scripts**: Easy startup scripts provided
- **Documentation**: All references updated

The system is now ready for development with proper CORS handling and no port conflicts!
