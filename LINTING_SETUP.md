# 🔍 Linting Setup Guide - IAM Services

Hướng dẫn cài đặt và sử dụng golangci-lint cho dự án IAM Services.

---

## 📦 Cài Đặt golangci-lint

### Cách 1: Sử dụng Go Install (Khuyến nghị)

```powershell
# Cài đặt phiên bản tương thích với Go 1.19
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
```

**Lưu ý về phiên bản:**
- **Go 1.19**: Sử dụng `v1.54.2` (khuyến nghị)
- **Go 1.20+**: Có thể dùng `v1.55.2` hoặc mới hơn
- Quá trình cài đặt có thể mất 2-5 phút để tải về tất cả dependencies.

### Cách 2: Download Binary (Khuyến nghị cho Go 1.19)

**Option 2a: Sử dụng Script Tự Động (Dễ nhất)**
```powershell
# Chạy script tự động download và cài đặt
.\scripts\install-golangci-lint.ps1 -Version "v1.54.2"

# Script sẽ:
# 1. Download binary từ GitHub
# 2. Giải nén vào GOPATH/bin
# 3. Verify cài đặt
# 4. Hướng dẫn thêm vào PATH
```

**Option 2b: Download Thủ Công**
Tải trực tiếp từ GitHub Releases:
- Go 1.19: https://github.com/golangci/golangci-lint/releases/tag/v1.54.2
- Go 1.20+: https://github.com/golangci/golangci-lint/releases/tag/v1.55.2
- Chọn file phù hợp với hệ điều hành (Windows: `windows-amd64.zip`)
- Giải nén và thêm vào PATH

### Kiểm Tra Cài Đặt

```powershell
# Kiểm tra version
golangci-lint version

# Nếu thấy: golangci-lint has version 1.54.2 (hoặc 1.55.2)
# => Cài đặt thành công!
```

---

## 🚀 Sử Dụng

### Cách 1: Sử Dụng Scripts (Dễ nhất)

```powershell
# Chạy lint tất cả
.\scripts\lint.ps1

# Chạy lint với auto-fix
.\scripts\lint.ps1 -Fix

# Chạy lint nhanh (bỏ qua slow linters)
.\scripts\lint.ps1 -Fast

# Chạy lint cho package cụ thể
.\scripts\lint.ps1 -Target model
.\scripts\lint.ps1 -Target handler
.\scripts\lint.ps1 -Target dao

# Chạy kiểm tra đầy đủ trước khi push
.\scripts\check-all.ps1
```

### Cách 2: Sử Dụng Makefile (Linux/Mac hoặc Windows với Make)

```bash
# Xem các lệnh có sẵn
make help

# Chạy lint
make lint

# Chạy lint với auto-fix
make lint-fix

# Chạy lint nhanh
make lint-fast

# Chạy lint cho package cụ thể
make lint-model
make lint-handler
make lint-dao

# Chạy kiểm tra đầy đủ
make check-all
```

### Cách 3: Chạy Trực Tiếp

```powershell
# Chạy lint với config
golangci-lint run --config .golangci.yml

# Chạy lint với auto-fix
golangci-lint run --config .golangci.yml --fix

# Chạy lint nhanh
golangci-lint run --config .golangci.yml --fast

# Chạy lint cho package cụ thể
golangci-lint run --config .golangci.yml internal/domain/model/...
```

---

## ⚙️ Configuration

File cấu hình: `.golangci.yml`

### Linters Enabled

- ✅ **revive** - Code quality & style (package comments, exported symbols)
- ✅ **errcheck** - Check error returns
- ✅ **gosec** - Security issues (G115 integer overflow)
- ✅ **goconst** - Repeated strings
- ✅ **gofmt** - Code formatting
- ✅ **goimports** - Import formatting
- ✅ **misspell** - Spelling mistakes
- ✅ **staticcheck** - Static analysis
- ✅ **typecheck** - Type errors
- ✅ **govet** - Go vet analysis

### Linters Settings

```yaml
revive:
  - exported: Check exported symbols have comments
  - package-comments: Check packages have comments

errcheck:
  - check-blank: true
  - Ignore rows.Close() in defer

gosec:
  - exclude G115 (we handle manually with #nosec)

goconst:
  - min-len: 3
  - min-occurrences: 3
  - ignore-tests: true
```

---

## 🎯 Workflow Khuyến Nghị

### Trước Khi Commit

```powershell
# Option 1: Chạy script tự động (khuyến nghị)
.\scripts\check-all.ps1

# Option 2: Chạy từng bước
.\scripts\lint.ps1
go build ./...
go test ./...
```

### Khi Fix Lỗi

```powershell
# 1. Chạy lint để xem lỗi
.\scripts\lint.ps1

# 2. Fix lỗi trong code

# 3. Chạy lint với auto-fix (nếu có thể)
.\scripts\lint.ps1 -Fix

# 4. Kiểm tra lại
.\scripts\lint.ps1
```

