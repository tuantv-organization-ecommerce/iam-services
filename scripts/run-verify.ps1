# Wrapper script to run verify-lint.ps1 with execution policy bypass
# This avoids the need to type -ExecutionPolicy Bypass every time

param(
    [switch]$Quick,
    [switch]$Verbose
)

# Get script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# Build arguments for verify-lint.ps1
$args = @()
if ($Quick) { $args += "-Quick" }
if ($Verbose) { $args += "-Verbose" }

# Run verify-lint.ps1 with execution policy bypass
Write-Host "Running verify with execution policy bypass..." -ForegroundColor Cyan
Write-Host ""

& powershell -ExecutionPolicy Bypass -File "$ScriptDir\verify-lint.ps1" @args
