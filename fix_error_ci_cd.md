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

## Liên hệ
- Nếu vẫn lỗi, đính kèm log step fail (trước và sau fix) để truy vết nhanh.
