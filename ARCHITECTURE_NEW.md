# IAM Service - Clean Architecture Structure

## Tổng quan

IAM Service được tái cấu trúc theo **Clean Architecture** với tách biệt rõ ràng giữa các layers.

## Cấu trúc thư mục mới

```
iam-services/
│
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
│
├── internal/
│   │
│   ├── adapter/                       # ADAPTER LAYER (Presentation)
│   │   ├── grpc/                      # gRPC handlers
│   │   │   ├── handler/              # gRPC service implementations
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── authorization_handler.go
│   │   │   │   ├── role_handler.go
│   │   │   │   ├── permission_handler.go
│   │   │   │   └── casbin_handler.go
│   │   │   └── interceptor/          # gRPC interceptors (auth, logging)
│   │   │       └── auth_interceptor.go
│   │   │
│   │   └── converter/                 # Data conversion between layers
│   │       ├── user_converter.go
│   │       ├── role_converter.go
│   │       └── casbin_converter.go
│   │
│   ├── application/                   # APPLICATION LAYER (Use Cases)
│   │   ├── usecase/                  # Use case implementations
│   │   │   ├── auth/
│   │   │   │   ├── register_usecase.go
│   │   │   │   ├── login_usecase.go
│   │   │   │   └── refresh_token_usecase.go
│   │   │   ├── authorization/
│   │   │   │   ├── assign_role_usecase.go
│   │   │   │   └── check_permission_usecase.go
│   │   │   ├── role/
│   │   │   │   ├── create_role_usecase.go
│   │   │   │   └── update_role_usecase.go
│   │   │   └── casbin/
│   │   │       ├── check_api_access_usecase.go
│   │   │       ├── check_cms_access_usecase.go
│   │   │       └── manage_cms_role_usecase.go
│   │   │
│   │   └── dto/                      # Data Transfer Objects
│   │       ├── auth_dto.go
│   │       ├── role_dto.go
│   │       └── casbin_dto.go
│   │
│   ├── domain/                        # DOMAIN LAYER (Business Logic)
│   │   ├── model/                    # Domain entities/aggregates
│   │   │   ├── user.go
│   │   │   ├── role.go
│   │   │   ├── permission.go
│   │   │   ├── cms_role.go
│   │   │   └── api_resource.go
│   │   │
│   │   ├── repository/               # Repository interfaces (Ports)
│   │   │   ├── user_repository.go
│   │   │   ├── role_repository.go
│   │   │   ├── permission_repository.go
│   │   │   ├── authorization_repository.go
│   │   │   ├── cms_repository.go
│   │   │   └── api_resource_repository.go
│   │   │
│   │   ├── service/                  # Domain services
│   │   │   ├── password_service.go
│   │   │   ├── token_service.go
│   │   │   └── authorization_service.go
│   │   │
│   │   └── valueobject/              # Value objects
│   │       ├── email.go
│   │       ├── username.go
│   │       └── token.go
│   │
│   └── infrastructure/                # INFRASTRUCTURE LAYER
│       ├── persistence/              # Data persistence
│       │   ├── postgres/             # PostgreSQL implementation
│       │   │   ├── dao/              # Data Access Objects
│       │   │   │   ├── user_dao.go
│       │   │   │   ├── role_dao.go
│       │   │   │   ├── permission_dao.go
│       │   │   │   ├── user_role_dao.go
│       │   │   │   ├── role_permission_dao.go
│       │   │   │   ├── cms_role_dao.go
│       │   │   │   ├── user_cms_role_dao.go
│       │   │   │   └── api_resource_dao.go
│       │   │   │
│       │   │   ├── repository/       # Repository implementations
│       │   │   │   ├── user_repository_impl.go
│       │   │   │   ├── role_repository_impl.go
│       │   │   │   ├── permission_repository_impl.go
│       │   │   │   ├── authorization_repository_impl.go
│       │   │   │   ├── cms_repository_impl.go
│       │   │   │   └── api_resource_repository_impl.go
│       │   │   │
│       │   │   └── connection.go     # Database connection
│       │   │
│       │   └── migrations/           # Database migrations
│       │
│       ├── casbin/                   # Casbin authorization
│       │   ├── enforcer.go
│       │   └── adapter.go
│       │
│       ├── security/                 # Security implementations
│       │   ├── jwt/
│       │   │   └── jwt_manager.go
│       │   └── password/
│       │       └── bcrypt_hasher.go
│       │
│       └── config/                   # Configuration
│           └── config.go
│
├── pkg/                               # PUBLIC PACKAGES
│   ├── proto/                        # Protocol buffers
│   │   └── iam.proto
│   │
│   ├── logger/                       # Logging utilities
│   │   └── logger.go
│   │
│   └── errors/                       # Custom errors
│       └── errors.go
│
├── configs/                           # CONFIGURATION FILES
│   ├── rbac_model.conf               # Casbin model
│   └── config.yaml                   # Application config
│
├── migrations/                        # DATABASE MIGRATIONS
│   ├── 001_init_schema.sql
│   ├── 002_seed_data.sql
│   ├── 003_casbin_tables.sql
│   └── 004_casbin_seed_data.sql
│
├── docs/                             # DOCUMENTATION
│   ├── ARCHITECTURE.md
│   ├── SETUP.md
│   ├── API.md
│   ├── DATABASE.md
│   └── CASBIN.md
│
├── scripts/                          # UTILITY SCRIPTS
│   ├── setup.sh
│   └── test-api.sh
│
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## Layers Explanation

### 1. **Adapter Layer** (Presentation Layer)

**Path**: `internal/adapter/`

**Responsibility**:
- Handle external communications (gRPC, REST API)
- Convert external requests to internal DTOs
- Convert internal results to external responses
- Input validation
- Error handling và response formatting

**Components**:
- **gRPC Handlers**: Implement gRPC service definitions
- **Converters**: Convert between Protocol Buffers and DTOs
- **Interceptors**: Authentication, logging, metrics

**Example**:
```go
// internal/adapter/grpc/handler/auth_handler.go
type AuthHandler struct {
    registerUseCase *usecase.RegisterUseCase
    loginUseCase    *usecase.LoginUseCase
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    // 1. Convert protobuf to DTO
    dto := converter.RegisterRequestToDTO(req)
    
    // 2. Call use case
    result, err := h.registerUseCase.Execute(ctx, dto)
    
    // 3. Convert result to protobuf
    return converter.RegisterResultToResponse(result), nil
}
```

---

### 2. **Application Layer** (Use Case Layer)

**Path**: `internal/application/`

**Responsibility**:
- Orchestrate business logic flow
- Coordinate domain entities and services
- Transaction management
- Application-specific business rules
- No dependency on external frameworks

**Components**:
- **Use Cases**: One use case per business operation
- **DTOs**: Data transfer between layers

**Example**:
```go
// internal/application/usecase/auth/register_usecase.go
type RegisterUseCase struct {
    userRepo       domain.UserRepository
    passwordSvc    domain.PasswordService
    eventPublisher EventPublisher
}

