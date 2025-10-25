# Kiến trúc IAM Service

## Tổng quan

IAM Service được xây dựng theo mô hình **Layered Architecture** (Kiến trúc phân lớp) với sự kết hợp của **DAO Pattern** và **Repository Pattern**. Kiến trúc này giúp tách biệt các mối quan tâm (separation of concerns), dễ bảo trì và mở rộng.

## Các Layer trong hệ thống

### 1. Presentation Layer (Handler Layer)

**Vị trí**: `internal/handler/`

**Chức năng**:
- Tiếp nhận request từ gRPC client
- Validate input data
- Gọi business logic layer
- Chuyển đổi response thành Protocol Buffer format
- Xử lý lỗi và trả về status code phù hợp

**Các file chính**:
- `grpc_handler.go`: Implements gRPC service interface
- `converter.go`: Chuyển đổi giữa domain models và Protocol Buffer messages

**Ví dụ**:

```go
func (h *GRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    user, tokenPair, err := h.authService.Login(ctx, req.Username, req.Password)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
    }
    // Convert and return response
}
```

### 2. Business Logic Layer (Service Layer)

**Vị trí**: `internal/service/`

**Chức năng**:
- Chứa toàn bộ business logic
- Orchestrate các operations phức tạp
- Validate business rules
- Quản lý transactions
- Không phụ thuộc vào database implementation

**Các service chính**:

#### AuthService (`auth_service.go`)
- User registration
- Login/Logout
- Token generation và verification
- Password management

#### AuthorizationService (`authorization_service.go`)
- Assign/remove roles
- Check permissions
- Get user roles

#### RoleService (`role_service.go`)
- CRUD operations cho roles
- Manage role-permission relationships

#### PermissionService (`permission_service.go`)
- CRUD operations cho permissions

**Ví dụ**:

```go
func (s *authService) Login(ctx context.Context, username, password string) (*domain.User, *domain.TokenPair, error) {
    // 1. Get user from repository
    user, err := s.userRepo.GetUserByUsername(ctx, username)
    
    // 2. Verify password
    if !s.passwordMgr.CheckPassword(password, user.PasswordHash) {
        return nil, nil, fmt.Errorf("invalid credentials")
    }
    
    // 3. Get user roles
    roles, err := s.authzRepo.GetUserRoles(ctx, user.ID)
    
    // 4. Generate tokens
    accessToken, _ := s.jwtManager.GenerateAccessToken(...)
    
    return user, tokenPair, nil
}
```

### 3. Data Access Layer

Bao gồm 2 sub-layers: Repository và DAO

#### 3.1 Repository Pattern

**Vị trí**: `internal/repository/`

**Chức năng**:
- Cung cấp abstraction cho data access
- Implements business-oriented data operations
- Có thể combine nhiều DAO operations
- Thêm business logic validation

**Các repository chính**:
- `user_repository.go`: User data operations
- `role_repository.go`: Role data operations
- `permission_repository.go`: Permission data operations
- `authorization_repository.go`: Authorization logic (user-role-permission relationships)

**Ví dụ**:

```go
func (r *authorizationRepository) GetUserPermissions(ctx context.Context, userID string) ([]*domain.Permission, error) {
    // 1. Get all roles for user
    roles, err := r.userRoleDAO.GetUserRoles(ctx, userID)
    
    // 2. Collect permissions from all roles
    permissionMap := make(map[string]*domain.Permission)
    for _, role := range roles {
        permissions, _ := r.rolePermissionDAO.GetRolePermissions(ctx, role.ID)
        for _, perm := range permissions {
            permissionMap[perm.ID] = perm
        }
    }
    
    return permissions, nil
}
```

#### 3.2 DAO Pattern (Data Access Object)

**Vị trí**: `internal/dao/`

**Chức năng**:
- Direct database access
- CRUD operations
- Simple queries
- No business logic

**Các DAO chính**:
- `user_dao.go`: User table operations
- `role_dao.go`: Role table operations
- `permission_dao.go`: Permission table operations
- `user_role_dao.go`: User-Role junction table
- `role_permission_dao.go`: Role-Permission junction table

**Ví dụ**:

