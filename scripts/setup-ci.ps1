# CI/CD Setup Helper Script for Windows PowerShell
# This script prepares the project for CI/CD pipeline

$ErrorActionPreference = "Stop"

Write-Host "🚀 IAM Service CI/CD Setup" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# Check if we're in the right directory
if (-not (Test-Path "go.mod")) {
    Write-Host "❌ Error: go.mod not found. Please run this script from iam-services directory." -ForegroundColor Red
    exit 1
}

Write-Host "📍 Current directory: $(Get-Location)" -ForegroundColor Yellow
Write-Host ""

# Step 1: Create .env.example if not exists
Write-Host "Step 1: Creating .env.example..." -ForegroundColor Yellow
if (-not (Test-Path ".env.example")) {
    $envContent = @"
# IAM Service Environment Variables
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
"@
    $envContent | Out-File -FilePath ".env.example" -Encoding UTF8
    Write-Host "✅ .env.example created" -ForegroundColor Green
} else {
    Write-Host "⚠️  .env.example already exists" -ForegroundColor Yellow
}
Write-Host ""

# Step 2: Download Go dependencies
Write-Host "Step 2: Downloading Go dependencies..." -ForegroundColor Yellow
go mod download
Write-Host "✅ Dependencies downloaded" -ForegroundColor Green
Write-Host ""

Write-Host "Step 3: Tidying Go modules..." -ForegroundColor Yellow
go mod tidy
Write-Host "✅ Go modules tidied" -ForegroundColor Green
Write-Host ""

# Step 4: Run local tests (optional)
Write-Host "Step 4: Running local tests..." -ForegroundColor Yellow
$runTests = Read-Host "Do you want to run tests locally? (y/n)"
if ($runTests -eq "y" -or $runTests -eq "Y") {
    Write-Host "Running tests that don't require database..." -ForegroundColor Cyan
    
    try {
        go test -v ./pkg/jwt/
        go test -v ./pkg/password/
        Write-Host "✅ Local tests completed" -ForegroundColor Green
    } catch {
        Write-Host "⚠️  Some tests failed (expected if dependencies missing)" -ForegroundColor Yellow
    }
} else {
    Write-Host "⏭️  Skipping local tests" -ForegroundColor Yellow
}
Write-Host ""

# Step 5: Check Git status
Write-Host "Step 5: Checking Git status..." -ForegroundColor Yellow
Write-Host ""
git status --short
Write-Host ""

# Step 6: Summary
Write-Host "======================================" -ForegroundColor Green
Write-Host "✅ Setup Complete!" -ForegroundColor Green
Write-Host ""
Write-Host "📝 Next steps:" -ForegroundColor Cyan
Write-Host "   1. Review the changes: git status"
Write-Host "   2. Stage changes: git add ."
Write-Host "   3. Commit: git commit -m 'ci: setup CI/CD pipeline with basic tests'"
Write-Host "   4. Create branch: git checkout -b feature/setup-cicd"
Write-Host "   5. Push: git push origin feature/setup-cicd"
Write-Host ""
Write-Host "Or run the quick commit script:" -ForegroundColor Yellow
Write-Host "   .\scripts\quick-commit.ps1"
Write-Host ""
Write-Host "⚠️  Remember to check GitHub Actions after pushing!" -ForegroundColor Yellow
Write-Host ""

