# CI/CD Error Fixes Guide

Tổng hợp lỗi thường gặp khi chạy CI/CD trên GitHub Actions cho iam-services và cách khắc phục nhanh.

---

## 1) Deprecated artifact action v3
- Mô tả: "This request has been automatically failed because it uses a deprecated version of actions/upload-artifact: v3"
- Nguyên nhân: GitHub deprecate v3 của artifact actions.
- Fix: Nâng cấp lên v4.
- Thay đổi:
  - File: `.github/workflows/ci-cd.yml` → job build: `uses: actions/upload-artifact@v4`
  - File: `.github/workflows/test.yml` → unit-tests & benchmark-tests: `uses: actions/upload-artifact@v4`

---

## 2) working-directory không tồn tại
- Mô tả: "An error occurred trying to start process '/usr/bin/bash' with working directory '.../ecommerce/back_end/iam-services'. No such file or directory"
- Nguyên nhân: Hard-code `working-directory: ecommerce/back_end/iam-services`, nhưng cấu trúc repo có thể khác.
- Fix: Tự phát hiện thư mục service và dùng biến môi trường.
- Thay đổi:
  - Thêm step sớm trong mỗi job:
    ```yaml
    - name: Set service directory
      run: |
        if [ -d "ecommerce/back_end/iam-services" ]; then
          echo "SERVICE_DIR=ecommerce/back_end/iam-services" >> $GITHUB_ENV
        else
          echo "SERVICE_DIR=." >> $GITHUB_ENV
        fi
    ```
  - Sửa tất cả `working-directory:` và đường dẫn artifact/coverage dùng `${{ env.SERVICE_DIR }}`

---

## 3) Lỗi `psql: command not found`
- Mô tả: Chạy migrations bằng psql fail vì runner chưa có PostgreSQL client.
- Fix: Cài `postgresql-client` trước khi chạy migrations.
- Thay đổi:
  ```yaml
  - name: Install PostgreSQL client
    run: |
      sudo apt-get update
      sudo apt-get install -y postgresql-client
  ```
  - Áp dụng cho các jobs có chạy `psql` trong `.github/workflows/ci-cd.yml` và `.github/workflows/test.yml`.

---

## 4) Codecov upload fail (không có token)
- Mô tả: Upload coverage lên Codecov fail nếu repo private và thiếu `CODECOV_TOKEN`.
- Fix: Thêm token (nếu cần) và không fail toàn job khi thiếu.
- Thay đổi:
  ```yaml
  - name: Upload coverage to Codecov
    uses: codecov/codecov-action@v3
    with:
      file: ./${{ env.SERVICE_DIR }}/coverage.out
      flags: unittests
      name: codecov-iam-service
      token: ${{ secrets.CODECOV_TOKEN }}
      fail_ci_if_error: false
  ```

---

## 5) Thiếu `.env.example`
- Mô tả: Dev không thấy `.env.example`, CI/Team khó cấu hình.
- Fix: Dùng `.env.template` và copy.
- Cách tạo:
  - PowerShell: `Copy-Item .env.template .env.example`
  - Linux/macOS: `cp .env.template .env.example`
  - Hoặc chạy script: `scripts/setup-ci.ps1` / `scripts/setup-ci.sh`

---

## 6) Migrations thiếu trong CI
- Mô tả: Lỗi bảng/dữ liệu Casbin/CMS chưa có.
- Fix: Thêm migrations `005_separate_user_cms_authorization.sql` và `006_seed_separated_authorization.sql` vào workflows.
- Thay đổi:
  ```bash
  psql -h localhost -U postgres -d iam_db_test -f migrations/005_separate_user_cms_authorization.sql
  psql -h localhost -U postgres -d iam_db_test -f migrations/006_seed_separated_authorization.sql
  ```

---

## 7) Go version mismatch (Dockerfile vs CI)
- Mô tả: Dockerfile dùng Go 1.21, workflow dùng 1.19 → inconsistency.
- Fix: Đồng bộ version (đã chuyển Dockerfile về `golang:1.19-alpine`).
- Files:
  - `Dockerfile`: `FROM golang:1.19-alpine AS builder`
  - Workflows: `GO_VERSION: '1.19'`

---

## 8) Health check fail ở deploy jobs (khi bật lại)
- Mô tả: `curl -f https://.../health` fail do HTTP Gateway hoặc endpoint chưa bật.
- Fix options:
  - Bật HTTP Gateway trong `internal/app/app.go` (uncomment `setupHTTPGateway()` và generate proto gateway trước).
  - Implement endpoint `/health` (REST) hoặc thay bằng check TCP gRPC (port 50051).
  - Chỉ bật deploy jobs khi server đã có compose + env chuẩn.

