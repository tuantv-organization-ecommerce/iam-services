# 🐳 Docker Deployment Guide

## 📋 Overview

IAM Service được containerized với Docker và orchestrated với Docker Compose. Configuration được quản lý thông qua `config-docker.yml`.

---

## 🚀 Quick Start

### **1. Prepare Configuration**

```bash
# Copy Docker config template
cp config-docker.yml.example config-docker.yml

# Edit if needed (optional, defaults work for local development)
vim config-docker.yml
```

### **2. Start All Services**

```bash
# Start all services (PostgreSQL + Redis + IAM Service)
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f iam-service
```

### **3. Verify Services**

```bash
# Health check
curl http://localhost:8080/health

# Test registration
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","email":"user@example.com","password":"Pass@123","full_name":"User One"}'
```

---

## 📦 Services Architecture

```
┌─────────────────┐
│   iam-service   │  Port: 8080 (HTTP), 50051 (gRPC)
│  (Go Binary)    │
└────────┬────────┘
         │
    ┌────┴─────┐
    │          │
┌───▼───┐  ┌──▼────┐
│ Postgres│  │ Redis │
│  :5432  │  │ :6379 │
└─────────┘  └───────┘
```

---

## 🔧 Configuration

### **Config File Hierarchy:**

1. **config-docker.yml** - Docker-specific config (mounted as volume)
2. **Environment variables** - Can override config.yml (prefix: `IAM_`)

### **config-docker.yml Structure:**

```yaml
server:
  host: "0.0.0.0"
  port: 50051
  http_host: "0.0.0.0"
  http_port: 8080

database:
  host: "postgres"      # ← Container name
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "iam_db"
  sslmode: "disable"

redis:
  host: "redis"         # ← Container name
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key-change-this-in-production"
  access_token_expiration_hours: 24
  refresh_token_expiration_hours: 168

log:
  level: "info"
  encoding: "json"
```

### **Override with Environment Variables:**

Edit `docker-compose.yml`:

```yaml
iam-service:
  environment:
    # Uncomment to override config-docker.yml
    IAM_LOG_LEVEL: debug
    IAM_JWT_SECRET: custom-secret-key
    IAM_DATABASE_PASSWORD: secure-password
```

---

## 📋 Docker Compose Commands

### **Start Services:**

```bash
# Start in detached mode
docker-compose up -d

# Start with rebuild
docker-compose up -d --build

# Start and view logs
docker-compose up
```

### **Stop Services:**

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (⚠️ deletes data!)
docker-compose down -v
```

### **View Logs:**

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f iam-service
docker-compose logs -f postgres
docker-compose logs -f redis

# Last 100 lines
docker-compose logs --tail=100 iam-service
```

### **Restart Services:**

```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart iam-service
```

### **Check Status:**

```bash
# Service status
docker-compose ps

# Health check
docker-compose exec iam-service wget -qO- http://localhost:8080/health
```

---

## 🔍 Troubleshooting

### **Problem 1: Port Already in Use**

**Symptom:**
```
Error starting userland proxy: listen tcp4 0.0.0.0:8080: bind: address already in use
```

**Solution:**
```bash
# Find process using port
netstat -ano | findstr :8080  # Windows
lsof -i :8080                 # Linux/Mac

# Stop the process or change port in docker-compose.yml
ports:
  - "9090:8080"  # External:Internal
```

---

### **Problem 2: Service Unhealthy**

**Symptom:**
```
docker-compose ps
iam-service   ... (unhealthy)
```

**Solution:**
```bash
# Check logs
docker-compose logs iam-service

# Check health
docker-compose exec iam-service wget -qO- http://localhost:8080/health

# Verify config
docker-compose exec iam-service cat /app/config.yml
```

---

### **Problem 3: Database Connection Failed**

**Symptom:**
```
failed to connect to database: dial tcp: lookup postgres
```

**Solution:**
```bash
# Ensure database host in config-docker.yml matches service name
database:
  host: "postgres"  # Must match docker-compose.yml service name

# Check postgres is running
docker-compose ps postgres

# Test connection
docker-compose exec iam-service ping postgres
```

---

### **Problem 4: Config Not Loaded**

**Symptom:**
```
WARNING: Failed to read config file
```

**Solution:**
```bash
# Verify config file is mounted
docker-compose exec iam-service ls -la /app/config.yml

# Check config content
docker-compose exec iam-service cat /app/config.yml

# Verify volume mount in docker-compose.yml
volumes:
  - ./config-docker.yml:/app/config.yml:ro
```

