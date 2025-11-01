# ✅ Verification Guide - golangci-lint for iam-services

Hướng dẫn verify golangci-lint đã hoạt động đúng cho iam-services.

---

## 🚀 Quick Verification

### Cách 1: Script Tự Động (Khuyến nghị)

```powershell
# Chạy verify script
.\scripts\verify-lint.ps1

# Quick mode (bỏ qua lint test)
.\scripts\verify-lint.ps1 -Quick

# Verbose mode (hiển thị chi tiết)
.\scripts\verify-lint.ps1 -Verbose
```

**Expected Output:**
```
========================================
  VERIFY GOLANGCI-LINT FOR IAM-SERVICES
========================================

1. Checking golangci-lint installation...
   ✅ golangci-lint has version 1.54.2

2. Checking Go version...
   ✅ go version go1.24.9 windows/amd64
   ✅ Go version is compatible

3. Checking go.mod compatibility...
   ✅ go 1.19
   ✅ go.mod compatible with Go 1.19+

4. Checking .golangci.yml config...
   ✅ .golangci.yml found
   ✅ Configuration is valid

5. Checking PATH...
   ✅ golangci-lint found in PATH

7. Checking lint scripts...
   ✅ All lint scripts present

========================================
  VERIFICATION SUMMARY
========================================
✅ ALL CHECKS PASSED!

golangci-lint is ready to use!
```

### Cách 2: Manual Verification

```powershell
# 1. Check golangci-lint version
golangci-lint version

# 2. Check Go version
go version

# 3. Check go.mod
Get-Content go.mod | Select-String "go "

# 4. Check config file
Test-Path .golangci.yml

# 5. Check PATH
Get-Command golangci-lint

# 6. Test lint on small package
golangci-lint run --fast internal/domain/model/...
```

---

## 🧪 Test Linting

### Quick Test

```powershell
# Test trên package nhỏ
.\scripts\lint.ps1 -Target model -Fast

# Test với auto-fix
.\scripts\lint.ps1 -Target model -Fix

# Test tất cả (có thể mất thời gian)
.\scripts\lint.ps1 -Fast
```

### Full Test

```powershell
# Chạy full check (lint + build + test)
.\scripts\check-all.ps1

# Hoặc từng bước
.\scripts\lint.ps1
go build ./...
go test ./...
```

---

## 🔍 What Gets Verified

### 1. Installation Check
- ✅ golangci-lint binary exists
- ✅ Version compatibility (1.54.2 for Go 1.19+)
- ✅ Binary is executable

### 2. Go Environment
- ✅ Go version (1.19+ recommended)
- ✅ go.mod compatibility
- ✅ Module structure

### 3. Configuration
- ✅ .golangci.yml exists
- ✅ Configuration is valid
- ✅ Linters are properly configured

### 4. PATH Setup
- ✅ golangci-lint in PATH
- ✅ Can be called from anywhere
- ✅ Correct binary location

### 5. Scripts
- ✅ lint.ps1 exists and executable
- ✅ check-all.ps1 exists
- ✅ install-golangci-lint.ps1 exists

### 6. Lint Test (Optional)
- ✅ Can run on codebase
- ✅ No critical errors
- ✅ Configuration works

---

## 🐛 Troubleshooting

### Issue: "golangci-lint not found"

**Solution:**
```powershell
# Add to PATH
$env:PATH += ";E:\go\src\bin"

# Or install if missing
.\scripts\install-golangci-lint.ps1 -Version "v1.54.2"
```

### Issue: "Configuration may have issues"

**Solution:**
```powershell
# Check config path
golangci-lint config path

# Validate config
golangci-lint run --config .golangci.yml --help
```

### Issue: "Go version incompatible"

**Solution:**
```powershell
# Check Go version
go version

# Update go.mod if needed
(Get-Content go.mod) -replace 'go 1\.\d+', 'go 1.19' | Set-Content go.mod
```

### Issue: "Lint test failed"

**Solution:**
```powershell
# Run with verbose to see errors
.\scripts\lint.ps1 -Verbose

# Fix specific issues
.\scripts\lint.ps1 -Fix

# Test on smaller scope
.\scripts\lint.ps1 -Target model
```

---

## 📊 Verification Levels

### Level 1: Basic (30 seconds)
```powershell
.\scripts\verify-lint.ps1 -Quick
```
- Installation check
- PATH check
- Config check

### Level 2: Standard (2-5 minutes)
```powershell
.\scripts\verify-lint.ps1
```
- All Level 1 checks
- Lint test on model package

### Level 3: Full (5-15 minutes)
```powershell
.\scripts\check-all.ps1
```
- All Level 2 checks
- Full lint on all packages
- Build test
- Unit tests

---

## 🎯 Success Criteria

### ✅ Ready to Use
- All verification checks pass
- golangci-lint responds to commands
- No critical configuration errors
- Scripts are executable

### ⚠️ Needs Attention
- Some checks fail but non-critical
- Configuration warnings
- PATH issues (temporary)

### ❌ Not Ready
- golangci-lint not installed
- Critical configuration errors
- Go version incompatible
- Scripts missing

---

## 📚 Related Commands

```powershell
# Verification
.\scripts\verify-lint.ps1              # Full verify
.\scripts\verify-lint.ps1 -Quick       # Quick verify
.\scripts\verify-lint.ps1 -Verbose     # Detailed verify

# Linting
.\scripts\lint.ps1                     # Lint all
.\scripts\lint.ps1 -Fast               # Fast lint
.\scripts\lint.ps1 -Fix                # Auto-fix
.\scripts\lint.ps1 -Target model       # Specific package

# Full Check
.\scripts\check-all.ps1                # Lint + Build + Test
.\scripts\check-all.ps1 -Fast          # Fast mode
.\scripts\check-all.ps1 -SkipTests     # Skip tests

# Installation
.\scripts\install-golangci-lint.ps1 -Version "v1.54.2"
```

---

## 📖 Documentation

- **`LINTING_SETUP.md`** - Complete setup guide
- **`INSTALLATION_GUIDE.md`** - Installation troubleshooting
- **`SETUP_COMPLETE.md`** - Setup summary
- **`fix_error_ci_cd.md`** - CI/CD fixes

---

**Last Updated:** 2024  
**Tested with:** Go 1.24.9 + golangci-lint v1.54.2
