# IAM Service - Identity and Access Management

**Version:** 1.0.0  
**Go Version:** 1.19  
**Web Framework:** Gin v1.9.1  
**Author:** E-commerce Platform Team

---

## 📖 Table of Contents

1. [Overview](#-overview)
2. [Features](#-features)
3. [Architecture](#-architecture)
4. [Tech Stack](#-tech-stack)
5. [Quick Start](#-quick-start)
6. [Configuration](#-configuration)
7. [API Documentation](#-api-documentation)
8. [Authorization (Casbin RBAC)](#-authorization-casbin-rbac)
9. [Database](#-database)
10. [Development](#-development)
11. [Testing](#-testing)
12. [CI/CD](#-cicd)
13. [Troubleshooting](#-troubleshooting)
14. [Best Practices](#-best-practices)

---

## 🎯 Overview

IAM Service is a comprehensive Identity and Access Management system for e-commerce platforms, built with **Clean Architecture** and **Casbin RBAC**.

### Key Features

- ✅ **Clean Architecture**: Clear separation of concerns
- ✅ **Dual Protocol**: gRPC (port 50051) + Gin HTTP (port 8080)
- ✅ **Multi-Domain RBAC**: Separated User/App and CMS authorization
- ✅ **JWT Authentication**: Access + Refresh tokens
- ✅ **High Performance**: Gin web framework
- ✅ **Swagger UI**: Protected with Basic Authentication
- ✅ **PostgreSQL**: With connection pooling and migrations
- ✅ **Structured Logging**: Uber Zap logger
- ✅ **Panic Recovery**: Multi-layered safety
- ✅ **CI/CD Ready**: GitHub Actions workflows

---

## 🚀 Features

### Authentication
- User registration, login, logout
- JWT token generation and validation
- Refresh token support
- Token verification

### Authorization
- Role-based access control (RBAC)
- Permission management
- User-role assignments
- Policy enforcement with Casbin

### Multi-Domain Authorization

#### 1. User/App Authorization
- **Purpose**: End users accessing web/mobile apps
- **Domains**: `user`, `api`
- **Table**: `casbin_rule_user`
- **Roles**: `user`, `premium_user`, `api_admin`

#### 2. CMS Authorization
- **Purpose**: Admin/staff accessing CMS
- **Structure**: Tab-based access control
- **Table**: `casbin_rule_cms`
- **Roles**: `cms_admin`, `cms_product_manager`, etc.
- **Tabs**: `product`, `inventory`, `order`, `user`, `report`, `setting`

---

## 🏗️ Architecture

### Clean Architecture Layers

```
┌────────────────────────────────────────────┐
│         CMD (Entry Point)                  │
│         cmd/server/main.go                 │
└──────────────┬─────────────────────────────┘
               ↓
┌──────────────────────────────────────────────┐
│      CONTAINER (Dependency Injection)        │
│      internal/container/container.go         │
└──────────────┬───────────────────────────────┘
               ↓
┌──────────────────────────────────────────────┐
│      ADAPTER LAYER (Handlers & Routes)       │
│      • Gin HTTP Handlers                     │
│      • gRPC Handlers                         │
│      • Middleware (Recovery, Logger, CORS)   │
└──────────────┬───────────────────────────────┘
               ↓
┌──────────────────────────────────────────────┐
│      APPLICATION LAYER (Use Cases)           │
│      • Auth Use Cases                        │
│      • Role Management                       │
│      • DTOs                                  │
└──────────────┬───────────────────────────────┘
               ↓
┌──────────────────────────────────────────────┐
│      DOMAIN LAYER (Business Logic)           │
│      • Domain Models                         │
│      • Repository Interfaces                 │
│      • Service Interfaces                    │
└──────────────┬───────────────────────────────┘
               ↓
┌──────────────────────────────────────────────┐
│      INFRASTRUCTURE LAYER                    │
│      • PostgreSQL Persistence                │
│      • JWT Security                          │
│      • Casbin Authorization                  │
└──────────────────────────────────────────────┘
```

### Request Flow

```
HTTP Request → Gin Router → Middleware → Gin Handler → Service → Repository → Database
   ↓                          ↓
gRPC Request              Recovery, Logger, CORS
```

---

## 🛠️ Tech Stack

### Core
- **Language**: Go 1.24
- **Web Framework**: Gin v1.9.1
- **RPC Framework**: gRPC
- **Database**: PostgreSQL 12+
- **Authorization**: Casbin v2
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Logging**: Uber Zap

### Key Dependencies
```go
github.com/gin-gonic/gin v1.9.1
github.com/gin-contrib/cors v1.5.0
github.com/casbin/casbin/v2
github.com/golang-jwt/jwt/v5
github.com/lib/pq
github.com/google/uuid
go.uber.org/zap
google.golang.org/grpc
gorm.io/gorm
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL 12+
- protoc (for gRPC development)

### 1. Clone & Install Dependencies

```bash
cd iam-services
go mod download
go mod tidy
```

### 2. Setup Database

```bash
# Create database
psql -U postgres
CREATE DATABASE iam_db;
\q

# Run migrations
psql -U postgres -d iam_db -f migrations/001_init_schema.sql
psql -U postgres -d iam_db -f migrations/002_seed_data.sql
psql -U postgres -d iam_db -f migrations/003_casbin_tables.sql
psql -U postgres -d iam_db -f migrations/004_casbin_seed_data.sql
psql -U postgres -d iam_db -f migrations/005_separate_user_cms_authorization.sql
psql -U postgres -d iam_db -f migrations/006_seed_separated_authorization.sql
```

### 3. Configure Environment

```bash
# Copy example config
cp .env.example .env

# Edit .env file
nano .env
```

Required environment variables:
```env
# Server
HTTP_PORT=8080
GRPC_PORT=50051

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=iam_db

# JWT (CRITICAL: Change in production!)
JWT_SECRET=your-secure-random-64-char-secret
JWT_EXPIRATION_HOURS=24

# Swagger (optional)
SWAGGER_ENABLED=true
SWAGGER_AUTH_USERNAME=admin
SWAGGER_AUTH_PASSWORD=changeme
```

### 4. Run the Service

```bash
# Development
go run cmd/server/main.go

# Or build and run
go build -o bin/iam-service cmd/server/main.go
./bin/iam-service
```

Service will start on:
- **gRPC**: `localhost:50051`
- **HTTP**: `http://localhost:8080`
- **Swagger UI**: `http://localhost:8080/swagger/` (requires auth)

---

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `HTTP_PORT` | HTTP server port | `8080` | Yes |
| `GRPC_PORT` | gRPC server port | `50051` | Yes |
| `DB_HOST` | PostgreSQL host | `localhost` | Yes |
| `DB_PORT` | PostgreSQL port | `5432` | Yes |
| `DB_USER` | Database user | `postgres` | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | `iam_db` | Yes |
| `JWT_SECRET` | JWT secret key (min 32 chars) | - | Yes |
| `JWT_EXPIRATION_HOURS` | Access token expiration | `24` | Yes |
| `CASBIN_MODEL_PATH` | Casbin model file | `./configs/rbac_model.conf` | Yes |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `info` | No |
| `SWAGGER_ENABLED` | Enable Swagger UI | `true` | No |
| `SWAGGER_AUTH_USERNAME` | Swagger username | `admin` | No |
| `SWAGGER_AUTH_PASSWORD` | Swagger password | `changeme` | No |

---

## 📚 API Documentation

### Gin HTTP API

**Base URL**: `http://localhost:8080`

#### Authentication
```bash
POST   /v1/auth/register     # Register new user
POST   /v1/auth/login        # Login
POST   /v1/auth/refresh      # Refresh access token
POST   /v1/auth/logout       # Logout
POST   /v1/auth/verify       # Verify token
```

#### Role Management
```bash
POST   /v1/roles             # Create role
GET    /v1/roles             # List roles (pagination)
GET    /v1/roles/:id         # Get role by ID
PUT    /v1/roles/:id         # Update role
DELETE /v1/roles/:id         # Delete role
POST   /v1/roles/assign      # Assign role to user
POST   /v1/roles/remove      # Remove role from user
GET    /v1/users/:user_id/roles  # Get user's roles
```

#### Permission Management
```bash
POST   /v1/permissions       # Create permission
GET    /v1/permissions       # List permissions
DELETE /v1/permissions/:id   # Delete permission
POST   /v1/permissions/check # Check permission
```

#### CMS Authorization
```bash
POST   /v1/cms/roles         # Create CMS role
GET    /v1/cms/roles         # List CMS roles
POST   /v1/cms/roles/assign  # Assign CMS role
POST   /v1/cms/roles/remove  # Remove CMS role
GET    /v1/cms/users/:user_id/tabs  # Get user's CMS tabs
```

#### Access Control
```bash
POST   /v1/access/api        # Check API access
POST   /v1/access/cms        # Check CMS access
POST   /v1/policies/enforce  # Enforce policy
```

#### API Resources
```bash
POST   /v1/api/resources     # Create API resource
GET    /v1/api/resources     # List API resources
```

#### Health Check
```bash
GET    /health               # Health check endpoint
```

### Example: Register & Login

```bash
# Register
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'

# Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecurePass123!"
  }'
```

### gRPC API

**Server**: `localhost:50051`

Use `grpcurl` for testing:
```bash
# List services
grpcurl -plaintext localhost:50051 list

# Call Login
grpcurl -plaintext -d '{
  "username": "johndoe",
  "password": "SecurePass123!"
}' localhost:50051 iam.IAMService/Login
```

### Swagger UI

**URL**: `http://localhost:8080/swagger/`  
**Auth**: Basic Authentication (admin/changeme by default)

Features:
- Interactive API documentation
- Try out APIs directly from browser
- View request/response schemas
- Auto-generated from proto files

---

## 🔐 Authorization (Casbin RBAC)

### Two Independent Authorization Systems

#### 1. User/App Authorization
- **Model**: `configs/rbac_user_model.conf`
- **Database**: `casbin_rule_user`
- **Domains**: `user`, `api`

**Example Policies**:
```
# Regular user - can browse products, create orders
p, user, user, /api/v1/products, GET
p, user, user, /api/v1/orders, (GET|POST)

# Premium user - more access
p, premium_user, api, /api/v1/products/**, (GET|POST)

# Admin - full access
p, api_admin, api, /api/v1/**, (GET|POST|PUT|DELETE)
```

**Check Access**:
```bash
curl -X POST http://localhost:8080/v1/access/api \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "api_path": "/api/v1/products",
    "method": "POST"
  }'
```

#### 2. CMS Authorization
- **Model**: `configs/rbac_cms_model.conf`
- **Database**: `casbin_rule_cms`
- **Structure**: Tab-based

**CMS Tabs**:
- `product` - Product management
- `inventory` - Inventory management
- `order` - Order management
- `user` - User management
- `report` - Reports & analytics
- `setting` - System settings

**Example Policies**:
```
# CMS Admin - full access to all tabs
p, cms_admin, product, /api/v1/products/*, (GET|POST|PUT|DELETE)
p, cms_admin, inventory, /api/v1/inventory/*, (GET|POST|PUT|DELETE)

# Product Manager - product & inventory only
p, cms_product_manager, product, /api/v1/products/*, (GET|POST|PUT)
p, cms_product_manager, inventory, /api/v1/inventory/*, (GET|POST|PUT)

# Viewer - read only
p, cms_viewer, product, /api/v1/products/*, GET
```

**Check CMS Access**:
```bash
curl -X POST http://localhost:8080/v1/access/cms \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "admin-789",
    "cms_tab": "product",
    "action": "POST"
  }'
```

### Pattern Matching

- **KeyMatch2**: `/api/v1/products/*` matches `/api/v1/products/123`
- **RegexMatch**: `(GET|POST)` matches "GET" OR "POST"

---

## 🗄️ Database

### Schema Overview

**Core Tables**:
- `users` - User accounts
- `roles` - User/app roles
- `permissions` - Permissions
- `user_roles` - User-role assignments
- `role_permissions` - Role-permission assignments

**Casbin Tables**:
- `casbin_rule_user` - User/app authorization policies
- `casbin_rule_cms` - CMS authorization policies

**CMS Tables**:
- `cms_roles` - CMS roles with tabs array
- `cms_tab_apis` - Maps tabs to APIs (many-to-many)
- `user_cms_roles` - User-CMS role assignments

**API Resources**:
- `api_resources` - Tracks API endpoints

### Migrations

Located in `migrations/` directory:
```
001_init_schema.sql                          # Core schema
002_seed_data.sql                            # Initial data
003_casbin_tables.sql                        # Casbin tables
004_casbin_seed_data.sql                     # Casbin seed
005_separate_user_cms_authorization.sql      # Separated auth architecture
006_seed_separated_authorization.sql         # Separated auth data
```

### Connection Pool

Configured for optimal performance:
```go
db.SetMaxOpenConns(25)              // Max open connections
db.SetMaxIdleConns(5)               // Max idle connections
db.SetConnMaxLifetime(5 * time.Minute)
```

---

## 💻 Development

### Project Structure

```
iam-services/
├── cmd/server/                # Entry point
├── internal/
│   ├── app/                   # Application lifecycle
│   ├── container/             # Dependency injection
│   ├── handler/               # HTTP/gRPC handlers
│   ├── router/                # Gin routes
│   ├── middleware/            # Middleware
│   ├── application/           # Use cases & DTOs
│   ├── domain/                # Business logic
│   └── infrastructure/        # Implementation details
├── pkg/
│   ├── jwt/                   # JWT utilities
│   ├── password/              # Password hashing
│   ├── casbin/                # Casbin enforcer
│   └── proto/                 # Proto definitions
├── configs/                   # Configuration files
├── migrations/                # SQL migrations
└── scripts/                   # Utility scripts
```

### Add New Endpoint

1. **Define DTO** in `internal/application/dto/`
2. **Create use case** in `internal/application/usecase/`
3. **Add handler method** in `internal/handler/gin_handler.go`
4. **Register route** in `internal/router/gin_router.go`
5. **Write tests**

### Code Quality Tools

#### Linting

```bash
# Install golangci-lint (Go 1.24 compatible)
.\scripts\install-golangci-lint.ps1 -Version "v1.54.2"

# Run lint
.\scripts\lint.ps1              # All checks
.\scripts\lint.ps1 -Fix         # Auto-fix issues
.\scripts\lint.ps1 -Fast        # Fast mode
```

#### Pre-Commit Checks

```bash
# Run all checks (lint + build + test)
.\scripts\check-all.ps1         # Windows
./scripts/check-all.sh          # Linux/Mac
```

#### Makefile Commands

```bash
make lint           # Run golangci-lint
make lint-fix       # Auto-fix issues
make build          # Build binary
make test           # Run tests
make check-all      # Lint + Build + Test
```

---

## 🧪 Testing

### Run Tests

```bash
# All tests
go test ./...

# Specific package
go test -v ./pkg/jwt/
go test -v ./internal/service/

# With coverage
go test -cover ./...

# Integration tests (requires database)
go test -v ./internal/dao/
```

### Test Files

- `pkg/jwt/jwt_manager_test.go` - JWT token tests
- `pkg/password/password_manager_test.go` - Password hashing tests
- `internal/dao/user_dao_test.go` - Database access tests
- `internal/service/auth_service_test.go` - Service layer tests

---

## 🔄 CI/CD

### GitHub Actions Workflow

Located in `.github/workflows/ci-cd.yml`

**Pipeline Stages**:
1. **Lint** - Code quality checks with golangci-lint
2. **Test** - Unit tests with PostgreSQL service
3. **Build** - Binary build and artifact upload
4. **Security** - Vulnerability scanning with Trivy
5. **Docker** - Build and push Docker image (main/develop only)
6. **Deploy** - Auto-deploy to staging/production (disabled by default)

### Quick CI Setup

```bash
# 1. Create .env.example
cp .env.template .env.example

# 2. Download dependencies
go mod download && go mod tidy

# 3. Create feature branch and push
git checkout -b feature/setup-ci
git add .
git commit -m "ci: setup CI/CD pipeline"
git push origin feature/setup-ci
```

**Expected CI Results**:
- ✅ Lint passes
- ✅ Tests pass with coverage report
- ✅ Build creates binary artifact
- ✅ Security scan completes

### Enable Deployment

To enable deployment (when servers are ready):
1. Configure GitHub Secrets:
   - `DOCKER_USERNAME`, `DOCKER_PASSWORD`
   - `STAGING_HOST`, `STAGING_USER`, `STAGING_SSH_KEY`
   - `PRODUCTION_HOST`, `PRODUCTION_USER`, `PRODUCTION_SSH_KEY`
2. Uncomment deploy jobs in `.github/workflows/ci-cd.yml`
3. Setup servers with Docker + Docker Compose

---

## 🐛 Troubleshooting

### Connection Refused

**Issue**: Cannot connect to service

**Solution**:
```bash
# Check if service is running
ps aux | grep iam-service

# Check ports
netstat -an | grep 8080
netstat -an | grep 50051
```

### Database Connection Failed

**Issue**: Cannot connect to PostgreSQL

**Solution**:
```bash
# Check PostgreSQL status
pg_isready

# Test connection
psql -U postgres -d iam_db

# Verify credentials in .env
cat .env | grep DB_
```

### Migrations Not Applied

**Issue**: Tables do not exist

**Solution**:
```bash
# Run all migrations in order
for i in {1..6}; do
  psql -U postgres -d iam_db -f migrations/00${i}_*.sql
done

# Verify tables
psql -U postgres -d iam_db -c "\dt"
```

### Token Invalid/Expired

**Issue**: JWT token rejected

**Solution**:
- Ensure `JWT_SECRET` matches between restarts
- Check token expiration time
- Login again to get fresh token

### Swagger UI Not Loading

**Issue**: 404 or blank page

**Solution**:
```bash
# Check if Swagger is enabled
echo $SWAGGER_ENABLED

# Verify spec file exists
ls pkg/proto/iam.swagger.json

# Regenerate proto files
.\scripts\generate-proto-simple.ps1
```

### Linting Fails

**Issue**: golangci-lint errors

**Solution**:
```bash
# Use compatible version for Go 1.24
.\scripts\install-golangci-lint.ps1 -Version "v1.54.2"

# Auto-fix issues
.\scripts\lint.ps1 -Fix

# Format code
go fmt ./...
```

---

## ✅ Best Practices

### Security

- ✅ Use strong JWT secrets (min 64 chars in production)
- ✅ Hash passwords with bcrypt (cost >= 10)
- ✅ Enable TLS in production
- ✅ Implement rate limiting
- ✅ Use parameterized SQL queries
- ✅ Audit critical operations
- ✅ Use secrets management (Vault, AWS Secrets Manager)

### Architecture

- ✅ Follow Clean Architecture principles
- ✅ Use dependency injection
- ✅ Keep layers independent
- ✅ Define focused interfaces
- ✅ Return domain errors, map in handlers

### Database

- ✅ Use connection pooling
- ✅ Version control migrations
- ✅ Add indexes for frequently queried columns
- ✅ Use transactions for multi-table operations
- ✅ Schedule regular backups

### Code Quality

- ✅ Run `golangci-lint` before commit
- ✅ Write tests for new features
- ✅ Maintain >70% code coverage
- ✅ Add comments for all exported symbols
- ✅ Use structured logging
- ✅ Handle errors properly

### Performance

- ✅ Cache frequently accessed data
- ✅ Paginate list operations
- ✅ Use lazy loading
- ✅ Batch database operations
- ✅ Monitor slow queries

---

## 📖 Additional Documentation

For more detailed information, see:

- **[CI/CD Setup Guide](CI_CD_SETUP_GUIDE.md)** - Complete CI/CD setup instructions
- **[Authorization Architecture](AUTHORIZATION_ARCHITECTURE.md)** - Detailed authorization design
- **[CI/CD Error Fixes](fix_error_ci_cd.md)** - Common CI/CD issues and solutions
- **[Linting Setup](LINTING_SETUP.md)** - golangci-lint configuration and usage
- **[Swagger Guide](SWAGGER_GUIDE.md)** - Swagger UI detailed documentation
- **[Gin Refactoring Summary](GIN_REFACTORING_SUMMARY.md)** - Gin framework migration details

---

## 🤝 Contributing

1. Create feature branch: `git checkout -b feature/my-feature`
2. Make changes following code quality standards
3. Run checks: `.\scripts\check-all.ps1`
4. Commit: `git commit -m "feat: description"`
5. Push: `git push origin feature/my-feature`
6. Create Pull Request

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

---

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-org/iam-services/issues)
- **Documentation**: See additional docs above
- **Check Logs**: Use structured logging to debug issues

---

## 📄 License

Copyright © 2024 E-commerce Platform Team

---

**Last Updated**: November 2024  
**Status**: ✅ Production Ready