### Khi Đẩy Code Lên GitHub

```powershell
# 1. Chạy kiểm tra đầy đủ
.\scripts\check-all.ps1

# 2. Nếu pass, commit và push
git add .
git commit -m "fix: resolve linting issues"
git push origin your-branch

# 3. GitHub Actions sẽ tự động chạy lại
```

---

## 🐛 Troubleshooting

### Lỗi: "golangci-lint not found"

**Nguyên nhân:** Chưa cài đặt hoặc chưa add vào PATH.

**Giải pháp:**
```powershell
# Cài đặt lại (chọn version phù hợp với Go version)
# Go 1.19:
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
# Go 1.20+:
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

# Kiểm tra GOPATH
go env GOPATH

# Thêm vào PATH (nếu cần)
# Windows: Thêm %GOPATH%\bin vào PATH
# Linux/Mac: export PATH=$PATH:$(go env GOPATH)/bin
```

### Lỗi: "invalid go version" hoặc "module requires Go 1.20"

**Nguyên nhân:** Phiên bản golangci-lint không tương thích với Go version.

**Giải pháp:**
```powershell
# Kiểm tra Go version
go version

# Option 1: Sử dụng script download binary (khuyến nghị)
.\scripts\install-golangci-lint.ps1 -Version "v1.54.2"

# Option 2: go install (có thể fail với Go 1.19)
# Go 1.19: Sử dụng v1.54.2
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

# Go 1.20+: Có thể dùng v1.55.2
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
```

### Lỗi: "go.work file requires go >= 1.21"

**Nguyên nhân:** File `go.mod` yêu cầu Go 1.21 nhưng bạn có Go 1.19.

**Giải pháp:**
```powershell
# Fix go.mod version
# Mở file go.mod và thay đổi:
# go 1.21 → go 1.19

# Hoặc dùng lệnh:
(Get-Content go.mod) -replace 'go 1\.21', 'go 1.19' | Set-Content go.mod
```

### Lỗi: Too many errors

**Nguyên nhân:** Có quá nhiều lỗi lint.

**Giải pháp:**
```powershell
# Fix từng package một
.\scripts\lint.ps1 -Target model -Fix
.\scripts\lint.ps1 -Target handler -Fix
.\scripts\lint.ps1 -Target dao -Fix

# Hoặc fix tất cả cùng lúc
.\scripts\lint.ps1 -Fix
```

### Lỗi: Execution Policy (PowerShell)

**Nguyên nhân:** PowerShell không cho phép chạy scripts.

**Giải pháp:**
```powershell
# Chạy với bypass
powershell -ExecutionPolicy Bypass -File .\scripts\lint.ps1

# Hoặc thay đổi policy
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

---

## 📊 Output Format

### Khi Có Lỗi

```
internal/dao/user_dao.go:95:2: Error return value of rows.Close is not checked (errcheck)
internal/domain/model/user.go:9:5: exported var ErrInvalidUsername should have comment or be unexported (revive)
internal/handler/grpc_handler.go:347:24: G115: integer overflow conversion int -> int32 (gosec)
```

**Format:** `file:line:column: message (linter)`

### Khi Pass

```
SUCCESS: Linting passed! Code is clean.
```

---

## 🚀 Quick Start Workflow

### First Time Setup

```powershell
# 1. Cài đặt golangci-lint (chọn 1 trong 2 cách)

# Cách 1: Script tự động (khuyến nghị)
.\scripts\install-golangci-lint.ps1 -Version "v1.54.2"

# Cách 2: go install (có thể fail với Go 1.19)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

# 2. Thêm vào PATH (nếu cần)
$env:PATH += ";E:\go\src\bin"  # Thay đổi theo GOPATH của bạn

# 3. Verify cài đặt
golangci-lint version
```

### Before Committing Code

```powershell
# Run full check
.\scripts\check-all.ps1

# If pass, commit
git add .
git commit -m "fix: your message"
git push
```

### During Development

```powershell
# Quick lint check
.\scripts\lint.ps1 -Fast

# Fix specific package
.\scripts\lint.ps1 -Target handler -Fix

# Run tests
go test ./...
```

---

## 📚 Resources

- [golangci-lint Documentation](https://golangci-lint.run/)
- [Enabled Linters](https://golangci-lint.run/usage/linters/)
- [Configuration Reference](https://golangci-lint.run/usage/configuration/)
- [GitHub Actions Integration](https://golangci-lint.run/usage/install/#github-actions)

---

## 🔗 Related Files

- `.golangci.yml` - Cấu hình lint
- `scripts/lint.ps1` - Script lint chính
- `scripts/check-all.ps1` - Script kiểm tra đầy đủ
- `Makefile` - Make commands
- `fix_error_ci_cd.md` - Log các lỗi đã fix

---

**Last Updated:** 2024  
**Maintainer:** IAM Service Team

