# Script to download and install golangci-lint binary for Windows
# This avoids the "go install" issue with Go 1.19

param(
    [string]$Version = "v1.53.3"
)

$ErrorActionPreference = "Stop"

Write-Host "================================" -ForegroundColor Cyan
Write-Host "  Installing golangci-lint" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# Get GOPATH
$GOPATH = go env GOPATH
$binDir = Join-Path $GOPATH "bin"

Write-Host "Version: $Version" -ForegroundColor Yellow
Write-Host "Install to: $binDir" -ForegroundColor Yellow
Write-Host ""

# Create bin directory if not exists
if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir -Force | Out-Null
    Write-Host "Created directory: $binDir" -ForegroundColor Green
}

# Download URL (format: golangci-lint-1.54.2-windows-amd64.zip not v1.54.2)
$versionNumber = $Version.TrimStart('v')
$fileName = "golangci-lint-$versionNumber-windows-amd64.zip"
$downloadUrl = "https://github.com/golangci/golangci-lint/releases/download/$Version/$fileName"
$tempDir = $env:TEMP
$zipPath = Join-Path $tempDir $fileName
$extractPath = Join-Path $tempDir "golangci-lint-$versionNumber-windows-amd64"

Write-Host "Downloading from: $downloadUrl" -ForegroundColor Cyan
Write-Host ""

try {
    # Download
    Write-Host "Downloading..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing
    Write-Host "Downloaded to: $zipPath" -ForegroundColor Green
    Write-Host ""

    # Extract
    Write-Host "Extracting..." -ForegroundColor Yellow
    if (Test-Path $extractPath) {
        Remove-Item -Path $extractPath -Recurse -Force
    }
    Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force
    Write-Host "Extracted to: $extractPath" -ForegroundColor Green
    Write-Host ""

    # Copy binary
    Write-Host "Installing binary..." -ForegroundColor Yellow
    $exePath = Join-Path $extractPath "golangci-lint.exe"
    $destPath = Join-Path $binDir "golangci-lint.exe"
    
    if (Test-Path $exePath) {
        Copy-Item -Path $exePath -Destination $destPath -Force
        Write-Host "Installed to: $destPath" -ForegroundColor Green
    } else {
        throw "Binary not found in extracted files"
    }
    Write-Host ""

    # Cleanup
    Write-Host "Cleaning up..." -ForegroundColor Yellow
    Remove-Item -Path $zipPath -Force
    Remove-Item -Path $extractPath -Recurse -Force
    Write-Host "Cleanup complete" -ForegroundColor Green
    Write-Host ""

    # Verify installation
    Write-Host "Verifying installation..." -ForegroundColor Yellow
    $env:PATH += ";$binDir"
    $version = & "$destPath" version
    Write-Host "$version" -ForegroundColor Green
    Write-Host ""

    Write-Host "================================" -ForegroundColor Green
    Write-Host "  SUCCESS!" -ForegroundColor Green
    Write-Host "================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "golangci-lint has been installed to:" -ForegroundColor White
    Write-Host "  $destPath" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "IMPORTANT: Make sure this path is in your PATH:" -ForegroundColor Yellow
    Write-Host "  $binDir" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "To add to PATH permanently:" -ForegroundColor Yellow
    Write-Host "  1. Open System Properties > Environment Variables" -ForegroundColor White
    Write-Host "  2. Edit 'Path' variable" -ForegroundColor White
    Write-Host "  3. Add: $binDir" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Or for current session only:" -ForegroundColor Yellow
    Write-Host "  `$env:PATH += `";$binDir`"" -ForegroundColor Cyan
    Write-Host ""

} catch {
    Write-Host ""
    Write-Host "ERROR: $_" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please try manual installation:" -ForegroundColor Yellow
    Write-Host "  1. Download from: https://github.com/golangci/golangci-lint/releases" -ForegroundColor White
    Write-Host "  2. Extract golangci-lint.exe" -ForegroundColor White
    Write-Host "  3. Copy to: $binDir" -ForegroundColor Cyan
    exit 1
}

