# IAM Service Refactoring Guide

## 📋 Tổng quan

IAM Service đang được refactor từ kiến trúc layered đơn giản sang **Clean Architecture** với structure rõ ràng hơn, phù hợp cho dự án lớn.

## 🎯 Mục tiêu

1. **Tách biệt concerns** rõ ràng giữa các layer
2. **Domain-driven design** với rich domain models
3. **Testability** - dễ dàng test từng layer
4. **Maintainability** - dễ maintain và mở rộng
5. **Framework independence** - Domain layer không phụ thuộc framework

## 📁 Cấu trúc mới

### Current Structure (Old)
```
internal/
├── dao/              # Data Access Objects
├── database/         # DB connection
├── domain/           # Simple entities
├── handler/          # gRPC handlers
├── repository/       # Repository pattern
└── service/          # Business logic
```

### New Structure (Clean Architecture)
```
internal/
├── adapter/                    # ADAPTER LAYER
│   ├── grpc/
│   │   ├── handler/           # gRPC handlers
│   │   └── interceptor/       # Middleware
│   └── converter/             # Data converters
│
├── application/               # APPLICATION LAYER
│   ├── usecase/              # Use cases
│   │   ├── auth/
│   │   ├── authorization/
│   │   ├── role/
│   │   └── casbin/
│   └── dto/                  # Data Transfer Objects
│
├── domain/                    # DOMAIN LAYER
│   ├── model/                # Rich domain models
│   ├── repository/           # Repository interfaces (Ports)
│   ├── service/              # Domain service interfaces
│   └── valueobject/          # Value objects
│
└── infrastructure/            # INFRASTRUCTURE LAYER
    ├── persistence/
    │   └── postgres/
    │       ├── dao/          # Database access
    │       └── repository/   # Repository implementations
    ├── casbin/               # Casbin integration
    ├── security/             # JWT, password hashing
    └── config/               # Configuration
```

## 🔄 Migration Progress

### ✅ Completed

1. **Domain Layer - Models**
   - ✅ `internal/domain/model/user.go`
   - ✅ `internal/domain/model/role.go`
   - ✅ `internal/domain/model/permission.go`
   - ✅ `internal/domain/model/cms_role.go`
   - ✅ `internal/domain/model/api_resource.go`

2. **Domain Layer - Repository Interfaces (Ports)**
   - ✅ `internal/domain/repository/user_repository.go`
   - ✅ `internal/domain/repository/role_repository.go`
   - ✅ `internal/domain/repository/permission_repository.go`
   - ✅ `internal/domain/repository/authorization_repository.go`
   - ✅ `internal/domain/repository/cms_repository.go`
   - ✅ `internal/domain/repository/api_resource_repository.go`

3. **Domain Layer - Service Interfaces**
   - ✅ `internal/domain/service/password_service.go`
   - ✅ `internal/domain/service/token_service.go`
   - ✅ `internal/domain/service/authorization_service.go`

4. **Application Layer - DTOs**
   - ✅ `internal/application/dto/auth_dto.go`
   - ✅ `internal/application/dto/casbin_dto.go`

5. **Application Layer - Use Cases (Examples)**
   - ✅ `internal/application/usecase/auth/register_usecase.go`
   - ✅ `internal/application/usecase/auth/login_usecase.go`
   - ✅ `internal/application/usecase/casbin/check_api_access_usecase.go`

### 🔄 In Progress

6. **Infrastructure Layer**
   - ⏳ Repository implementations
   - ⏳ DAO adapters
   - ⏳ Casbin adapter
   - ⏳ JWT & Password implementations

7. **Adapter Layer**
   - ⏳ gRPC handlers
   - ⏳ Converters
   - ⏳ Interceptors

### 📝 Pending

8. **Main Application**
   - ⏳ Dependency injection
   - ⏳ Wire all layers together

9. **Documentation**
   - ⏳ Update architecture docs
   - ⏳ Add usage examples

## 🔧 Key Differences

### Old Way (Service Layer)
```go
// internal/service/auth_service.go
type AuthService interface {
    Register(...) (*domain.User, error)
    Login(...) (*domain.User, *domain.TokenPair, error)
}

type authService struct {
    userRepo    repository.UserRepository
    jwtManager  *jwt.JWTManager
    passwordMgr *password.PasswordManager
}

func (s *authService) Register(...) {
    // Mix of business logic and infrastructure concerns
    hashedPassword, _ := s.passwordMgr.HashPassword(password)
    user := &domain.User{...}
    s.userRepo.CreateUser(ctx, user)
}
```

