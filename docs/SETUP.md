# Hướng dẫn Cài đặt IAM Service

## Yêu cầu hệ thống

### Phần mềm cần thiết

1. **Go** - Version 1.21 hoặc cao hơn
   - Download: https://golang.org/dl/
   - Kiểm tra version: `go version`

2. **PostgreSQL** - Version 12 hoặc cao hơn
   - Download: https://www.postgresql.org/download/
   - Kiểm tra version: `psql --version`

3. **Protocol Buffer Compiler (protoc)** - Version 3.x
   - Download: https://github.com/protocolbuffers/protobuf/releases
   - Kiểm tra version: `protoc --version`

4. **grpcurl** (Optional - dùng để test)
   - Install: `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`

### Cài đặt Go plugins cho protoc

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Đảm bảo rằng `$GOPATH/bin` nằm trong `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Các bước cài đặt

### Bước 1: Clone hoặc Navigate đến project

```bash
cd ecommerce/back_end/iam-services
```

### Bước 2: Cài đặt Go dependencies

```bash
go mod download
go mod tidy
```

### Bước 3: Setup PostgreSQL Database

#### 3.1. Tạo database

```bash
# Đăng nhập PostgreSQL
psql -U postgres

# Tạo database
CREATE DATABASE iam_db;

# Tạo user (optional)
CREATE USER iam_user WITH PASSWORD 'your_password';

# Grant permissions
GRANT ALL PRIVILEGES ON DATABASE iam_db TO iam_user;

# Exit
\q
```

#### 3.2. Chạy migrations

Chạy schema migration:

```bash
psql -U postgres -d iam_db -f migrations/001_init_schema.sql
```

Chạy seed data:

```bash
psql -U postgres -d iam_db -f migrations/002_seed_data.sql
```

#### 3.3. Verify database

```bash
psql -U postgres -d iam_db

# List tables
\dt

# Check sample data
SELECT * FROM roles;
SELECT * FROM permissions;

\q
```

### Bước 4: Cấu hình Environment Variables

Copy file `.env.example` sang `.env`:

```bash
cp .env.example .env
```

Chỉnh sửa file `.env`:

```env
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=50051

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=iam_db
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=change-this-to-a-strong-random-secret-key
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_EXPIRATION_HOURS=168

# Log Configuration
LOG_LEVEL=info
LOG_ENCODING=json
```

**⚠️ Lưu ý bảo mật:**
- Đổi `JWT_SECRET` thành một chuỗi ngẫu nhiên mạnh
- Không commit file `.env` lên Git
- Trong production, sử dụng secret management tools

### Bước 5: Generate Protocol Buffer Code

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       pkg/proto/iam.proto
```

Hoặc nếu bạn trên Windows PowerShell:

```powershell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/iam.proto
```

Điều này sẽ tạo ra:
- `pkg/proto/iam.pb.go`
- `pkg/proto/iam_grpc.pb.go`

### Bước 6: Build và Run

#### Option 1: Run trực tiếp với Go

```bash
go run cmd/server/main.go
```

#### Option 2: Build binary rồi chạy

```bash
# Build
go build -o bin/iam-service cmd/server/main.go

# Run
./bin/iam-service
```

Hoặc trên Windows:

```powershell
# Build
go build -o bin/iam-service.exe cmd/server/main.go

# Run
.\bin\iam-service.exe
```

Nếu thành công, bạn sẽ thấy:

```
{"level":"info","ts":...,"msg":"Starting IAM Service..."}
{"level":"info","ts":...,"msg":"Database connected successfully"}
{"level":"info","ts":...,"msg":"IAM Service is running","address":"0.0.0.0:50051"}
```

## Kiểm tra Service

### Test 1: List services

```bash
grpcurl -plaintext localhost:50051 list
```

Output:
```
grpc.reflection.v1alpha.ServerReflection
iam.IAMService
```

### Test 2: Register user

```bash
grpcurl -plaintext -d '{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "full_name": "Test User"
}' localhost:50051 iam.IAMService/Register
```

### Test 3: Login

```bash
grpcurl -plaintext -d '{
  "username": "testuser",
  "password": "password123"
}' localhost:50051 iam.IAMService/Login
```

Nếu thành công, bạn sẽ nhận được access token và refresh token.

### Test 4: Verify token

Thay `YOUR_TOKEN` bằng access token từ login response:

```bash
grpcurl -plaintext -d '{
  "token": "YOUR_TOKEN"
}' localhost:50051 iam.IAMService/VerifyToken
```

## Troubleshooting

### Issue 1: "connection refused"

**Nguyên nhân**: Service chưa chạy hoặc port bị conflict

**Giải pháp**:
- Kiểm tra service đã chạy chưa
- Kiểm tra port 50051 có bị chiếm bởi process khác không
- Thử đổi port trong `.env`

### Issue 2: "failed to connect to database"

**Nguyên nhân**: PostgreSQL không chạy hoặc config sai

**Giải pháp**:
- Kiểm tra PostgreSQL đã chạy: `pg_isready`
- Kiểm tra credentials trong `.env`
- Kiểm tra database đã được tạo chưa

### Issue 3: "table does not exist"

**Nguyên nhân**: Migrations chưa chạy

**Giải pháp**:
```bash
psql -U postgres -d iam_db -f migrations/001_init_schema.sql
psql -U postgres -d iam_db -f migrations/002_seed_data.sql
```

### Issue 4: "invalid token"

**Nguyên nhân**: JWT_SECRET không khớp hoặc token expired

**Giải pháp**:
- Đảm bảo JWT_SECRET trong `.env` không thay đổi
- Login lại để lấy token mới

### Issue 5: Protocol Buffer files không generate

**Nguyên nhân**: protoc hoặc plugins chưa cài đúng

**Giải pháp**:
```bash
# Kiểm tra protoc
protoc --version

# Cài lại plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Kiểm tra PATH
echo $PATH | grep $(go env GOPATH)/bin
```

## Deployment Production

### Checklist

- [ ] Đổi `JWT_SECRET` thành random strong key
- [ ] Set `DB_SSL_MODE=require` cho production database
- [ ] Set `LOG_LEVEL=warn` hoặc `error`
- [ ] Enable database backup
- [ ] Setup monitoring và logging
- [ ] Use connection pooling
- [ ] Setup reverse proxy (nginx/envoy)
- [ ] Enable TLS for gRPC
- [ ] Rate limiting
- [ ] Health checks

### Docker (Optional)

Tạo file `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o bin/iam-service cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/iam-service .
COPY --from=builder /app/.env .

EXPOSE 50051
CMD ["./iam-service"]
```

Build và run:

```bash
docker build -t iam-service .
docker run -p 50051:50051 --env-file .env iam-service
```

## Tài liệu tham khảo

- [Architecture Documentation](ARCHITECTURE.md)
- [API Documentation](API.md)
- [Database Schema](DATABASE.md)
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/quickstart/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

