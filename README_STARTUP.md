# 🚀 IAM Services - Quick Start Guide

## 📋 Prerequisites

- ✅ Docker Desktop installed and running
- ✅ Go 1.21+ installed
- ✅ PostgreSQL and Redis (via Docker)

---

## ⚡ Quick Start (Automated)

### **Step 1: Start Services**

```powershell
cd C:\Users\tvttt\OneDrive\Desktop\go_workspace\my_project\ecommerce\back_end\iam-services

# Run startup script
.\start-services.ps1
```

**Choose option:**
- **Option 1**: Docker Compose (Full containerized)
- **Option 2**: Local Go ⭐ **RECOMMENDED** (Faster, easier debug)

---

### **Step 2: Verify Services**

In a **new PowerShell window**:

```powershell
cd C:\Users\tvttt\OneDrive\Desktop\go_workspace\my_project\ecommerce\back_end\iam-services

# Run verification script
.\verify-services.ps1
```

**Should see:**
- ✅ PostgreSQL is listening on port 5432
- ✅ Redis is listening on port 6379
- ✅ gRPC Server is listening on port 50051
- ✅ HTTP Server is listening on port 8080
- ✅ Health check passed!
- ✅ OPTIONS preflight passed!

---

## 🛠️ Manual Start (If Scripts Don't Work)

### **Terminal 1: Database**

```powershell
cd C:\Users\tvttt\OneDrive\Desktop\go_workspace\my_project\ecommerce\back_end\iam-services

# Start PostgreSQL and Redis
docker-compose up -d postgres redis

# Check status
docker-compose ps

# Should see both "Up"
```

---

### **Terminal 2: Backend (Choose One)**

#### **Option A: Local Go (Recommended)**

```powershell
cd C:\Users\tvttt\OneDrive\Desktop\go_workspace\my_project\ecommerce\back_end\iam-services

# Run backend
go run cmd/server/main.go
```

**Wait for:**
```
INFO  gRPC server is running      address=0.0.0.0:50051
INFO  Gin HTTP server is running  address=0.0.0.0:8080
                                   health=http://0.0.0.0:8080/health
```

---

#### **Option B: Docker Compose**

```powershell
cd C:\Users\tvttt\OneDrive\Desktop\go_workspace\my_project\ecommerce\back_end\iam-services

# Build and start
docker-compose build iam-service
docker-compose up -d iam-service

# View logs
docker-compose logs -f iam-service
```

---

## ✅ Quick Tests

### **Test 1: Health Check**

```powershell
curl http://localhost:8080/health
```

**Expected:**
```json
{"status":"ok"}
```

---

### **Test 2: OPTIONS Preflight (CORS)**

```powershell
curl -X OPTIONS http://localhost:8080/v1/auth/register `
  -H "Origin: http://localhost:3000" `
  -H "Access-Control-Request-Method: POST" `
  -v
```

**Expected:**
```
< HTTP/1.1 204 No Content
< Access-Control-Allow-Origin: *
< Access-Control-Allow-Methods: POST, OPTIONS, GET, PUT, DELETE, PATCH
```

---

### **Test 3: Register API**

```powershell
curl -X POST http://localhost:8080/v1/auth/register `
  -H "Content-Type: application/json" `
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test1234",
    "full_name": "Test User"
  }'
```

**Expected:**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "...",
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

---

## 🐛 Troubleshooting

### **Issue: Port already in use**

```powershell
# Find and kill process
netstat -ano | findstr ":8080"
taskkill /PID <PID> /F
```

---

### **Issue: Docker not running**

1. Open Docker Desktop
2. Wait for it to fully start
3. Verify: `docker ps`

---

### **Issue: Database connection failed**

```powershell
# Restart database
docker-compose restart postgres
Start-Sleep -Seconds 5
go run cmd/server/main.go
```

---

### **Issue: Container keeps restarting**

```powershell
# View logs
docker-compose logs iam-service

# Rebuild
docker-compose down
docker-compose build --no-cache iam-service
docker-compose up -d
```

---

## 📊 Services Overview

| Service | Port | Check Command |
|---------|------|---------------|
| PostgreSQL | 5432 | `docker-compose ps postgres` |
| Redis | 6379 | `docker-compose ps redis` |
| gRPC Server | 50051 | `netstat -ano \| findstr :50051` |
| HTTP Server | 8080 | `curl http://localhost:8080/health` |

---

## 🎯 Next Steps

After services are running:

1. **Start Frontend:**
   ```powershell
   cd C:\Users\tvttt\OneDrive\Desktop\go_workspace\my_project\ecommerce\web\web_shop
   npm run dev
   ```

2. **Test Registration:**
   - Open: http://localhost:3000/auth/signup
   - Fill the form
   - Submit

3. **Check Logs:**
   - Backend logs show the request
   - Frontend Network tab shows 200 OK

---

## 💡 Tips

### **Development Mode**
- Use **Local Go** (Option 2) for faster iteration
- Hot reload with `air` or `reflex`
- Easy to add breakpoints

### **Production-like Testing**
- Use **Docker Compose** (Option 1)
- Tests containerization
- Closer to production setup

### **Quick Restart**
```powershell
# Stop backend (Ctrl+C)
# Start again
go run cmd/server/main.go
```

---

## 📞 Support

If issues persist:

1. Run `.\verify-services.ps1` to diagnose
2. Check logs: `docker-compose logs iam-service`
3. Verify Docker Desktop is running
4. Check firewall/antivirus not blocking ports

---

**Last Updated:** 2024-11-06

