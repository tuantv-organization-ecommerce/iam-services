# Verify golangci-lint setup for iam-services
# Run this to check if everything is working correctly

param(
    [switch]$Quick,
    [switch]$Verbose
)

$ErrorActionPreference = "Continue"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  VERIFY GOLANGCI-LINT FOR IAM-SERVICES" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$allPassed = $true

# Step 1: Check golangci-lint installation
Write-Host "1. Checking golangci-lint installation..." -ForegroundColor Yellow
try {
    $lintVersion = & golangci-lint version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "   ✅ $lintVersion" -ForegroundColor Green
    } else {
        Write-Host "   ❌ golangci-lint not found or error" -ForegroundColor Red
        Write-Host "   $lintVersion" -ForegroundColor Red
        $allPassed = $false
    }
} catch {
    Write-Host "   ❌ golangci-lint not found" -ForegroundColor Red
    Write-Host "   Please install: .\scripts\install-golangci-lint.ps1 -Version 'v1.54.2'" -ForegroundColor Yellow
    $allPassed = $false
}
Write-Host ""

# Step 2: Check Go version
Write-Host "2. Checking Go version..." -ForegroundColor Yellow
try {
    $goVersion = & go version 2>&1
    Write-Host "   ✅ $goVersion" -ForegroundColor Green
    
    # Extract Go version number
    if ($goVersion -match "go(\d+\.\d+)") {
        $goVersionNum = [version]$matches[1]
        if ($goVersionNum -ge [version]"1.19") {
            Write-Host "   ✅ Go version is compatible" -ForegroundColor Green
        } else {
            Write-Host "   ⚠️  Go version may be too old (recommend 1.19+)" -ForegroundColor Yellow
        }
    }
} catch {
    Write-Host "   ❌ Go not found" -ForegroundColor Red
    $allPassed = $false
}
Write-Host ""

# Step 3: Check go.mod compatibility
Write-Host "3. Checking go.mod compatibility..." -ForegroundColor Yellow
if (Test-Path "go.mod") {
    $goModVersion = Get-Content go.mod | Select-String "go " | ForEach-Object { $_.Line.Trim() }
    Write-Host "   ✅ $goModVersion" -ForegroundColor Green
    
    if ($goModVersion -match "go 1\.(\d+)") {
        $modVersion = [int]$matches[1]
        if ($modVersion -le 19) {
            Write-Host "   ✅ go.mod compatible with Go 1.19+" -ForegroundColor Green
        } else {
            Write-Host "   ⚠️  go.mod requires Go $modVersion+ (may need update)" -ForegroundColor Yellow
        }
    }
} else {
    Write-Host "   ❌ go.mod not found" -ForegroundColor Red
    $allPassed = $false
}
Write-Host ""

# Step 4: Check .golangci.yml config
Write-Host "4. Checking .golangci.yml config..." -ForegroundColor Yellow
if (Test-Path ".golangci.yml") {
    Write-Host "   ✅ .golangci.yml found" -ForegroundColor Green
    
    # Check if config is valid
    try {
        $configTest = & golangci-lint config path 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "   ✅ Configuration is valid" -ForegroundColor Green
        } else {
            Write-Host "   ⚠️  Configuration may have issues" -ForegroundColor Yellow
            if ($Verbose) {
                Write-Host "   $configTest" -ForegroundColor Gray
            }
        }
    } catch {
        Write-Host "   ⚠️  Could not validate config" -ForegroundColor Yellow
    }
} else {
    Write-Host "   ❌ .golangci.yml not found" -ForegroundColor Red
    $allPassed = $false
}
Write-Host ""

# Step 5: Check PATH
Write-Host "5. Checking PATH..." -ForegroundColor Yellow
$golangciInPath = $false
$golangciPath = ""

try {
    $null = Get-Command golangci-lint -ErrorAction Stop
    $golangciInPath = $true
    $golangciPath = (Get-Command golangci-lint).Source
    Write-Host "   ✅ golangci-lint found in PATH" -ForegroundColor Green
    Write-Host "   Location: $golangciPath" -ForegroundColor Gray
} catch {
    Write-Host "   ❌ golangci-lint not in PATH" -ForegroundColor Red
    Write-Host "   Please add to PATH: E:\go\src\bin" -ForegroundColor Yellow
    $allPassed = $false
}
Write-Host ""

# Step 6: Quick lint test
if (-not $Quick) {
    Write-Host "6. Running quick lint test..." -ForegroundColor Yellow
    
    try {
        Write-Host "   Testing on internal/domain/model/..." -ForegroundColor Gray
        $lintResult = & golangci-lint run --config .golangci.yml --fast internal/domain/model/... 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "   ✅ Lint test passed (no issues found)" -ForegroundColor Green
        } else {
            Write-Host "   ⚠️  Lint test found issues (this is normal)" -ForegroundColor Yellow
            if ($Verbose) {
                Write-Host "   Issues found:" -ForegroundColor Gray
                $lintResult | Select-Object -First 5 | ForEach-Object { Write-Host "   $_" -ForegroundColor Gray }
                if ($lintResult.Count -gt 5) {
                    Write-Host "   ... and $($lintResult.Count - 5) more" -ForegroundColor Gray
                }
            }
        }
    } catch {
        Write-Host "   ❌ Lint test failed" -ForegroundColor Red
        Write-Host "   Error: $_" -ForegroundColor Red
        $allPassed = $false
    }
    Write-Host ""
}

# Step 7: Check scripts
Write-Host "7. Checking lint scripts..." -ForegroundColor Yellow
$scripts = @("lint.ps1", "check-all.ps1", "install-golangci-lint.ps1")
$scriptsOk = $true

foreach ($script in $scripts) {
    if (Test-Path "scripts\$script") {
        Write-Host "   ✅ scripts\$script" -ForegroundColor Green
    } else {
        Write-Host "   ❌ scripts\$script missing" -ForegroundColor Red
        $scriptsOk = $false
        $allPassed = $false
    }
}

if ($scriptsOk) {
    Write-Host "   ✅ All lint scripts present" -ForegroundColor Green
} else {
    Write-Host "   ❌ Some scripts missing" -ForegroundColor Red
}
Write-Host ""

# Summary
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  VERIFICATION SUMMARY" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

if ($allPassed) {
    Write-Host "✅ ALL CHECKS PASSED!" -ForegroundColor Green
    Write-Host ""
    Write-Host "golangci-lint is ready to use!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "  .\scripts\lint.ps1              # Run lint" -ForegroundColor White
    Write-Host "  .\scripts\lint.ps1 -Fix         # Auto-fix" -ForegroundColor White
    Write-Host "  .\scripts\check-all.ps1         # Full check" -ForegroundColor White
} else {
    Write-Host "❌ SOME CHECKS FAILED!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please fix the issues above before using golangci-lint." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Common fixes:" -ForegroundColor Yellow
    Write-Host "  .\scripts\install-golangci-lint.ps1 -Version 'v1.54.2'" -ForegroundColor White
    Write-Host "  `$env:PATH += ';E:\go\src\bin'" -ForegroundColor White
}

Write-Host ""
Write-Host "For detailed help, see:" -ForegroundColor Cyan
Write-Host "  LINTING_SETUP.md" -ForegroundColor White
Write-Host "  INSTALLATION_GUIDE.md" -ForegroundColor White

exit $(if ($allPassed) { 0 } else { 1 })
