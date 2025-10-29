# Wrapper script to run lint.ps1 with execution policy bypass
# This avoids the need to type -ExecutionPolicy Bypass every time

param(
    [string]$Target = "all",
    [switch]$Fix,
    [switch]$Fast,
    [switch]$Verbose
)

# Get script directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# Build arguments for lint.ps1
$args = @()
if ($Target -ne "all") { $args += "-Target", $Target }
if ($Fix) { $args += "-Fix" }
if ($Fast) { $args += "-Fast" }
if ($Verbose) { $args += "-Verbose" }

# Run lint.ps1 with execution policy bypass
Write-Host "Running lint with execution policy bypass..." -ForegroundColor Cyan
Write-Host ""

& powershell -ExecutionPolicy Bypass -File "$ScriptDir\lint.ps1" @args
