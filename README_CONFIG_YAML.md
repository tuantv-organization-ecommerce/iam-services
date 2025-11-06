# 📝 YAML Configuration Guide

## ✨ New: YAML-based Configuration

IAM service now uses **YAML configuration** instead of `.env` files for better structure and readability.

---

## 🚀 Quick Start

### **Step 1: Copy example config**

```powershell
cd C:\Users\...\iam-services

# Copy example to config.yml
Copy-Item config.example.yml config.yml
```

### **Step 2: Edit config.yml**

```powershell
# Edit with your preferred editor
notepad config.yml
# or
code config.yml
```

### **Step 3: Start service**

```powershell
go run cmd/server/main.go
```

---

## 📁 Config File Structure

```yaml
# config.yml

server:
  host: "0.0.0.0"
  port: 50051
  http:
    host: "0.0.0.0"      # ← HTTP REST API host
    port: 8080            # ← HTTP REST API port

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "iam_db"
  sslmode: "disable"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  access_token_duration: 24h
  refresh_token_duration: 168h

log:
  level: "info"          # debug, info, warn, error
  encoding: "json"       # json or console

swagger:
  enabled: true
  base_path: "/swagger/"
  auth:
    username: "admin"
    password: "changeme"
```

---

## 🔧 Configuration Loading

### **Search Paths**

Config file is searched in these locations (in order):

1. `./config.yml` (current directory) ✅
2. `./configs/config.yml`
3. `../../config.yml` (from cmd/server/)
4. `/etc/iam-services/config.yml` (Linux system config)

### **Logs**

When service starts:

```
Current working directory: /path/to/iam-services
SUCCESS: Config file loaded from: /path/to/iam-services/config.yml
✓ Server Configuration:
  - gRPC: 0.0.0.0:50051
  - HTTP: 0.0.0.0:8080
✓ Database: postgres@localhost:5432/iam_db
✓ Redis: localhost:6379 (DB:0)
✓ Log Level: info
```

---

## 🌍 Environment Variables Override

Environment variables can **override** config file values.

**Format:** `IAM_<SECTION>_<KEY>`

```powershell
# Override HTTP port
$env:IAM_SERVER_HTTP_PORT="9000"

# Override database password
$env:IAM_DATABASE_PASSWORD="secret123"

# Start service
go run cmd/server/main.go

# Will use port 9000, not 8080 from config.yml
```

---

## 🎯 Configuration Priority

1. **Environment variables** (highest - overrides everything)
2. **config.yml file** (middle)
3. **Default values in code** (lowest - fallback only)

---

## 📝 Configuration Sections

### **server**

gRPC and HTTP server settings.

```yaml
server:
  host: "0.0.0.0"        # gRPC host
  port: 50051            # gRPC port
  http:
    host: "0.0.0.0"      # HTTP host
    port: 8080           # HTTP port
```

### **database**

PostgreSQL connection settings.

```yaml
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "iam_db"
  sslmode: "disable"     # disable, require, verify-full
```

### **redis**

Redis cache settings.

```yaml
redis:
  host: "localhost"
  port: 6379
  password: ""           # Empty if no auth
  db: 0                  # Redis database number
```

### **jwt**

JWT token configuration.

```yaml
jwt:
  secret: "change-this-in-production"
  access_token_duration: 24h      # 24 hours
  refresh_token_duration: 168h    # 7 days
```

**Duration format:**
- `1h` = 1 hour
- `24h` = 24 hours  
- `30m` = 30 minutes
- `168h` = 7 days

### **log**

Logging configuration.

```yaml
log:
  level: "info"          # debug, info, warn, error
  encoding: "json"       # json or console
```

### **swagger**

API documentation settings.

```yaml
swagger:
  enabled: true
  base_path: "/swagger/"
  auth:
    username: "admin"
    password: "changeme"
```

---

## 🔒 Security Best Practices

### **Development**