---

## 🔐 Production Deployment

### **1. Secure Configuration**

Create `config-docker.prod.yml`:

```yaml
server:
  host: "0.0.0.0"
  port: 50051
  http_host: "0.0.0.0"
  http_port: 8080

database:
  host: "prod-postgres.internal"
  port: 5432
  user: "iam_user"
  password: "${DB_PASSWORD}"  # From environment
  dbname: "iam_prod"
  sslmode: "require"

redis:
  host: "prod-redis.internal"
  port: 6379
  password: "${REDIS_PASSWORD}"
  db: 0

jwt:
  secret: "${JWT_SECRET}"  # From secrets manager
  access_token_expiration_hours: 1
  refresh_token_expiration_hours: 24

log:
  level: "warn"
  encoding: "json"
```

### **2. Use Docker Secrets**

```yaml
# docker-compose.prod.yml
services:
  iam-service:
    secrets:
      - jwt_secret
      - db_password
    environment:
      IAM_JWT_SECRET_FILE: /run/secrets/jwt_secret
      IAM_DATABASE_PASSWORD_FILE: /run/secrets/db_password

secrets:
  jwt_secret:
    external: true
  db_password:
    external: true
```

### **3. Enable TLS**

```yaml
# Use reverse proxy (nginx/traefik)
nginx:
  image: nginx:alpine
  ports:
    - "443:443"
  volumes:
    - ./nginx.conf:/etc/nginx/nginx.conf
    - ./certs:/etc/nginx/certs
```

---

## 📊 Monitoring

### **Health Checks:**

```bash
# Manual check
curl http://localhost:8080/health

# Docker health status
docker-compose ps
docker inspect iam-service --format='{{.State.Health.Status}}'
```

### **Logs:**

```bash
# Stream logs with timestamps
docker-compose logs -f --timestamps iam-service

# Export logs
docker-compose logs --no-color > iam-service.log
```

### **Metrics:**

```bash
# Container stats
docker stats iam-service

# Detailed info
docker-compose top iam-service
```

---

## 🔄 Updates & Maintenance

### **Update Service:**

```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose up -d --build iam-service

# Check logs
docker-compose logs -f iam-service
```

### **Database Migrations:**

```bash
# Run migrations manually
docker-compose exec iam-service ./iam-service migrate

# Or use init scripts (already mounted in postgres)
# Migrations run automatically on first start
```

### **Backup Data:**

```bash
# Backup PostgreSQL
docker-compose exec postgres pg_dump -U postgres iam_db > backup.sql

# Backup Redis
docker-compose exec redis redis-cli SAVE
docker cp iam-redis:/data/dump.rdb ./redis-backup.rdb
```

### **Restore Data:**

```bash
# Restore PostgreSQL
cat backup.sql | docker-compose exec -T postgres psql -U postgres iam_db

# Restore Redis
docker cp ./redis-backup.rdb iam-redis:/data/dump.rdb
docker-compose restart redis
```

---

## 📝 Configuration Examples

### **Development (default):**

```yaml
# config-docker.yml
database:
  host: "postgres"
  password: "postgres"
log:
  level: "debug"
jwt:
  access_token_expiration_hours: 24
```

### **Staging:**

```yaml
# config-docker.staging.yml
database:
  host: "staging-db.internal"
  password: "${DB_PASSWORD}"
  sslmode: "require"
log:
  level: "info"
jwt:
  access_token_expiration_hours: 8
```

### **Production:**

```yaml
# config-docker.prod.yml
database:
  host: "prod-db.internal"
  password: "${DB_PASSWORD}"
  sslmode: "verify-full"
log:
  level: "warn"
jwt:
  access_token_expiration_hours: 1
```

---

## 🆘 Getting Help

### **Debug Mode:**

```bash
# Enable debug logs
docker-compose down
# Edit config-docker.yml: log.level = "debug"
docker-compose up

# Or override with env var
docker-compose run -e IAM_LOG_LEVEL=debug iam-service
```

### **Interactive Shell:**

```bash
# Access container
docker-compose exec iam-service sh

# Test internal connectivity
wget -qO- http://localhost:8080/health
ping postgres
ping redis
```

### **Clean Restart:**

```bash
# Stop and remove everything
docker-compose down -v

# Remove images
docker-compose down --rmi all

# Start fresh
docker-compose up -d --build
```

---

**Last Updated:** 2025-11-06  
**Docker Compose Version:** 3.8  
**Required Docker Version:** 20.10+

