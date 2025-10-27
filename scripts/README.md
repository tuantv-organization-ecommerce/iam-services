# Scripts Directory

Helper scripts for IAM Service development and CI/CD setup.

## 📜 Available Scripts

### CI/CD Setup Scripts

#### `setup-ci.sh` (Linux/macOS)
Automated setup script for CI/CD pipeline.

**Usage:**
```bash
chmod +x scripts/setup-ci.sh
./scripts/setup-ci.sh
```

**What it does:**
- Creates `.env.example` from template
- Downloads Go dependencies (`go mod download`)
- Tidies Go modules (`go mod tidy`)
- Optionally runs local tests
- Shows Git status

#### `setup-ci.ps1` (Windows PowerShell)
Windows version of the setup script.

**Usage:**
```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\setup-ci.ps1
```

### Quick Commit Scripts

#### `quick-commit.sh` (Linux/macOS)
Automated commit and push script for CI/CD setup.

**Usage:**
```bash
chmod +x scripts/quick-commit.sh
./scripts/quick-commit.sh
```

**What it does:**
- Shows current Git status
- Stages all changes (`git add .`)
- Creates commit with detailed message
- Optionally creates feature branch
- Pushes to GitHub
- Shows next steps

#### `quick-commit.ps1` (Windows PowerShell)
Windows version of the quick commit script.

**Usage:**
```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\quick-commit.ps1
```

### Proto Generation Scripts

#### `setup-proto.ps1`
Generates Protocol Buffer files for gRPC and REST API.

**Usage:**
```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\setup-proto.ps1
```

#### `generate-proto.ps1`
Alternative proto generation script.

**Usage:**
```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\generate-proto.ps1
```

#### `setup.sh` (Linux/macOS)
Unix version of proto generation.

**Usage:**
```bash
chmod +x scripts/setup.sh
./scripts/setup.sh
```

### Testing Scripts

#### `test-api.sh`
Tests API endpoints.

**Usage:**
```bash
chmod +x scripts/test-api.sh
./scripts/test-api.sh
```

---

## 🚀 Quick Start Workflow

### For CI/CD Setup (First Time)

**Linux/macOS:**
```bash
# Setup environment
./scripts/setup-ci.sh

# Commit and push
./scripts/quick-commit.sh
```

**Windows:**
```powershell
# Setup environment
.\scripts\setup-ci.ps1

# Commit and push
.\scripts\quick-commit.ps1
```

### Manual Steps (Alternative)

```bash
# 1. Create .env.example
cat > .env.example << 'EOF'
SERVER_HOST=0.0.0.0
SERVER_PORT=50051
HTTP_HOST=0.0.0.0
HTTP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=iam_db
DB_SSL_MODE=disable
JWT_SECRET=CHANGE-THIS-TO-64-CHAR-RANDOM-STRING
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_EXPIRATION_HOURS=168
CASBIN_MODEL_PATH=./configs/rbac_model.conf
LOG_LEVEL=info
LOG_ENCODING=json
EOF

# 2. Download dependencies
go mod download
go mod tidy

# 3. Commit and push
git add .
git commit -m "ci: setup CI/CD pipeline with basic tests"
git checkout -b feature/setup-cicd
git push origin feature/setup-cicd
```

---

## 📝 Script Permissions (Linux/macOS)

Make scripts executable:

```bash
chmod +x scripts/setup-ci.sh
chmod +x scripts/quick-commit.sh
chmod +x scripts/setup.sh
chmod +x scripts/test-api.sh
```

Or all at once:

```bash
chmod +x scripts/*.sh
```

---

## 🐛 Troubleshooting

### PowerShell Execution Policy Error

**Error:**
```
scripts\setup-ci.ps1 cannot be loaded because running scripts is disabled
```

**Solution:**
```powershell
# Run as Administrator
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser

# Or bypass for single execution
powershell -ExecutionPolicy Bypass -File .\scripts\setup-ci.ps1
```

### Permission Denied (Linux/macOS)

**Error:**
```
bash: ./scripts/setup-ci.sh: Permission denied
```

**Solution:**
```bash
chmod +x scripts/setup-ci.sh
```

### Git Not Found

**Error:**
```
git: command not found
```

**Solution:**
- Install Git: https://git-scm.com/downloads
- Verify: `git --version`

---

## 📚 Additional Resources

- **CI/CD Quick Start**: [../CI_QUICK_START.md](../CI_QUICK_START.md)
- **Full Setup Guide**: [../CI_CD_SETUP_GUIDE.md](../CI_CD_SETUP_GUIDE.md)
- **Main README**: [../README.md](../README.md)

---

**Last Updated**: 2024  
**Maintainer**: IAM Service Team