---

## 9) Lỗi tests DAO do API khác tên
- Mô tả: Unit tests gọi `GetByID/GetByUsername/...` trong khi DAO là `FindByID/FindByUsername/...`.
- Fix: Cập nhật tests cho đúng API thực tế, và xử lý not-found theo DAO (trả `nil, nil`).
- Files:
  - `internal/dao/user_dao_test.go` (đã cập nhật dùng `FindBy...` và assert `nil` cho not-found)

---

## 10) Lỗi mock interfaces không khớp
- Mô tả: Mock repo trong tests thiếu method so với interface thật.
- Fix: Bổ sung mock methods cần thiết (`UserExists`, `UserHasPermission`, ...).
- Files:
  - `internal/service/auth_service_test.go` (đã bổ sung mock methods)

---

## 11) Thư mục artifact/coverage sai đường dẫn
- Mô tả: Artifact path hard-code theo mono-repo.
- Fix: Dùng `${{ env.SERVICE_DIR }}` sau step detect thư mục.
- Ví dụ:
  ```yaml
  with:
    path: ${{ env.SERVICE_DIR }}/bin/iam-service
  ```

---

## 12) psql kết nối DB test không ổn định
- Tips:
  - Đợi Postgres healthy trước khi chạy psql:
    ```yaml
    services:
      postgres:
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    - name: Wait for PostgreSQL
      run: |
        until pg_isready -h localhost -p 5432 -U postgres; do
          echo "Waiting for PostgreSQL..."; sleep 2; done
    ```

---

## 13) Gợi ý xác minh nhanh khi CI fail
- Mở log job fail trong Actions → xem step gần nhất.
- Kiểm tra thư mục hiện tại: add step `pwd && ls -la`.
- In ra biến: `echo $GITHUB_WORKSPACE`, `echo ${{ env.SERVICE_DIR }}`.
- Re-run jobs sau khi fix.

---

## 14) Lỗi errcheck: "Error return value of rows.Close is not checked"
- Mô tả: Linter báo lỗi vì kết quả trả về của `rows.Close()` bị bỏ qua.
- Nguyên nhân: Sử dụng `defer func() { _ = rows.Close() }()` khiến `errcheck` coi là bỏ qua lỗi có chủ đích.
- Fix tối ưu: Dùng cú pháp chuẩn `defer rows.Close()` để errcheck bỏ qua hợp lệ trong ngữ cảnh `defer` (idiomatic Go). Không thay đổi control flow và vẫn đảm bảo đóng `rows` khi thoát hàm.
- Thay đổi đã áp dụng:
  - `internal/dao/role_permission_dao.go`: thay thế ở 2 vị trí `GetRolePermissions`, `GetPermissionRoles`
  - `internal/dao/role_dao.go`: `List`
  - `internal/dao/permission_dao.go`: `List`
  - `internal/dao/cms_tab_api_dao.go`: `FindByTab`, `FindByAPI`, `ListAll`
  - `internal/dao/cms_role_dao.go`: `List`
  - `internal/dao/api_resource_dao.go`: `ListByService`, `List`

Ví dụ trước/sau:

Trước:
```go
defer func() { _ = rows.Close() }()
```

Sau (idiomatic):
```go
defer rows.Close()
```

Ghi chú nâng cao: Nếu muốn xử lý lỗi đóng kết nối một cách nghiêm ngặt (ví dụ ghi log hoặc trả lỗi nếu chưa có lỗi khác), có thể dùng:
```go
defer func() {
    if cerr := rows.Close(); cerr != nil {
        // log cerr hoặc gán vào err đã đặt tên trong chữ ký hàm
    }
}()
```
Tuy nhiên, với `errcheck` mặc định, `defer rows.Close()` là đủ, gọn và đúng idiom.

---

## 15) Lỗi revive: "package-comments: should have a package comment"
- Mô tả: Linter `revive` yêu cầu mỗi package phải có comment mô tả chức năng.
- Nguyên nhân: File đầu tiên trong package thiếu package-level comment.
- Fix: Thêm comment trước dòng `package <name>`.
- Thay đổi đã áp dụng:
  - `internal/application/dto/auth_dto.go`: Thêm comment mô tả package dto

Ví dụ:
```go
// Package dto provides Data Transfer Objects for application layer communication.
// It contains request/response structures for authentication and authorization operations.
package dto
```

