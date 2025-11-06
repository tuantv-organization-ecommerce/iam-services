# 🔧 IAM Services Troubleshooting Guide

## 🔴 Common Errors and Solutions

### **Error 1: Redis RDB Format Version Error**

```
Can't handle RDB format version 12
Error reading the RDB base file appendonly.aof.1.base.rdb
```

**Cause:** Redis data from incompatible version or corrupted data file.

**Solution:**

```powershell
# Quick fix - Reset Redis data
.\fix-redis.ps1

# Or manually:
docker-compose down -v
docker-compose up -d postgres redis
```

**⚠️ Warning:** This will delete all cached data and require users to re-login.

---

### **Error 2: Socket Hang Up / Connection Refused**

```
Error: connect ECONNREFUSED 127.0.0.1:8080
Error: socket hang up
```

**Cause:** Backend service is not running or crashed.

**Check:**
```powershell
# Check if port 8080 is listening
netstat -ano | findstr ":8080"

# Check container status
docker-compose ps
```

**Solution:**
```powershell
# Start backend
.\start-services.ps1

# Or manually:
go run cmd/server/main.go
```

---

### **Error 3: Port Already in Use**

```
bind: Only one usage of each socket address (protocol/network address/port) is normally permitted
```

**Solution:**
```powershell
# Find process using port 8080
netstat -ano | findstr ":8080"

# Kill the process (replace <PID> with actual PID)
taskkill /PID <PID> /F

# Restart service
go run cmd/server/main.go
```

---

### **Error 4: Database Connection Failed**

```
ERROR Failed to connect to database: dial tcp [::1]:5432
```

**Check:**
```powershell
# Check if PostgreSQL is running
docker-compose ps postgres
```

**Solution:**
```powershell
# Restart PostgreSQL
docker-compose restart postgres

# Wait 5 seconds
Start-Sleep -Seconds 5

# Retry backend
go run cmd/server/main.go
```

---

### **Error 5: Casbin Enforcer Failed**

```
ERROR failed to initialize Casbin enforcer
```

**Cause:** Missing config files or database not ready.

**Check:**
```powershell
# Verify config files exist
dir configs\rbac_model.conf
dir configs\rbac_cms_model.conf
dir configs\rbac_user_model.conf
```

**Solution:**
```powershell
# If files missing, restore from git
git checkout configs/

# Rebuild if using Docker
docker-compose build --no-cache iam-service
docker-compose up -d iam-service
```

---

### **Error 6: Docker Daemon Not Running**

```
error during connect: This error may indicate that the docker daemon is not running
```

**Solution:**
1. Open **Docker Desktop**
2. Wait for it to fully start (green icon in system tray)
3. Verify: `docker ps`

---

### **Error 7: CORS Preflight Failed**

```
Access to fetch at 'http://localhost:8080' from origin 'http://localhost:3000' 
has been blocked by CORS policy
```

**Check Backend Logs:**
```powershell
# Should see OPTIONS request logged
docker-compose logs -f iam-service
```

**Verify CORS Middleware:**
```powershell
# Test OPTIONS manually
curl -X OPTIONS http://localhost:8080/v1/auth/register `
  -H "Origin: http://localhost:3000" `
  -v
```

**Expected Response:**
```
HTTP/1.1 204 No Content
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: POST, OPTIONS, GET, PUT, DELETE, PATCH
```

---

### **Error 8: Module Not Found (Go)**

```
module github.com/tvttt/iam-services: cannot find module
```

**Solution:**
```powershell
# Download dependencies
go mod download

# Tidy up
go mod tidy

# Retry
go run cmd/server/main.go
```

---

### **Error 9: Container Keeps Restarting**

```powershell
# Check logs for crash reason
docker-compose logs iam-service

# Common causes:
# - Database not ready
# - Port conflict
# - Panic in code
# - Missing environment variables
```

**Solution:**
```powershell
# View full logs
docker-compose logs --tail=100 iam-service

# Rebuild if needed
docker-compose down
docker-compose build --no-cache iam-service
docker-compose up -d
```

---

## 🛠️ Diagnostic Commands

### **Check All Services**

```powershell
# Run verification script
.\verify-services.ps1

# Or manually check each:

# PostgreSQL
docker-compose ps postgres
netstat -ano | findstr ":5432"

# Redis
docker-compose ps redis  
netstat -ano | findstr ":6379"

# Backend HTTP
netstat -ano | findstr ":8080"
curl http://localhost:8080/health

# Backend gRPC
netstat -ano | findstr ":50051"
```

---

### **View Logs**

```powershell
# All services
docker-compose logs

# Specific service
docker-compose logs postgres
docker-compose logs redis
docker-compose logs iam-service

# Follow logs real-time
docker-compose logs -f iam-service

# Last N lines
docker-compose logs --tail=50 iam-service
```

---

### **Clean Restart**

```powershell
# Nuclear option - Clean everything and start fresh

# 1. Stop and remove everything
docker-compose down -v

# 2. Remove images (optional)
docker-compose down --rmi all -v

# 3. Rebuild from scratch
docker-compose build --no-cache

# 4. Start services
docker-compose up -d

# 5. View logs
docker-compose logs -f
```

---

## 🔍 Debug Checklist

When debugging issues, check these in order:

- [ ] Docker Desktop is running
- [ ] `docker ps` shows containers
- [ ] PostgreSQL is Up (port 5432)
- [ ] Redis is Up (port 6379)
- [ ] Backend is listening on 8080
- [ ] `curl http://localhost:8080/health` returns 200
- [ ] OPTIONS request returns 204 with CORS headers
- [ ] No firewall blocking ports
- [ ] Frontend .env.local has correct URL

---

## 📊 Service Status Matrix

| Service | Port | Check | Expected |
|---------|------|-------|----------|
| PostgreSQL | 5432 | `netstat -ano \| findstr :5432` | LISTENING |
| Redis | 6379 | `netstat -ano \| findstr :6379` | LISTENING |
| gRPC | 50051 | `netstat -ano \| findstr :50051` | LISTENING |
| HTTP | 8080 | `curl localhost:8080/health` | `{"status":"ok"}` |
| Frontend | 3000 | Browser: `localhost:3000` | Home page |

---

## 🆘 Still Having Issues?

1. **Run all diagnostic scripts:**
   ```powershell
   .\verify-services.ps1
   ```

2. **Check logs for errors:**
   ```powershell
   docker-compose logs --tail=100
   ```

3. **Try clean restart:**
   ```powershell
   .\fix-redis.ps1
   .\start-services.ps1
   ```

4. **Test each component individually:**
   - Database: `docker-compose ps`
   - Backend: `go run cmd/server/main.go`
   - Frontend: `npm run dev`

---

**Last Updated:** 2024-11-06

