# ============================================
# Test Configuration Loading
# ============================================

Write-Host "`n╔════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║   Configuration Test Script            ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════╝`n" -ForegroundColor Cyan

# Step 1: Check .env file
Write-Host "🔍 Step 1: Checking .env file..." -ForegroundColor Yellow
if (Test-Path ".env") {
    Write-Host "✅ .env file exists" -ForegroundColor Green
    Write-Host "`n📄 Content:" -ForegroundColor Cyan
    Get-Content .env | Select-String "HTTP_HOST|HTTP_PORT|SERVER_"
    Write-Host ""
} else {
    Write-Host "❌ .env file NOT found!" -ForegroundColor Red
    Write-Host "Run .\create-env.ps1 to create it`n" -ForegroundColor Yellow
    exit 1
}

# Step 2: Check current directory
Write-Host "🔍 Step 2: Current directory..." -ForegroundColor Yellow
$currentDir = Get-Location
Write-Host "📁 $currentDir" -ForegroundColor Cyan
Write-Host ""

# Step 3: Test environment variables (without .env)
Write-Host "🔍 Step 3: Testing without environment variables..." -ForegroundColor Yellow
Write-Host "HTTP_HOST = [$env:HTTP_HOST]" -ForegroundColor Gray
Write-Host "HTTP_PORT = [$env:HTTP_PORT]" -ForegroundColor Gray
Write-Host ""

# Step 4: Build and run service in test mode
Write-Host "🔍 Step 4: Starting service to test config loading..." -ForegroundColor Yellow
Write-Host "Press Ctrl+C after seeing the logs`n" -ForegroundColor Cyan

# Run the service and capture first few seconds of output
Write-Host "════════════════════════════════════════" -ForegroundColor Gray
$job = Start-Job -ScriptBlock {
    Set-Location $using:currentDir
    go run cmd/server/main.go 2>&1
}

# Wait and show output
Start-Sleep -Seconds 5
$output = Receive-Job -Job $job

# Stop the job
Stop-Job -Job $job
Remove-Job -Job $job

# Display output
$output | Select-String "Current working directory|SUCCESS:|WARNING:|CRITICAL:|FINAL CONFIG:|HTTP server" | ForEach-Object {
    $line = $_.Line
    if ($line -match "SUCCESS") {
        Write-Host $line -ForegroundColor Green
    } elseif ($line -match "WARNING|CRITICAL") {
        Write-Host $line -ForegroundColor Red
    } elseif ($line -match "FINAL CONFIG") {
        Write-Host $line -ForegroundColor Yellow
    } else {
        Write-Host $line -ForegroundColor Cyan
    }
}
Write-Host "════════════════════════════════════════" -ForegroundColor Gray

# Analysis
Write-Host "`n📊 Analysis:" -ForegroundColor Cyan
if ($output -match "FINAL CONFIG.*HTTP_HOST=0\.0\.0\.0.*HTTP_PORT=8080") {
    Write-Host "✅ Configuration loaded correctly!" -ForegroundColor Green
    Write-Host "   HTTP_HOST = 0.0.0.0" -ForegroundColor Green
    Write-Host "   HTTP_PORT = 8080" -ForegroundColor Green
} elseif ($output -match "FINAL CONFIG.*HTTP_HOST=.*HTTP_PORT=") {
    Write-Host "⚠️  Configuration loaded but values might be wrong" -ForegroundColor Yellow
    $output | Select-String "FINAL CONFIG"
} else {
    Write-Host "❌ Could not detect configuration" -ForegroundColor Red
}

Write-Host "`n💡 To start the service normally:" -ForegroundColor Cyan
Write-Host "   go run cmd/server/main.go`n" -ForegroundColor Yellow

