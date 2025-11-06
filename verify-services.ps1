# ============================================
# IAM Services Verification Script
# ============================================

Write-Host "`n╔════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║   IAM Services Verification            ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════╝`n" -ForegroundColor Cyan

# Check PostgreSQL
Write-Host "🔍 Checking PostgreSQL (port 5432)..." -ForegroundColor Yellow
$postgres = netstat -ano | Select-String ":5432.*LISTENING"
if ($postgres) {
    Write-Host "✅ PostgreSQL is listening on port 5432" -ForegroundColor Green
} else {
    Write-Host "❌ PostgreSQL is NOT running" -ForegroundColor Red
}

# Check Redis
Write-Host "`n🔍 Checking Redis (port 6379)..." -ForegroundColor Yellow
$redis = netstat -ano | Select-String ":6379.*LISTENING"
if ($redis) {
    Write-Host "✅ Redis is listening on port 6379" -ForegroundColor Green
} else {
    Write-Host "❌ Redis is NOT running" -ForegroundColor Red
}

# Check gRPC Server
Write-Host "`n🔍 Checking gRPC Server (port 50051)..." -ForegroundColor Yellow
$grpc = netstat -ano | Select-String ":50051.*LISTENING"
if ($grpc) {
    Write-Host "✅ gRPC Server is listening on port 50051" -ForegroundColor Green
} else {
    Write-Host "❌ gRPC Server is NOT running" -ForegroundColor Red
}

# Check HTTP Server
Write-Host "`n🔍 Checking HTTP Server (port 8080)..." -ForegroundColor Yellow
$http = netstat -ano | Select-String ":8080.*LISTENING"
if ($http) {
    Write-Host "✅ HTTP Server is listening on port 8080" -ForegroundColor Green
} else {
    Write-Host "❌ HTTP Server is NOT running" -ForegroundColor Red
}

# Test Health Endpoint
Write-Host "`n🔍 Testing Health Endpoint..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -UseBasicParsing -TimeoutSec 5
    if ($response.StatusCode -eq 200) {
        Write-Host "✅ Health check passed!" -ForegroundColor Green
        Write-Host "   Response: $($response.Content)" -ForegroundColor Cyan
    }
} catch {
    Write-Host "❌ Health check failed!" -ForegroundColor Red
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
}

# Test Register Endpoint
Write-Host "`n🔍 Testing Register Endpoint (OPTIONS)..." -ForegroundColor Yellow
try {
    $headers = @{
        "Origin" = "http://localhost:3000"
        "Access-Control-Request-Method" = "POST"
        "Access-Control-Request-Headers" = "content-type"
    }
    
    $response = Invoke-WebRequest -Uri "http://localhost:8080/v1/auth/register" -Method Options -Headers $headers -UseBasicParsing -TimeoutSec 5
    if ($response.StatusCode -eq 204) {
        Write-Host "✅ OPTIONS preflight passed!" -ForegroundColor Green
        Write-Host "   CORS Headers:" -ForegroundColor Cyan
        Write-Host "   - Allow-Origin: $($response.Headers['Access-Control-Allow-Origin'])" -ForegroundColor Cyan
        Write-Host "   - Allow-Methods: $($response.Headers['Access-Control-Allow-Methods'])" -ForegroundColor Cyan
    }
} catch {
    Write-Host "❌ OPTIONS preflight failed!" -ForegroundColor Red
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
}

# Summary
Write-Host "`n════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "SUMMARY" -ForegroundColor Yellow
Write-Host "════════════════════════════════════════" -ForegroundColor Cyan

$allGood = $postgres -and $redis -and $http
if ($allGood) {
    Write-Host "✅ All services are running correctly!" -ForegroundColor Green
    Write-Host "`n🎉 You can now test the frontend:" -ForegroundColor Green
    Write-Host "   http://localhost:3000/auth/signup" -ForegroundColor Cyan
} else {
    Write-Host "⚠️  Some services are not running." -ForegroundColor Yellow
    Write-Host "`nPlease run: .\start-services.ps1" -ForegroundColor Cyan
}

Write-Host "`n"

