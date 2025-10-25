# gRPC Gateway - REST API Documentation

## Tổng Quan

IAM Service hỗ trợ cả gRPC và REST API thông qua **gRPC Gateway**. Bạn có thể sử dụng REST API trực tiếp mà không cần client gRPC.

## Cấu Hình

### Environment Variables

```bash
# gRPC Server (mặc định)
SERVER_HOST=0.0.0.0
SERVER_PORT=50051

# HTTP Gateway (mặc định)
HTTP_HOST=0.0.0.0
HTTP_PORT=8080
```

### File .env

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=50051
HTTP_HOST=0.0.0.0
HTTP_PORT=8080
```

## Khởi Động Service

### 1. Generate Proto Files

Chạy script để generate cả gRPC và Gateway code:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\setup-proto.ps1
```

Script này sẽ:
- Cài đặt protoc-gen-go, protoc-gen-go-grpc
- Cài đặt protoc-gen-grpc-gateway, protoc-gen-openapiv2
- Tải Google API proto files
- Generate gRPC code
- Generate Gateway code
- Generate OpenAPI/Swagger documentation

### 2. Run Service

```bash
go run cmd/server/main.go
```

Service sẽ khởi động 2 servers:
- **gRPC Server**: `0.0.0.0:50051`
- **HTTP Gateway**: `0.0.0.0:8080`

## REST API Endpoints

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

#### 1. Register User

**Endpoint:** `POST /api/v1/auth/register`

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe"
}
```

**Response:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john_doe",
  "email": "john@example.com",
  "message": "User registered successfully"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'
```

#### 2. Login

**Endpoint:** `POST /api/v1/auth/login`

**Request Body:**
```json
{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "is_active": true
  }
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "SecurePass123!"
  }'
```

#### 3. Refresh Token

**Endpoint:** `POST /api/v1/auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

#### 4. Verify Token

**Endpoint:** `POST /api/v1/auth/verify`

**Request Body:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:**
```json
{
  "valid": true,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "roles": ["admin", "user"],
  "message": "Token is valid"
}
```

#### 5. Logout

**Endpoint:** `POST /api/v1/auth/logout`

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Authorization

#### 1. Assign Role to User

**Endpoint:** `POST /api/v1/authorization/assign-role`

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "role_id": "role-admin-id"
}
```

#### 2. Remove Role from User

**Endpoint:** `DELETE /api/v1/authorization/users/{user_id}/roles/{role_id}`

**cURL Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/authorization/users/550e8400-e29b-41d4-a716-446655440000/roles/role-admin-id
```

#### 3. Get User Roles

**Endpoint:** `GET /api/v1/authorization/users/{user_id}/roles`

**cURL Example:**
```bash
curl http://localhost:8080/api/v1/authorization/users/550e8400-e29b-41d4-a716-446655440000/roles
```

#### 4. Check Permission

**Endpoint:** `POST /api/v1/authorization/check-permission`

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "resource": "users",
  "action": "read"
}
```

### Role Management

#### 1. Create Role

**Endpoint:** `POST /api/v1/roles`

**Request Body:**
```json
{
  "name": "editor",
  "description": "Editor role with limited permissions",
  "permission_ids": ["perm-1", "perm-2"]
}
```

#### 2. Update Role

**Endpoint:** `PUT /api/v1/roles/{role_id}`

**Request Body:**
```json
{
  "role_id": "role-editor-id",
  "name": "senior_editor",
  "description": "Senior editor with more permissions",
  "permission_ids": ["perm-1", "perm-2", "perm-3"]
}
```

#### 3. Get Role

**Endpoint:** `GET /api/v1/roles/{role_id}`

**cURL Example:**
```bash
curl http://localhost:8080/api/v1/roles/role-editor-id
```

#### 4. List Roles

**Endpoint:** `GET /api/v1/roles?page=1&page_size=10`

**cURL Example:**
```bash
curl "http://localhost:8080/api/v1/roles?page=1&page_size=10"
```

#### 5. Delete Role

**Endpoint:** `DELETE /api/v1/roles/{role_id}`

### Permission Management

#### 1. Create Permission

**Endpoint:** `POST /api/v1/permissions`

**Request Body:**
```json
{
  "name": "read_users",
  "resource": "users",
  "action": "read",
  "description": "Permission to read user data"
}
```

#### 2. List Permissions

**Endpoint:** `GET /api/v1/permissions?page=1&page_size=10`

#### 3. Delete Permission

**Endpoint:** `DELETE /api/v1/permissions/{permission_id}`

### Casbin Authorization

#### 1. Check API Access