```go
func (d *userDAO) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
    query := `
        SELECT id, username, email, password_hash, full_name, is_active, created_at, updated_at
        FROM users
        WHERE username = $1
    `
    user := &domain.User{}
    err := d.db.QueryRowContext(ctx, query, username).Scan(...)
    return user, err
}
```

### 4. Domain Layer

**Vị trí**: `internal/domain/`

**Chức năng**:
- Định nghĩa domain entities
- Business objects
- Không có dependencies

**Các entities**:
- `User`: User information
- `Role`: Role definition
- `Permission`: Permission definition
- `UserRole`: User-Role relationship
- `RolePermission`: Role-Permission relationship
- `TokenPair`: JWT token pair

## Mô hình dữ liệu

### Entity Relationship Diagram

```
┌─────────────┐       ┌──────────────┐       ┌──────────────┐
│    Users    │       │  User_Roles  │       │    Roles     │
├─────────────┤       ├──────────────┤       ├──────────────┤
│ id (PK)     │──────<│ user_id (FK) │       │ id (PK)      │
│ username    │       │ role_id (FK) │>──────│ name         │
│ email       │       │ created_at   │       │ description  │
│ password_h  │       └──────────────┘       │ created_at   │
│ full_name   │                              │ updated_at   │
│ is_active   │                              └──────┬───────┘
│ created_at  │                                     │
│ updated_at  │                                     │
└─────────────┘                                     │
                                                    │
                      ┌─────────────────────┐       │
                      │ Role_Permissions    │       │
                      ├─────────────────────┤       │
                      │ role_id (FK)        │<──────┘
                      │ permission_id (FK)  │>──────┐
                      │ created_at          │       │
                      └─────────────────────┘       │
                                                    │
                                                    │
                      ┌─────────────────────┐       │
                      │   Permissions       │       │
                      ├─────────────────────┤       │
                      │ id (PK)             │<──────┘
                      │ name                │
                      │ resource            │
                      │ action              │
                      │ description         │
                      │ created_at          │
                      │ updated_at          │
                      └─────────────────────┘
```

## Flow của một Request

### Ví dụ: Login Flow

```
1. Client gửi LoginRequest qua gRPC
         ↓
2. GRPCHandler.Login() nhận request
         ↓
3. Validate input (username, password)
         ↓
4. Call AuthService.Login(username, password)
         ↓
5. AuthService:
   a. Call UserRepository.GetUserByUsername()
         ↓
   b. UserRepository call UserDAO.FindByUsername()
         ↓
   c. UserDAO query database
         ↓
   d. Return User entity
         ↓
   e. Verify password using PasswordManager
         ↓
   f. Get user roles from AuthorizationRepository
         ↓
   g. Generate JWT tokens using JWTManager
         ↓
   h. Return User and TokenPair
         ↓
6. GRPCHandler convert to Protocol Buffer
         ↓
7. Return LoginResponse to client
```

## Các Pattern được áp dụng

### 1. Layered Architecture
- Tách biệt concerns
- Mỗi layer có trách nhiệm riêng
- Dependencies chỉ đi một chiều (top-down)

### 2. Repository Pattern
- Abstraction layer giữa business logic và data access
- Dễ dàng thay đổi data source
- Testable

### 3. DAO Pattern
- Encapsulate database access
- CRUD operations
- Simple and focused

### 4. Dependency Injection
- Loose coupling
- Testability
- Flexibility

## Ưu điểm của kiến trúc này

1. **Separation of Concerns**: Mỗi layer có trách nhiệm riêng biệt
2. **Maintainability**: Dễ bảo trì và mở rộng
3. **Testability**: Có thể test từng layer độc lập
4. **Reusability**: Code có thể tái sử dụng
5. **Flexibility**: Dễ dàng thay đổi implementation
6. **Scalability**: Dễ dàng scale từng component

## Best Practices

1. **Handler Layer**:
   - Chỉ xử lý HTTP/gRPC concerns
   - Không chứa business logic
   - Validate input cơ bản

2. **Service Layer**:
   - Chứa toàn bộ business logic
   - Orchestrate multiple operations
   - Independent of transport layer

3. **Repository Layer**:
   - Business-oriented queries
   - Combine multiple DAO operations
   - Add validation logic

4. **DAO Layer**:
   - Simple CRUD operations
   - Direct database access
   - No business logic

5. **Domain Layer**:
   - Pure business objects
   - No dependencies
   - Reusable across layers

