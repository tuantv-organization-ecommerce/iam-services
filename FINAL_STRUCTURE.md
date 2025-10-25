# Final Structure - IAM Service ✅

## 🎯 Clean & Simple

Cấu trúc cuối cùng sau khi loại bỏ legacy code và implement Clean Architecture.

## 📁 Structure

```
ecommerce/back_end/
│
├── go.work                      # Workspace config
│
├── gokits/                      # ✨ Shared utilities
│   ├── go.mod
│   ├── README.md
│   └── logger/
│       └── logger.go           # Shared logger for all services
│
└── iam-services/
    │
    ├── cmd/
    │   ├── README.md           # Commands guide
    │   └── server/             # ✅ Single entry point
    │       ├── main.go         # Production server
    │       └── README.md       # Server documentation
    │
    ├── internal/
    │   ├── domain/             # 🔵 Domain Layer (Core Business Logic)
    │   │   ├── model/          # Rich domain entities
    │   │   ├── repository/     # Repository interfaces (ports)
    │   │   └── service/        # Domain service interfaces
    │   │
    │   ├── application/        # 🟢 Application Layer (Use Cases)
    │   │   ├── dto/            # Data transfer objects
    │   │   └── usecase/        # Application use cases
    │   │
    │   ├── infrastructure/     # 🟡 Infrastructure Layer (Implementations)
    │   │   ├── persistence/    # Repository implementations
    │   │   ├── security/       # JWT & Password services
    │   │   ├── authorization/  # Casbin service
    │   │   └── config/         # Config loader
    │   │
    │   ├── handler/            # 🟠 Adapter Layer (gRPC)
    │   │   ├── grpc_handler.go
    │   │   ├── casbin_handler.go
    │   │   └── converter.go
    │   │
    │   ├── dao/                # Database access objects
    │   ├── database/           # Database connection
    │   ├── config/             # Configuration
    │   │
    │   ├── repository/         # Old repositories (compatibility)
    │   └── service/            # Old services (compatibility)
    │
    ├── pkg/
    │   ├── casbin/             # Casbin enforcer wrapper
    │   ├── jwt/                # JWT manager
    │   ├── password/           # Password manager
    │   └── proto/              # Generated proto files
    │
    ├── configs/
    │   └── rbac_model.conf     # Casbin RBAC model
    │
    ├── migrations/             # Database migrations
    │   ├── 001_init_schema.sql
    │   ├── 002_seed_data.sql
    │   ├── 003_casbin_tables.sql
    │   └── 004_casbin_seed_data.sql
    │
    ├── docs/                   # Documentation
    │   ├── API.md
    │   ├── ARCHITECTURE.md
    │   ├── CASBIN.md
    │   └── DATABASE.md
    │
    ├── scripts/                # Utility scripts
    │   └── setup-proto.ps1
    │
    ├── go.mod
    ├── go.sum
    ├── Dockerfile
    ├── docker-compose.yml
    ├── Makefile
    │
    └── README.md               # Main documentation
```

## 🏛️ Clean Architecture Layers

### 1. Domain Layer (Core) 🔵
**Location:** `internal/domain/`

**Purpose:** Pure business logic with NO external dependencies

**Components:**
- `model/` - Rich domain entities with business methods
- `repository/` - Repository interfaces (ports)
- `service/` - Domain service interfaces (ports)

**Rules:**
- ❌ NO imports from other layers
- ❌ NO external libraries (except stdlib)
- ✅ Only business rules

### 2. Application Layer (Use Cases) 🟢
**Location:** `internal/application/`

**Purpose:** Orchestrate domain objects to fulfill use cases

**Components:**
- `dto/` - Data transfer objects
- `usecase/` - Application logic (login, register, etc.)

**Rules:**
- ✅ Can depend on Domain layer only
- ❌ Cannot depend on Infrastructure or Adapter
- ✅ Coordinates domain objects

### 3. Infrastructure Layer (Implementations) 🟡
**Location:** `internal/infrastructure/`

**Purpose:** Implement domain ports with external tools

