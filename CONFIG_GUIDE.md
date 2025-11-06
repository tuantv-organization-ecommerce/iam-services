# 📋 IAM Service Configuration Guide

## 🎯 Overview

IAM Service sử dụng **config.yml** để quản lý cấu hình thay vì `.env` files. Configuration được load bằng [Viper](https://github.com/spf13/viper) với hỗ trợ:
- ✅ YAML configuration files
- ✅ Environment variables override
- ✅ Default values
- ✅ Multiple search paths

---

## 📁 Configuration File Location

Viper sẽ tìm `config.yml` ở các vị trí sau (theo thứ tự ưu tiên):

1. **Current directory** (`.`)
2. **Config directory** (`./configs/`)
3. **Parent directories** (`../../`, `../../../`)
4. **System config** (`/etc/iam-services/`) - Linux only

---

## 🔧 Configuration Structure

### **config.yml Example**

```yaml
# ============================================
# IAM Service Configuration
# ============================================

# Server Configuration
server:
  host: "0.0.0.0"           # gRPC server host
  port: 50051               # gRPC server port
  http_host: "0.0.0.0"      # HTTP REST API host
  http_port: 8080           # HTTP REST API port

# Database Configuration
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "iam_db"
  sslmode: "disable"

# Redis Configuration
redis:
  host: "localhost"
  port: 6379
  password: ""              # Leave empty if no password
  db: 0

# JWT Configuration
jwt:
  secret: "your-secret-key-change-this-in-production"
  access_token_expiration_hours: 24
  refresh_token_expiration_hours: 168

# Logging Configuration
log:
  level: "info"             # Options: debug, info, warn, error
  encoding: "json"          # Options: json, console

# Swagger Configuration
swagger:
  enabled: true
  base_path: "/swagger/"
  spec_path: "/swagger.json"
  title: "IAM Service API Documentation"
  auth:
    username: "admin"
    password: "changeme"
    realm: "IAM Service API Documentation"
```

---

## 🌍 Environment Variables Override

Bạn có thể override bất kỳ config nào bằng environment variables với prefix `IAM_`:

### **Format:**
```
IAM_<SECTION>_<KEY>=value
```

### **Examples:**

```bash
# Override HTTP port
export IAM_SERVER_HTTP_PORT=9090

# Override database password
export IAM_DATABASE_PASSWORD=secret123

# Override JWT secret
export IAM_JWT_SECRET=super-secret-key

# Override log level
export IAM_LOG_LEVEL=debug
```

---

## 🚀 Quick Start

### **1. Local Development**

```bash
# Copy config template
cp config.yml.example config.yml

# Edit configuration
vim config.yml

# Run service
go run cmd/server/main.go
```

### **2. Docker Development**

```bash
# Create config.yml in project root
cat > config.yml << EOF
server:
  host: "0.0.0.0"
  port: 50051
  http_host: "0.0.0.0"
  http_port: 8080
database:
  host: "postgres"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "iam_db"
  sslmode: "disable"
redis:
  host: "redis"
  port: 6379
  password: ""
  db: 0
EOF

# Start services
docker-compose up -d
```

### **3. Production Deployment**

```bash
# Place config in system location
sudo mkdir -p /etc/iam-services
sudo cp config.yml /etc/iam-services/

# OR use environment variables
export IAM_SERVER_HOST=0.0.0.0
export IAM_SERVER_PORT=50051
export IAM_SERVER_HTTP_HOST=0.0.0.0
export IAM_SERVER_HTTP_PORT=8080
export IAM_DATABASE_HOST=prod-postgres.example.com
export IAM_DATABASE_PASSWORD=prod-secret
export IAM_REDIS_HOST=prod-redis.example.com
export IAM_JWT_SECRET=prod-jwt-secret

# Start service
./iam-service
```

---

## 🔍 Configuration Validation

Service sẽ validate configuration khi khởi động:

### **Required Fields:**
- ✅ `server.port` - gRPC port
- ✅ `server.http_port` - HTTP REST API port
- ✅ `database.host` - Database host
- ✅ `database.dbname` - Database name
- ✅ `jwt.secret` - JWT secret key

### **Validation Checks:**
```go
✓ Server Configuration:
  - gRPC: 0.0.0.0:50051
  - HTTP: 0.0.0.0:8080
✓ Database: postgres@localhost:5432/iam_db
✓ Redis: localhost:6379 (DB:0)
✓ Log Level: info
```

---

## 📊 Default Values

Nếu config.yml không tồn tại hoặc thiếu fields, Viper sẽ dùng default values:

| Section | Key | Default Value |
|---------|-----|---------------|
| `server.host` | gRPC Host | `0.0.0.0` |
| `server.port` | gRPC Port | `50051` |
| `server.http_host` | HTTP Host | `0.0.0.0` |
| `server.http_port` | HTTP Port | `8080` |
| `database.host` | DB Host | `localhost` |
| `database.port` | DB Port | `5432` |
| `database.user` | DB User | `postgres` |
| `database.password` | DB Password | `postgres` |
| `database.dbname` | DB Name | `iam_db` |
| `database.sslmode` | SSL Mode | `disable` |
| `redis.host` | Redis Host | `localhost` |
| `redis.port` | Redis Port | `6379` |
| `redis.password` | Redis Password | `` (empty) |
| `redis.db` | Redis DB | `0` |
| `jwt.secret` | JWT Secret | `your-secret-key-change-this-in-production` |
| `jwt.access_token_expiration_hours` | Access Token TTL | `24` hours |
| `jwt.refresh_token_expiration_hours` | Refresh Token TTL | `168` hours (7 days) |
| `log.level` | Log Level | `info` |
| `log.encoding` | Log Encoding | `json` |
| `swagger.enabled` | Swagger Enabled | `true` |

---

## 🐛 Troubleshooting

### **Problem: Config file not found**

**Symptom:**
```
WARNING: Failed to read config file: open config.yml: no such file or directory
Will use default values
```

**Solution:**
1. Create `config.yml` in the service root directory
2. Or set environment variables
3. Or place config in `/etc/iam-services/config.yml`

---

### **Problem: HTTP server not starting**

**Symptom:**
```
✓ Server Configuration:
  - gRPC: 0.0.0.0:50051
  - HTTP: :8080  # ← Missing host!
```

**Solution:**
Ensure `config.yml` has `http_host` and `http_port`:
```yaml
server:
  http_host: "0.0.0.0"
  http_port: 8080
```

---

### **Problem: Database connection fails**

**Symptom:**
```
failed to connect to database: pq: password authentication failed
```

**Solution:**
1. Check database credentials in `config.yml`
2. Or override with environment variables:
```bash
export IAM_DATABASE_PASSWORD=correct_password
```

---

### **Problem: Redis connection warning**

**Symptom:**
```
Failed to connect to Redis, token storage will be disabled
```

**Solution:**
1. Start Redis: `docker-compose up redis -d`
2. Update `config.yml`:
```yaml
redis:
  host: "localhost"  # or "redis" for Docker
  port: 6379
```

---

## 🔐 Security Best Practices

### **1. Protect JWT Secret**
```yaml
# ❌ BAD
jwt:
  secret: "your-secret-key"

# ✅ GOOD
jwt:
  secret: "ZXhhbXBsZS1zZWN1cmUtcmFuZG9tLXNlY3JldC1rZXk="  # Base64 encoded, 32+ chars
```

### **2. Don't Commit Secrets**
```bash
# Add to .gitignore
echo "config.yml" >> .gitignore
echo "*.local.yml" >> .gitignore
```

### **3. Use Environment Variables in Production**
```bash
# Use secrets management (Vault, AWS Secrets Manager, etc.)
export IAM_JWT_SECRET=$(vault read -field=jwt_secret secret/iam)
export IAM_DATABASE_PASSWORD=$(aws secretsmanager get-secret-value --secret-id iam-db-password --query SecretString --output text)
```

### **4. Enable SSL for Production**
```yaml
database:
  sslmode: "require"  # or "verify-full"
```

---

## 📚 Additional Resources

- **Viper Documentation**: https://github.com/spf13/viper
- **Configuration Code**: `internal/config/config.go`
- **Validation Code**: `internal/infrastructure/config/config_loader.go`

---

## 🆘 Need Help?

1. Check logs: `go run cmd/server/main.go`
2. Verify config: Look for "✓ Server Configuration" in logs
3. Test endpoints:
   ```bash
   curl http://localhost:8080/health
   curl http://localhost:8080/v1/auth/register
   ```

---

**Last Updated:** 2025-11-06  
**Version:** 1.0.0