**Endpoint:** `POST /api/v1/casbin/check-api-access`

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "api_path": "/api/v1/users",
  "method": "GET"
}
```

**Response:**
```json
{
  "allowed": true,
  "message": "Access granted"
}
```

#### 2. Check CMS Access

**Endpoint:** `POST /api/v1/casbin/check-cms-access`

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "cms_tab": "product",
  "action": "read"
}
```

**Response:**
```json
{
  "allowed": true,
  "message": "Access granted to CMS tab",
  "accessible_tabs": ["product", "inventory", "report"]
}
```

#### 3. Enforce Policy

**Endpoint:** `POST /api/v1/casbin/enforce`

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "domain": "cms",
  "resource": "product",
  "action": "write"
}
```

### CMS Role Management

#### 1. Create CMS Role

**Endpoint:** `POST /api/v1/cms/roles`

**Request Body:**
```json
{
  "name": "product_manager",
  "description": "Product management role",
  "tabs": ["product", "inventory"]
}
```

#### 2. Assign CMS Role

**Endpoint:** `POST /api/v1/cms/assign-role`

**Request Body:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "cms_role_id": "cms-role-id"
}
```

#### 3. Get User CMS Tabs

**Endpoint:** `GET /api/v1/cms/users/{user_id}/tabs`

**Response:**
```json
{
  "tabs": ["product", "inventory", "report"],
  "roles": [
    {
      "id": "cms-role-id",
      "name": "product_manager",
      "description": "Product management role",
      "tabs": ["product", "inventory"]
    }
  ]
}
```

#### 4. List CMS Roles

**Endpoint:** `GET /api/v1/cms/roles?page=1&page_size=10`

#### 5. Remove CMS Role

**Endpoint:** `DELETE /api/v1/cms/users/{user_id}/roles/{cms_role_id}`

### API Resource Management

#### 1. Create API Resource

**Endpoint:** `POST /api/v1/resources`

**Request Body:**
```json
{
  "path": "/api/v1/users",
  "method": "GET",
  "service": "iam-service",
  "description": "Get all users"
}
```

#### 2. List API Resources

**Endpoint:** `GET /api/v1/resources?service=iam-service&page=1&page_size=10`

## CORS Support

HTTP Gateway hỗ trợ CORS với các headers sau:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Accept, Authorization, Content-Type, X-CSRF-Token
```

## OpenAPI/Swagger Documentation

Sau khi generate proto files, bạn sẽ có file OpenAPI specification:

```
pkg/proto/iam_gateway.swagger.json
```

Bạn có thể sử dụng file này với Swagger UI hoặc các công cụ tương tự để xem API documentation tương tác.

## Testing với Postman

1. Import OpenAPI file vào Postman
2. Set base URL: `http://localhost:8080`
3. Test các endpoints

## Testing với cURL

### Example: Complete Flow

```bash
# 1. Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test123!",
    "full_name": "Test User"
  }'

# 2. Login
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123!"
  }' | jq -r '.access_token')

# 3. Verify Token
curl -X POST http://localhost:8080/api/v1/auth/verify \
  -H "Content-Type: application/json" \
  -d "{\"token\": \"$TOKEN\"}"

# 4. List Roles
curl http://localhost:8080/api/v1/roles
```

## So Sánh gRPC vs REST

| Feature | gRPC | REST |
|---------|------|------|
| Port | 50051 | 8080 |
| Protocol | HTTP/2 | HTTP/1.1 |
| Format | Protobuf (binary) | JSON |
| Performance | Cao hơn | Thấp hơn |
| Client | Cần gRPC client | Bất kỳ HTTP client |
| Browser Support | Giới hạn | Đầy đủ |
| Debugging | Khó hơn | Dễ hơn |

## Best Practices

1. **Authentication**: Sử dụng JWT tokens trong header `Authorization: Bearer <token>`
2. **Error Handling**: REST API trả về HTTP status codes chuẩn
3. **Pagination**: Sử dụng `page` và `page_size` query parameters
4. **Filtering**: Sử dụng query parameters cho filtering
5. **HTTPS**: Trong production, luôn sử dụng HTTPS

## Troubleshooting

### Gateway không khởi động

```bash
# Check if gRPC server is running
netstat -an | findstr :50051

# Check if HTTP gateway is running
netstat -an | findstr :8080
```

### Proto generation fails

```bash
# Re-run setup script
powershell -ExecutionPolicy Bypass -File .\scripts\setup-proto.ps1

# Or manually generate
protoc -I. -I./third_party --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative pkg/proto/iam_gateway.proto
```

### CORS errors

Nếu gặp CORS errors từ browser, kiểm tra `corsMiddleware` trong `main.go` đã được apply đúng.

## Next Steps

- [ ] Add authentication middleware cho REST API
- [ ] Add rate limiting
- [ ] Add request logging
- [ ] Add metrics endpoint
- [ ] Add health check endpoint
- [ ] Deploy với TLS/SSL

