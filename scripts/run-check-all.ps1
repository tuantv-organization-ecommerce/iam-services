# Wrapper script to run check-all.ps1 with execution policy bypass
# This avoids the need to type -ExecutionPolicy Bypass every time

param(
    [switch]$SkipTests,
    [switch]$Fast
)

# Get script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# Build arguments for check-all.ps1
$args = @()
if ($SkipTests) { $args += "-SkipTests" }
if ($Fast) { $args += "-Fast" }

# Run check-all.ps1 with execution policy bypass
Write-Host "Running check-all with execution policy bypass..." -ForegroundColor Cyan
Write-Host ""

& powershell -ExecutionPolicy Bypass -File "$ScriptDir\check-all.ps1" @args