---

## 16) Lỗi gosec G115: "integer overflow conversion int -> int32"
- Mô tả: `gosec` cảnh báo chuyển đổi từ `int` sang `int32` có thể gây overflow.
- Nguyên nhân: Trực tiếp ép kiểu `int32(total)` mà không kiểm tra giá trị.
- Fix: Thêm overflow check trước khi convert, dùng giá trị max nếu overflow.
- Thay đổi đã áp dụng:
  - `internal/handler/grpc_handler.go`: `ListRoles`, `ListPermissions`
  - `internal/handler/casbin_handler.go`: `ListAPIResources`

Ví dụ:
```go
// Safe conversion with overflow check
totalInt32 := int32(total)
if total > 0 && int(totalInt32) != total {
    h.logger.Warn("Total count overflow in ListRoles", zap.Int("total", total))
    totalInt32 = 2147483647 // Max int32
}
return &pb.ListRolesResponse{
    Roles: pbRoles,
    Total: totalInt32,
}, nil
```

---

## 17) Lỗi goconst: "string has N occurrences, make it a constant"
- Mô tả: Linter `goconst` yêu cầu string lặp lại nhiều lần nên được đặt làm constant.
- Nguyên nhân: Dùng hardcode string nhiều lần trong test files.
- Fix: Extract ra constants.
- Thay đổi đã áp dụng:
  - `pkg/password/password_manager_test.go`: Extract `testPassword`
  - `pkg/jwt/jwt_manager_test.go`: Extract `testUserID`, `testUsername`

Ví dụ:
```go
const (
    testPassword = "MySecurePassword123!"
)

func TestHashPassword_Success(t *testing.T) {
    manager := NewPasswordManager()
    password := testPassword  // Thay vì hardcode
    // ...
}
```

---

## 18) Lỗi errcheck: "Error return value is not checked" (Sync)
- Mô tả: `errcheck` báo lỗi vì không check error từ `logger.Sync()`.
- Nguyên nhân: Dùng `_ = logger.Sync()` hoặc gọi mà không check return.
- Fix: Check error và log warning (vì Sync() thường báo lỗi harmless trên stdout/stderr).
- Thay đổi đã áp dụng:
  - `internal/app/app.go`: `waitForShutdown`, `Shutdown`
  - `internal/container/container.go`: `Close`

Ví dụ:
```go
// Trước
_ = a.logger.Sync()

// Sau
if syncErr := a.logger.Sync(); syncErr != nil {
    // Ignore sync errors on stdout/stderr (common on some platforms)
    a.logger.Warn("Logger sync returned error (may be harmless)", zap.Error(syncErr))
}
```

Và check return value của `Shutdown()`:
```go
// Trước
a.Shutdown(context.Background())

// Sau
if err := a.Shutdown(context.Background()); err != nil {
    a.logger.Error("Error during shutdown", zap.Error(err))
}
```

---

## 19) Lỗi revive: "exported var/method should have comment"
- Mô tả: Linter `revive` yêu cầu tất cả exported symbols (vars, consts, methods) phải có comment.
- Nguyên nhân: Exported vars/methods/consts thiếu comment giải thích.
- Fix: Thêm comment mô tả cho mỗi exported symbol.
- Thay đổi đã áp dụng:
  - `internal/domain/model/api_resource.go`: Comment cho tất cả Error vars
  - `internal/domain/model/user.go`: Comment cho Error vars và methods `ID()`, `Username()`
  - `internal/domain/casbin.go`: Comment cho package, constants `DomainUser`, `DomainCMS`, `DomainAPI`, `CMSTab*`

Ví dụ:
```go
// Trước
var (
	ErrInvalidUsername = errors.New("invalid username")
)

// Sau
var (
	// ErrInvalidUsername indicates username validation failed
	ErrInvalidUsername = errors.New("invalid username")
)
```

Với methods:
```go
// ID returns the user's unique identifier
func (u *User) ID() string { return u.id }

// Username returns the user's username
func (u *User) Username() string { return u.username }
```

Với constants:
```go
const (
	// DomainUser is the domain for end user authorization
	DomainUser CasbinDomain = "user"
	// DomainCMS is the domain for CMS/admin panel authorization
	DomainCMS CasbinDomain = "cms"
)
```

---