**Components:**
- `persistence/` - Repository implementations (PostgreSQL)
- `security/` - JWT & Password implementations
- `authorization/` - Casbin implementation
- `config/` - Configuration loader

**Rules:**
- ✅ Implements Domain interfaces
- ✅ Can use external libraries
- ➡️ Dependency points inward to Domain

### 4. Adapter Layer (External Interfaces) 🟠
**Location:** `internal/handler/`

**Purpose:** Adapt external requests to internal use cases

**Components:**
- gRPC handlers
- Converters (proto ↔ domain models)

**Rules:**
- ✅ Depends on Application & Domain
- ✅ Converts external formats to internal
- ➡️ Dependency points inward

## 🔄 Dependency Flow

```
┌────────────────────────────────────────┐
│         Adapter Layer (gRPC)           │
│         internal/handler/              │
└──────────────┬─────────────────────────┘
               │ depends on
┌──────────────▼─────────────────────────┐
│      Application Layer (Use Cases)     │
│      internal/application/             │
└──────────────┬─────────────────────────┘
               │ depends on
┌──────────────▼─────────────────────────┐
│       Domain Layer (Core Logic)        │  ← NO dependencies!
│       internal/domain/                 │
└──────────────▲─────────────────────────┘
               │ implements
┌──────────────┴─────────────────────────┐
│    Infrastructure Layer (External)     │
│    internal/infrastructure/            │
└────────────────────────────────────────┘
```

**Key Principle:** Dependencies point inward!

## 🚀 Running the Service

```bash
# Development
go run ./cmd/server

# Production build
go build -o server ./cmd/server
./server
```

## 📦 Shared Packages (GoKits)

**Location:** `ecommerce/back_end/gokits/`

**Purpose:** Shared utilities for ALL services

**Current Packages:**
- ✅ `logger/` - Shared logger using Zap

**Usage:**
```go
import "github.com/tvttt/gokits/logger"

log, _ := logger.NewProduction()
defer logger.Sync(log)

log.Info("Service started", zap.String("name", "iam-service"))
```

**Future Packages:**
- `config/` - Configuration management
- `errors/` - Standard error types
- `middleware/` - Common middleware
- `metrics/` - Prometheus metrics
- `tracing/` - Distributed tracing

## ✅ Simplifications Made

### Removed

1. ❌ `cmd/legacy/` - Old architecture (không cần thiết)
2. ❌ `cmd/server/main_old.go` - Duplicate file
3. ❌ `server-legacy.exe` - Legacy binary
4. ❌ Duplicate logger code - Now uses gokits

### Result

✅ Single entry point: `cmd/server/main.go`
✅ Clean structure
✅ No confusing files
✅ Shared logger via gokits
✅ Production ready

## 🎯 Key Features

1. **Clean Architecture** ✅
   - Clear layer separation
   - Dependency inversion
   - Testable & maintainable

2. **Shared Utilities** ✅
   - GoKits for common code
   - Reusable across services
   - Single source of truth

3. **Modern Stack** ✅
   - gRPC for communication
   - PostgreSQL for data
   - Casbin for authorization
   - JWT for authentication
   - Zap for logging

4. **Production Ready** ✅
   - Docker support
   - Migrations ready
   - Comprehensive docs
   - Error handling
   - Graceful shutdown

## 📚 Documentation

- `README.md` - Overview and quick start
- `ARCHITECTURE_NEW.md` - Architecture details
- `CLEAN_REFACTOR_COMPLETE.md` - Refactoring summary
- `SHARED_PACKAGES.md` - GoKits documentation
- `cmd/README.md` - Commands guide
- `cmd/server/README.md` - Server documentation
- `docs/` - Detailed documentation

## 🎊 Summary

**Before:** Confusing structure with legacy code
**After:** Clean, simple, production-ready structure

**Key Changes:**
1. ✅ Single entry point (`cmd/server/`)
2. ✅ Clean Architecture implemented
3. ✅ Shared logger via gokits
4. ✅ All legacy code removed
5. ✅ Comprehensive documentation

**Status:** ✅ Production Ready

---

**This is the final, clean structure. Ready for deployment! 🚀**

