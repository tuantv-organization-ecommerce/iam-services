# Migration Guide: Clean Architecture

## Overview

This guide explains how to migrate from the old architecture (`main.go`) to the new Clean Architecture (`main_new.go`).

## Architecture Comparison

### Old Architecture (Layered)
```
Handler -> Service -> Repository -> DAO -> Database
```

### New Architecture (Clean/Hexagonal)
```
Adapter (Handler) -> Application (Use Case) -> Domain (Ports) <- Infrastructure (Adapters)
```

## Key Differences

### 1. **Dependency Direction**

**Old:**
- Dependencies point inward: Handler depends on Service, Service depends on Repository
- Concrete implementations everywhere

**New:**
- Dependencies point to domain: All layers depend on domain interfaces (ports)
- Infrastructure implements domain ports
- Domain has NO dependencies on other layers

### 2. **File Structure**

**Old Structure:**
```
internal/
├── handler/          # gRPC handlers
├── service/          # Business logic
├── repository/       # Data access abstraction
└── dao/             # Database operations
```

**New Structure:**
```
internal/
├── domain/
│   ├── model/       # Rich domain models
│   ├── repository/  # Repository interfaces (ports)
│   └── service/     # Domain service interfaces
├── application/
│   ├── dto/         # Data transfer objects
│   └── usecase/     # Application use cases
├── infrastructure/
│   ├── persistence/ # Repository implementations
│   ├── security/    # JWT, Password services
│   ├── authorization/ # Casbin service
│   └── config/      # Config loader
└── adapter/
    └── grpc/        # gRPC handlers (adapters)
```

### 3. **Dependency Injection**

**Old main.go:**
```go
// Create DAOs
userDAO := dao.NewUserDAO(db)

// Create Repositories
userRepo := repository.NewUserRepository(userDAO)

// Create Services
authService := service.NewAuthService(userRepo, jwtManager, passwordMgr)

// Create Handlers
grpcHandler := handler.NewGRPCHandler(authService, logger)
```

**New main_new.go:**
```go
// 1. DAOs (Infrastructure - Database)
userDAO := dao.NewUserDAO(db)

// 2. Infrastructure Layer (Implements Domain Ports)
// 2.1 Persistence
userRepo := persistence.NewUserRepository(userDAO)

// 2.2 Security
tokenService := security.NewJWTService(jwtManager)
passwordService := security.NewPasswordService(passwordMgr)

// 2.3 Authorization
casbinAuthzService := authorization.NewCasbinService(...)

// 3. Application Layer (Use Cases)
registerUseCase := authUseCase.NewRegisterUseCase(userRepo, passwordService)
loginUseCase := authUseCase.NewLoginUseCase(userRepo, tokenService, passwordService)

// 4. Adapter Layer (gRPC Handlers)
grpcHandler := handler.NewGRPCHandler(registerUseCase, loginUseCase, logger)
```

## Migration Steps

### Step 1: Run with Old Architecture (Current)

The current `main.go` is still functional. Use it while migrating:

```bash
go run cmd/server/main.go
```

### Step 2: Test New Architecture

To use the new Clean Architecture:

1. Rename `main.go` to `main_old.go`:
```bash
cd ecommerce/back_end/iam-services
mv cmd/server/main.go cmd/server/main_old.go
```

2. Rename `main_new.go` to `main.go`:
```bash
mv cmd/server/main_new.go cmd/server/main.go
```

3. Update the main function name:
```go
// Change from:
func mainNew() {

// To:
func main() {
```

4. Run the service:
```bash
go run cmd/server/main.go
```

### Step 3: Update Handlers to Use Use Cases

Currently, `main_new.go` still uses old services for compatibility. To fully migrate:

1. **Create new gRPC handler implementations** that use use cases:

