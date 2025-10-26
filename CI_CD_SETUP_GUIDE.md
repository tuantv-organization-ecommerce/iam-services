# CI/CD Setup Guide - IAM Service

## 📋 Quick Start Guide

Hướng dẫn chi tiết để thiết lập CI/CD cho IAM Service từ đầu.

---

## 🔧 Prerequisites

### 1. GitHub Repository
- [ ] Repository đã được tạo
- [ ] Code đã được push lên GitHub
- [ ] Có quyền admin để setup Actions và Secrets

### 2. Docker Hub Account
- [ ] Tạo account tại https://hub.docker.com
- [ ] Tạo repository: `your-username/iam-service`
- [ ] Generate access token: Account Settings → Security → New Access Token

### 3. Server Infrastructure
- [ ] Staging server (Ubuntu/Debian recommended)
- [ ] Production server (Ubuntu/Debian recommended)
- [ ] SSH access với key-based authentication
- [ ] Docker & Docker Compose đã cài đặt

### 4. Optional Services
- [ ] Slack workspace (cho notifications)
- [ ] Codecov account (cho coverage tracking)

---

## 📝 Step-by-Step Setup

### Step 1: GitHub Repository Setup

#### 1.1 Enable GitHub Actions
```bash
# Actions tab sẽ tự động hiện sau khi push workflow files
# Verify: GitHub repo → Actions tab
```

#### 1.2 Configure Branch Protection
```
Settings → Branches → Add rule

Branch name pattern: main
✅ Require a pull request before merging
✅ Require status checks to pass before merging
   - Select: lint, test, build, security
✅ Require branches to be up to date before merging
```

### Step 2: Configure GitHub Secrets

Đi tới: `Settings → Secrets and variables → Actions → New repository secret`

#### 2.1 Docker Hub Credentials
```
Name: DOCKER_USERNAME
Secret: your-docker-username

Name: DOCKER_PASSWORD
Secret: your-docker-access-token
```

#### 2.2 Staging Server
```bash
# Generate SSH key trên máy local
ssh-keygen -t ed25519 -C "github-actions-staging" -f ~/.ssh/github_staging

# Copy public key to staging server
ssh-copy-id -i ~/.ssh/github_staging.pub deploy@staging.example.com

# Verify connection
ssh -i ~/.ssh/github_staging deploy@staging.example.com

# Add private key to GitHub Secrets
cat ~/.ssh/github_staging | pbcopy  # Copy to clipboard
```

```
Name: STAGING_HOST
Secret: staging.example.com

Name: STAGING_USER
Secret: deploy

Name: STAGING_SSH_KEY
Secret: <paste-private-key-here>
```

#### 2.3 Production Server
```bash
# Generate SSH key
ssh-keygen -t ed25519 -C "github-actions-production" -f ~/.ssh/github_production

# Copy to production server
ssh-copy-id -i ~/.ssh/github_production.pub deploy@production.example.com

# Add to GitHub Secrets
cat ~/.ssh/github_production | pbcopy
```

```
Name: PRODUCTION_HOST
Secret: production.example.com

Name: PRODUCTION_USER
Secret: deploy

Name: PRODUCTION_SSH_KEY
Secret: <paste-private-key-here>
```

#### 2.4 Slack Notification (Optional)
```
1. Create Slack App: https://api.slack.com/apps
2. Add Incoming Webhooks
3. Copy Webhook URL

Name: SLACK_WEBHOOK
Secret: https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXX
```

#### 2.5 Codecov (Optional)
```
1. Sign up at https://codecov.io with GitHub
2. Add repository
3. Copy token

Name: CODECOV_TOKEN
Secret: your-codecov-token
```

### Step 3: Server Setup

#### 3.1 Staging Server Setup
```bash
# SSH into staging server
ssh deploy@staging.example.com

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify
docker --version
docker-compose --version

# Create deployment directory
mkdir -p /app/iam-service
cd /app/iam-service

# Create .env file
nano .env
```

