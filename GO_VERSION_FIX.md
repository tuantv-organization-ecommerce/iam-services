# 🔧 Go Version Compatibility Fix

Hướng dẫn fix lỗi compatibility khi sử dụng Go 1.24 với golangci-lint.

---

## 🚨 Vấn Đề Với Go 1.24

**Lỗi thường gặp:**
```
level=error msg="Running error: context loading failed: failed to load packages: 
failed to load with go/packages: err: exit status 1: stderr: go: module . listed 
in go.work file requires go >= 1.24, but go.work lists go 1.19; to update it:\n\tgo work use\n"
```

**Nguyên nhân:** Go version mới (1.24) nhưng go.work file vẫn yêu cầu Go 1.19.

---

## ✅ Giải Pháp

### 1. Update go.work File

```bash
# File: ecommerce/back_end/go.work
# Change from:
go 1.19

# To:
go 1.24
```

### 2. Update go.mod Files

```bash
# File: ecommerce/back_end/iam-services/go.mod
# Change from:
go 1.19

# To:
go 1.24
```

### 3. Install Latest golangci-lint

```powershell
# Go 1.24 supports latest golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 4. Fix Batch Files PATH

Updated batch files now automatically add GOPATH/bin to PATH:

```batch
REM Add GOPATH/bin to PATH for this session
for /f "tokens=*" %%i in ('go env GOPATH') do set GOPATH=%%i
set PATH=%PATH%;%GOPATH%\bin
```

---

## 📦 Files Updated

### Version Files
- ✅ `ecommerce/back_end/go.work` - Updated to go 1.24
- ✅ `ecommerce/back_end/iam-services/go.mod` - Updated to go 1.24

### Batch Files (Auto PATH)
- ✅ `scripts/lint.bat` - Auto-adds GOPATH/bin to PATH
- ✅ `scripts/check-all.bat` - Auto-adds GOPATH/bin to PATH
- ✅ `scripts/verify.bat` - Auto-adds GOPATH/bin to PATH

---

## 🚀 Usage After Fix

### Test Linting

```cmd
# Quick lint test
.\scripts\lint.bat -Fast

# Full lint
.\scripts\lint.bat

# Auto-fix
.\scripts\lint.bat -Fix
```

### Verify Setup

```cmd
# Verify everything works
.\scripts\verify.bat -Quick
```

### Full Check

```cmd
# Pre-push check
.\scripts\check-all.bat
```

---

## 🔍 Troubleshooting

### Issue: "golangci-lint not found"

**Solution:**
```cmd
# Batch files now auto-add PATH, but if still failing:
go env GOPATH
# Add the result\bin to your system PATH
```

### Issue: "go.work file requires go >= 1.24"

**Solution:**
```bash
# Update go.work file
# Change: go 1.19
# To:     go 1.24
```

### Issue: "module requires go >= 1.24"

**Solution:**
```bash
# Update go.mod file
# Change: go 1.19  
# To:     go 1.24
```

---

## 📊 Version Compatibility

| Go Version | golangci-lint | go.mod | go.work | Status |
|------------|---------------|--------|---------|--------|
| 1.19       | v1.54.2       | 1.19   | 1.19    | ✅ OK  |
| 1.24       | latest        | 1.24   | 1.24    | ✅ OK  |
| 1.24       | v1.54.2       | 1.19   | 1.19    | ❌ FAIL |
| 1.24       | latest        | 1.19   | 1.19    | ❌ FAIL |

---

## 🎯 Quick Fix Commands

```powershell
# 1. Update go.work
(Get-Content "..\go.work") -replace 'go 1\.19', 'go 1.24' | Set-Content "..\go.work"

# 2. Update go.mod  
(Get-Content "go.mod") -replace 'go 1\.19', 'go 1.24' | Set-Content "go.mod"

# 3. Install latest golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 4. Test
.\scripts\lint.bat -Fast
```

---

## 📚 Related Files

- `go.work` - Workspace configuration
- `go.mod` - Module configuration  
- `scripts/lint.bat` - Lint runner with auto PATH
- `scripts/verify.bat` - Setup verification
- `EXECUTION_POLICY_FIX.md` - PowerShell issues

---

**Last Updated:** 2024  
**Tested with:** Go 1.24 + golangci-lint latest
