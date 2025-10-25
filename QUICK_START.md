# Quick Start Guide - IAM Service

## Yêu cầu

- **Go**: 1.19 hoặc cao hơn
- **PostgreSQL**: 12 hoặc cao hơn  
- **Protoc**: 3.x (Protocol Buffer Compiler)

## Setup nhanh (Windows PowerShell)

### Option 1: Sử dụng Setup Script (Khuyến nghị)

```powershell
# Di chuyển đến thư mục project
cd ecommerce\back_end\iam-services

# Chạy setup script
.\scripts\setup-proto.ps1
```

Script sẽ tự động:
1. ✅ Cài đặt protoc-gen-go (v1.28.1)
2. ✅ Cài đặt protoc-gen-go-grpc (v1.2.0)
3. ✅ Generate proto files
4. ✅ Chạy go mod tidy

### Option 2: Setup thủ công

#### Bước 1: Cài đặt Go plugins

```powershell
# Cài protoc-gen-go (version tương thích với Go 1.19)
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

# Cài protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
```

#### Bước 2: Thêm GOPATH/bin vào PATH (nếu chưa có)

```powershell
# Kiểm tra GOPATH
go env GOPATH

# Thêm vào PATH (PowerShell - session hiện tại)
$env:PATH += ";$(go env GOPATH)\bin"

# Hoặc thêm vĩnh viễn qua System Properties > Environment Variables
# Thêm: C:\Users\<YourUsername>\go\bin vào PATH
```

#### Bước 3: Generate proto files

```powershell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/iam.proto
```

Sẽ tạo ra:
- `pkg/proto/iam.pb.go`
- `pkg/proto/iam_grpc.pb.go`

#### Bước 4: Download dependencies

```powershell
go mod download
go mod tidy
```

## Setup Database

### Tạo Database

```sql
-- Kết nối PostgreSQL
psql -U postgres

-- Tạo database
CREATE DATABASE iam_db;

-- Tạo user (optional)
CREATE USER iam_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE iam_db TO iam_user;

-- Exit
\q
```

### Chạy Migrations

```powershell
# Schema migration
psql -U postgres -d iam_db -f migrations\001_init_schema.sql

# Seed data
psql -U postgres -d iam_db -f migrations\002_seed_data.sql

# Casbin tables
psql -U postgres -d iam_db -f migrations\003_casbin_tables.sql

# Casbin seed data
psql -U postgres -d iam_db -f migrations\004_casbin_seed_data.sql
```

## Configuration

Copy và chỉnh sửa file config:

```powershell
# Copy template
copy .env.example .env

# Chỉnh sửa .env
notepad .env
```

Cập nhật các giá trị:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=iam_db

# JWT
JWT_SECRET=your-strong-secret-key-change-this
```

## Build & Run

### Run Development Mode

```powershell
go run cmd/server/main.go
```

### Build Binary

```powershell
# Build
go build -o bin\iam-service.exe cmd\server\main.go

# Run
.\bin\iam-service.exe
```

## Testing

### Test với grpcurl

```powershell
# Install grpcurl (if not installed)
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List services
grpcurl -plaintext localhost:50051 list

# Register user
grpcurl -plaintext -d '{\"username\": \"testuser\", \"email\": \"test@example.com\", \"password\": \"password123\", \"full_name\": \"Test User\"}' localhost:50051 iam.IAMService/Register

# Login
grpcurl -plaintext -d '{\"username\": \"testuser\", \"password\": \"password123\"}' localhost:50051 iam.IAMService/Login
```

## Troubleshooting

### Lỗi: "protoc-gen-go: Plugin failed"

**Nguyên nhân**: Chưa cài protoc-gen-go hoặc không có trong PATH

**Giải pháp**:
```powershell
# Cài lại plugin
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

# Kiểm tra PATH
$env:PATH -split ';' | Select-String 'go\\bin'

# Thêm vào PATH nếu chưa có
$env:PATH += ";$(go env GOPATH)\bin"
```

### Lỗi: "go.mod file indicates go 1.21, but maximum version supported is 1.19"

**Giải pháp**: Đã fix trong go.mod (go 1.19)

### Lỗi: "cannot find module providing package pkg/proto"

**Nguyên nhân**: Proto files chưa được generate

**Giải pháp**: Generate proto files trước (xem bước 3)

### Lỗi: "connection refused" khi test

**Nguyên nhân**: Service chưa chạy hoặc port conflict

**Giải pháp**:
```powershell
# Kiểm tra service đang chạy
netstat -ano | findstr :50051

# Kill process nếu cần
taskkill /PID <PID> /F

# Chạy lại service
go run cmd/server/main.go
```

## Docker Setup (Optional)

```powershell
# Build image
docker build -t iam-service .

# Run with docker-compose
docker-compose up -d

# Check logs
docker-compose logs -f iam-service
```

## Makefile Commands

```powershell
# Nếu có Make trên Windows (hoặc dùng Git Bash)

make proto          # Generate proto files
make build          # Build binary
make run            # Run service
make test           # Run tests
make clean          # Clean build artifacts
make db-migrate     # Run migrations
```

## Next Steps

1. ✅ Setup xong? → Đọc [API Documentation](docs/API.md)
2. 🔐 Setup Casbin? → Đọc [Casbin Guide](docs/CASBIN.md)
3. 🏗️ Refactor code? → Đọc [Refactoring Guide](REFACTORING_GUIDE.md)
4. 📚 Hiểu architecture? → Đọc [Architecture](ARCHITECTURE_NEW.md)

## Support

- Architecture: `ARCHITECTURE_NEW.md`
- Refactoring: `REFACTORING_GUIDE.md`
- Casbin: `docs/CASBIN.md`
- Database: `docs/DATABASE.md`

