# Cấu trúc thư mục IAM Service

```
iam-services/
│
├── cmd/                                # Entry points
│   └── server/
│       └── main.go                     # Main application entry point
│
├── internal/                           # Private application code
│   │
│   ├── config/                         # Configuration management
│   │   └── config.go                   # Load config from env vars
│   │
│   ├── dao/                            # Data Access Objects (Database layer)
│   │   ├── user_dao.go                 # User table CRUD operations
│   │   ├── role_dao.go                 # Role table CRUD operations
│   │   ├── permission_dao.go           # Permission table CRUD operations
│   │   ├── user_role_dao.go            # User-Role relationship operations
│   │   └── role_permission_dao.go      # Role-Permission relationship operations
│   │
│   ├── database/                       # Database connection
│   │   └── database.go                 # DB connection and configuration
│   │
│   ├── domain/                         # Domain entities/models
│   │   └── user.go                     # User, Role, Permission entities
│   │
│   ├── handler/                        # Presentation layer (gRPC)
│   │   ├── grpc_handler.go             # gRPC service implementation
│   │   └── converter.go                # Convert between domain and protobuf
│   │
│   ├── repository/                     # Repository pattern
│   │   ├── user_repository.go          # User business operations
│   │   ├── role_repository.go          # Role business operations
│   │   ├── permission_repository.go    # Permission business operations
│   │   └── authorization_repository.go # Authorization logic
│   │
│   └── service/                        # Business logic layer
│       ├── auth_service.go             # Authentication service
│       ├── authorization_service.go    # Authorization service
│       ├── role_service.go             # Role management service
│       └── permission_service.go       # Permission management service
│
├── pkg/                                # Public packages (reusable)
│   ├── jwt/                            # JWT token management
│   │   └── jwt_manager.go              # Generate and verify JWT tokens
│   │
│   ├── password/                       # Password utilities
│   │   └── password_manager.go         # Hash and verify passwords
│   │
│   └── proto/                          # Protocol Buffer definitions
│       └── iam.proto                   # gRPC service definitions
│
├── migrations/                         # Database migrations
│   ├── 001_init_schema.sql            # Initial schema creation
│   └── 002_seed_data.sql              # Seed default data
│
├── docs/                               # Documentation
│   ├── ARCHITECTURE.md                 # Architecture documentation
│   ├── API.md                          # API documentation
│   ├── DATABASE.md                     # Database schema documentation
│   └── SETUP.md                        # Setup guide
│
├── scripts/                            # Utility scripts
│   ├── setup.sh                        # Initial setup script
│   └── test-api.sh                     # API testing script
│
├── .env.example                        # Environment variables template
├── .gitignore                          # Git ignore rules
├── docker-compose.yml                  # Docker Compose configuration
├── Dockerfile                          # Docker image definition
├── go.mod                              # Go module definition
├── Makefile                            # Build automation
├── README.md                           # Main documentation
└── STRUCTURE.md                        # This file
```

## Layer Responsibilities

### 1. Presentation Layer (Handler)
- **Path**: `internal/handler/`
- **Purpose**: Handle gRPC requests and responses
- **Dependencies**: Service layer

### 2. Business Logic Layer (Service)
- **Path**: `internal/service/`
- **Purpose**: Implement business rules and orchestrate operations
- **Dependencies**: Repository layer, utilities (JWT, Password)

### 3. Data Access Layer (Repository + DAO)
- **Path**: `internal/repository/`, `internal/dao/`
- **Purpose**: 
  - Repository: Business-oriented data operations
  - DAO: Direct database CRUD operations
- **Dependencies**: Domain layer, Database

### 4. Domain Layer
- **Path**: `internal/domain/`
- **Purpose**: Define business entities
- **Dependencies**: None (pure domain models)

## Key Files

### Configuration
- `internal/config/config.go`: Load configuration from environment
- `.env.example`: Template for environment variables

### Entry Point
- `cmd/server/main.go`: Application startup and dependency injection

### Database
- `internal/database/database.go`: Database connection setup
- `migrations/`: SQL migration files

### API Definition
- `pkg/proto/iam.proto`: gRPC service and message definitions

### Documentation
- `README.md`: Main documentation
- `docs/ARCHITECTURE.md`: Detailed architecture explanation
- `docs/API.md`: API usage examples
- `docs/DATABASE.md`: Database schema
- `docs/SETUP.md`: Installation guide

### Build & Deploy
- `Makefile`: Build commands
- `Dockerfile`: Container image
- `docker-compose.yml`: Multi-container setup
- `scripts/setup.sh`: Automated setup

## Dependencies Flow

```
main.go
  ↓
Handler (gRPC)
  ↓
Service (Business Logic)
  ↓
Repository (Data Operations)
  ↓
DAO (Database Access)
  ↓
PostgreSQL
```

## How to Navigate

1. **Start with**: `README.md` for overview
2. **Understanding Architecture**: `docs/ARCHITECTURE.md`
3. **Setup**: `docs/SETUP.md`
4. **API Usage**: `docs/API.md`
5. **Database**: `docs/DATABASE.md`
6. **Code Entry**: `cmd/server/main.go`

## Development Workflow

1. Define API in `pkg/proto/iam.proto`
2. Generate code: `make proto`
3. Add domain models in `internal/domain/`
4. Implement DAO in `internal/dao/`
5. Implement Repository in `internal/repository/`
6. Implement Service in `internal/service/`
7. Implement Handler in `internal/handler/`
8. Wire everything in `cmd/server/main.go`

## Testing Workflow

1. Start database: `make db-create db-migrate`
2. Run service: `make run`
3. Test API: `scripts/test-api.sh`
4. Or use grpcurl manually