```yaml
jwt:
  secret: "dev-secret-key"

database:
  password: "postgres"

swagger:
  enabled: true
  auth:
    password: "changeme"
```

### **Production**

```yaml
jwt:
  secret: "${JWT_SECRET}"    # From environment variable

database:
  password: "${DB_PASSWORD}" # From environment variable

swagger:
  enabled: false             # Disable in production
  auth:
    password: "${SWAGGER_PASSWORD}"
```

**Set via environment:**

```bash
export IAM_JWT_SECRET="super-secret-key-xyz123"
export IAM_DATABASE_PASSWORD="complex-password"
export IAM_SWAGGER_AUTH_PASSWORD="admin-password"
```

---

## 📋 Multiple Environment Configs

### **Development**

```yaml
# config.yml (default for dev)
database:
  host: "localhost"
  
log:
  level: "debug"
```

### **Staging**

```yaml
# config.staging.yml
database:
  host: "staging-db.example.com"
  
log:
  level: "info"
```

### **Production**

```yaml
# config.production.yml
database:
  host: "prod-db.example.com"
  sslmode: "require"
  
log:
  level: "warn"

swagger:
  enabled: false
```

**Load specific config:**

```bash
# Set config name via env var
export CONFIG_NAME="production"
go run cmd/server/main.go
```

---

## ✅ Validation

Config is validated on load:

```go
// Ensures HTTP host/port are never empty
if config.Server.HTTPHost == "" {
    log.Printf("WARNING: using default: 0.0.0.0")
    config.Server.HTTPHost = "0.0.0.0"
}
```

---

## 🐛 Troubleshooting

### **Issue: Config file not found**

```
WARNING: Failed to read config file: Config File "config" Not Found
```

**Solution:**

```powershell
# Check file exists
Get-ChildItem config.yml

# Copy from example
Copy-Item config.example.yml config.yml
```

---

### **Issue: Values not loading**

**Check YAML syntax:**

```powershell
# Validate YAML online
# https://www.yamllint.com/

# Check for:
# - Correct indentation (2 spaces)
# - No tabs
# - Quotes around special characters
```

**Correct:**
```yaml
server:
  http:
    host: "0.0.0.0"
    port: 8080
```

**Wrong:**
```yaml
server:
http:              # ❌ Wrong indentation
    host: 0.0.0.0  # ❌ Missing quotes
	port: 8080     # ❌ Tab instead of spaces
```

---

### **Issue: Environment override not working**

**Format must be:** `IAM_<SECTION>_<SUBSECTION>_<KEY>`

```powershell
# ❌ Wrong
$env:HTTP_PORT="9000"

# ✅ Correct
$env:IAM_SERVER_HTTP_PORT="9000"
```

---

## 📊 Comparison: .env vs YAML

| Feature | .env | YAML ✅ |
|---------|------|---------|
| Structure | Flat | Hierarchical |
| Comments | Basic | Rich |
| Arrays/Lists | No | Yes |
| Type support | Strings only | Strings, numbers, bools, durations |
| Readability | Basic | Excellent |
| Validation | Manual | Built-in |
| IDE support | Limited | Excellent |

---

## 🎯 Migration from .env

### **Old (.env):**

```env
HTTP_HOST=0.0.0.0
HTTP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
```

### **New (config.yml):**

```yaml
server:
  http:
    host: "0.0.0.0"
    port: 8080

database:
  host: "localhost"
  port: 5432
```

**Benefits:**
- ✅ Better organized
- ✅ Easier to read
- ✅ Type-safe
- ✅ Comments and documentation
- ✅ Environment-specific configs

---

## 💡 Tips

1. **Always use quotes** for string values
2. **Use 2 spaces** for indentation (not tabs)
3. **Document your settings** with comments
4. **Keep secrets in environment variables**, not in config file
5. **Use config.example.yml** as template
6. **Don't commit config.yml** to git (.gitignore)

---

**Last Updated:** 2024-11-06