**Staging .env file:**
```env
# Database
DB_USER=postgres
DB_PASSWORD=<generate-strong-password>
DB_NAME=iam_db_staging

# JWT
JWT_SECRET=<generate-random-64-char-string>

# Docker
DOCKER_USERNAME=your-docker-username

# Logs
LOG_LEVEL=info
```

```bash
# Upload docker-compose file
# From local machine:
scp docker-compose.staging.yml deploy@staging.example.com:/app/iam-service/

# Test deployment
cd /app/iam-service
docker-compose -f docker-compose.staging.yml pull
docker-compose -f docker-compose.staging.yml up -d
docker-compose -f docker-compose.staging.yml logs -f
```

#### 3.2 Production Server Setup
```bash
# SSH into production server
ssh deploy@production.example.com

# Install Docker (same as staging)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
newgrp docker

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create deployment directory
mkdir -p /app/iam-service
cd /app/iam-service

# Create .env file
nano .env
```

**Production .env file:**
```env
# Database
DB_USER=iam_user
DB_PASSWORD=<very-strong-password-min-32-chars>
DB_NAME=iam_db

# JWT - IMPORTANT: Use strong random key
JWT_SECRET=<generate-with-openssl-rand-base64-64>
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_EXPIRATION_HOURS=168

# Docker
DOCKER_USERNAME=your-docker-username

# Logs
LOG_LEVEL=warn
```

```bash
# Generate strong JWT secret
openssl rand -base64 64 | tr -d '\n'

# Upload docker-compose file
# From local machine:
scp docker-compose.prod.yml deploy@production.example.com:/app/iam-service/

# Test deployment
cd /app/iam-service
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d
docker-compose -f docker-compose.prod.yml logs -f
```

### Step 4: Database Setup on Servers

#### 4.1 Run Migrations on Staging
```bash
ssh deploy@staging.example.com
cd /app/iam-service

# Copy migrations
mkdir -p migrations
# Upload migration files via scp or git clone

# Run migrations
docker exec -i iam-postgres-staging psql -U postgres -d iam_db_staging < migrations/001_init_schema.sql
docker exec -i iam-postgres-staging psql -U postgres -d iam_db_staging < migrations/002_seed_data.sql
docker exec -i iam-postgres-staging psql -U postgres -d iam_db_staging < migrations/003_casbin_tables.sql
docker exec -i iam-postgres-staging psql -U postgres -d iam_db_staging < migrations/004_casbin_seed_data.sql

# Verify
docker exec -it iam-postgres-staging psql -U postgres -d iam_db_staging -c "\dt"
```

#### 4.2 Run Migrations on Production
```bash
ssh deploy@production.example.com
cd /app/iam-service

# Run migrations (same as staging)
docker exec -i iam-postgres-prod psql -U iam_user -d iam_db < migrations/001_init_schema.sql
docker exec -i iam-postgres-prod psql -U iam_user -d iam_db < migrations/002_seed_data.sql
docker exec -i iam-postgres-prod psql -U iam_user -d iam_db < migrations/003_casbin_tables.sql
docker exec -i iam-postgres-prod psql -U iam_user -d iam_db < migrations/004_casbin_seed_data.sql

# Setup automated backups
docker-compose -f docker-compose.prod.yml up -d postgres-backup
```

### Step 5: Test CI/CD Pipeline

#### 5.1 Test Feature Branch
```bash
# Create feature branch
git checkout -b feature/test-cicd

# Make a small change
echo "# CI/CD Test" >> TEST.md
git add TEST.md
git commit -m "test: CI/CD pipeline test"

# Push to GitHub
git push origin feature/test-cicd

# Check GitHub Actions
# Go to: Repository → Actions tab
# Verify: lint, test, build, security jobs run successfully
```

#### 5.2 Test Staging Deployment
```bash
# Create PR to develop
# On GitHub: Create pull request from feature/test-cicd to develop

# After PR approved, merge to develop
git checkout develop
git pull origin develop

# Verify deployment
# GitHub Actions should auto-deploy to staging

# Check staging server
ssh deploy@staging.example.com
cd /app/iam-service
docker-compose -f docker-compose.staging.yml logs -f iam-service

# Test staging API
curl http://staging.example.com:8080/health
```