### New Way (Use Case + Domain Model)
```go
// internal/application/usecase/auth/register_usecase.go
type RegisterUseCase struct {
    userRepo    repository.UserRepository  // Domain interface
    passwordSvc service.PasswordService    // Domain service interface
}

func (uc *RegisterUseCase) Execute(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
    // 1. Create rich domain model
    user := model.NewUser(id, req.Username, req.Email, req.FullName)
    
    // 2. Domain validation
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    // 3. Use domain service
    hashedPassword, _ := uc.passwordSvc.Hash(req.Password)
    user.SetPasswordHash(hashedPassword)
    
    // 4. Save through repository interface
    uc.userRepo.Save(ctx, user)
}

// internal/domain/model/user.go (Rich domain model)
type User struct {
    id           string
    username     string
    isActive     bool
}

func (u *User) Activate() error {
    if u.isActive {
        return ErrUserAlreadyActive
    }
    u.isActive = true
    return nil
}

func (u *User) Validate() error {
    if u.username == "" {
        return ErrInvalidUsername
    }
    return nil
}
```

## 📊 Benefits

### 1. Clear Separation
- **Domain** = Business logic only
- **Application** = Use case orchestration
- **Infrastructure** = Technical implementation
- **Adapter** = External interface

### 2. Testability
```go
// Test domain model (no mocks needed)
func TestUser_Activate(t *testing.T) {
    user := model.NewUser("1", "john", "john@test.com", "John")
    user.Deactivate()
    
    err := user.Activate()
    assert.NoError(t, err)
    assert.True(t, user.IsActive())
}

// Test use case (mock repositories)
func TestRegisterUseCase_Execute(t *testing.T) {
    mockRepo := &MockUserRepository{}
    mockPasswordSvc := &MockPasswordService{}
    
    uc := NewRegisterUseCase(mockRepo, mockPasswordSvc)
    // ... test
}
```

### 3. Framework Independence
- Domain layer không import `database/sql`, `gRPC`, `JWT library`
- Dễ dàng thay đổi framework
- Business logic tái sử dụng

### 4. Scalability
- Thêm use case mới dễ dàng
- Multiple teams work parallel
- Microservices ready

## 🚀 Next Steps

### Phase 1: Complete Infrastructure Layer
- [ ] Implement repository adapters
- [ ] Move DAOs to infrastructure
- [ ] Implement domain services (JWT, Password, Casbin)

### Phase 2: Complete Adapter Layer
- [ ] Create gRPC handlers using use cases
- [ ] Create converters (Proto ↔ DTO ↔ Domain)
- [ ] Add interceptors

### Phase 3: Dependency Injection
- [ ] Wire layers in `main.go`
- [ ] Setup proper DI container
- [ ] Remove old code

### Phase 4: Testing & Documentation
- [ ] Add unit tests for each layer
- [ ] Add integration tests
- [ ] Update documentation

## 📚 Examples

### Example: Register Flow

**Old Code**:
```
gRPC Request 
→ Handler 
→ Service (business logic + infra mixed) 
→ Repository 
→ DAO 
→ Database
```

**New Code**:
```
gRPC Request 
→ Handler (Adapter)
   ↓ convert to DTO
→ Use Case (Application)
   ↓ use domain model
→ Domain Model + Domain Services
   ↓ through repository interface
→ Repository Implementation (Infrastructure)
   ↓ use DAO
→ DAO (Infrastructure)
   ↓
→ Database
```

### Example: Check API Access Flow

```go
// 1. gRPC Handler (Adapter)
func (h *AuthHandler) CheckAPIAccess(ctx, req *pb.CheckAPIAccessRequest) {
    // Convert proto to DTO
    dto := converter.ToCheckAPIAccessDTO(req)
    
    // Call use case
    result := h.checkAPIAccessUC.Execute(ctx, dto)
    
    // Convert back to proto
    return converter.ToProtoResponse(result)
}

// 2. Use Case (Application)
func (uc *CheckAPIAccessUseCase) Execute(ctx, dto *dto.CheckAPIAccessRequest) {
    // Validate
    uc.validateRequest(dto)
    
    // Use domain service
    allowed := uc.authzSvc.Enforce(ctx, dto.UserID, domain.API, dto.Path, dto.Method)
    
    return &dto.CheckAPIAccessResponse{Allowed: allowed}
}

// 3. Domain Service Interface (Domain)
type AuthorizationService interface {
    Enforce(ctx, subject, domain, object, action string) (bool, error)
}

// 4. Casbin Implementation (Infrastructure)
type CasbinAuthorizationService struct {
    enforcer *casbin.Enforcer
}

func (s *CasbinAuthorizationService) Enforce(...) (bool, error) {
    return s.enforcer.Enforce(subject, domain, object, action)
}
```

## ⚠️ Important Notes

1. **Keep old code** until refactoring is complete
2. **Test thoroughly** before removing old code
3. **Update incrementally** - one layer at a time
4. **Document as you go** - update docs with changes
5. **Review with team** before major changes

## 🤝 Contributing

Khi refactor:
1. Follow existing patterns in new structure
2. Write unit tests for new code
3. Update this guide with progress
4. Keep old code until fully tested

## 📞 Questions?

See:
- `ARCHITECTURE_NEW.md` - Detailed architecture explanation
- `docs/ARCHITECTURE.md` - Original architecture
- Domain model files - Rich examples

