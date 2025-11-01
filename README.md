# IAM Service - Identity and Access Management

**Version:** 1.0.0  
**Go Version:** 1.24  
**Web Framework:** Gin v1.9.1  
**Author:** E-commerce Platform Team

---

## 📖 Mục lục

1. [Tổng quan](#-tổng-quan)
2. [Tính năng](#-tính-năng)
3. [Kiến trúc](#-kiến-trúc)
4. [Cấu trúc Project](#-cấu-trúc-project)
5. [Công nghệ](#️-công-nghệ)
6. [Cài đặt](#-cài-đặt)
7. [Cấu hình](#️-cấu-hình)
8. [Chạy Service](#-chạy-service)
9. [API Documentation](#-api-documentation)
10. [Database Schema](#️-database-schema)
11. [Casbin Authorization](#-casbin-authorization-rbac)
12. [gRPC Gateway (REST API)](#-grpc-gateway-rest-api)
13. [Testing](#-testing)
14. [CI/CD Pipeline](#-cicd-pipeline)
15. [Deployment](#-deployment)
16. [Troubleshooting](#-troubleshooting)
17. [Best Practices](#-best-practices)

---

## 🎯 Tổng quan

IAM Service là một hệ thống quản lý danh tính và phân quyền người dùng toàn diện cho nền tảng e-commerce, được xây dựng theo **Clean Architecture** với **Casbin RBAC** và hỗ trợ cả **gRPC** và **REST API**.

### Điểm nổi bật

- ✅ **Clean Architecture**: Tách biệt rõ ràng giữa các layers
- ✅ **Gin Web Framework**: High-performance HTTP server với rich middleware
- ✅ **Multi-Domain RBAC**: User, CMS, API domains với Casbin
- ✅ **Dual Protocol**: gRPC (port 50051) + Gin HTTP (port 8080)
- ✅ **JWT Authentication**: Access token + Refresh token
- ✅ **PostgreSQL**: Với connection pooling và migrations
- ✅ **Logging**: Uber Zap với structured logging
- ✅ **Panic Recovery**: Multi-layered recovery system
- ✅ **Swagger UI**: Protected với Basic Authentication
- ✅ **Shared GoKits**: Reusable infrastructure library

---

## 🚀 Tính năng

### Authentication (Xác thực)
- ✅ **Register**: Đăng ký người dùng mới
- ✅ **Login**: Đăng nhập với username/password
- ✅ **Refresh Token**: Làm mới access token
- ✅ **Verify Token**: Xác minh tính hợp lệ của token
- ✅ **Logout**: Đăng xuất người dùng

### Authorization (Phân quyền)
- ✅ **Role Management**: CRUD operations cho roles
- ✅ **Permission Management**: CRUD operations cho permissions
- ✅ **User-Role Assignment**: Gán/xóa roles cho users
- ✅ **Permission Checking**: Kiểm tra quyền truy cập

### Casbin RBAC (Advanced Authorization)
- ✅ **Multi-Domain**: User, CMS, API domains
- ✅ **CMS Roles**: Phân quyền theo tabs (product, inventory, order, report, ...)
- ✅ **API Resources**: Track và phân quyền chi tiết cho API endpoints
- ✅ **Pattern Matching**: Flexible permissions với wildcards và regex
- ✅ **Policy Enforcement**: Real-time policy evaluation

### gRPC Gateway
- ✅ **REST API**: Tự động generate từ gRPC definitions
- ✅ **OpenAPI/Swagger**: Documentation tự động
- ✅ **CORS Support**: Hỗ trợ cross-origin requests
- ✅ **Dual Protocol**: Chạy song song gRPC và HTTP

---

## 🏗️ Kiến trúc

### Clean Architecture Layers

```
┌────────────────────────────────────────────────────────┐
│                  CMD (Main Entry)                      │
│                   cmd/server/main.go                   │
└────────────────────────────────────────────────────────┘
                            ↓
┌────────────────────────────────────────────────────────┐
│              APPLICATION CONTAINER                      │
│            (Dependency Injection)                       │
│            internal/container/container.go              │
└────────────────────────────────────────────────────────┘
                            ↓
┌────────────────────────────────────────────────────────┐
│                   ADAPTER LAYER                        │
│              (Presentation/Interface)                   │
│           • gRPC Handlers (grpc_handler.go)            │
│           • Gin HTTP Handlers (gin_handler.go)         │
│           • Gin Router (gin_router.go)                 │
│           • Converters (converter.go)                  │
│           • Middleware (recovery.go, gin_middleware.go)│
└────────────────────────────────────────────────────────┘
                            ↓
┌────────────────────────────────────────────────────────┐
│                APPLICATION LAYER                       │
│                  (Use Cases)                           │
│           • Auth Use Cases (register, login)           │
│           • Role Use Cases (CRUD)                      │
│           • Casbin Use Cases (check access)            │
│           • DTOs (Data Transfer Objects)               │
└────────────────────────────────────────────────────────┘
                            ↓
┌────────────────────────────────────────────────────────┐
│                   DOMAIN LAYER                         │
│              (Business Logic Core)                     │
│           • Domain Models (entities)                   │
│           • Repository Interfaces (ports)              │
│           • Domain Services (interfaces)               │
└────────────────────────────────────────────────────────┘
                            ↓
┌────────────────────────────────────────────────────────┐
│              INFRASTRUCTURE LAYER                      │
│               (Implementation Details)                  │
│           • Persistence (PostgreSQL, DAOs)             │
│           • Security (JWT, Password)                   │
│           • Authorization (Casbin)                     │
│           • Config (Environment variables)             │
└────────────────────────────────────────────────────────┘
```

### Request Flow Example: Login

#### gRPC Flow
```
1. gRPC Request → GRPCHandler.Login()
                        ↓
2. Validate Input
                        ↓
3. AuthService.Login()
   ├── UserRepository.GetByUsername()
   ├── PasswordService.CheckPassword()
   ├── AuthorizationRepository.GetUserRoles()
   └── TokenService.GenerateAccessToken()
                        ↓
4. Convert to Protocol Buffer
                        ↓
5. Response → Client
```

#### Gin HTTP Flow (NEW)
```
1. HTTP Request → Gin Router → GinHandler.Login()
                        ↓
2. Middleware Stack (Logger, CORS, Recovery)
                        ↓
3. Request Binding & Validation
                        ↓
4. AuthService.Login()
   ├── UserRepository.GetByUsername()
   ├── PasswordService.CheckPassword()
   ├── AuthorizationRepository.GetUserRoles()
   └── TokenService.GenerateAccessToken()
                        ↓
5. Convert to JSON Response
                        ↓
6. Response → Client
```

---

## 📁 Cấu trúc Project

```
iam-services/
├── cmd/
│   └── server/
│       └── main.go                          # Entry point
│
├── internal/
│   ├── app/
│   │   └── app.go                           # Application lifecycle
│   │
│   ├── container/
│   │   └── container.go                     # Dependency injection
│   │
│   ├── application/                         # Use cases & DTOs
│   │   ├── dto/
│   │   │   ├── auth_dto.go
│   │   │   └── casbin_dto.go
│   │   └── usecase/
│   │       ├── auth/
│   │       │   ├── login_usecase.go
│   │       │   └── register_usecase.go
│   │       └── casbin/
│   │           └── check_api_access_usecase.go
│   │
│   ├── domain/                              # Business logic core
│   │   ├── model/                           # Rich domain models
│   │   │   ├── user.go
│   │   │   ├── role.go
│   │   │   ├── permission.go
│   │   │   ├── cms_role.go
│   │   │   └── api_resource.go
│   │   ├── repository/                      # Repository interfaces (ports)
│   │   │   ├── user_repository.go
│   │   │   ├── role_repository.go
│   │   │   ├── authorization_repository.go
│   │   │   ├── cms_repository.go
│   │   │   └── api_resource_repository.go
│   │   └── service/                         # Domain service interfaces
│   │       ├── password_service.go
│   │       ├── token_service.go
│   │       └── authorization_service.go
│   │
│   ├── infrastructure/                      # Implementation details
│   │   ├── persistence/                     # Repository implementations
│   │   │   ├── user_repository_impl.go
│   │   │   ├── role_repository_impl.go
│   │   │   ├── authorization_repository_impl.go
│   │   │   ├── cms_repository_impl.go
│   │   │   └── api_resource_repository_impl.go
│   │   ├── security/                        # Security implementations
│   │   │   ├── jwt_service_impl.go
│   │   │   └── password_service_impl.go
│   │   ├── authorization/                   # Casbin implementation
│   │   │   └── casbin_service_impl.go
│   │   └── config/
│   │       └── config_loader.go
│   │
│   ├── dao/                                 # Data Access Objects
│   │   ├── user_dao.go
│   │   ├── role_dao.go
│   │   ├── permission_dao.go
│   │   ├── user_role_dao.go
│   │   ├── role_permission_dao.go
│   │   ├── cms_role_dao.go
│   │   ├── user_cms_role_dao.go
│   │   └── api_resource_dao.go
│   │
│   ├── handler/                             # Request handlers
│   │   ├── grpc_handler.go                  # gRPC handlers
│   │   ├── gin_handler.go                   # Gin HTTP handlers (NEW)
│   │   ├── casbin_handler.go
│   │   └── converter.go
│   │
│   ├── router/                              # Gin routing (NEW)
│   │   └── gin_router.go
│   │
│   ├── middleware/
│   │   ├── recovery.go                      # gRPC panic recovery
│   │   └── gin_middleware.go                # Gin middleware (NEW)
│   │
│   ├── config/
│   │   └── config.go
│   │
│   └── database/
│       └── database.go                      # Database connection
│
├── pkg/
│   ├── proto/                               # Protocol Buffer definitions
│   │   ├── iam.proto
│   │   ├── iam_gateway.proto
│   │   ├── iam.pb.go                        # Generated
│   │   ├── iam_grpc.pb.go                   # Generated
│   │   └── iam.pb.gw.go                     # Generated (Gateway)
│   ├── jwt/
│   │   └── jwt_manager.go
│   ├── password/
│   │   └── password_manager.go
│   └── casbin/
│       └── enforcer.go
│
├── configs/
│   └── rbac_model.conf                      # Casbin model
│
├── migrations/                              # SQL migrations
│   ├── 001_init_schema.sql
│   ├── 002_seed_data.sql
│   ├── 003_casbin_tables.sql
│   └── 004_casbin_seed_data.sql
│
├── scripts/
│   ├── setup-proto.ps1                      # Proto generation script
│   ├── setup.sh
│   └── test-api.sh
│
├── .env.example
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── docker-compose.yml
└── README.md (this file)
```

---

## 🛠️ Công nghệ

### Core Technologies
- **Language**: Go 1.24
- **Web Framework**: Gin v1.9.1 (NEW)
- **RPC Framework**: gRPC
- **Database**: PostgreSQL 12+
- **Authorization**: Casbin v2
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Logging**: Uber Zap
- **Config**: godotenv

### Libraries & Dependencies
```
github.com/gin-gonic/gin                     # Web framework (NEW)
github.com/gin-contrib/cors                  # CORS middleware (NEW)
github.com/casbin/casbin/v2                  # RBAC engine
github.com/casbin/gorm-adapter/v3            # Casbin adapter
github.com/golang-jwt/jwt/v5                 # JWT
github.com/lib/pq                            # PostgreSQL driver
github.com/google/uuid                       # UUID generation
go.uber.org/zap                              # Logging
google.golang.org/grpc                       # gRPC
google.golang.org/protobuf                   # Protocol Buffers
gorm.io/gorm                                 # ORM
gorm.io/driver/postgres                      # PostgreSQL driver
github.com/grpc-ecosystem/grpc-gateway/v2    # Gateway (for Swagger spec)
github.com/tvttt/gokits                      # Shared utilities (local)
```

---

## 📦 Cài đặt

### Yêu cầu hệ thống

1. **Go** - Version 1.19 hoặc cao hơn
   ```bash
   go version  # Kiểm tra version
   ```

2. **PostgreSQL** - Version 12 hoặc cao hơn
   ```bash
   psql --version
   ```

3. **Protocol Buffer Compiler (protoc)**
   - Download: https://github.com/protocolbuffers/protobuf/releases
   ```bash
   protoc --version
   ```

4. **Go Protoc Plugins**
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
   go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
   go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
   ```

5. **grpcurl** (Optional - để test)
   ```bash
   go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
   ```

Đảm bảo `$GOPATH/bin` trong PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Các bước cài đặt

#### 1. Clone/Navigate đến project
```bash
cd ecommerce/back_end/iam-services
```

#### 2. Cài đặt Go dependencies
```bash
go mod download
go mod tidy
```

#### 3. Setup PostgreSQL Database

**Tạo database:**
```sql
-- Đăng nhập PostgreSQL
psql -U postgres

-- Tạo database
CREATE DATABASE iam_db;

-- Tạo user (optional)
CREATE USER iam_user WITH PASSWORD 'your_password';

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE iam_db TO iam_user;

-- Exit
\q
```

**Chạy migrations:**
```bash
# Schema migration
psql -U postgres -d iam_db -f migrations/001_init_schema.sql

# Seed data
psql -U postgres -d iam_db -f migrations/002_seed_data.sql

# Casbin tables
psql -U postgres -d iam_db -f migrations/003_casbin_tables.sql

# Casbin seed data
psql -U postgres -d iam_db -f migrations/004_casbin_seed_data.sql
```

**Verify database:**
```bash
psql -U postgres -d iam_db

\dt  # List tables
SELECT * FROM roles;
SELECT * FROM permissions;
SELECT * FROM casbin_rule LIMIT 10;
\q
```

#### 4. Cấu hình Environment Variables

Copy `.env.example` sang `.env`:
```bash
cp .env.example .env
```

Chỉnh sửa `.env`:
```env
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=50051

# HTTP Gateway Configuration
HTTP_HOST=0.0.0.0
HTTP_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=iam_db
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=change-this-to-a-strong-random-secret-key-minimum-32-chars
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_EXPIRATION_HOURS=168

# Casbin Configuration
CASBIN_MODEL_PATH=./configs/rbac_model.conf

# Log Configuration
LOG_LEVEL=info
LOG_ENCODING=json
```

**⚠️ Lưu ý bảo mật:**
- Đổi `JWT_SECRET` thành chuỗi ngẫu nhiên mạnh (ít nhất 32 ký tự)
- Không commit file `.env` lên Git
- Trong production, sử dụng secret management tools (Vault, AWS Secrets Manager)

#### 5. Generate Protocol Buffer Code

**Windows PowerShell:**
```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\setup-proto.ps1
```

**Linux/MacOS:**
```bash
chmod +x scripts/setup.sh
./scripts/setup.sh
```

Hoặc manual:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
       --openapiv2_out=. \
       pkg/proto/iam.proto pkg/proto/iam_gateway.proto
```

---

## ⚙️ Cấu hình

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVER_HOST` | gRPC server host | `0.0.0.0` | Yes |
| `SERVER_PORT` | gRPC server port | `50051` | Yes |
| `HTTP_HOST` | HTTP gateway host | `0.0.0.0` | Yes |
| `HTTP_PORT` | HTTP gateway port | `8080` | Yes |
| `DB_HOST` | PostgreSQL host | `localhost` | Yes |
| `DB_PORT` | PostgreSQL port | `5432` | Yes |
| `DB_USER` | Database user | `postgres` | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | `iam_db` | Yes |
| `DB_SSL_MODE` | SSL mode | `disable` | Yes |
| `JWT_SECRET` | JWT secret key | - | Yes |
| `JWT_EXPIRATION_HOURS` | Access token expiration | `24` | Yes |
| `JWT_REFRESH_EXPIRATION_HOURS` | Refresh token expiration | `168` | Yes |
| `CASBIN_MODEL_PATH` | Casbin model file path | `./configs/rbac_model.conf` | Yes |
| `LOG_LEVEL` | Logging level | `info` | No |
| `LOG_ENCODING` | Log format (json/console) | `json` | No |

### Database Connection Pool

Configured in `internal/database/database.go`:
```go
db.SetMaxOpenConns(25)          // Maximum open connections
db.SetMaxIdleConns(5)           // Maximum idle connections
db.SetConnMaxLifetime(5 * time.Minute)
```

---

## 🎬 Chạy Service

### Option 1: Run trực tiếp
```bash
go run cmd/server/main.go
```

### Option 2: Build binary
```bash
# Build
go build -o bin/iam-service cmd/server/main.go

# Run
./bin/iam-service
```

### Option 3: Docker
```bash
# Build image
docker build -t iam-service .

# Run container
docker run -p 50051:50051 -p 8080:8080 --env-file .env iam-service
```

### Option 4: Docker Compose
```bash
docker-compose up -d
```

### Success Output

Nếu thành công, bạn sẽ thấy:
```
{"level":"info","ts":...,"msg":"Starting IAM Service..."}
{"level":"info","ts":...,"msg":"Database connected successfully"}
{"level":"info","ts":...,"msg":"Casbin enforcer initialized"}
{"level":"info","ts":...,"msg":"gRPC server starting","address":"0.0.0.0:50051"}
{"level":"info","ts":...,"msg":"HTTP gateway starting","address":"0.0.0.0:8080"}
```

Service chạy trên:
- **gRPC**: `0.0.0.0:50051`
- **REST API**: `http://localhost:8080`

---

## 📚 API Documentation

### gRPC API

**Server Address**: `localhost:50051`

#### Authentication APIs

##### 1. Register
```bash
grpcurl -plaintext -d '{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe"
}' localhost:50051 iam.IAMService/Register
```

##### 2. Login
```bash
grpcurl -plaintext -d '{
  "username": "johndoe",
  "password": "SecurePass123!"
}' localhost:50051 iam.IAMService/Login
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 86400,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "is_active": true
  }
}
```

##### 3. Refresh Token
```bash
grpcurl -plaintext -d '{
  "refresh_token": "your-refresh-token"
}' localhost:50051 iam.IAMService/RefreshToken
```

##### 4. Verify Token
```bash
grpcurl -plaintext -d '{
  "token": "your-access-token"
}' localhost:50051 iam.IAMService/VerifyToken
```

##### 5. Logout
```bash
grpcurl -plaintext -d '{
  "user_id": "user-id",
  "token": "access-token"
}' localhost:50051 iam.IAMService/Logout
```

#### Authorization APIs

##### 6. Assign Role
```bash
grpcurl -plaintext -d '{
  "user_id": "user-123",
  "role_id": "role-admin"
}' localhost:50051 iam.IAMService/AssignRole
```

##### 7. Remove Role
```bash
grpcurl -plaintext -d '{
  "user_id": "user-123",
  "role_id": "role-admin"
}' localhost:50051 iam.IAMService/RemoveRole
```

##### 8. Get User Roles
```bash
grpcurl -plaintext -d '{
  "user_id": "user-123"
}' localhost:50051 iam.IAMService/GetUserRoles
```

##### 9. Check Permission
```bash
grpcurl -plaintext -d '{
  "user_id": "user-123",
  "resource": "users",
  "action": "create"
}' localhost:50051 iam.IAMService/CheckPermission
```

#### Role Management APIs

##### 10. Create Role
```bash
grpcurl -plaintext -d '{
  "name": "editor",
  "description": "Content editor role",
  "permission_ids": ["perm-001", "perm-002"]
}' localhost:50051 iam.IAMService/CreateRole
```

##### 11-14. Update/Delete/Get/List Roles
Similar format...

#### Permission Management APIs

##### 15. Create Permission
```bash
grpcurl -plaintext -d '{
  "name": "Product Create",
  "resource": "product",
  "action": "create",
  "description": "Permission to create products"
}' localhost:50051 iam.IAMService/CreatePermission
```

#### Casbin APIs (See Casbin section for details)

---

### REST API

**Base URL**: `http://localhost:8080/api/v1`

#### Authentication Endpoints

##### POST /api/v1/auth/register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'
```

##### POST /api/v1/auth/login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecurePass123!"
  }'
```

##### POST /api/v1/auth/refresh
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "your-refresh-token"}'
```

##### POST /api/v1/auth/verify
```bash
curl -X POST http://localhost:8080/api/v1/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"token": "your-access-token"}'
```

##### POST /api/v1/auth/logout
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-id",
    "token": "access-token"
  }'
```

#### Authorization Endpoints

##### POST /api/v1/authorization/assign-role
##### DELETE /api/v1/authorization/users/{user_id}/roles/{role_id}
##### GET /api/v1/authorization/users/{user_id}/roles
##### POST /api/v1/authorization/check-permission

#### Role Management Endpoints

##### POST /api/v1/roles
##### PUT /api/v1/roles/{role_id}
##### GET /api/v1/roles/{role_id}
##### GET /api/v1/roles?page=1&page_size=10
##### DELETE /api/v1/roles/{role_id}

#### Permission Management Endpoints

##### POST /api/v1/permissions
##### GET /api/v1/permissions?page=1&page_size=10
##### DELETE /api/v1/permissions/{permission_id}

#### Casbin Endpoints (See Casbin section)

---

### Error Codes

| Code | Description |
|------|-------------|
| `OK` (200) | Success |
| `INVALID_ARGUMENT` (400) | Invalid input parameters |
| `UNAUTHENTICATED` (401) | Invalid credentials or token |
| `PERMISSION_DENIED` (403) | User doesn't have required permission |
| `NOT_FOUND` (404) | Resource not found |
| `ALREADY_EXISTS` (409) | Resource already exists |
| `INTERNAL` (500) | Internal server error |

---

## 🗄️ Database Schema

### Entity Relationship Diagram

```
┌─────────────────────────────────┐
│           users                  │
├─────────────────────────────────┤
│ id              VARCHAR(36) PK   │
│ username        VARCHAR(100) UK  │
│ email           VARCHAR(255) UK  │
│ password_hash   VARCHAR(255)     │
│ full_name       VARCHAR(255)     │
│ is_active       BOOLEAN          │
│ created_at      TIMESTAMP        │
│ updated_at      TIMESTAMP        │
└────────────┬────────────────────┘
             │ 1:N
┌────────────▼────────────────────┐
│       user_roles                 │
├─────────────────────────────────┤
│ user_id         VARCHAR(36) PK,FK│
│ role_id         VARCHAR(36) PK,FK│
│ created_at      TIMESTAMP        │
└────────────┬────────────────────┘
             │ N:1
┌────────────▼────────────────────┐
│           roles                  │
├─────────────────────────────────┤
│ id              VARCHAR(36) PK   │
│ name            VARCHAR(100) UK  │
│ description     TEXT             │
│ created_at      TIMESTAMP        │
│ updated_at      TIMESTAMP        │
└────────────┬────────────────────┘
             │ 1:N
┌────────────▼────────────────────┐
│    role_permissions              │
├─────────────────────────────────┤
│ role_id         VARCHAR(36) PK,FK│
│ permission_id   VARCHAR(36) PK,FK│
│ created_at      TIMESTAMP        │
└────────────┬────────────────────┘
             │ N:1
┌────────────▼────────────────────┐
│       permissions                │
├─────────────────────────────────┤
│ id              VARCHAR(36) PK   │
│ name            VARCHAR(100) UK  │
│ resource        VARCHAR(100)     │
│ action          VARCHAR(50)      │
│ description     TEXT             │
│ created_at      TIMESTAMP        │
│ updated_at      TIMESTAMP        │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│        casbin_rule               │
├─────────────────────────────────┤
│ id              SERIAL PK        │
│ ptype           VARCHAR(100)     │
│ v0              VARCHAR(100)     │
│ v1              VARCHAR(100)     │
│ v2              VARCHAR(100)     │
│ v3              VARCHAR(100)     │
│ v4              VARCHAR(100)     │
│ v5              VARCHAR(100)     │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│         cms_roles                │
├─────────────────────────────────┤
│ id              VARCHAR(36) PK   │
│ name            VARCHAR(100) UK  │
│ description     TEXT             │
│ tabs            TEXT[]           │
│ created_at      TIMESTAMP        │
│ updated_at      TIMESTAMP        │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│      api_resources               │
├─────────────────────────────────┤
│ id              VARCHAR(36) PK   │
│ path            VARCHAR(500)     │
│ method          VARCHAR(20)      │
│ service         VARCHAR(100)     │
│ description     TEXT             │
│ created_at      TIMESTAMP        │
│ updated_at      TIMESTAMP        │
│ UNIQUE(path, method)             │
└─────────────────────────────────┘
```

### Tables

#### 1. users
Lưu trữ thông tin người dùng.

**Columns**:
- `id`: UUID - Primary key
- `username`: Tên đăng nhập (unique)
- `email`: Email (unique)
- `password_hash`: Password đã hash bằng bcrypt (cost 10)
- `full_name`: Tên đầy đủ
- `is_active`: Trạng thái active
- `created_at`, `updated_at`: Timestamps

**Indexes**:
- `idx_users_username`
- `idx_users_email`
- `idx_users_is_active`

#### 2. roles
Lưu trữ vai trò.

**Default Roles**:
- `admin`: Full access
- `user`: Basic access
- `moderator`: Intermediate access

#### 3. permissions
Lưu trữ quyền hạn.

**Format**: `resource:action`
- Examples: `user:read`, `user:create`, `product:update`

#### 4. user_roles
Junction table (Many-to-Many): Users ↔ Roles

#### 5. role_permissions
Junction table (Many-to-Many): Roles ↔ Permissions

#### 6. casbin_rule
Casbin policy storage (RBAC policies)

#### 7. cms_roles
CMS-specific roles với phân quyền theo tabs

#### 8. api_resources
Tracking API endpoints cho phân quyền chi tiết

#### 9. user_cms_roles
Junction table: Users ↔ CMS Roles

### Common Queries

#### Get all roles của user
```sql
SELECT r.*
FROM roles r
INNER JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = 'user-id';
```

#### Get all permissions của user
```sql
SELECT DISTINCT p.*
FROM permissions p
INNER JOIN role_permissions rp ON p.id = rp.permission_id
INNER JOIN user_roles ur ON rp.role_id = ur.role_id
WHERE ur.user_id = 'user-id';
```

#### Check if user has permission
```sql
SELECT EXISTS(
    SELECT 1
    FROM permissions p
    INNER JOIN role_permissions rp ON p.id = rp.permission_id
    INNER JOIN user_roles ur ON rp.role_id = ur.role_id
    WHERE ur.user_id = 'user-id'
      AND p.resource = 'users'
      AND p.action = 'create'
) AS has_permission;
```

---

## 🔐 Casbin Authorization (RBAC)

### 🎯 Tổng quan

IAM Service sử dụng **2 hệ thống Casbin RBAC độc lập** với architecture tách biệt:

| System | Purpose | Model | Database | Structure |
|--------|---------|-------|----------|-----------|
| **User/App Authorization** | End users on web/app | `rbac_user_model.conf` | `casbin_rule_user` | Domain-based (user, api) |
| **CMS Authorization** | Admin/staff on CMS | `rbac_cms_model.conf` | `casbin_rule_cms` | Tab-based (product, inventory...) |

📖 **[Chi tiết Architecture](AUTHORIZATION_ARCHITECTURE.md)** - Xem document đầy đủ về separated authorization architecture

---

### 1️⃣ User/App Authorization (`roles` table)

**Purpose**: Phân quyền cho end users truy cập web/app

**Model**: `configs/rbac_user_model.conf`

```ini
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
```

**Domains**:
- `user`: User-specific resources
- `api`: API access control

**Default Roles**:
- `user`: Regular user (browse products, create orders)
- `premium_user`: Premium features
- `api_admin`: Full API access

**Example Policies**:
```
# Regular user policies
p, user, user, /api/v1/products, GET
p, user, user, /api/v1/products/*, GET
p, user, user, /api/v1/orders, (GET|POST)
p, user, user, /api/v1/users/self, (GET|PUT)

# Premium user policies
p, premium_user, api, /api/v1/products/**, (GET|POST)
p, premium_user, api, /api/v1/orders/**, (GET|POST|PUT)

# Admin policies
p, api_admin, api, /api/v1/**, (GET|POST|PUT|DELETE)

# Role assignments
g, user-123, user, user
g, user-456, premium_user, api
```

---

### 2️⃣ CMS Authorization (`cms_roles` table)

**Purpose**: Phân quyền cho admin/staff trên CMS với tab-based access control

**Model**: `configs/rbac_cms_model.conf`

```ini
[request_definition]
r = sub, tab, api, act

[policy_definition]
p = sub, tab, api, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.tab == p.tab && keyMatch2(r.api, p.api) && regexMatch(r.act, p.act)
```

**CMS Tabs**:

| Tab | Description | Example APIs |
|-----|-------------|--------------|
| `product` | Product management | `/api/v1/products`, `/api/v1/categories` |
| `inventory` | Inventory management | `/api/v1/inventory/*`, `/api/v1/warehouses` |
| `order` | Order management | `/api/v1/orders/*`, `/api/v1/orders/*/ship` |
| `user` | User management | `/api/v1/users/*`, `/api/v1/users/*/roles` |
| `report` | Reports & analytics | `/api/v1/reports/sales`, `/api/v1/reports/revenue` |
| `setting` | System settings | `/api/v1/settings/*`, `/api/v1/roles` |

**Tab-API Mapping** (`cms_tab_apis` table):
```
tab_name='product'    | api_path='/api/v1/products'     | method='GET'
tab_name='product'    | api_path='/api/v1/products'     | method='POST'
tab_name='inventory'  | api_path='/api/v1/inventory/*'  | method='GET'
tab_name='inventory'  | api_path='/api/v1/products'     | method='GET'  ← Shared API
```

**Default CMS Roles**:
- `cms_admin`: Full access to all tabs
- `cms_product_manager`: Product + Inventory tabs only
- `cms_order_manager`: Order tab only
- `cms_viewer`: Read-only access to selected tabs

**Example Policies**:
```
# CMS Admin - full access to all tabs
p, cms_admin, product, /api/v1/products/*, (GET|POST|PUT|DELETE)
p, cms_admin, inventory, /api/v1/inventory/*, (GET|POST|PUT|DELETE)
p, cms_admin, order, /api/v1/orders/*, (GET|POST|PUT|DELETE)

# Product Manager - product & inventory tabs only
p, cms_product_manager, product, /api/v1/products/*, (GET|POST|PUT|DELETE)
p, cms_product_manager, inventory, /api/v1/inventory/*, (GET|POST|PUT)

# Viewer - read only
p, cms_viewer, product, /api/v1/products/*, GET
p, cms_viewer, inventory, /api/v1/inventory/*, GET

# Role assignments (no domain, only role name)
g, user-789, cms_admin
g, user-456, cms_product_manager
```

---

### 📊 Database Tables

#### User/App Tables
- `roles`: User/app roles (user, premium_user, api_admin)
- `casbin_rule_user`: Policies for user/app authorization
- `user_roles`: User-to-role assignments

#### CMS Tables
- `cms_roles`: CMS roles with tabs array
- `cms_tab_apis`: Maps tabs to APIs (many-to-many)
- `casbin_rule_cms`: Policies for CMS authorization
- `user_cms_roles`: User-to-CMS-role assignments

---

### 🔄 Authorization Flows

#### User/App Flow
```
Request → Extract (user_id, domain, api_path, method)
    ↓
Check User Enforcer (rbac_user_model.conf)
    ↓
Get user roles from casbin_rule_user
    ↓
Match policies (keyMatch2 + regexMatch)
    ↓
Allow/Deny
```

#### CMS Flow
```
Request → Extract (user_id, tab, api_path, method)
    ↓
Verify tab-API mapping in cms_tab_apis
    ↓
Check CMS Enforcer (rbac_cms_model.conf)
    ↓
Get user CMS roles from casbin_rule_cms
    ↓
Match policies for that tab
    ↓
Allow/Deny
```

### Casbin APIs

#### Check API Access
```bash
# gRPC
grpcurl -plaintext -d '{
  "user_id": "user-123",
  "api_path": "/api/v1/products",
  "method": "POST"
}' localhost:50051 iam.IAMService/CheckAPIAccess

# REST
curl -X POST http://localhost:8080/api/v1/casbin/check-api-access \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "api_path": "/api/v1/products",
    "method": "POST"
  }'
```

#### Check CMS Access
```bash
# gRPC
grpcurl -plaintext -d '{
  "user_id": "user-456",
  "cms_tab": "product",
  "action": "POST"
}' localhost:50051 iam.IAMService/CheckCMSAccess

# REST
curl -X POST http://localhost:8080/api/v1/casbin/check-cms-access \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-456",
    "cms_tab": "product",
    "action": "POST"
  }'
```

#### Enforce Policy (General)
```bash
grpcurl -plaintext -d '{
  "user_id": "user-789",
  "domain": "api",
  "resource": "/api/v1/users/123",
  "action": "DELETE"
}' localhost:50051 iam.IAMService/EnforcePolicy
```

#### Create CMS Role
```bash
grpcurl -plaintext -d '{
  "name": "cms_content_editor",
  "description": "Content editor role",
  "tabs": ["product", "inventory"]
}' localhost:50051 iam.IAMService/CreateCMSRole
```

#### Assign CMS Role
```bash
grpcurl -plaintext -d '{
  "user_id": "user-123",
  "cms_role_id": "cms-role-001"
}' localhost:50051 iam.IAMService/AssignCMSRole
```

#### Get User CMS Tabs
```bash
grpcurl -plaintext -d '{
  "user_id": "user-123"
}' localhost:50051 iam.IAMService/GetUserCMSTabs
```

#### Create API Resource
```bash
grpcurl -plaintext -d '{
  "path": "/api/v1/products",
  "method": "POST",
  "service": "product-service",
  "description": "Create new product"
}' localhost:50051 iam.IAMService/CreateAPIResource
```

#### List API Resources
```bash
grpcurl -plaintext -d '{
  "service": "product-service",
  "page": 1,
  "page_size": 10
}' localhost:50051 iam.IAMService/ListAPIResources
```

### Pattern Matching

#### KeyMatch2 (cho resource paths)
- `/api/v1/products` matches `/api/v1/products`
- `/api/v1/products/*` matches `/api/v1/products/123`
- `/api/v1/**` matches tất cả sub-paths

#### RegexMatch (cho actions)
- `GET` matches exact "GET"
- `(GET|POST)` matches "GET" OR "POST"
- `(GET|POST|PUT|DELETE)` matches tất cả CRUD operations

### Workflow Examples

#### Setup CMS Admin
```go
// 1. Create CMS role
CreateCMSRole("cms_admin", "Full CMS access", 
    []string{"product", "inventory", "order", "user", "report", "setting"})

// 2. Add policies
AddPolicy("cms_admin", "cms", "/cms/product/*", "(GET|POST|PUT|DELETE)")
AddPolicy("cms_admin", "cms", "/cms/inventory/*", "(GET|POST|PUT|DELETE)")

// 3. Assign to user
AssignCMSRole("user-123", "cms-role-admin-001")

// 4. Check access
CheckCMSAccess("user-123", "product", "POST")  // true
```

#### Multi-Domain User
```go
// User có roles ở nhiều domains
AssignUserRole("user-789", "user", "user")        // End user
AssignCMSRole("user-789", "cms-product-manager")  // CMS staff
AssignUserRole("user-789", "moderator", "api")    // API access

// Check domains khác nhau
CheckAPIAccess("user-789", "/api/v1/products", "GET")   // API domain
CheckCMSAccess("user-789", "product", "POST")            // CMS domain
```

---

### 🔄 Migration to Separated Architecture

Nếu bạn đang migrate từ unified `casbin_rule` table sang separated architecture:

#### Step 1: Run Migrations
```bash
# Create new tables and migrate data
psql -U postgres -d iam_db -f migrations/005_separate_user_cms_authorization.sql

# Seed initial data
psql -U postgres -d iam_db -f migrations/006_seed_separated_authorization.sql
```

#### Step 2: Verify Migration
```sql
-- Check User/App policies
SELECT COUNT(*) FROM casbin_rule_user;
SELECT * FROM casbin_rule_user WHERE ptype = 'p' LIMIT 5;

-- Check CMS policies
SELECT COUNT(*) FROM casbin_rule_cms;
SELECT * FROM casbin_rule_cms WHERE ptype = 'p' LIMIT 5;

-- Check tab-API mappings
SELECT COUNT(*) FROM cms_tab_apis;
SELECT * FROM cms_tab_apis WHERE tab_name = 'product';

-- Old data backed up
SELECT COUNT(*) FROM casbin_rule_old_backup;
```

#### Step 3: Update Application Code
- Initialize 2 Casbin enforcers (user + CMS)
- Update authorization checks to use appropriate enforcer
- Register tab-API mappings for new APIs

#### Step 4: Benefits
✅ **Clearer separation**: User vs CMS authorization  
✅ **Tab-based control**: Flexible CMS access management  
✅ **API sharing**: Same API can belong to multiple tabs  
✅ **Better security**: Granular control per domain  
✅ **Easier maintenance**: Independent policy management

📖 **Chi tiết**: Xem [AUTHORIZATION_ARCHITECTURE.md](AUTHORIZATION_ARCHITECTURE.md)

---

## 🌐 gRPC Gateway (REST API)

### Tổng quan

IAM Service hỗ trợ **cả gRPC và REST API** thông qua gRPC Gateway.

### Servers

- **gRPC**: `localhost:50051`
- **REST API**: `http://localhost:8080`

### Generate Gateway Code

**Windows:**
```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\setup-proto.ps1
```

**Linux/MacOS:**
```bash
./scripts/setup.sh
```

Script sẽ:
1. Install `protoc-gen-grpc-gateway` và `protoc-gen-openapiv2`
2. Download Google API proto files
3. Generate Gateway code (`.pb.gw.go`)
4. Generate OpenAPI/Swagger documentation (`.swagger.json`)

### CORS Support

Gateway hỗ trợ CORS với config:
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Accept, Authorization, Content-Type, X-CSRF-Token
```

### OpenAPI/Swagger

Sau khi generate, file OpenAPI spec:
```
pkg/proto/iam_gateway.swagger.json
```

Import vào Swagger UI hoặc Postman để xem interactive API documentation.

### Quick Test

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test123!",
    "full_name": "Test User"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123!"
  }'
```

### gRPC vs REST Comparison

| Feature | gRPC | REST (Gateway) | REST (Gin) |
|---------|------|----------------|------------|
| Port | 50051 | 8080 | 8080 |
| Protocol | HTTP/2 | HTTP/1.1 | HTTP/1.1 |
| Format | Protobuf | JSON | JSON |
| Performance | Cao nhất | Trung bình | Cao |
| Client | gRPC client | HTTP client | HTTP client |
| Browser Support | Giới hạn | Đầy đủ | Đầy đủ |
| Routing | Proto | Proto+Gateway | Gin Router |
| Middleware | gRPC | Limited | Đầy đủ |

---

## 🎨 Gin HTTP API (NEW)

### Tổng quan

IAM Service hiện sử dụng **Gin Web Framework** để handle REST API, thay thế cho gRPC Gateway. Gin mang lại:

- ✅ **Performance cao**: Nhanh nhất trong Go web frameworks
- ✅ **Flexible routing**: RESTful routing với params, query, body binding
- ✅ **Rich middleware**: Logger, CORS, Recovery, Authentication
- ✅ **Better error handling**: Consistent JSON error responses
- ✅ **Easy testing**: Standard HTTP testing với httptest

### Server Configuration

- **HTTP Server**: `http://localhost:8080`
- **gRPC Server**: `localhost:50051` (vẫn hoạt động song song)
- **Health Check**: `GET http://localhost:8080/health`
- **Swagger UI**: `http://localhost:8080/swagger/` (với Basic Auth)

### API Endpoints

#### Authentication
```
POST   /v1/auth/register     - Đăng ký user mới
POST   /v1/auth/login        - Đăng nhập
POST   /v1/auth/refresh      - Refresh access token
POST   /v1/auth/logout       - Đăng xuất
POST   /v1/auth/verify       - Verify token
```

#### Role Management
```
POST   /v1/roles             - Tạo role mới
GET    /v1/roles             - List roles (page, page_size query params)
GET    /v1/roles/:id         - Get role by ID
PUT    /v1/roles/:id         - Update role
DELETE /v1/roles/:id         - Delete role
POST   /v1/roles/assign      - Assign role to user
POST   /v1/roles/remove      - Remove role from user
```

#### User Roles
```
GET    /v1/users/:user_id/roles  - Get user's roles
```

#### Permission Management
```
POST   /v1/permissions           - Create permission
GET    /v1/permissions           - List permissions
DELETE /v1/permissions/:id       - Delete permission
POST   /v1/permissions/check     - Check permission
```

#### Access Control
```
POST   /v1/access/api        - Check API access
POST   /v1/access/cms        - Check CMS access
```

#### Policy Management
```
POST   /v1/policies/enforce  - Enforce authorization policy
```

#### CMS Management
```
POST   /v1/cms/roles         - Create CMS role
GET    /v1/cms/roles         - List CMS roles
POST   /v1/cms/roles/assign  - Assign CMS role
POST   /v1/cms/roles/remove  - Remove CMS role
GET    /v1/cms/users/:user_id/tabs  - Get user's CMS tabs
```

#### API Resources
```
POST   /v1/api/resources     - Create API resource
GET    /v1/api/resources     - List API resources
```

### Middleware Stack

1. **GinRecovery**: Panic recovery với logging
2. **GinLogger**: Structured logging cho mỗi request
3. **GinCORS**: CORS headers tự động
4. **GinAuth**: JWT authentication (optional per route)

### Request/Response Format

#### Success Response
```json
{
  "data": {
    // Response data object
  },
  "message": "Optional success message"
}
```

#### Error Response
```json
{
  "error": "Error message",
  "message": "Detailed explanation",
  "code": 400
}
```

### Example Usage

#### Register và Login
```bash
# Register
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'

# Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "SecurePass123!"
  }'
```

#### Role Management
```bash
# Create role
curl -X POST http://localhost:8080/v1/roles \
  -H "Content-Type: application/json" \
  -d '{
    "name": "product_manager",
    "description": "Can manage products",
    "permission_ids": ["perm-1", "perm-2"]
  }'

# List roles with pagination
curl "http://localhost:8080/v1/roles?page=1&page_size=10"

# Get specific role
curl "http://localhost:8080/v1/roles/role-123"

# Assign role to user
curl -X POST http://localhost:8080/v1/roles/assign \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "role_id": "role-456"
  }'
```

#### CMS Access Control
```bash
# Create CMS role
curl -X POST http://localhost:8080/v1/cms/roles \
  -H "Content-Type: application/json" \
  -d '{
    "name": "cms_editor",
    "description": "CMS Editor role",
    "tabs": ["product", "inventory", "order"]
  }'

# Check CMS access
curl -X POST http://localhost:8080/v1/access/cms \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "cms_tab": "product",
    "action": "POST"
  }'

# Get user's CMS tabs
curl "http://localhost:8080/v1/cms/users/user-123/tabs"
```

### Integration với Swagger UI

Swagger UI vẫn hoạt động và được serve bởi Gin router:

```bash
# Access Swagger UI (requires Basic Auth)
# Username: admin (default, configurable via SWAGGER_AUTH_USERNAME)
# Password: changeme (default, configurable via SWAGGER_AUTH_PASSWORD)
open http://localhost:8080/swagger/
```

### Configuration

Gin server được cấu hình qua environment variables:

```bash
# Server configuration
HTTP_PORT=8080                    # Gin HTTP server port
GRPC_PORT=50051                   # gRPC server port (vẫn hoạt động)

# Swagger configuration  
SWAGGER_ENABLED=true
SWAGGER_BASE_PATH=/swagger/
SWAGGER_SPEC_PATH=/swagger.json
SWAGGER_AUTH_USERNAME=admin
SWAGGER_AUTH_PASSWORD=changeme
SWAGGER_AUTH_REALM="IAM Service API Documentation"

# Log configuration
LOG_LEVEL=debug                   # Gin runs in debug mode if LOG_LEVEL=debug
```

### Benefits của Gin

1. **Performance**: Gin is one of the fastest Go web frameworks
2. **Clean Code**: Intuitive API, easy to maintain
3. **Rich Ecosystem**: Extensive middleware collection
4. **Better Testing**: Standard Go HTTP testing
5. **Flexible Routing**: Path parameters, query params, request binding
6. **Error Handling**: Consistent error responses
7. **Production Ready**: Used by many companies in production

### Dual Protocol Support

Service vẫn hỗ trợ cả gRPC và HTTP:

- **gRPC** (port 50051): For microservice-to-microservice communication
- **Gin HTTP** (port 8080): For external clients, web apps, mobile apps

Cả hai protocol đều sử dụng chung business logic (services, repositories).

---

## 🧪 Testing

### List Services
```bash
grpcurl -plaintext localhost:50051 list
```

### Describe Service
```bash
grpcurl -plaintext localhost:50051 describe iam.IAMService
```

### Test Registration Flow
```bash
# 1. Register
grpcurl -plaintext -d '{
  "username": "testuser",
  "email": "test@example.com",
  "password": "Password123!",
  "full_name": "Test User"
}' localhost:50051 iam.IAMService/Register

# 2. Login
grpcurl -plaintext -d '{
  "username": "testuser",
  "password": "Password123!"
}' localhost:50051 iam.IAMService/Login

# Save token from response
TOKEN="eyJhbGciOiJIUzI1NiIs..."

# 3. Verify Token
grpcurl -plaintext -d '{
  "token": "'$TOKEN'"
}' localhost:50051 iam.IAMService/VerifyToken

# 4. Get User Roles
grpcurl -plaintext -d '{
  "user_id": "user-id-from-login"
}' localhost:50051 iam.IAMService/GetUserRoles
```

### Test REST API Flow
```bash
# Complete flow với cURL
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }' | jq -r '.access_token')

echo "Token: $TOKEN"

# Verify
curl -X POST http://localhost:8080/api/v1/auth/verify \
  -H "Content-Type: application/json" \
  -d "{\"token\": \"$TOKEN\"}"

# List roles
curl "http://localhost:8080/api/v1/roles?page=1&page_size=10"
```

### Test Script
```bash
chmod +x scripts/test-api.sh
./scripts/test-api.sh
```

---

## 🔄 CI/CD Pipeline

IAM Service sử dụng **GitHub Actions** để tự động hóa quy trình build, test, và deployment.

### Pipeline Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                      CI/CD Pipeline Flow                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  Push/PR → Lint → Test → Build → Security Scan → Docker → Deploy   │
│              ↓       ↓       ↓          ↓            ↓        ↓     │
│           Passed  Passed  Passed    Passed      Pushed   Staging/   │
│                                                          Production   │
└─────────────────────────────────────────────────────────────────────┘
```

### Pipeline Stages

#### 1. 🔍 Lint & Code Quality
- **Công cụ**: golangci-lint
- **Kiểm tra**: 
  - Code formatting (`gofmt`)
  - Code quality và best practices
  - Cyclomatic complexity
  - Security issues (`gosec`)
  - Dead code detection

#### 2. 🧪 Automated Testing
- **Unit Tests**: Test các function và method riêng lẻ
- **Integration Tests**: Test tích hợp với database và Casbin
- **Coverage**: Báo cáo code coverage (upload lên Codecov)
- **Race Detection**: Phát hiện race conditions
- **Benchmark Tests**: Đánh giá performance (chỉ chạy khi PR)

#### 3. 🏗️ Build Application
- **Binary Build**: Compile Go binary cho Linux/AMD64
- **Optimization**: Strip symbols (`-ldflags="-w -s"`)
- **Artifact**: Upload binary để reuse trong các stage sau

#### 4. 🔒 Security Scan
- **Trivy**: Scan vulnerabilities trong dependencies và filesystem
- **Gosec**: Scan security issues trong Go code
- **SARIF Report**: Upload kết quả lên GitHub Security tab

#### 5. 🐳 Docker Build & Push
- **Trigger**: Chỉ chạy khi push lên `main` hoặc `develop`
- **Multi-stage Build**: Tối ưu image size
- **Registry**: Push lên Docker Hub
- **Tags**: 
  - `latest` (từ main)
  - `develop` (từ develop)
  - `<branch>-<sha>` (commit-specific)
- **Cache**: Sử dụng registry cache để tăng tốc build

#### 6. 🚀 Deployment
- **Staging**: Auto-deploy khi push lên `develop`
  - Environment: `staging`
  - URL: `https://iam-staging.example.com`
- **Production**: Auto-deploy khi push lên `main`
  - Environment: `production`
  - URL: `https://iam.example.com`
  - Create GitHub Release với version tag
- **Health Check**: Verify service sau khi deploy

#### 7. 📢 Notification
- **Slack**: Thông báo kết quả deployment
- **Status**: Success/Failure với chi tiết commit

---

### GitHub Workflows

#### Main CI/CD Pipeline
**File**: `.github/workflows/ci-cd.yml`

**Triggers**:
```yaml
on:
  push:
    branches: [main, develop, feature/**, hotfix/**]
  pull_request:
    branches: [main, develop]
  workflow_dispatch:
```

**Jobs**:
1. `lint` - Code quality checks
2. `test` - Unit tests với PostgreSQL service
3. `build` - Build binary artifact
4. `security` - Security scanning
5. `docker` - Build & push Docker image
6. `deploy-staging` - Deploy to staging (develop only)
7. `deploy-production` - Deploy to production (main only)
8. `notify` - Send notifications

#### Automated Testing
**File**: `.github/workflows/test.yml`

**Triggers**:
```yaml
on:
  push:
    branches: ['**']
  pull_request:
    branches: ['**']
  schedule:
    - cron: '0 2 * * *'  # Daily at 2 AM UTC
```

**Jobs**:
1. `unit-tests` - Unit tests với coverage report
2. `integration-tests` - Integration tests
3. `benchmark-tests` - Performance benchmarks (PR only)

---

### Setup CI/CD

#### 1. GitHub Secrets

Cấu hình các secrets trong GitHub Repository Settings:

```bash
# Docker Hub
DOCKER_USERNAME=your-docker-username
DOCKER_PASSWORD=your-docker-password

# Staging Server
STAGING_HOST=staging.example.com
STAGING_USER=deploy
STAGING_SSH_KEY=<private-ssh-key>

# Production Server
PRODUCTION_HOST=production.example.com
PRODUCTION_USER=deploy
PRODUCTION_SSH_KEY=<private-ssh-key>

# Notifications
SLACK_WEBHOOK=https://hooks.slack.com/services/...
```

#### 2. Server Setup

**Cài đặt trên Staging/Production servers**:
```bash
# Install Docker & Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Create deployment directory
mkdir -p /app/iam-service
cd /app/iam-service

# Setup docker-compose files
# - docker-compose.staging.yml
# - docker-compose.prod.yml
```

#### 3. Database Migrations

**Tự động chạy migrations**:
```yaml
# docker-compose.prod.yml
services:
  iam-service:
    image: your-username/iam-service:latest
    depends_on:
      - postgres
    environment:
      - DB_HOST=postgres
      - DB_NAME=iam_db
      # ... other env vars
    command: |
      sh -c "
        /app/iam-service migrate up &&
        /app/iam-service
      "
```

#### 4. Environment Variables

**Quản lý env vars**:
- Development: `.env` file (gitignored)
- Staging/Production: Docker secrets hoặc external config management (Vault, AWS Secrets Manager)

---

### Branch Strategy

```
main (production)
  ↑
  └─ develop (staging)
       ↑
       ├─ feature/user-authentication
       ├─ feature/casbin-integration
       └─ hotfix/critical-bug
```

**Workflow**:
1. **Feature Development**: 
   - Branch từ `develop`: `feature/feature-name`
   - Push → CI runs (lint, test, build, security)
   - Create PR to `develop` → Full pipeline runs
   - Merge → Auto-deploy to staging

2. **Release to Production**:
   - Create PR từ `develop` → `main`
   - Merge → Auto-deploy to production
   - GitHub Release được tạo tự động

3. **Hotfix**:
   - Branch từ `main`: `hotfix/issue-name`
   - Create PR to `main` → Deploy to production
   - Merge back to `develop`

---

### Monitoring CI/CD

#### GitHub Actions
- **Workflow Runs**: Repository → Actions tab
- **Build History**: Xem lịch sử và logs của mỗi run
- **Artifacts**: Download build artifacts (binaries, coverage reports)

#### Code Coverage
- **Codecov**: https://codecov.io/gh/your-org/iam-services
- **Trend**: Theo dõi coverage trend qua commits

#### Security
- **GitHub Security**: Repository → Security tab
- **Dependabot**: Auto PR để update vulnerable dependencies
- **Code Scanning**: Trivy và Gosec results

#### Docker Registry
- **Docker Hub**: https://hub.docker.com/r/your-username/iam-service
- **Image Tags**: Xem các version đã build
- **Image Size**: Monitor image size optimization

---

### Performance Optimization

#### Cache Strategy
```yaml
# Go modules cache
- uses: actions/setup-go@v4
  with:
    cache: true  # Cache go modules

# Docker layer cache
- uses: docker/build-push-action@v5
  with:
    cache-from: type=registry,ref=...
    cache-to: type=registry,ref=...
```

#### Parallel Jobs
- Lint và Security scan chạy song song
- Test chạy độc lập với Build
- Deploy jobs chỉ chạy khi cần

#### Smart Triggers
- PR: Chạy tất cả checks
- Push to feature branches: Chỉ lint + test
- Push to main/develop: Full pipeline + deploy

---

### Troubleshooting CI/CD

#### Build Fails

**Check logs**:
```bash
# Go to Actions tab → Click failed workflow → View logs
```

**Common issues**:
- Go module conflicts: `go mod tidy`
- Test failures: Check PostgreSQL service connectivity
- Docker build: Check Dockerfile syntax

#### Deployment Fails

**SSH connection issues**:
```bash
# Verify SSH key in GitHub Secrets
# Test SSH connection manually
ssh -i private-key deploy@production.example.com
```

**Docker pull fails**:
```bash
# Check Docker Hub credentials
# Verify image exists: docker pull username/iam-service:latest
```

**Health check fails**:
```bash
# Check service logs on server
docker-compose logs iam-service

# Verify environment variables
docker-compose config
```

---

### Local CI Testing

**Sử dụng Act để test workflows locally**:
```bash
# Install Act
brew install act  # macOS
# or download from https://github.com/nektos/act

# Run CI workflow locally
act -j lint
act -j test
act -j build

# Run entire workflow
act push
```

**Docker Compose cho local testing**:
```bash
# Build và test local
docker-compose -f docker-compose.local.yml up --build

# Verify health
curl http://localhost:8080/health
```

---

### Best Practices

#### 1. Commit Messages
```bash
# Format: <type>(<scope>): <subject>
git commit -m "feat(auth): add JWT refresh token support"
git commit -m "fix(casbin): resolve policy loading issue"
git commit -m "docs(readme): update CI/CD section"
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

#### 2. Pull Request
- **Title**: Clear và descriptive
- **Description**: Explain what, why, how
- **Labels**: `feature`, `bugfix`, `documentation`, `enhancement`
- **Review**: Request review từ team members

#### 3. Testing
- Viết tests cho mọi feature mới
- Maintain coverage > 80%
- Test cả success và error cases
- Add integration tests cho critical flows

#### 4. Security
- Never commit secrets hoặc passwords
- Use GitHub Secrets cho sensitive data
- Scan images trước khi deploy
- Review security alerts promptly

#### 5. Versioning
- Use semantic versioning: `v1.2.3`
- Tag releases: `git tag v1.2.3`
- Document breaking changes trong release notes

---

## 🚀 Deployment

### Production Checklist

- [ ] Đổi `JWT_SECRET` thành random strong key (ít nhất 64 chars)
- [ ] Set `DB_SSL_MODE=require` cho PostgreSQL
- [ ] Set `LOG_LEVEL=warn` hoặc `error`
- [ ] Enable database backup schedules
- [ ] Setup monitoring (Prometheus, Grafana)
- [ ] Setup centralized logging (ELK, Loki)
- [ ] Use connection pooling
- [ ] Setup reverse proxy (nginx, envoy)
- [ ] Enable TLS for gRPC
- [ ] Add rate limiting
- [ ] Add health checks
- [ ] Use secrets management (Vault, AWS Secrets Manager)

### Docker Production Build

```dockerfile
# Dockerfile
FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o iam-service cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/iam-service .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/.env .

EXPOSE 50051 8080
CMD ["./iam-service"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: iam-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: iam-service
  template:
    metadata:
      labels:
        app: iam-service
    spec:
      containers:
      - name: iam-service
        image: iam-service:latest
        ports:
        - containerPort: 50051
          name: grpc
        - containerPort: 8080
          name: http
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: iam-secrets
              key: db-password
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: iam-secrets
              key: jwt-secret
        livenessProbe:
          tcpSocket:
            port: 50051
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          tcpSocket:
            port: 50051
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Health Checks

Add health check endpoint (future enhancement):
```go
// internal/handler/health_handler.go
func (h *HealthHandler) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
    // Check database
    if err := h.db.Ping(); err != nil {
        return &pb.HealthCheckResponse{Status: "unhealthy"}, nil
    }
    
    return &pb.HealthCheckResponse{Status: "healthy"}, nil
}
```

---

## 🐛 Troubleshooting

### Issue 1: Connection refused
**Nguyên nhân**: Service chưa chạy hoặc port conflict

**Giải pháp**:
```bash
# Check if service is running
ps aux | grep iam-service

# Check ports
netstat -an | grep 50051
netstat -an | grep 8080

# Try different ports in .env
```

### Issue 2: Failed to connect to database
**Nguyên nhân**: PostgreSQL không chạy hoặc config sai

**Giải pháp**:
```bash
# Check PostgreSQL
pg_isready

# Check credentials
psql -U postgres -d iam_db

# Verify .env config
cat .env | grep DB_
```

### Issue 3: Table does not exist
**Nguyên nhân**: Migrations chưa chạy

**Giải pháp**:
```bash
# Run all migrations
psql -U postgres -d iam_db -f migrations/001_init_schema.sql
psql -U postgres -d iam_db -f migrations/002_seed_data.sql
psql -U postgres -d iam_db -f migrations/003_casbin_tables.sql
psql -U postgres -d iam_db -f migrations/004_casbin_seed_data.sql

# Verify tables
psql -U postgres -d iam_db -c "\dt"
```

### Issue 4: Invalid token
**Nguyên nhân**: JWT_SECRET không khớp hoặc token expired

**Giải pháp**:
```bash
# Đảm bảo JWT_SECRET không thay đổi
# Login lại để lấy token mới
# Check token expiration trong .env
```

### Issue 5: Protocol Buffer generation fails
**Nguyên nhân**: protoc hoặc plugins chưa cài

**Giải pháp**:
```bash
# Check protoc
protoc --version

# Reinstall plugins (for Go 1.19)
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

# Check PATH
echo $PATH | grep $(go env GOPATH)/bin

# Add to PATH if missing
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Issue 6: CORS errors in browser
**Nguyên nhân**: CORS middleware chưa config đúng

**Giải pháp**:
- Kiểm tra `corsMiddleware` trong `internal/app/app.go`
- Verify headers được set đúng
- Test với Postman trước (bypass CORS)

### Issue 7: Casbin authorization always denies
**Nguyên nhân**: User chưa được gán role hoặc policy chưa đúng

**Giải pháp**:
```sql
-- Check user role assignments
SELECT * FROM casbin_rule WHERE ptype = 'g' AND v0 = 'user-id';

-- Check policies for role
SELECT * FROM casbin_rule WHERE ptype = 'p' AND v0 = 'role-name';

-- Verify domain
SELECT * FROM casbin_rule WHERE v1 = 'cms' OR v1 = 'api' OR v1 = 'user';
```

### Issue 8: Go version mismatch errors
**Nguyên nhân**: Packages require newer Go version

**Giải pháp**:
```bash
# Check Go version
go version

# For Go 1.19, use compatible versions:
# google.golang.org/grpc v1.53.0
# github.com/redis/go-redis/v9 v9.0.5
# go.mongodb.org/mongo-driver v1.11.9

# Update go.mod và run
go mod tidy
```

---

## ✅ Best Practices

### Security
1. **Passwords**: Always hash với bcrypt (cost >= 10)
2. **JWT**: Use strong secret (minimum 32 chars), rotate regularly
3. **SQL Injection**: Always use parameterized queries
4. **TLS**: Enable trong production
5. **Rate Limiting**: Implement để prevent brute force
6. **Audit Logging**: Log critical operations (role assignments, permission changes)
7. **Secrets Management**: Use Vault hoặc AWS Secrets Manager trong production

### Architecture
1. **Separation of Concerns**: Keep layers independent
2. **Dependency Injection**: Use container pattern
3. **Interface Segregation**: Define focused interfaces
4. **Error Handling**: Return domain errors, map to gRPC/HTTP codes in handlers
5. **Logging**: Use structured logging with context
6. **Configuration**: Use environment variables, never hardcode

### Database
1. **Connection Pooling**: Configure appropriate pool size
2. **Migrations**: Version control all schema changes
3. **Indexes**: Add indexes for foreign keys và frequently queried columns
4. **Transactions**: Use for multi-table operations
5. **Backups**: Schedule regular backups
6. **Monitoring**: Track slow queries

### Casbin
1. **Least Privilege**: Gán quyền tối thiểu
2. **Domain Separation**: Tách biệt user/cms/api domains
3. **Pattern Matching**: Use wildcards carefully
4. **Testing**: Test cả positive và negative cases
5. **Audit**: Review policies regularly

### Performance
1. **Caching**: Cache frequently accessed data (roles, permissions)
2. **Pagination**: Always paginate list operations
3. **Lazy Loading**: Load related entities only when needed
4. **Batch Operations**: Batch database operations khi có thể
5. **Connection Reuse**: Use HTTP/2 for gRPC, connection pooling for DB

---

## 📖 Development Guide

### Add New Use Case

1. Define DTO in `internal/application/dto/`
2. Create use case in `internal/application/usecase/`
3. Implement business logic using domain repositories
4. Add handler method in `internal/handler/`
5. Update proto file và regenerate
6. Add tests

### Add New Domain Entity

1. Create entity in `internal/domain/model/`
2. Define repository interface in `internal/domain/repository/`
3. Create DAO in `internal/dao/`
4. Implement repository in `internal/infrastructure/persistence/`
5. Create migration file
6. Add seed data if needed

### Add New Casbin Policy

1. Update policy in database:
```sql
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES ('p', 'role_name', 'domain', '/resource/*', 'GET');
```

2. Or via API:
```go
enforcer.AddPolicy("role_name", "domain", "/resource/*", "GET")
```

3. Reload policies:
```go
enforcer.LoadPolicy()
```

---

## 📝 License

Copyright © 2024 E-commerce Platform Team

---

## 👥 Contributors

- **IAM Service Team**
- **GoKits Library Team**

---

## 📞 Support

For issues or questions:
1. Check this README
2. Check logs: `{"level":"error",...}`
3. Review database state
4. Contact development team

---

**Last Updated**: 2024-01-XX  
**Version**: 1.0.0
