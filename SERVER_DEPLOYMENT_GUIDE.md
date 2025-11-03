# Server Deployment Guide - IAM Services

Hướng dẫn đầy đủ về cách có server/IP để deploy IAM Service API.

**Date**: November 2024  
**Status**: ✅ Complete Guide

---

## 📋 Table of Contents

1. [Quick Comparison](#quick-comparison)
2. [Free Options (Miễn Phí)](#free-options-miễn-phí)
3. [Cloud VPS (Trả Phí)](#cloud-vps-trả-phí)
4. [Local Network (Tạm Thời)](#local-network-tạm-thời)
5. [Production Setup](#production-setup)
6. [Enable CI/CD Auto-Deploy](#enable-cicd-auto-deploy)

---

## Quick Comparison

| Option | Cost | Setup Time | Best For | Public IP | Production Ready |
|--------|------|------------|----------|-----------|------------------|
| **ngrok** | Free tier | 5 mins | Quick demo | Temporary | ❌ No |
| **Render** | Free tier | 15 mins | Test/Demo | ✅ Yes | ⚠️ Limited |
| **Railway** | Free tier | 15 mins | Test/Demo | ✅ Yes | ⚠️ Limited |
| **DigitalOcean** | $4-6/month | 30 mins | Production | ✅ Yes | ✅ Yes |
| **AWS EC2** | ~$5/month | 45 mins | Production | ✅ Yes | ✅ Yes |
| **Google Cloud** | Free tier | 45 mins | Production | ✅ Yes | ✅ Yes |
| **Azure** | Free tier | 45 mins | Production | ✅ Yes | ✅ Yes |
| **Local Network** | Free | 5 mins | Team test | ⚠️ Local only | ❌ No |

---

## Free Options (Miễn Phí)

### Option 1: ngrok (Nhanh Nhất - 5 Phút)

**Ưu điểm**: 
- ✅ Cực kỳ nhanh (5 phút)
- ✅ Không cần config phức tạp
- ✅ Tốt cho demo/test tạm thời

**Nhược điểm**:
- ❌ URL thay đổi mỗi lần restart (free tier)
- ❌ Không dùng cho production
- ❌ Giới hạn bandwidth

#### Bước 1: Install ngrok

**Windows**:
```powershell
# Download từ: https://ngrok.com/download
# Hoặc dùng Chocolatey:
choco install ngrok

# Verify
ngrok version
```

**Linux/macOS**:
```bash
# Download và install
curl -s https://ngrok-agent.s3.amazonaws.com/ngrok.asc | \
  sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null && \
  echo "deb https://ngrok-agent.s3.amazonaws.com buster main" | \
  sudo tee /etc/apt/sources.list.d/ngrok.list && \
  sudo apt update && sudo apt install ngrok

# Verify
ngrok version
```

#### Bước 2: Signup & Get Token

1. Đăng ký tại: https://dashboard.ngrok.com/signup
2. Copy authtoken từ: https://dashboard.ngrok.com/get-started/your-authtoken
3. Configure:
```powershell
ngrok config add-authtoken YOUR_TOKEN_HERE
```

#### Bước 3: Start Service & ngrok

**Terminal 1** (Service):
```powershell
cd iam-services
go run cmd/server/main.go
```

**Terminal 2** (ngrok):
```powershell
ngrok http 8080
```

#### Bước 4: Get Public URL

ngrok sẽ hiển thị:
```
Session Status    online
Account           your-email@example.com
Version           3.x.x
Region            United States (us)
Forwarding        https://abc123.ngrok-free.app -> http://localhost:8080
```

**Public URL**: `https://abc123.ngrok-free.app`

#### Bước 5: Test API

```powershell
# Health check
curl https://abc123.ngrok-free.app/health

# Swagger UI
start https://abc123.ngrok-free.app/swagger/

# Share with team
# Gửi URL cho team members: https://abc123.ngrok-free.app
```

**Lưu ý**: URL này chỉ hoạt động khi service và ngrok đang chạy!

---

### Option 2: Render.com (Free Tier - 15 Phút)

**Ưu điểm**:
- ✅ Free tier generous (750 hours/month)
- ✅ Auto SSL certificate
- ✅ PostgreSQL database free
- ✅ URL cố định
- ✅ Auto deploy từ GitHub

**Nhược điểm**:
- ⚠️ Service sleep sau 15 phút không dùng
- ⚠️ Cold start ~30s khi wake up
- ⚠️ Giới hạn resources

#### Bước 1: Signup Render

1. Đăng ký tại: https://render.com
2. Chọn "Sign up with GitHub"

#### Bước 2: Create PostgreSQL Database

1. Dashboard → **New +** → **PostgreSQL**
2. Name: `iam-db`
3. Database: `iam_db`
4. User: `postgres`
5. Region: Choose closest to you
6. Plan: **Free**
7. Click **Create Database**
8. **Copy Internal Database URL** (dạng: `postgres://...`)

#### Bước 3: Run Migrations

Trên Render PostgreSQL dashboard:
1. Click **Connect** → **External**
2. Copy PSQL command
3. Chạy trên máy local:

```bash
# Connect to Render PostgreSQL
psql YOUR_EXTERNAL_CONNECTION_STRING

# Run migrations
\i C:/path/to/iam-services/migrations/001_init_schema.sql
\i C:/path/to/iam-services/migrations/002_seed_data.sql
\i C:/path/to/iam-services/migrations/003_casbin_tables.sql
\i C:/path/to/iam-services/migrations/004_casbin_seed_data.sql
\i C:/path/to/iam-services/migrations/005_separate_user_cms_authorization.sql
\i C:/path/to/iam-services/migrations/006_seed_separated_authorization.sql
\q
```

#### Bước 4: Create Web Service

1. Dashboard → **New +** → **Web Service**
2. Connect your GitHub repository
3. Select repository: `your-repo/iam-services`
4. Settings:
   - **Name**: `iam-service`
   - **Region**: Same as database
   - **Branch**: `main`
   - **Root Directory**: `ecommerce/back_end/iam-services` (nếu có)
   - **Runtime**: `Docker`
   - **Plan**: **Free**

#### Bước 5: Configure Environment Variables

Trong Web Service settings → **Environment**:

```bash
# Database (từ Render PostgreSQL)
DATABASE_URL=postgres://...  # Internal Database URL

# Server
HTTP_PORT=8080
GRPC_PORT=50051

# JWT
JWT_SECRET=your-super-secure-random-64-char-secret-key-here
JWT_EXPIRATION_HOURS=24

# Swagger
SWAGGER_ENABLED=true
SWAGGER_AUTH_USERNAME=admin
SWAGGER_AUTH_PASSWORD=change-this-in-production

# Log
LOG_LEVEL=info
```

#### Bước 6: Update Dockerfile (Nếu Cần)

Đảm bảo Dockerfile expose đúng ports:

```dockerfile
# iam-services/Dockerfile
FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o iam-service cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/iam-service .
COPY --from=builder /app/configs ./configs

EXPOSE 8080 50051

CMD ["./iam-service"]
```

#### Bước 7: Deploy

1. Click **Create Web Service**
2. Đợi build & deploy (~5-10 phút)
3. Service sẽ available tại: `https://iam-service-xxxx.onrender.com`

#### Bước 8: Test

```powershell
# Get your Render URL
$RENDER_URL = "https://iam-service-xxxx.onrender.com"

# Health check
curl "$RENDER_URL/health"

# Swagger UI
start "$RENDER_URL/swagger/"

# Share with team
Write-Host "API URL: $RENDER_URL"
```

**Auto Deploy**: Mỗi khi push lên GitHub `main` branch → Render tự động deploy!

---

### Option 3: Railway.app (Free Tier - 15 Phút)

**Ưu điểm**:
- ✅ Free $5 credit/month
- ✅ Không sleep như Render
- ✅ PostgreSQL included
- ✅ Deploy từ GitHub
- ✅ Better performance than Render free

**Nhược điểm**:
- ⚠️ $5/month chỉ đủ ~500 hours
- ⚠️ Cần credit card (không charge nếu không vượt limit)

#### Setup Tương Tự Render

1. Signup: https://railway.app
2. **New Project** → **Deploy from GitHub**
3. Select repository
4. Add **PostgreSQL** database
5. Configure environment variables
6. Deploy!

**URL**: `https://iam-service-production-xxxx.up.railway.app`

---

### Option 4: Google Cloud Free Tier

**Ưu điểm**:
- ✅ Free tier rất generous
- ✅ Không expire (always free)
- ✅ Compute Engine: 1 f1-micro instance
- ✅ Production grade
- ✅ Public static IP free

**Nhược điểm**:
- ⚠️ Cần credit card
- ⚠️ Setup phức tạp hơn
- ⚠️ f1-micro có performance thấp

#### Bước 1: Setup Google Cloud

1. Đăng ký: https://cloud.google.com/free
2. Tạo project mới: `iam-service`
3. Enable Compute Engine API

#### Bước 2: Create VM Instance

**Console** → **Compute Engine** → **VM Instances** → **Create Instance**:

```
Name: iam-service-vm
Region: us-central1 (hoặc gần bạn)
Zone: us-central1-a
Machine type: e2-micro (free tier)
Boot disk: Ubuntu 22.04 LTS, 30GB
Firewall: 
  ✅ Allow HTTP traffic
  ✅ Allow HTTPS traffic
```

#### Bước 3: Configure Firewall

**VPC Network** → **Firewall** → **Create Firewall Rule**:

```yaml
Name: allow-iam-service
Direction: Ingress
Targets: All instances
Source IP ranges: 0.0.0.0/0
Protocols/ports:
  tcp: 8080, 50051
```

#### Bước 4: SSH & Setup

Click **SSH** button trên VM instance:

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Logout and login again
exit
```

#### Bước 5: Deploy Service

SSH lại vào VM:

```bash
# Create directory
mkdir -p ~/iam-service
cd ~/iam-service

# Create docker-compose.yml
cat > docker-compose.yml <<'EOF'
version: '3.8'

services:
  postgres:
    image: postgres:14-alpine
    environment:
      POSTGRES_DB: iam_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: your_secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  iam-service:
    image: YOUR_DOCKER_USERNAME/iam-service:latest
    ports:
      - "8080:8080"
      - "50051:50051"
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: your_secure_password
      DB_NAME: iam_db
      JWT_SECRET: your-super-secure-random-64-char-secret
      HTTP_PORT: 8080
      GRPC_PORT: 50051
      SWAGGER_ENABLED: true
    depends_on:
      - postgres
    restart: unless-stopped

volumes:
  postgres_data:
EOF

# Pull and start
docker-compose pull
docker-compose up -d

# Check logs
docker-compose logs -f
```

#### Bước 6: Get External IP

```bash
# On VM
curl ifconfig.me

# Or from GCP Console
# VM Instances → External IP
```

#### Bước 7: Test

```powershell
# Replace with your VM's external IP
$GCP_IP = "34.123.45.67"

# Health check
curl "http://${GCP_IP}:8080/health"

# Swagger UI
start "http://${GCP_IP}:8080/swagger/"
```

**Lưu ý**: Free tier e2-micro có performance hạn chế nhưng đủ cho test/demo!

---

## Cloud VPS (Trả Phí)

### Option 5: DigitalOcean Droplet ($4-6/month)

**Ưu điểm**:
- ✅ Giá rẻ nhất ($4/month)
- ✅ Setup đơn giản
- ✅ Performance tốt
- ✅ Static IP
- ✅ Good documentation

**Nhược điểm**:
- ❌ Không có free tier
- ❌ Phải trả phí monthly

#### Pricing

| Plan | RAM | CPU | Storage | Transfer | Price |
|------|-----|-----|---------|----------|-------|
| Basic | 512MB | 1 CPU | 10GB | 500GB | $4/mo |
| Basic | 1GB | 1 CPU | 25GB | 1TB | $6/mo |
| Basic | 2GB | 1 CPU | 50GB | 2TB | $12/mo |

**Recommend**: 1GB plan ($6/month) - Đủ cho production nhỏ

#### Setup Steps

1. **Signup**: https://www.digitalocean.com (referral bonus $200)
2. **Create Droplet**:
   - Image: Ubuntu 22.04 LTS
   - Plan: Basic $6/month (1GB RAM)
   - Region: Singapore/US (gần users)
   - Authentication: SSH key
3. **Get IP**: Copy droplet IP (VD: `159.89.123.45`)
4. **SSH & Setup**: (Giống Google Cloud steps above)
5. **Deploy**: Docker Compose

**Public URL**: `http://YOUR_DROPLET_IP:8080`

#### Setup Domain (Optional)

Nếu muốn domain name (VD: `api.yourcompany.com`):

1. Buy domain: Namecheap, GoDaddy (~$10/year)
2. DigitalOcean → **Networking** → **Domains** → Add domain
3. Update nameservers tại domain registrar:
   ```
   ns1.digitalocean.com
   ns2.digitalocean.com
   ns3.digitalocean.com
   ```
4. Create A record: `api.yourcompany.com` → `YOUR_DROPLET_IP`
5. Wait DNS propagation (5-30 mins)

**Public URL**: `http://api.yourcompany.com:8080`

---

### Option 6: AWS EC2 (~$5/month)

**Ưu điểm**:
- ✅ Free tier 12 tháng đầu
- ✅ Scalable
- ✅ Enterprise grade
- ✅ Many regions

**Nhược điểm**:
- ⚠️ Setup phức tạp
- ⚠️ UI overwhelming
- ⚠️ Dễ vượt budget nếu không cẩn thận

#### Free Tier

- **EC2**: t2.micro (1GB RAM) - Free 750 hours/month (12 months)
- **RDS**: db.t2.micro - Free 750 hours/month (12 months)

#### Quick Setup

1. **Signup**: https://aws.amazon.com/free
2. **Launch EC2**:
   - AMI: Ubuntu Server 22.04 LTS
   - Instance type: t2.micro (free tier)
   - Configure Security Group:
     - SSH: 22 (your IP)
     - Custom TCP: 8080 (0.0.0.0/0)
     - Custom TCP: 50051 (0.0.0.0/0)
3. **Elastic IP**: Allocate để có static IP
4. **SSH & Deploy**: (Giống steps trên)

**Cost after free tier**: ~$5-10/month

---

### Option 7: Azure VM (Free Tier Available)

**Ưu điểm**:
- ✅ Free tier 12 tháng
- ✅ Good if using Microsoft stack
- ✅ Enterprise support

**Nhược điểm**:
- ⚠️ Setup phức tạp
- ⚠️ Portal confusing

#### Free Tier

- **VM**: B1S (1 vCPU, 1GB RAM) - Free 750 hours/month (12 months)
- **Database**: PostgreSQL Basic tier

Similar setup như AWS EC2.

---

## Local Network (Tạm Thời)

### Option 8: Local IP (Cho Team Nội Bộ)

**Use case**: Team members cùng WiFi/LAN

#### Bước 1: Get Local IP

**Windows**:
```powershell
ipconfig
# Tìm dòng "IPv4 Address" (VD: 192.168.1.100)
```

**Linux/macOS**:
```bash
ip addr show
# hoặc
ifconfig
```

#### Bước 2: Start Service

```powershell
cd iam-services
go run cmd/server/main.go
```

#### Bước 3: Configure Firewall

**Windows Firewall**:
```powershell
# Allow port 8080
netsh advfirewall firewall add rule name="IAM Service HTTP" dir=in action=allow protocol=TCP localport=8080

# Allow port 50051
netsh advfirewall firewall add rule name="IAM Service gRPC" dir=in action=allow protocol=TCP localport=50051
```

#### Bước 4: Share with Team

```
API URL: http://192.168.1.100:8080
Swagger: http://192.168.1.100:8080/swagger/
```

**Lưu ý**: Chỉ work khi:
- Team cùng mạng WiFi/LAN
- Service đang chạy trên máy bạn
- Firewall đã mở ports

---

## Production Setup

### Recommend Stack for Production

**Server**: DigitalOcean Droplet (1GB, $6/month)  
**Database**: DigitalOcean Managed PostgreSQL (hoặc same droplet)  
**Domain**: Namecheap (~$10/year)  
**SSL**: Let's Encrypt (Free)  
**Reverse Proxy**: Nginx  

**Total Cost**: ~$16/month + $10/year domain

### Complete Production Setup

#### Step 1: DigitalOcean Droplet

```bash
# SSH to droplet
ssh root@YOUR_DROPLET_IP

# Update
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install Nginx
sudo apt install nginx -y

# Install Certbot (SSL)
sudo apt install certbot python3-certbot-nginx -y
```

#### Step 2: Setup Domain

1. Buy domain: `yourcompany.com`
2. Add A record: `api.yourcompany.com` → `YOUR_DROPLET_IP`
3. Wait DNS propagation

#### Step 3: Deploy Service

```bash
# Create directory
mkdir -p /app/iam-service
cd /app/iam-service

# Create docker-compose.yml
cat > docker-compose.yml <<'EOF'
version: '3.8'

services:
  postgres:
    image: postgres:14-alpine
    environment:
      POSTGRES_DB: iam_db
      POSTGRES_USER: iam_user
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - /app/iam-service/backups:/backups
    restart: unless-stopped

  iam-service:
    image: YOUR_DOCKER_USERNAME/iam-service:latest
    ports:
      - "127.0.0.1:8080:8080"  # Only localhost
      - "127.0.0.1:50051:50051"
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: iam_user
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: iam_db
      JWT_SECRET: ${JWT_SECRET}
      HTTP_PORT: 8080
      GRPC_PORT: 50051
      LOG_LEVEL: info
    depends_on:
      - postgres
    restart: unless-stopped

volumes:
  postgres_data:
EOF

# Create .env file
cat > .env <<'EOF'
DB_PASSWORD=your-super-secure-database-password
JWT_SECRET=your-super-secure-64-character-jwt-secret-key-here
EOF

# Secure .env
chmod 600 .env

# Start services
docker-compose up -d
```

#### Step 4: Configure Nginx

```bash
# Create Nginx config
sudo nano /etc/nginx/sites-available/iam-service
```

Paste:
```nginx
server {
    listen 80;
    server_name api.yourcompany.com;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
    
    location / {
        limit_req zone=api_limit burst=20 nodelay;
        
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/iam-service /etc/nginx/sites-enabled/

# Test config
sudo nginx -t

# Restart Nginx
sudo systemctl restart nginx
```

#### Step 5: Setup SSL (HTTPS)

```bash
# Get SSL certificate
sudo certbot --nginx -d api.yourcompany.com

# Follow prompts
# Select: Redirect HTTP to HTTPS (option 2)

# Auto-renewal (certbot tự động setup)
sudo certbot renew --dry-run
```

#### Step 6: Test Production API

```powershell
# Health check (HTTPS!)
curl https://api.yourcompany.com/health

# Swagger UI
start https://api.yourcompany.com/swagger/

# Test from mobile
# Mở browser trên phone: https://api.yourcompany.com/swagger/
```

**Public URLs**:
- **API**: `https://api.yourcompany.com`
- **Swagger**: `https://api.yourcompany.com/swagger/`

#### Step 7: Monitoring & Backups

**Database Backup** (automated):
```bash
# Create backup script
cat > /app/iam-service/backup.sh <<'EOF'
#!/bin/bash
BACKUP_DIR="/app/iam-service/backups"
DATE=$(date +%Y%m%d_%H%M%S)
docker exec iam-service-postgres-1 pg_dump -U iam_user iam_db > $BACKUP_DIR/iam_db_$DATE.sql
# Keep only last 7 days
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
EOF

chmod +x /app/iam-service/backup.sh

# Add to crontab (daily at 2 AM)
sudo crontab -e
# Add: 0 2 * * * /app/iam-service/backup.sh
```

**Log Rotation**:
```bash
sudo nano /etc/logrotate.d/iam-service
```

```
/var/log/nginx/access.log /var/log/nginx/error.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    sharedscripts
    postrotate
        systemctl reload nginx
    endscript
}
```

---

## Enable CI/CD Auto-Deploy

Sau khi có server, enable auto-deploy từ GitHub:

### Step 1: Setup SSH Key for GitHub Actions

**Trên server**:
```bash
# Create deploy user
sudo adduser deploy
sudo usermod -aG docker deploy

# Switch to deploy user
sudo su - deploy

# Generate SSH key
ssh-keygen -t ed25519 -C "github-actions" -f ~/.ssh/github_actions
# Press Enter (no passphrase)

# Show public key
cat ~/.ssh/github_actions.pub
# Copy and add to: ~/.ssh/authorized_keys
```

**Copy private key**:
```bash
cat ~/.ssh/github_actions
# Copy entire content (including BEGIN/END lines)
```

### Step 2: Add GitHub Secrets

Repository → Settings → Secrets → Actions → New secret:

```
DOCKER_USERNAME=your-dockerhub-username
DOCKER_PASSWORD=your-dockerhub-token

PRODUCTION_HOST=api.yourcompany.com
PRODUCTION_USER=deploy
PRODUCTION_SSH_KEY=<paste-private-key-here>
```

### Step 3: Enable Deploy Jobs

Edit `.github/workflows/ci-cd.yml`:

**Line 388** (deploy-production):
```yaml
# Change from:
if: false  # Disabled temporarily

# To:
if: github.ref == 'refs/heads/main'
```

**Uncomment deploy-production job** (remove all `#` from lines 384-426)

### Step 4: Test Auto-Deploy

```bash
# Make a change
echo "# Test deploy" >> README.md

# Commit and push to main
git add .
git commit -m "test: trigger auto-deploy"
git push origin main

# Check GitHub Actions
# Repository → Actions → Latest workflow
# Verify: deploy-production job runs ✓
```

**Result**: 
- Push to `main` → Auto build → Auto deploy → Service restart
- Zero downtime với blue-green deployment (advanced)

---

## Summary & Recommendations

### For Quick Demo (Hôm Nay)

**Use**: ngrok (Free, 5 minutes)
```powershell
go run cmd/server/main.go
ngrok http 8080
# Share URL: https://abc123.ngrok-free.app
```

### For Team Testing (1-2 Tuần)

**Use**: Render.com or Railway (Free tier)
- Auto deploy từ GitHub
- PostgreSQL included
- URL cố định
- Cost: $0

### For Production (Lâu Dài)

**Use**: DigitalOcean Droplet + Domain + SSL
- Cost: ~$16/month + $10/year domain
- Static IP
- Full control
- Production ready

### Budget Comparison

| Use Case | Option | Cost | Setup Time |
|----------|--------|------|------------|
| Demo today | ngrok | $0 | 5 mins |
| Team test | Render | $0 | 15 mins |
| Small production | DigitalOcean | $6/mo | 1 hour |
| Full production | DigitalOcean + Domain | $16/mo + $10/year | 2 hours |
| Enterprise | AWS/GCP/Azure | $50+/mo | 1 day |

---

## Quick Start Commands

### ngrok (Fastest)
```powershell
# Install
choco install ngrok

# Setup
ngrok config add-authtoken YOUR_TOKEN

# Run (service must be running on localhost:8080)
ngrok http 8080
```

### Render.com
```
1. Signup → Connect GitHub
2. New Web Service → Select repo
3. Add PostgreSQL
4. Deploy
```

### DigitalOcean
```bash
# Create droplet → SSH
ssh root@YOUR_IP

# Install Docker
curl -fsSL https://get.docker.com | sh

# Deploy
cd /app
docker-compose up -d
```

---

## Support Resources

- **ngrok**: https://ngrok.com/docs
- **Render**: https://render.com/docs
- **Railway**: https://docs.railway.app
- **DigitalOcean**: https://www.digitalocean.com/community/tutorials
- **AWS**: https://aws.amazon.com/getting-started/
- **Google Cloud**: https://cloud.google.com/docs

---

**Questions?**

1. **Cần demo ngay?** → Use ngrok
2. **Test với team?** → Use Render/Railway
3. **Production thật?** → Use DigitalOcean
4. **Budget unlimited?** → Use AWS/GCP

Cho tôi biết bạn chọn option nào, tôi sẽ hướng dẫn chi tiết! 🚀

---

**Last Updated**: November 2024  
**Status**: ✅ Complete Guide

