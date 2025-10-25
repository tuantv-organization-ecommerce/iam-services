# IAM Service API Documentation

## Tổng quan

IAM Service cung cấp gRPC API cho các chức năng xác thực và phân quyền.

**Server Address**: `localhost:50051` (default)

## Authentication APIs

### 1. Register

Đăng ký người dùng mới.

**Method**: `iam.IAMService/Register`

**Request**:
```protobuf
message RegisterRequest {
  string username = 1;
  string email = 2;
  string password = 3;
  string full_name = 4;
}
```

**Response**:
```protobuf
message RegisterResponse {
  string user_id = 1;
  string username = 2;
  string email = 3;
  string message = 4;
}
```

**Example (grpcurl)**:
```bash
grpcurl -plaintext -d '{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepass123",
  "full_name": "John Doe"
}' localhost:50051 iam.IAMService/Register
```

---

### 2. Login

Đăng nhập và nhận access token.

**Method**: `iam.IAMService/Login`

**Request**:
```protobuf
message LoginRequest {
  string username = 1;
  string password = 2;
}
```

**Response**:
```protobuf
message LoginResponse {
  string access_token = 1;
  string refresh_token = 2;
  string token_type = 3;
  int64 expires_in = 4;
  User user = 5;
}
```

**Example**:
```bash
grpcurl -plaintext -d '{
  "username": "johndoe",
  "password": "securepass123"
}' localhost:50051 iam.IAMService/Login
```

---

### 3. RefreshToken

Làm mới access token sử dụng refresh token.

**Method**: `iam.IAMService/RefreshToken`

**Request**:
```protobuf
message RefreshTokenRequest {
  string refresh_token = 1;
}
```

**Response**:
```protobuf
message RefreshTokenResponse {
  string access_token = 1;
  string refresh_token = 2;
  string token_type = 3;
  int64 expires_in = 4;
}
```

**Example**:
```bash
grpcurl -plaintext -d '{
  "refresh_token": "your-refresh-token"
}' localhost:50051 iam.IAMService/RefreshToken
```

---

### 4. VerifyToken

Xác minh tính hợp lệ của token.

**Method**: `iam.IAMService/VerifyToken`

**Request**:
```protobuf
message VerifyTokenRequest {
  string token = 1;
}
```

**Response**:
```protobuf
message VerifyTokenResponse {
  bool valid = 1;
  string user_id = 2;
  repeated string roles = 3;
  string message = 4;
}
```

**Example**:
```bash
grpcurl -plaintext -d '{
  "token": "your-access-token"
}' localhost:50051 iam.IAMService/VerifyToken
```

---

### 5. Logout

Đăng xuất người dùng.

**Method**: `iam.IAMService/Logout`

**Request**:
```protobuf
message LogoutRequest {
  string user_id = 1;
  string token = 2;
}
```

**Response**:
```protobuf
message LogoutResponse {
  string message = 1;
}
```

---

## Authorization APIs

### 6. AssignRole

Gán vai trò cho người dùng.

**Method**: `iam.IAMService/AssignRole`

**Request**:
```protobuf
message AssignRoleRequest {
  string user_id = 1;
  string role_id = 2;
}
```

**Response**:
```protobuf
message AssignRoleResponse {
  string message = 1;
}
```

**Example**:
```bash
grpcurl -plaintext -d '{
  "user_id": "user-123",
  "role_id": "role-001"
}' localhost:50051 iam.IAMService/AssignRole
```

---

### 7. RemoveRole

Xóa vai trò khỏi người dùng.

**Method**: `iam.IAMService/RemoveRole`

**Request**:
```protobuf
message RemoveRoleRequest {
  string user_id = 1;
  string role_id = 2;
}
```

**Response**:
```protobuf
message RemoveRoleResponse {
  string message = 1;
}
```

---

### 8. GetUserRoles

Lấy danh sách vai trò của người dùng.

**Method**: `iam.IAMService/GetUserRoles`

**Request**:
```protobuf
message GetUserRolesRequest {
  string user_id = 1;
}
```

**Response**:
```protobuf
message GetUserRolesResponse {
  repeated Role roles = 1;
}
```

**Example**:
```bash
grpcurl -plaintext -d '{
  "user_id": "user-123"
}' localhost:50051 iam.IAMService/GetUserRoles
```

---

### 9. CheckPermission

Kiểm tra quyền của người dùng.

**Method**: `iam.IAMService/CheckPermission`

**Request**:
```protobuf
message CheckPermissionRequest {
  string user_id = 1;
  string resource = 2;
  string action = 3;
}
```

**Response**:
```protobuf
message CheckPermissionResponse {
  bool allowed = 1;
  string message = 2;
}
```

**Example**:
```bash
grpcurl -plaintext -d '{
  "user_id": "user-123",
  "resource": "user",
  "action": "create"
}' localhost:50051 iam.IAMService/CheckPermission
```

