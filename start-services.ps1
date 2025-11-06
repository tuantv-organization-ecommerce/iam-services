# ============================================
# IAM Services Startup Script
# ============================================

Write-Host "`n╔════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║   IAM Services Startup Script          ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════╝`n" -ForegroundColor Cyan

$ErrorActionPreference = "Stop"

# Check if Docker is running
Write-Host "🔍 Step 1: Checking Docker..." -ForegroundColor Yellow
try {
    docker --version | Out-Null
    Write-Host "✅ Docker is available`n" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker not found. Please install Docker Desktop and start it." -ForegroundColor Red
    Write-Host "Download from: https://www.docker.com/products/docker-desktop/`n" -ForegroundColor Yellow
    exit 1
}

# Stop any existing containers
Write-Host "🛑 Step 2: Stopping existing containers..." -ForegroundColor Yellow
docker-compose down 2>&1 | Out-Null
Write-Host "✅ Containers stopped`n" -ForegroundColor Green

# Start PostgreSQL and Redis
Write-Host "🚀 Step 3: Starting PostgreSQL and Redis..." -ForegroundColor Yellow
docker-compose up -d postgres redis
Start-Sleep -Seconds 3

# Check database status
Write-Host "`n📊 Step 4: Checking database status..." -ForegroundColor Yellow
$dbStatus = docker-compose ps postgres
if ($dbStatus -match "Up") {
    Write-Host "✅ PostgreSQL is running" -ForegroundColor Green
} else {
    Write-Host "⚠️  PostgreSQL status unknown" -ForegroundColor Yellow
}

$redisStatus = docker-compose ps redis
if ($redisStatus -match "Up") {
    Write-Host "✅ Redis is running`n" -ForegroundColor Green
} else {
    Write-Host "⚠️  Redis status unknown`n" -ForegroundColor Yellow
}

# Wait for database to be ready
Write-Host "⏳ Step 5: Waiting for database to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 8
Write-Host "✅ Database should be ready now`n" -ForegroundColor Green

# Ask user how to run backend
Write-Host "════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "How do you want to run the backend?" -ForegroundColor Yellow
Write-Host "════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "[1] Docker Compose (Full containerized)" -ForegroundColor White
Write-Host "[2] Local Go (Faster, easier to debug) ⭐ RECOMMENDED" -ForegroundColor Green
Write-Host "════════════════════════════════════════`n" -ForegroundColor Cyan

$choice = Read-Host "Enter your choice (1 or 2)"

if ($choice -eq "1") {
    # Docker Compose method
    Write-Host "`n🔨 Building IAM service image..." -ForegroundColor Yellow
    docker-compose build iam-service
    
    Write-Host "`n🚀 Starting IAM service in Docker..." -ForegroundColor Yellow
    docker-compose up -d iam-service
    
    Write-Host "`n📋 Showing logs (Press Ctrl+C to stop viewing logs)..." -ForegroundColor Yellow
    Write-Host "Service will continue running in background`n" -ForegroundColor Cyan
    
    Start-Sleep -Seconds 3
    docker-compose logs -f iam-service
    
} elseif ($choice -eq "2") {
    # Local Go method
    Write-Host "`n🚀 Starting IAM service locally with Go..." -ForegroundColor Yellow
    Write-Host "Press Ctrl+C to stop the service`n" -ForegroundColor Cyan
    
    # Run Go in foreground so user can see logs
    go run cmd/server/main.go
    
} else {
    Write-Host "`n❌ Invalid choice. Please run the script again." -ForegroundColor Red
    exit 1
}