#### 5.3 Test Production Deployment
```bash
# Create PR from develop to main
# After review and approval, merge to main

# Verify production deployment
ssh deploy@production.example.com
cd /app/iam-service
docker-compose -f docker-compose.prod.yml logs -f iam-service

# Test production API
curl https://production.example.com/health

# Check GitHub Release
# Repository → Releases → Should see new release with version tag
```

### Step 6: Configure Monitoring (Optional)

#### 6.1 Setup Health Checks
```bash
# On servers, setup cron job for health monitoring
crontab -e

# Add:
*/5 * * * * curl -f http://localhost:8080/health || echo "Service down!" | mail -s "IAM Service Alert" admin@example.com
```

#### 6.2 Setup Log Aggregation
```bash
# Consider using:
# - ELK Stack (Elasticsearch, Logstash, Kibana)
# - Grafana Loki
# - AWS CloudWatch
# - Datadog

# Example: View logs
docker-compose -f docker-compose.prod.yml logs -f --tail=100 iam-service
```

---

## 🎯 Verification Checklist

### GitHub Actions
- [ ] Lint job passes
- [ ] Test job passes với coverage report
- [ ] Build job creates artifact
- [ ] Security scan completes
- [ ] Docker image builds và pushes
- [ ] Staging auto-deploys on develop merge
- [ ] Production auto-deploys on main merge
- [ ] Notifications are sent

### Staging Environment
- [ ] Service accessible via HTTP/gRPC
- [ ] Database migrations applied
- [ ] Health check endpoint works
- [ ] Can register/login users
- [ ] Logs are readable

### Production Environment
- [ ] Service accessible via HTTP/gRPC
- [ ] SSL/TLS configured (if applicable)
- [ ] Database migrations applied
- [ ] Backups configured
- [ ] Health check endpoint works
- [ ] Performance is acceptable
- [ ] Logs are aggregated

---

## 🔍 Troubleshooting Common Issues

### Issue: GitHub Actions fails to connect to Docker Hub
**Solution:**
```bash
# Verify credentials
echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin

# Check secrets in GitHub
Settings → Secrets → Verify DOCKER_USERNAME and DOCKER_PASSWORD
```

### Issue: SSH deployment fails
**Solution:**
```bash
# Test SSH connection manually
ssh -i ~/.ssh/github_staging deploy@staging.example.com

# Verify SSH key in GitHub Secrets
# Ensure private key includes:
-----BEGIN OPENSSH PRIVATE KEY-----
...
-----END OPENSSH PRIVATE KEY-----

# Check server SSH config
sudo nano /etc/ssh/sshd_config
# Ensure: PubkeyAuthentication yes
```

### Issue: Docker pull fails on server
**Solution:**
```bash
# Login to Docker Hub on server
docker login -u your-username

# Verify image exists
docker pull your-username/iam-service:latest

# Check docker-compose file
docker-compose -f docker-compose.staging.yml config
```

### Issue: Database migration fails
**Solution:**
```bash
# Check PostgreSQL logs
docker logs iam-postgres-staging

# Verify database exists
docker exec -it iam-postgres-staging psql -U postgres -l

# Re-run migrations manually
docker exec -i iam-postgres-staging psql -U postgres -d iam_db_staging < migrations/001_init_schema.sql
```

---

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Backup Best Practices](https://www.postgresql.org/docs/current/backup.html)
- [IAM Service README](./README.md)

---

## 🎉 Success!

Nếu tất cả checklist items đều pass, CI/CD pipeline của bạn đã sẵn sàng!

**Next Steps:**
1. Monitor first few deployments closely
2. Setup alerts và monitoring
3. Document runbooks cho team
4. Review và optimize pipeline performance
5. Train team members on CI/CD process

---

**Last Updated**: 2024-01-XX  
**Version**: 1.0.0