func (uc *RegisterUseCase) Execute(ctx context.Context, dto RegisterDTO) (*RegisterResult, error) {
    // 1. Business validation
    if err := uc.validateRegistration(dto); err != nil {
        return nil, err
    }
    
    // 2. Create domain entity
    user := domain.NewUser(dto.Username, dto.Email, dto.FullName)
    
    // 3. Hash password using domain service
    hashedPassword, err := uc.passwordSvc.Hash(dto.Password)
    user.SetPassword(hashedPassword)
    
    // 4. Save to repository
    if err := uc.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // 5. Publish event
    uc.eventPublisher.Publish(UserRegisteredEvent{UserID: user.ID})
    
    return &RegisterResult{UserID: user.ID}, nil
}
```

---

### 3. **Domain Layer** (Business Logic Layer)

**Path**: `internal/domain/`

**Responsibility**:
- Core business logic
- Business rules và validations
- Domain entities và aggregates
- Domain events
- **100% framework independent**

**Components**:
- **Models**: Domain entities with business logic
- **Repository Interfaces**: Define data access contracts (Ports)
- **Domain Services**: Business logic spanning multiple entities
- **Value Objects**: Immutable objects representing concepts

**Example**:
```go
// internal/domain/model/user.go
type User struct {
    id           string
    username     Username      // Value object
    email        Email         // Value object
    passwordHash string
    isActive     bool
    roles        []Role
}

// Business logic trong entity
func (u *User) Activate() error {
    if u.isActive {
        return ErrUserAlreadyActive
    }
    u.isActive = true
    return nil
}

func (u *User) HasRole(roleName string) bool {
    for _, role := range u.roles {
        if role.Name() == roleName {
            return true
        }
    }
    return false
}

// internal/domain/repository/user_repository.go (Interface/Port)
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    FindByUsername(ctx context.Context, username string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}
```

---

### 4. **Infrastructure Layer**

**Path**: `internal/infrastructure/`

**Responsibility**:
- External dependencies implementation
- Database access
- Third-party integrations
- Framework-specific code

**Components**:
- **Persistence**: Database implementations
  - **DAO**: Raw database operations
  - **Repository Implementations**: Implement domain repository interfaces
- **Casbin**: Authorization engine
- **Security**: JWT, password hashing
- **Config**: Configuration loading

**Example**:
```go
// internal/infrastructure/persistence/postgres/dao/user_dao.go
type UserDAO struct {
    db *sql.DB
}