```go
// Example: internal/adapter/grpc/auth_handler.go
type AuthHandler struct {
    pb.UnimplementedIAMServiceServer
    registerUseCase *authUseCase.RegisterUseCase
    loginUseCase    *authUseCase.LoginUseCase
    logger          *zap.Logger
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    // Use DTO
    input := &dto.RegisterRequest{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
        FullName: req.FullName,
    }
    
    // Call use case
    output, err := h.registerUseCase.Execute(ctx, input)
    if err != nil {
        return nil, err
    }
    
    // Convert to gRPC response
    return &pb.RegisterResponse{
        UserId:  output.UserID,
        Message: output.Message,
    }, nil
}
```

2. **Update main.go** to use new handlers:

```go
// Remove old services
// authService := service.NewAuthService(...)

// Use new handlers
authHandler := adapter.NewAuthHandler(registerUseCase, loginUseCase, logger)
pb.RegisterIAMServiceServer(grpcServer, authHandler)
```

## Benefits of New Architecture

### 1. **Testability**
- Mock interfaces instead of concrete implementations
- Test business logic without database or external dependencies

```go
// Example test
func TestRegisterUseCase(t *testing.T) {
    // Mock repository
    mockRepo := &mockUserRepository{}
    mockPasswordService := &mockPasswordService{}
    
    // Create use case with mocks
    useCase := authUseCase.NewRegisterUseCase(mockRepo, mockPasswordService)
    
    // Test without real database
    result, err := useCase.Execute(context.Background(), input)
    assert.NoError(t, err)
}
```

### 2. **Flexibility**
- Swap implementations without changing business logic
- Example: Switch from PostgreSQL to MongoDB by changing persistence layer

### 3. **Maintainability**
- Clear separation of concerns
- Domain logic is isolated and pure
- Changes in one layer don't affect others

### 4. **Scalability**
- Easy to add new features
- Use cases are independent and composable

## Example: Adding a New Feature

### Old Architecture
```go
// 1. Add method to service
func (s *AuthService) ResetPassword(...) { }

// 2. Add method to handler
func (h *Handler) ResetPassword(...) { }

// Service depends on Repository, Repository depends on DAO
// Changes ripple through layers
```

### New Architecture
```go
// 1. Define use case (independent)
type ResetPasswordUseCase struct {
    userRepo        repository.UserRepository  // Domain port
    passwordService service.PasswordService    // Domain port
    emailService    service.EmailService       // Domain port
}

func (uc *ResetPasswordUseCase) Execute(ctx context.Context, input *dto.ResetPasswordRequest) error {
    // Pure business logic
    // No dependencies on external concerns
}

// 2. Wire in main.go
resetPasswordUseCase := authUseCase.NewResetPasswordUseCase(
    userRepo,
    passwordService,
    emailService,
)

// 3. Add handler
authHandler := adapter.NewAuthHandler(..., resetPasswordUseCase, ...)
```

## Compatibility

The new architecture is **100% compatible** with the old one during migration:

- Both `main.go` and `main_new.go` can coexist
- Old services still work
- Migrate one feature at a time
- No downtime required

## Rollback Plan

If issues occur:

1. Switch back to old main.go:
```bash
mv cmd/server/main.go cmd/server/main_new.go
mv cmd/server/main_old.go cmd/server/main.go
go run cmd/server/main.go
```

2. Old architecture files are intact:
   - `internal/service/`
   - `internal/repository/`
   - `internal/handler/`

## Next Steps

1. ✅ Infrastructure layer created
2. ✅ main_new.go with DI created
3. ⏳ Create new gRPC handlers using use cases
4. ⏳ Migrate handlers one by one
5. ⏳ Remove old service layer
6. ⏳ Update documentation

## Questions?

Refer to:
- [ARCHITECTURE_NEW.md](./ARCHITECTURE_NEW.md) - New architecture details
- [REFACTORING_GUIDE.md](./REFACTORING_GUIDE.md) - Refactoring process
- [QUICK_START.md](./QUICK_START.md) - Getting started guide

