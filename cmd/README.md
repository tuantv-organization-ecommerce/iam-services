# Commands

Entry point for the IAM service using Clean Architecture.

## Directory Structure

```
cmd/
└── server/          # ✅ Production server (Clean Architecture)
    ├── main.go
    └── README.md
```

## Quick Start

### Run Server

```bash
# Run directly
go run ./cmd/server

# Or build and run
go build -o server ./cmd/server
./server
```

## Features

✅ **Clean Architecture** - Domain-driven design with proper dependency inversion
✅ **Shared Logger** - Uses `github.com/tvttt/gokits/logger` for consistent logging
✅ **Infrastructure Layer** - Persistence, security, and authorization implementations
✅ **gRPC Server** - High-performance gRPC communication
✅ **Casbin RBAC** - Role-based access control with Casbin
✅ **JWT Authentication** - Secure token-based authentication
✅ **PostgreSQL** - Reliable data persistence

## Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=iam_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=168h

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=50051

# Logging
LOG_LEVEL=info
LOG_ENCODING=json
```

## Build

```bash
# Development
go run ./cmd/server

# Production build
go build -o server ./cmd/server

# With optimizations
go build -ldflags="-s -w" -o server ./cmd/server

# Run binary
./server
```

## Docker

```bash
# Build image
docker build -t iam-service:latest .

# Run container
docker run -p 50051:50051 \
  -e DB_HOST=postgres \
  -e JWT_SECRET=your-secret \
  iam-service:latest
```

## Documentation

- **Architecture**: `../ARCHITECTURE_NEW.md`
- **Clean Architecture**: `../CLEAN_REFACTOR_COMPLETE.md`
- **Shared Packages**: `../../SHARED_PACKAGES.md`
- **Setup Guide**: `../QUICK_START.md`
- **API Documentation**: `../docs/API.md`

## Troubleshooting

### Build Issues

```bash
# Regenerate proto files
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       pkg/proto/iam.proto

# Update dependencies
go mod tidy
```

### Database Issues

```bash
# Run migrations
psql -U postgres -d iam_db -f migrations/001_init_schema.sql
psql -U postgres -d iam_db -f migrations/002_seed_data.sql
psql -U postgres -d iam_db -f migrations/003_casbin_tables.sql
psql -U postgres -d iam_db -f migrations/004_casbin_seed_data.sql
```

### Logger Issues

The service uses shared logger from `gokits`:
```go
import "github.com/tvttt/gokits/logger"

log, _ := logger.NewProduction()
defer logger.Sync(log)
```

## Development

### Adding New Features

1. Define in **domain layer** (interfaces/ports)
2. Implement **use cases** in application layer
3. Add **infrastructure** implementations
4. Wire dependencies in `main.go`

### Project Structure

```
cmd/server/
├── main.go              # Entry point with DI
└── README.md           # This file

internal/
├── domain/             # Business logic (ports)
├── application/        # Use cases & DTOs
├── infrastructure/     # Implementations (adapters)
└── handler/           # gRPC handlers

pkg/
├── casbin/            # Casbin enforcer
├── jwt/               # JWT manager
├── password/          # Password manager
└── proto/             # Generated proto files
```

---

**Production Ready** ✅ - This version is tested and ready for deployment.
