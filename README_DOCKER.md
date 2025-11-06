# 🐳 Docker Quick Reference

## 🚀 Quick Start

```bash
# 1. Ensure config-docker.yml exists (or it will use defaults)
cp config-docker.yml.example config-docker.yml

# 2. Start all services
docker-compose up -d

# 3. Check status
docker-compose ps

# 4. View logs
docker-compose logs -f iam-service
```

## ✅ Verify Running

```bash
# Health check
curl http://localhost:8080/health
# Expected: {"service":"iam-service","status":"healthy"}

# Test registration
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","email":"user@example.com","password":"Pass@123","full_name":"User One"}'
```

## 📊 Service Endpoints

| Service | Type | Port | Endpoint |
|---------|------|------|----------|
| IAM HTTP | REST API | 8080 | http://localhost:8080 |
| IAM gRPC | gRPC | 50051 | localhost:50051 |
| PostgreSQL | Database | 5432 | localhost:5432 |
| Redis | Cache | 6379 | localhost:6379 |

## 🔧 Common Commands

```bash
# Stop all services
docker-compose down

# Restart a service
docker-compose restart iam-service

# View logs (last 100 lines)
docker-compose logs --tail=100 iam-service

# Rebuild after code changes
docker-compose up -d --build iam-service

# Clean start (⚠️ deletes data!)
docker-compose down -v && docker-compose up -d
```

## 📖 Full Documentation

See [DOCKER_GUIDE.md](./DOCKER_GUIDE.md) for comprehensive documentation including:
- Configuration options
- Environment variables
- Troubleshooting
- Production deployment
- Monitoring & maintenance

## 🔧 Configuration

### **Files:**
- `config-docker.yml` - Active config (gitignored, create from example)
- `config-docker.yml.example` - Template (committed to git)
- `docker-compose.yml` - Service orchestration

### **Override Config:**

Edit `docker-compose.yml`:

```yaml
iam-service:
  environment:
    IAM_LOG_LEVEL: debug          # Override log level
    IAM_JWT_SECRET: custom-key    # Override JWT secret
```

## 🐛 Troubleshooting

### Port conflict:
```bash
# Check what's using the port
netstat -ano | findstr :8080  # Windows
lsof -i :8080                 # Linux/Mac
```

### Service unhealthy:
```bash
docker-compose logs iam-service
docker-compose exec iam-service wget -qO- http://localhost:8080/health
```

### Config not loaded:
```bash
docker-compose exec iam-service cat /app/config.yml
```

## 📋 What's Updated

✅ **docker-compose.yml:**
- Mount `config-docker.yml` as volume
- Simplified environment variables (config.yml handles defaults)
- Added health check for iam-service
- Environment variables can override config values

✅ **Dockerfile:**
- Added `wget` for health checks
- Optimized for config.yml support

✅ **Configuration:**
- `config-docker.yml` - Docker-specific config
- Uses container names (`postgres`, `redis`)
- Viper loads config with defaults

## 🆘 Need Help?

1. **Check logs:** `docker-compose logs -f iam-service`
2. **Verify services:** `docker-compose ps`
3. **Test endpoints:** `curl http://localhost:8080/health`
4. **Read full guide:** [DOCKER_GUIDE.md](./DOCKER_GUIDE.md)

---

**Happy Dockering! 🐳**

