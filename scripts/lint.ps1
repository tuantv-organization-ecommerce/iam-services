# PowerShell script to run golangci-lint
# Usage:
#   .\scripts\lint.ps1              # Lint all
#   .\scripts\lint.ps1 -Fix         # Lint with auto-fix
#   .\scripts\lint.ps1 -Fast        # Fast lint
#   .\scripts\lint.ps1 -Target model  # Lint specific package

param(
    [string]$Target = "all",
    [switch]$Fix,
    [switch]$Fast,
    [switch]$Verbose
)

# Get script directory and project root
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir

Write-Host "================================" -ForegroundColor Cyan
Write-Host "  golangci-lint for iam-services" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# Change to project directory
Set-Location $ProjectRoot

# Check if golangci-lint is installed
try {
    $null = Get-Command golangci-lint -ErrorAction Stop
} catch {
    Write-Host "ERROR: golangci-lint not found!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install golangci-lint:"
    Write-Host "  Go 1.19: go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2"
    Write-Host "  Go 1.20+: go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2"
    Write-Host ""
    Write-Host "Or download from: https://github.com/golangci/golangci-lint/releases"
    exit 1
}

# Build command
$cmd = "golangci-lint run"

if ($Fast) {
    $cmd += " --fast"
    Write-Host "Fast mode enabled (skipping some slow linters)" -ForegroundColor Yellow
}

if ($Fix) {
    $cmd += " --fix"
    Write-Host "Auto-fix mode enabled" -ForegroundColor Yellow
}

if ($Verbose) {
    $cmd += " --verbose"
}

# Add config file
if (Test-Path ".golangci.yml") {
    $cmd += " --config .golangci.yml"
}

# Target specific package
switch ($Target) {
    "model" { 
        $cmd += " internal/domain/model/..."
        Write-Host "Target: domain/model" -ForegroundColor Cyan
    }
    "handler" { 
        $cmd += " internal/handler/..."
        Write-Host "Target: handler" -ForegroundColor Cyan
    }
    "dao" { 
        $cmd += " internal/dao/..."
        Write-Host "Target: dao" -ForegroundColor Cyan
    }
    "service" { 
        $cmd += " internal/service/..."
        Write-Host "Target: service" -ForegroundColor Cyan
    }
    "all" { 
        Write-Host "Target: all packages" -ForegroundColor Cyan
    }
    default { 
        $cmd += " $Target"
        Write-Host "Target: $Target" -ForegroundColor Cyan
    }
}

Write-Host ""
Write-Host "Running: $cmd" -ForegroundColor Green
Write-Host ""

# Run the command
Invoke-Expression $cmd
$exitCode = $LASTEXITCODE

Write-Host ""
if ($exitCode -eq 0) {
    Write-Host "SUCCESS: Linting passed! Code is clean." -ForegroundColor Green
} else {
    Write-Host "FAILED: Linting failed. Please fix the issues above." -ForegroundColor Red
}

exit $exitCode
