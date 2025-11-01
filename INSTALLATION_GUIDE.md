# 📦 Installation Guide - golangci-lint for Go 1.19

Hướng dẫn cài đặt golangci-lint cho Go 1.19 (tránh lỗi "module requires Go 1.20").

---

## 🚨 Vấn Đề Với Go 1.19

Khi sử dụng `go install` với Go 1.19, bạn có thể gặp lỗi:
```
module requires Go 1.20
```

**Nguyên nhân:** Các phiên bản golangci-lint gần đây có dependencies yêu cầu Go 1.20+.

**Giải pháp:** Sử dụng **binary download** thay vì `go install`.

---

## ✅ Cách Cài Đặt (Khuyến Nghị)

### Option 1: Script Tự Động (Dễ nhất)

```powershell
# 1. Chạy script install
.\scripts\install-golangci-lint.ps1 -Version "v1.54.2"

# 2. Thêm vào PATH (nếu cần)
$env:PATH += ";E:\go\src\bin"  # Thay đổi theo GOPATH của bạn

# 3. Verify
golangci-lint version
```

**Script sẽ tự động:**
- Download binary từ GitHub releases
- Giải nén vào `GOPATH/bin`
- Verify cài đặt
- Hướng dẫn thêm vào PATH

### Option 2: Download Thủ Công

```powershell
# 1. Download từ GitHub
# https://github.com/golangci/golangci-lint/releases/tag/v1.54.2
# File: golangci-lint-1.54.2-windows-amd64.zip

# 2. Giải nén vào thư mục
# Ví dụ: C:\golangci-lint\

# 3. Thêm vào PATH
# System Properties > Environment Variables > Path > Add: C:\golangci-lint\

# 4. Verify
golangci-lint version
```

### Option 3: go install (Có thể fail)

```powershell
# Chỉ dùng nếu không gặp lỗi "module requires Go 1.20"
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
```

---

## 🔧 Troubleshooting

### Lỗi: "module requires Go 1.20"

**Nguyên nhân:** Dependencies của golangci-lint yêu cầu Go 1.20+.

**Giải pháp:**
```powershell
# Sử dụng binary download thay vì go install
.\scripts\install-golangci-lint.ps1 -Version "v1.54.2"
```

### Lỗi: "go.work file requires go >= 1.21"

**Nguyên nhân:** File `go.mod` yêu cầu Go 1.21.

**Giải pháp:**
```powershell
# Fix go.mod version
(Get-Content go.mod) -replace 'go 1\.21', 'go 1.19' | Set-Content go.mod
```

### Lỗi: "golangci-lint not found"

**Nguyên nhân:** Chưa thêm vào PATH.

**Giải pháp:**
```powershell
# Thêm vào PATH cho session hiện tại
$env:PATH += ";E:\go\src\bin"  # Thay đổi theo GOPATH

# Hoặc thêm vĩnh viễn:
# System Properties > Environment Variables > Path > Add: E:\go\src\bin
```

---

## 📋 Version Compatibility

| Go Version | golangci-lint Version | Method |
|------------|----------------------|---------|
| Go 1.19    | v1.54.2              | Binary download (khuyến nghị) |
| Go 1.19    | v1.54.2              | go install (có thể fail) |
| Go 1.20+   | v1.55.2+             | go install (OK) |
| Go 1.20+   | v1.55.2+             | Binary download (OK) |

---

## 🚀 Quick Start

```powershell
# 1. Install
.\scripts\install-golangci-lint.ps1 -Version "v1.54.2"

# 2. Add to PATH
$env:PATH += ";E:\go\src\bin"

# 3. Test
golangci-lint version

# 4. Run lint
.\scripts\lint.ps1
```

---

## 📚 Related Files

- `scripts/install-golangci-lint.ps1` - Auto install script
- `scripts/lint.ps1` - Lint runner
- `scripts/check-all.ps1` - Pre-push check
- `LINTING_SETUP.md` - Complete linting guide
- `SETUP_COMPLETE.md` - Setup summary

---

**Last Updated:** 2024  
**Tested with:** Go 1.19 + golangci-lint v1.54.2
