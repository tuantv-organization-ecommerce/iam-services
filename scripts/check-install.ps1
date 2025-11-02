# Quick check script to verify golangci-lint installation

Write-Host "================================" -ForegroundColor Cyan
Write-Host "  Checking golangci-lint Setup" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# Check Go version
Write-Host "1. Checking Go version..." -ForegroundColor Yellow
$goVersion = go version
Write-Host "   $goVersion" -ForegroundColor Green
Write-Host ""

# Check GOPATH
Write-Host "2. Checking GOPATH..." -ForegroundColor Yellow
$goPath = go env GOPATH
Write-Host "   GOPATH: $goPath" -ForegroundColor Green
Write-Host ""

# Check if golangci-lint is installed
Write-Host "3. Checking golangci-lint..." -ForegroundColor Yellow
try {
    $lintVersion = & golangci-lint version 2>&1
    Write-Host "   $lintVersion" -ForegroundColor Green
    Write-Host ""
    Write-Host "SUCCESS: golangci-lint is installed and ready!" -ForegroundColor Green
    exit 0
} catch {
    Write-Host "   NOT FOUND!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install golangci-lint:" -ForegroundColor Yellow
    Write-Host "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2"
    Write-Host ""
    Write-Host "Or check if $goPath\bin is in your PATH" -ForegroundColor Yellow
    exit 1
}