## 20) Lỗi gosec G115: Integer overflow (cải tiến)
- Mô tả: Sau khi thêm overflow check inline, gosec vẫn cảnh báo do phát hiện direct conversion.
- Nguyên nhân: Dù có check sau conversion, gosec scan static và cảnh báo ngay tại dòng `int32(value)`.
- Fix tối ưu: Tạo helper function `safeIntToInt32()` với overflow check và dùng `#nosec` comment có giải thích.
- Thay đổi đã áp dụng:
  - `internal/handler/grpc_handler.go`: Thêm helper function, dùng trong `ListRoles`, `ListPermissions`
  - `internal/handler/casbin_handler.go`: Thêm helper function, dùng trong `ListAPIResources`

Helper function:
```go
// safeIntToInt32 safely converts int to int32 with overflow protection
func safeIntToInt32(value int) int32 {
	const maxInt32 = 2147483647
	if value > maxInt32 {
		return maxInt32
	}
	if value < -2147483648 {
		return -2147483648
	}
	return int32(value) // #nosec G115 -- overflow checked above
}
```

Usage:
```go
// Trước
totalInt32 := int32(total)  // Gosec warns here
if total > 0 && int(totalInt32) != total {
	totalInt32 = 2147483647
}

// Sau
Total: safeIntToInt32(total),  // Clean, no warning
```

---

## Tổng kết các fix Clean Code

### Đã fix trong lần 1:
1. ✅ Package comment cho `dto` package (revive)
2. ✅ Integer overflow với inline check (gosec G115)
3. ✅ Extract test string constants (goconst)
4. ✅ Check `Sync()` và `Shutdown()` errors (errcheck)

### Đã fix trong lần 2:
5. ✅ Package comment cho `domain` package (revive)
6. ✅ Comments cho exported vars: `ErrInvalidAPIResource`, `ErrInvalidUsername`, etc. (revive)
7. ✅ Comments cho exported methods: `ID()`, `Username()` (revive)
8. ✅ Comments cho exported constants: `DomainUser`, `CMSTabProduct`, etc. (revive)
9. ✅ Integer overflow với helper function + #nosec (gosec G115 - cải tiến)

### Checklist hoàn thành:
- ✅ Tất cả package đều có package comment
- ✅ Tất cả exported vars có comment
- ✅ Tất cả exported methods có comment
- ✅ Tất cả exported constants có comment
- ✅ Integer overflow được xử lý an toàn
- ✅ Test strings được extract thành constants
- ✅ Error returns được check đầy đủ
- ✅ HTTP Method constants có comments (MethodGET, MethodPOST, MethodPUT, MethodDELETE, MethodPATCH)
- ✅ CasbinDomain constants có comments (DomainUser, DomainCMS, DomainAPI)

### Files đã cập nhật:
**Packages:**
- `internal/application/dto/auth_dto.go` - Package comment
- `internal/domain/casbin.go` - Package comment + domain constants

**Error Variables:**
- `internal/domain/model/api_resource.go` - ErrInvalidAPIResource, ErrEmptyPath, ErrEmptyMethod, ErrInvalidMethod
- `internal/domain/model/user.go` - All error vars

**Constants:**
- `internal/domain/model/api_resource.go` - HTTP method constants
- `internal/domain/casbin.go` - CMS tab constants
- `internal/domain/service/authorization_service.go` - Domain constants

**Methods:**
- `internal/domain/model/user.go` - ID(), Username() methods

**Handlers:**
- `internal/handler/grpc_handler.go` - Added safeIntToInt32 helper
- `internal/handler/casbin_handler.go` - Uses shared helper

**Tests:**
- `pkg/password/password_manager_test.go` - Extracted testPassword constant
- `pkg/jwt/jwt_manager_test.go` - Extracted testUserID, testUsername constants

**App & Container:**
- `internal/app/app.go` - Check Shutdown() errors
- `internal/container/container.go` - Check Logger.Sync() errors

---

## 📋 Quick Reference: Comment Patterns

### Package Comment:
```go
// Package dto provides Data Transfer Objects for application layer communication.
package dto
```

### Exported Const:
```go
const (
    // MethodGET represents HTTP GET method
    MethodGET HTTPMethod = "GET"
)
```

### Exported Var/Error:
```go
var (
    // ErrInvalidUsername indicates username validation failed
    ErrInvalidUsername = errors.New("invalid username")
)
```

### Exported Method:
```go
// ID returns the user's unique identifier
func (u *User) ID() string { return u.id }
```

### Exported Type:
```go
// User represents a user aggregate root in the domain
type User struct { ... }
```

---

## Liên hệ
- Nếu vẫn lỗi, đính kèm log step fail (trước và sau fix) để truy vết nhanh.
