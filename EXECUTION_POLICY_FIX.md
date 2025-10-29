# 🔧 PowerShell Execution Policy Fix

Hướng dẫn fix lỗi "running scripts is disabled on this system" khi chạy PowerShell scripts.

---

## 🚨 Lỗi Thường Gặp

```
.\scripts\lint.ps1 : File ...\lint.ps1 cannot be loaded because running 
scripts is disabled on this system. For more information, see 
about_Execution_Policies at https:/go.microsoft.com/fwlink/?LinkID=135170.
```

**Nguyên nhân:** PowerShell Execution Policy mặc định là `Restricted`, không cho phép chạy scripts.

---

## ✅ Giải Pháp

### Cách 1: Sử Dụng Batch Files (Khuyến nghị - Dễ nhất)

```cmd
# Thay vì chạy .ps1 files, dùng .bat files
.\scripts\lint.bat
.\scripts\lint.bat -Fast
.\scripts\lint.bat -Fix
.\scripts\check-all.bat
.\scripts\verify.bat
```

**Ưu điểm:**
- ✅ Không cần thay đổi system settings
- ✅ Hoạt động ngay lập tức
- ✅ An toàn (không ảnh hưởng system-wide)

### Cách 2: Bypass cho PowerShell Commands

```powershell
# Thêm -ExecutionPolicy Bypass vào mỗi lệnh
powershell -ExecutionPolicy Bypass -File .\scripts\lint.ps1
powershell -ExecutionPolicy Bypass -File .\scripts\lint.ps1 -Fast
powershell -ExecutionPolicy Bypass -File .\scripts\check-all.ps1
```

### Cách 3: Thay Đổi Execution Policy (Không khuyến nghị)

```powershell
# Chỉ cho current user (tương đối an toàn)
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser

# Sau đó có thể chạy bình thường
.\scripts\lint.ps1
```

**⚠️ Lưu ý:** Cách này thay đổi system settings, có thể ảnh hưởng bảo mật.

---

## 📦 Files Đã Tạo

### Batch Files (Khuyến nghị sử dụng)
- `scripts/lint.bat` - Run lint.ps1
- `scripts/check-all.bat` - Run check-all.ps1  
- `scripts/verify.bat` - Run verify-lint.ps1

### PowerShell Wrapper Scripts
- `scripts/run-lint.ps1` - Wrapper cho lint.ps1
- `scripts/run-check-all.ps1` - Wrapper cho check-all.ps1
- `scripts/run-verify.ps1` - Wrapper cho verify-lint.ps1

---

## 🚀 Usage Examples

### Sử Dụng Batch Files

```cmd
# Basic lint
.\scripts\lint.bat

# Fast lint
.\scripts\lint.bat -Fast

# Lint with auto-fix
.\scripts\lint.bat -Fix

# Lint specific package
.\scripts\lint.bat -Target model

# Full check
.\scripts\check-all.bat

# Verify setup
.\scripts\verify.bat
```

### Sử Dụng PowerShell với Bypass

```powershell
# Lint commands
powershell -ExecutionPolicy Bypass -File .\scripts\lint.ps1
powershell -ExecutionPolicy Bypass -File .\scripts\lint.ps1 -Fast
powershell -ExecutionPolicy Bypass -File .\scripts\lint.ps1 -Fix

# Check all
powershell -ExecutionPolicy Bypass -File .\scripts\check-all.ps1

# Verify
powershell -ExecutionPolicy Bypass -File .\scripts\verify-lint.ps1
```

### Sử Dụng PowerShell Wrappers

```powershell
# Lint commands
.\scripts\run-lint.ps1
.\scripts\run-lint.ps1 -Fast
.\scripts\run-lint.ps1 -Fix

# Check all
.\scripts\run-check-all.ps1

# Verify
.\scripts\run-verify.ps1
```

---

## 🔍 Kiểm Tra Execution Policy

```powershell
# Xem current policy
Get-ExecutionPolicy -List

# Xem policy cho current scope
Get-ExecutionPolicy

# Thay đổi policy (nếu cần)
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

---

## 📋 Comparison Table

| Method | Ease | Security | System Impact | Recommended |
|--------|------|----------|---------------|-------------|
| Batch Files | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | None | ✅ Yes |
| PowerShell Bypass | ⭐⭐⭐ | ⭐⭐⭐⭐ | None | ✅ Yes |
| PowerShell Wrappers | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | None | ✅ Yes |
| Change Policy | ⭐⭐ | ⭐⭐ | System-wide | ❌ No |

---

## 🎯 Quick Reference

### For Daily Use (Recommended)
```cmd
.\scripts\lint.bat -Fast          # Quick lint
.\scripts\check-all.bat           # Full check before push
.\scripts\verify.bat -Quick       # Verify setup
```

### For Development
```cmd
.\scripts\lint.bat -Fix           # Auto-fix issues
.\scripts\lint.bat -Target model  # Lint specific package
.\scripts\lint.bat -Verbose       # Detailed output
```

### For Troubleshooting
```cmd
.\scripts\verify.bat -Verbose     # Detailed verification
.\scripts\check-all.bat -SkipTests # Skip tests
```

---

## 🔗 Related Files

- `scripts/lint.bat` - Main lint runner
- `scripts/check-all.bat` - Pre-push check
- `scripts/verify.bat` - Setup verification
- `LINTING_SETUP.md` - Complete setup guide
- `VERIFICATION_GUIDE.md` - Verification guide

---

**Last Updated:** 2024  
**Tested with:** Windows PowerShell + Execution Policy Restricted
