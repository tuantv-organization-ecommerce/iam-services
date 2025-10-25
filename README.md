# IAM Service (Identity and Access Management)

Service quản lý xác thực và phân quyền người dùng cho hệ thống e-commerce với **Casbin RBAC** integration.

## 🆕 What's New: gRPC Gateway (REST API)

IAM Service hiện hỗ trợ **cả gRPC và REST API** thông qua gRPC Gateway!

- **gRPC Server**: `localhost:50051`
- **REST API**: `http://localhost:8080`

Quick test REST API:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

👉 **[Quick Start Guide](GATEWAY_QUICKSTART.md)**  
👉 **[Full REST API Documentation](docs/GRPC_GATEWAY.md)**

## Tính năng chính

### Authentication (Xác thực)
- ✅ Đăng ký người dùng mới
- ✅ Đăng nhập với username/password
- ✅ Refresh token
- ✅ Xác minh token
- ✅ Đăng xuất

### Authorization (Phân quyền)
- ✅ Gán vai trò cho người dùng
- ✅ Xóa vai trò khỏi người dùng
- ✅ Lấy danh sách vai trò của người dùng
- ✅ Kiểm tra quyền truy cập

### **🆕 Casbin RBAC Authorization**
- ✅ **Multi-domain authorization** (user, cms, api)
- ✅ **CMS role management** với phân quyền theo tabs
- ✅ **API resource tracking** và phân quyền chi tiết
- ✅ **Pattern matching** cho flexible permissions
- ✅ **Policy enforcement** với Casbin engine

### Role Management (Quản lý vai trò)
- ✅ Tạo vai trò mới
- ✅ Cập nhật vai trò
- ✅ Xóa vai trò
- ✅ Lấy thông tin vai trò
- ✅ Liệt kê tất cả vai trò

### Permission Management (Quản lý quyền)
- ✅ Tạo quyền mới
- ✅ Xóa quyền
- ✅ Liệt kê tất cả quyền

## Kiến trúc

Service được xây dựng theo **Layered Architecture** với các layer:

```
┌─────────────────────────────────────┐
│      Presentation Layer (gRPC)      │
│         (Handler Layer)              │
├─────────────────────────────────────┤
│      Business Logic Layer            │
│         (Service Layer)              │
├─────────────────────────────────────┤
│      Data Access Layer               │
│    (Repository & DAO Pattern)        │
├─────────────────────────────────────┤
│         Domain Layer                 │
│       (Entities/Models)              │
└─────────────────────────────────────┘
```

## Cấu trúc thư mục

```
iam-services/
├── cmd/
│   └── server/              # Entry point của application
│       └── main.go
├── internal/
│   ├── config/             # Configuration management
│   ├── dao/                # Data Access Objects
│   ├── database/           # Database connection
│   ├── domain/             # Domain entities
│   ├── handler/            # gRPC handlers
│   ├── repository/         # Repository pattern
│   └── service/            # Business logic
├── pkg/
│   ├── jwt/                # JWT token management
│   ├── password/           # Password hashing
│   └── proto/              # Protocol Buffer definitions
├── migrations/             # Database migrations
├── docs/                   # Documentation
├── .env.example           # Environment variables template
├── go.mod                 # Go module file
└── README.md
```

## Công nghệ sử dụng

- **Language**: Go 1.21+
- **RPC Framework**: gRPC
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Logging**: Uber Zap

## Dependencies

- `google.golang.org/grpc` - gRPC framework
- `google.golang.org/protobuf` - Protocol Buffers
- `github.com/golang-jwt/jwt/v5` - JWT implementation
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/google/uuid` - UUID generation
- `golang.org/x/crypto` - Password hashing
- `go.uber.org/zap` - Structured logging

## Hướng dẫn cài đặt

### Yêu cầu

- Go 1.21 hoặc cao hơn
- PostgreSQL 12 hoặc cao hơn
- Protocol Buffer compiler (protoc)

### Các bước cài đặt

1. **Clone repository và di chuyển vào thư mục**

```bash
cd ecommerce/back_end/iam-services
```

2. **Cài đặt dependencies**

```bash
go mod download
```

3. **Cấu hình môi trường**

Copy file `.env.example` thành `.env` và cập nhật các giá trị:

```bash
cp .env.example .env
```

Chỉnh sửa file `.env`:

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=50051

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=iam_db
DB_SSL_MODE=disable

JWT_SECRET=your-secret-key-here
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_EXPIRATION_HOURS=168
```

4. **Tạo database**

```sql
CREATE DATABASE iam_db;
```

5. **Chạy migrations**

```bash
psql -U postgres -d iam_db -f migrations/001_init_schema.sql
psql -U postgres -d iam_db -f migrations/002_seed_data.sql
```

6. **Generate Protocol Buffer code**

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       pkg/proto/iam.proto
```

7. **Build và chạy service**

```bash
go run cmd/server/main.go
```

Hoặc build binary:

```bash
go build -o bin/iam-service cmd/server/main.go
./bin/iam-service
```

## Testing với grpcurl

### Đăng ký người dùng

```bash
grpcurl -plaintext -d '{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "full_name": "Test User"
}' localhost:50051 iam.IAMService/Register
```

### Đăng nhập

```bash
grpcurl -plaintext -d '{
  "username": "testuser",
  "password": "password123"
}' localhost:50051 iam.IAMService/Login
```

### Xác minh token

```bash
grpcurl -plaintext -d '{
  "token": "your-access-token-here"
}' localhost:50051 iam.IAMService/VerifyToken
```

## Tài liệu

- [Kiến trúc chi tiết](docs/ARCHITECTURE.md)
- [Hướng dẫn cài đặt](docs/SETUP.md)
- [API Documentation](docs/API.md)
- [Database Schema](docs/DATABASE.md)
- **[🆕 Casbin RBAC Guide](docs/CASBIN.md)** - Chi tiết về Casbin authorization
- **[🆕 Casbin Quick Start](README_CASBIN.md)** - Quick start guide

## License

Copyright © 2024