---

## Role Management APIs

### 10. CreateRole

Tạo vai trò mới.

**Method**: `iam.IAMService/CreateRole`

**Request**:
```protobuf
message CreateRoleRequest {
  string name = 1;
  string description = 2;
  repeated string permission_ids = 3;
}
```

**Response**:
```protobuf
message CreateRoleResponse {
  string role_id = 1;
  string message = 2;
}
```

**Example**:
```bash
grpcurl -plaintext -d '{
  "name": "editor",
  "description": "Content editor role",
  "permission_ids": ["perm-001", "perm-003"]
}' localhost:50051 iam.IAMService/CreateRole
```

---

### 11. UpdateRole

Cập nhật vai trò.

**Method**: `iam.IAMService/UpdateRole`

**Request**:
```protobuf
message UpdateRoleRequest {
  string role_id = 1;
  string name = 2;
  string description = 3;
  repeated string permission_ids = 4;
}
```

**Response**:
```protobuf
message UpdateRoleResponse {
  string message = 1;
}
```

---

### 12. DeleteRole

Xóa vai trò.

**Method**: `iam.IAMService/DeleteRole`

**Request**:
```protobuf
message DeleteRoleRequest {
  string role_id = 1;
}
```

**Response**:
```protobuf
message DeleteRoleResponse {
  string message = 1;
}
```

---

### 13. GetRole

Lấy thông tin vai trò.

**Method**: `iam.IAMService/GetRole`

**Request**:
```protobuf
message GetRoleRequest {
  string role_id = 1;
}
```

**Response**:
```protobuf
message GetRoleResponse {
  Role role = 1;
}
```

---

### 14. ListRoles

Liệt kê tất cả vai trò.

**Method**: `iam.IAMService/ListRoles`

**Request**:
```protobuf
message ListRolesRequest {
  int32 page = 1;
  int32 page_size = 2;
}
```

**Response**:
```protobuf
message ListRolesResponse {
  repeated Role roles = 1;
  int32 total = 2;
}
```

**Example**:
```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:50051 iam.IAMService/ListRoles
```

---

## Permission Management APIs

### 15. CreatePermission

Tạo quyền mới.

**Method**: `iam.IAMService/CreatePermission`

**Request**:
```protobuf
message CreatePermissionRequest {
  string name = 1;
  string resource = 2;
  string action = 3;
  string description = 4;
}
```

**Response**:
```protobuf
message CreatePermissionResponse {
  string permission_id = 1;
  string message = 2;
}
```

**Example**:
```bash
grpcurl -plaintext -d '{
  "name": "Product Create",
  "resource": "product",
  "action": "create",
  "description": "Permission to create products"
}' localhost:50051 iam.IAMService/CreatePermission
```

---

### 16. DeletePermission

Xóa quyền.

**Method**: `iam.IAMService/DeletePermission`

**Request**:
```protobuf
message DeletePermissionRequest {
  string permission_id = 1;
}
```

**Response**:
```protobuf
message DeletePermissionResponse {
  string message = 1;
}
```

---

### 17. ListPermissions

Liệt kê tất cả quyền.

**Method**: `iam.IAMService/ListPermissions`

**Request**:
```protobuf
message ListPermissionsRequest {
  int32 page = 1;
  int32 page_size = 2;
}
```

**Response**:
```protobuf
message ListPermissionsResponse {
  repeated Permission permissions = 1;
  int32 total = 2;
}
```

**Example**:
```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:50051 iam.IAMService/ListPermissions
```

---

## Data Models

### User
```protobuf
message User {
  string id = 1;
  string username = 2;
  string email = 3;
  string full_name = 4;
  bool is_active = 5;
  string created_at = 6;
  string updated_at = 7;
}
```

### Role
```protobuf
message Role {
  string id = 1;
  string name = 2;
  string description = 3;
  repeated Permission permissions = 4;
  string created_at = 5;
  string updated_at = 6;
}
```

### Permission
```protobuf
message Permission {
  string id = 1;
  string name = 2;
  string resource = 3;
  string action = 4;
  string description = 5;
  string created_at = 6;
  string updated_at = 7;
}
```

## Error Codes

| Code | Description |
|------|-------------|
| `OK` | Success |
| `INVALID_ARGUMENT` | Invalid input parameters |
| `UNAUTHENTICATED` | Invalid credentials or token |
| `PERMISSION_DENIED` | User doesn't have required permission |
| `NOT_FOUND` | Resource not found |
| `ALREADY_EXISTS` | Resource already exists |
| `INTERNAL` | Internal server error |

## Testing với grpcurl

### List all services
```bash
grpcurl -plaintext localhost:50051 list
```

### Describe a service
```bash
grpcurl -plaintext localhost:50051 describe iam.IAMService
```

### Describe a method
```bash
grpcurl -plaintext localhost:50051 describe iam.IAMService.Login
```

