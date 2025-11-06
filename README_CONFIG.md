# 🔧 Configuration Guide

## 📋 Environment Variables

IAM service uses environment variables for configuration, loaded from `.env` file.

---

## ⚡ Quick Fix

### **If HTTP_HOST and HTTP_PORT are not loading:**

**Step 1: Create .env file**
```powershell
.\create-env.ps1
```

**Step 2: Verify file**
```powershell
Get-Content .env
```

**Step 3: Test configuration**
```powershell
.\test-config.ps1
```

**Step 4: Restart service**
```powershell
go run cmd/server/main.go
```

---

## 🔍 Configuration Loading Process

The service tries to load `.env` from multiple locations:

1. `.env` (current directory)
2. `../../.env` (from cmd/server/)
3. `../../../.env` (fallback)

**Logs show:**
```
Current working directory: /path/to/iam-services
SUCCESS: .env file loaded from: .env
ENV: HTTP_HOST = 0.0.0.0
ENV: HTTP_PORT = 8080
FINAL CONFIG: HTTP_HOST=0.0.0.0, HTTP_PORT=8080
```

---

## ✅ Expected Configuration

### **Server**
- `SERVER_HOST=0.0.0.0` - gRPC server host
- `SERVER_PORT=50051` - gRPC server port
- `HTTP_HOST=0.0.0.0` - HTTP/REST server host  ← IMPORTANT
- `HTTP_PORT=8080` - HTTP/REST server port      ← IMPORTANT

### **Database**
- `DB_HOST=localhost`
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=postgres`
- `DB_NAME=iam_db`
- `DB_SSL_MODE=disable`

### **Redis**
- `REDIS_HOST=localhost`
- `REDIS_PORT=6379`
- `REDIS_PASSWORD=` (empty)
- `REDIS_DB=0`

### **JWT**
- `JWT_SECRET=your-secret-key-change-this-in-production`
- `JWT_EXPIRATION_HOURS=24`
- `JWT_REFRESH_EXPIRATION_HOURS=168`

---

## 🐛 Troubleshooting

### **Problem: HTTP server address shows ":"**

**Symptom:**
```
Gin HTTP server is running address=: health=http://:/health
```

**Cause:** `HTTP_HOST` and `HTTP_PORT` are empty

**Solution:**
```powershell
# 1. Check .env exists
Get-ChildItem .env

# 2. Check content
Get-Content .env | Select-String "HTTP"

# 3. Recreate if needed
.\create-env.ps1

# 4. Restart service
go run cmd/server/main.go
```

---

### **Problem: .env not loading**

**Symptom:**
```
WARNING: No .env file found in any location
```

**Solution 1: Ensure file is in correct location**
```powershell
# Should be in iam-services/ directory
cd C:\Users\...\iam-services
dir .env
```

**Solution 2: Set environment variables directly**
```powershell
$env:HTTP_HOST="0.0.0.0"
$env:HTTP_PORT="8080"
go run cmd/server/main.go
```

---

### **Problem: Values show as empty in logs**

**Symptom:**
```
ENV: HTTP_HOST not set, using default: 0.0.0.0
ENV: HTTP_PORT not set, using default: 8080
```

**This is OK!** The service will use the default values.

**But if you see:**
```
CRITICAL: HTTP_HOST is empty, forcing default: 0.0.0.0
```

This means getEnv() returned empty string. Check .env file format.

---

## 📝 .env File Format

**Correct format:**
```env
HTTP_HOST=0.0.0.0
HTTP_PORT=8080
```

**Wrong formats:**
```env
# ❌ Spaces around =
HTTP_HOST = 0.0.0.0

# ❌ Quotes
HTTP_HOST="0.0.0.0"

# ❌ Comments on same line
HTTP_HOST=0.0.0.0 # host
```

---

## 🔍 Debug Logs

When service starts, check for these logs:

### **✅ Good:**
```
SUCCESS: .env file loaded from: .env
ENV: HTTP_HOST = 0.0.0.0
ENV: HTTP_PORT = 8080
FINAL CONFIG: HTTP_HOST=0.0.0.0, HTTP_PORT=8080
Gin HTTP server is running address=0.0.0.0:8080
```

### **⚠️ Warning (but OK):**
```
WARNING: No .env file found in any location
ENV: HTTP_HOST not set, using default: 0.0.0.0
ENV: HTTP_PORT not set, using default: 8080
FINAL CONFIG: HTTP_HOST=0.0.0.0, HTTP_PORT=8080
Gin HTTP server is running address=0.0.0.0:8080
```

### **❌ Bad:**
```
WARNING: No .env file found
CRITICAL: HTTP_HOST is empty, forcing default: 0.0.0.0
CRITICAL: HTTP_PORT is empty, forcing default: 8080
Gin HTTP server is running address=:
```

---

## 🎯 Testing Configuration

### **Manual Test:**
```powershell
# 1. Start service
go run cmd/server/main.go

# 2. In another terminal, test
curl http://localhost:8080/health

# Expected: {"status":"ok"}
```

### **Automated Test:**
```powershell
.\test-config.ps1
```

---

## 📊 Configuration Priority

1. **Environment variables** (highest priority)
2. **.env file** (middle priority)
3. **Default values in code** (lowest priority)

Example:
```powershell
# Set env var (overrides .env)
$env:HTTP_PORT="9000"

# Start service
go run cmd/server/main.go

# Will use port 9000, not 8080 from .env
```

---

## 💡 Best Practices

### **Development**
- Use `.env` file for local config
- Don't commit `.env` to git (in .gitignore)
- Use `.env.example` as template

### **Production**
- Use environment variables (not .env)
- Set via Docker/Kubernetes secrets
- Never expose secrets in logs

### **Docker**
- Environment variables set in docker-compose.yml
- .env file not needed in container

---

## 🆘 Still Not Working?

1. **Run diagnostic:**
   ```powershell
   .\test-config.ps1
   ```

2. **Check full logs:**
   ```powershell
   go run cmd/server/main.go 2>&1 | Select-String "CONFIG|ENV:|HTTP"
   ```

3. **Force environment variables:**
   ```powershell
   $env:HTTP_HOST="0.0.0.0"
   $env:HTTP_PORT="8080"
   go run cmd/server/main.go
   ```

4. **Verify with:**
   ```powershell
   curl http://localhost:8080/health
   ```

---

**Last Updated:** 2024-11-06

