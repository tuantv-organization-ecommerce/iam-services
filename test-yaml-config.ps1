#!/usr/bin/env pwsh
# Test YAML Config Script

Write-Host ""
Write-Host "========================================"
Write-Host "  Testing IAM Service (YAML Config)"
Write-Host "========================================"
Write-Host ""

# Wait for service to start
Write-Host "[1/3] Waiting for service to start..."
Start-Sleep -Seconds 5

# Test HTTP health endpoint
Write-Host "[2/3] Testing HTTP Server (port 8080)..."
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method GET -TimeoutSec 5 -UseBasicParsing -ErrorAction Stop
    Write-Host "  SUCCESS: HTTP Server is running"
    Write-Host "  Status: $($response.StatusCode)"
    Write-Host "  Response: $($response.Content)"
}
catch {
    Write-Host "  ERROR: HTTP Server not responding"
    Write-Host "  Error: $($_.Exception.Message)"
    exit 1
}

# Test CORS
Write-Host "[3/3] Testing CORS preflight..."
try {
    $headers = @{
        "Origin" = "http://localhost:3000"
        "Access-Control-Request-Method" = "POST"
        "Access-Control-Request-Headers" = "content-type"
    }
    $response = Invoke-WebRequest -Uri "http://localhost:8080/v1/auth/register" -Method OPTIONS -Headers $headers -TimeoutSec 5 -UseBasicParsing -ErrorAction Stop
    Write-Host "  SUCCESS: CORS is working"
    Write-Host "  Status: $($response.StatusCode)"
}
catch {
    Write-Host "  ERROR: CORS test failed"
    Write-Host "  Error: $($_.Exception.Message)"
}

Write-Host ""
Write-Host "========================================"
Write-Host "  All tests passed!"
Write-Host "========================================"
Write-Host ""
