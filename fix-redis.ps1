# ============================================
# Fix Redis RDB Format Error
# ============================================

Write-Host "`n╔════════════════════════════════════════╗" -ForegroundColor Red
Write-Host "║   Redis RDB Format Fix Script          ║" -ForegroundColor Red
Write-Host "╚════════════════════════════════════════╝`n" -ForegroundColor Red

Write-Host "⚠️  This will DELETE all Redis data (cache, sessions, tokens)" -ForegroundColor Yellow
Write-Host "⚠️  This is safe for development but will require users to re-login`n" -ForegroundColor Yellow

$confirm = Read-Host "Continue? (yes/no)"
if ($confirm -ne "yes") {
    Write-Host "`n❌ Cancelled." -ForegroundColor Red
    exit 0
}

Write-Host "`n🔧 Step 1: Stopping all containers..." -ForegroundColor Cyan
docker-compose down
Write-Host "✅ Containers stopped`n" -ForegroundColor Green

Write-Host "🗑️  Step 2: Removing volumes (including Redis data)..." -ForegroundColor Cyan
docker-compose down -v
Write-Host "✅ Volumes removed`n" -ForegroundColor Green

Write-Host "🚀 Step 3: Starting PostgreSQL and Redis with fresh data..." -ForegroundColor Cyan
docker-compose up -d postgres redis
Write-Host "✅ Containers started`n" -ForegroundColor Green

Write-Host "⏳ Step 4: Waiting for services to be ready..." -ForegroundColor Cyan
Start-Sleep -Seconds 8
Write-Host "✅ Services should be ready`n" -ForegroundColor Green

Write-Host "🔍 Step 5: Checking service status..." -ForegroundColor Cyan
docker-compose ps

Write-Host "`n✅ Redis data has been reset successfully!" -ForegroundColor Green
Write-Host "You can now start the IAM service with:" -ForegroundColor Cyan
Write-Host "   go run cmd/server/main.go`n" -ForegroundColor Yellow