func (dao *UserDAO) Insert(ctx context.Context, user *UserEntity) error {
    query := `INSERT INTO users ...`
    _, err := dao.db.ExecContext(ctx, query, ...)
    return err
}

// internal/infrastructure/persistence/postgres/repository/user_repository_impl.go
type UserRepositoryImpl struct {
    dao *dao.UserDAO
}

// Implement domain.UserRepository interface
func (r *UserRepositoryImpl) Save(ctx context.Context, user *domain.User) error {
    // Convert domain model to DAO entity
    entity := r.toEntity(user)
    
    // Use DAO to persist
    return r.dao.Insert(ctx, entity)
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*domain.User, error) {
    entity, err := r.dao.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Convert DAO entity to domain model
    return r.toDomainModel(entity), nil
}
```

---

## Dependency Flow

```
┌─────────────────────────────────────────┐
│     Adapter Layer (gRPC Handlers)       │
│   - Convert external requests/responses │
│   - Input validation                    │
└───────────────┬─────────────────────────┘
                │ depends on ↓
┌───────────────▼─────────────────────────┐
│     Application Layer (Use Cases)       │
│   - Orchestrate business logic          │
│   - Transaction management              │
└───────────────┬─────────────────────────┘
                │ depends on ↓
┌───────────────▼─────────────────────────┐
│     Domain Layer (Business Logic)       │
│   - Core business rules                 │
│   - Entities, Value Objects             │
│   - Repository Interfaces (Ports)       │
└───────────────┬─────────────────────────┘
                │ implemented by ↓
┌───────────────▼─────────────────────────┐
│     Infrastructure Layer                │
│   - Database implementations            │
│   - External service integrations       │
└─────────────────────────────────────────┘
```

**Key Principles**:
1. **Dependency Inversion**: Inner layers define interfaces, outer layers implement them
2. **No circular dependencies**: Dependencies only flow inward
3. **Framework independence**: Domain layer has NO framework dependencies
4. **Testability**: Each layer can be tested independently

---

## Benefits of This Structure

### 1. **Clear Separation of Concerns**
- Mỗi layer có responsibility riêng biệt
- Dễ tìm code, dễ maintain

### 2. **Testability**
- Mock dependencies dễ dàng
- Test từng layer độc lập
- Domain logic test không cần database

### 3. **Flexibility**
- Thay đổi database dễ dàng (chỉ sửa infrastructure)
- Thay đổi presentation layer (REST → gRPC)
- Add new adapters (GraphQL, WebSocket)

### 4. **Scalability**
- Add new features theo use case
- Multiple teams work parallel
- Microservices ready

### 5. **Domain-Driven Design Ready**
- Business logic tập trung ở domain layer
- Ubiquitous language
- Rich domain models

---

## Migration Strategy

### Phase 1: Setup New Structure
1. Tạo folder structure mới
2. Move proto files
3. Setup base infrastructure

### Phase 2: Migrate Infrastructure
1. Move DAOs → `infrastructure/persistence/postgres/dao/`
2. Move repositories → `infrastructure/persistence/postgres/repository/`
3. Move Casbin → `infrastructure/casbin/`
4. Move config → `infrastructure/config/`

### Phase 3: Extract Domain
1. Move entities → `domain/model/`
2. Create repository interfaces → `domain/repository/`
3. Extract domain services → `domain/service/`
4. Create value objects → `domain/valueobject/`

### Phase 4: Create Use Cases
1. Extract use cases from current services
2. Create DTOs
3. Wire dependencies

### Phase 5: Refactor Adapters
1. Move handlers → `adapter/grpc/handler/`
2. Create converters → `adapter/converter/`
3. Update main.go with DI

---

## Example: Complete Flow

```
1. gRPC Request
   ↓
2. AuthHandler (Adapter)
   - Validate input
   - Convert to DTO
   ↓
3. RegisterUseCase (Application)
   - Check business rules
   - Create User entity
   - Call domain services
   ↓
4. User (Domain Entity)
   - Business logic
   - Validation
   ↓
5. UserRepository Interface (Domain Port)
   ↓
6. UserRepositoryImpl (Infrastructure)
   - Convert to DAO entity
   - Call UserDAO
   ↓
7. UserDAO (Infrastructure)
   - Execute SQL
   - Return result
   ↓
8. Results flow back up through layers
   ↓
9. gRPC Response
```

---

## Next Steps

1. Review và approve cấu trúc
2. Begin migration process
3. Update documentation
4. Setup testing strategy

