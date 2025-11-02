# Complete check script: lint + build + test
# Run this before pushing to ensure everything passes

param(
    [switch]$SkipTests,
    [switch]$Fast
)

# Get script directory and project root
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  IAM-Services Pre-Push Checks" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Set-Location $ProjectRoot

$allPassed = $true

# Step 1: Run linter
Write-Host "Step 1/3: Running golangci-lint..." -ForegroundColor Cyan
Write-Host ""

if ($Fast) {
    & "$ScriptDir\lint.ps1" -Fast
} else {
    & "$ScriptDir\lint.ps1"
}

if ($LASTEXITCODE -ne 0) {
    $allPassed = $false
    Write-Host ""
    Write-Host "FAILED: Linting failed!" -ForegroundColor Red
} else {
    Write-Host "PASSED: Linting passed!" -ForegroundColor Green
}

Write-Host ""
Write-Host "----------------------------------------"
Write-Host ""

# Step 2: Build
Write-Host "Step 2/3: Building project..." -ForegroundColor Cyan
Write-Host ""

go build -v ./...
if ($LASTEXITCODE -ne 0) {
    $allPassed = $false
    Write-Host ""
    Write-Host "FAILED: Build failed!" -ForegroundColor Red
} else {
    Write-Host ""
    Write-Host "PASSED: Build passed!" -ForegroundColor Green
}

Write-Host ""
Write-Host "----------------------------------------"
Write-Host ""

# Step 3: Tests
if (-not $SkipTests) {
    Write-Host "Step 3/3: Running tests..." -ForegroundColor Cyan
    Write-Host ""
    
    go test -v ./...
    if ($LASTEXITCODE -ne 0) {
        $allPassed = $false
        Write-Host ""
        Write-Host "FAILED: Tests failed!" -ForegroundColor Red
    } else {
        Write-Host ""
        Write-Host "PASSED: Tests passed!" -ForegroundColor Green
    }
} else {
    Write-Host "Step 3/3: Tests skipped" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========================================"
Write-Host ""

if ($allPassed) {
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  SUCCESS: All checks passed!" -ForegroundColor Green
    Write-Host "  Ready to push to GitHub" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    exit 0
} else {
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "  FAILED: Some checks failed!" -ForegroundColor Red
    Write-Host "  Please fix issues before pushing" -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
    exit 1
}
