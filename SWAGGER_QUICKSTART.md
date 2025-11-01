# Swagger UI Quick Start Guide

## Tổng quan

IAM Service đã được tích hợp Swagger UI để hiển thị API documentation. Swagger UI được bảo vệ bằng HTTP Basic Authentication.

## Quick Start (3 bước)

### 1. Generate Proto Files (Chỉ cần 1 lần hoặc khi proto files thay đổi)

```powershell
cd iam-services
.\scripts\generate-proto-simple.ps1
```

**Output:**
- ✅ `pkg/proto/iam.pb.go` - gRPC message definitions
- ✅ `pkg/proto/iam_grpc.pb.go` - gRPC service definitions
- ✅ `pkg/proto/iam.pb.gw.go` - HTTP gateway handlers
- ✅ `pkg/proto/iam.swagger.json` - OpenAPI specification

### 2. Start Service

```powershell
go run cmd/server/main.go
```

**Logs sẽ hiển thị:**
```
INFO    gRPC server is running  {"address": "0.0.0.0:50051"}
INFO    Swagger UI enabled with Basic Authentication
INFO    HTTP Gateway is running {"address": "0.0.0.0:8080", "swagger": "http://0.0.0.0:8080/swagger/"}
```

### 3. Access Swagger UI

1. Mở browser: **http://localhost:8080/swagger/**
2. Nhập credentials:
   - **Username**: `admin`
   - **Password**: `changeme`
3. Explore & test APIs! 🎉

## Các Endpoints Có Sẵn

### Authentication
- `POST /v1/auth/register` - Đăng ký user mới
- `POST /v1/auth/login` - Đăng nhập
- `POST /v1/auth/refresh` - Refresh access token
- `POST /v1/auth/logout` - Đăng xuất
- `POST /v1/auth/verify` - Verify token

### Authorization
- `POST /v1/roles/assign` - Gán role cho user
- `POST /v1/roles/remove` - Xóa role khỏi user
- `GET /v1/users/{user_id}/roles` - Lấy roles của user
- `POST /v1/permissions/check` - Kiểm tra permission

### Role Management
- `POST /v1/roles` - Tạo role mới
- `PUT /v1/roles/{role_id}` - Cập nhật role
- `DELETE /v1/roles/{role_id}` - Xóa role
- `GET /v1/roles/{role_id}` - Lấy thông tin role
- `GET /v1/roles` - List tất cả roles

### Permission Management
- `POST /v1/permissions` - Tạo permission mới
- `DELETE /v1/permissions/{permission_id}` - Xóa permission
- `GET /v1/permissions` - List tất cả permissions

### Casbin Authorization
- `POST /v1/access/api` - Kiểm tra API access
- `POST /v1/access/cms` - Kiểm tra CMS access
- `POST /v1/policies/enforce` - Enforce policy

### CMS Role Management
- `POST /v1/cms/roles` - Tạo CMS role
- `POST /v1/cms/roles/assign` - Gán CMS role
- `POST /v1/cms/roles/remove` - Xóa CMS role
- `GET /v1/cms/users/{user_id}/tabs` - Lấy CMS tabs của user
- `GET /v1/cms/roles` - List CMS roles

### API Resource Management
- `POST /v1/api/resources` - Tạo API resource
- `GET /v1/api/resources` - List API resources

## Cấu hình

### Environment Variables

Tạo file `.env` trong thư mục `iam-services`:

```bash
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=50051
HTTP_HOST=0.0.0.0
HTTP_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=iam_db
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_EXPIRATION_HOURS=168

# Swagger Configuration
SWAGGER_ENABLED=true
SWAGGER_BASE_PATH=/swagger/
SWAGGER_SPEC_PATH=/swagger.json
SWAGGER_TITLE=IAM Service API Documentation
SWAGGER_AUTH_USERNAME=admin
SWAGGER_AUTH_PASSWORD=changeme
SWAGGER_AUTH_REALM=IAM Service API Documentation
```

### Đổi Swagger Credentials

#### Option 1: Environment Variables
```bash
SWAGGER_AUTH_USERNAME=myusername
SWAGGER_AUTH_PASSWORD=mypassword
```

#### Option 2: Update Config File
Edit `internal/config/config.go`:
```go
AuthUsername: getEnv("SWAGGER_AUTH_USERNAME", "myusername"),
AuthPassword: getEnv("SWAGGER_AUTH_PASSWORD", "mypassword"),
```

### Disable Swagger (Production)

```bash
SWAGGER_ENABLED=false
```

Hoặc trong code:
```go
Swagger: SwaggerConfig{
    Enabled: false,
    // ...
}
```

## Features

### ✅ Interactive API Testing
- Click "Try it out" button
- Fill in request parameters
- Execute API calls
- View response data

### ✅ Request/Response Schemas
- View data structures
- See field types và constraints
- Example values

### ✅ Security
- Basic Authentication protection
- Constant-time password comparison (prevents timing attacks)
- Configurable credentials

### ✅ Auto-Generated Documentation
- Sync tự động với proto files
- Không cần maintain docs manually
- Support cả gRPC và HTTP/REST

## Troubleshooting

### Issue: "protoc-gen-go: not found"

**Solution:**
```powershell
# Add GOPATH/bin to PATH
$env:PATH += ";E:\go\src\bin"  # Adjust based on your GOPATH

# Or install plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.18.1
```

### Issue: Swagger UI không load

**Solution:**
```powershell
# 1. Check if swagger.json exists
ls pkg/proto/iam.swagger.json

# 2. Regenerate proto files
.\scripts\generate-proto-simple.ps1

# 3. Restart service
go run cmd/server/main.go
```

### Issue: 401 Unauthorized loop

**Solution:**
1. Check credentials trong `.env` file
2. Clear browser cache/cookies
3. Try incognito mode
4. Verify service logs for "Swagger UI enabled" message

### Issue: Proto generation fails

**Common causes:**
1. Protoc not installed
2. Protoc plugins not in PATH
3. Missing google proto files

**Solution:**
```powershell
# Check protoc version
protoc --version

# Check plugins
where protoc-gen-go
where protoc-gen-grpc-gateway

# Regenerate with error output
.\scripts\generate-proto-simple.ps1
```

## Development Workflow

### 1. Thêm API mới

**Step 1:** Update proto file `pkg/proto/iam.proto`
```protobuf
// Add new RPC method
rpc MyNewAPI(MyRequest) returns (MyResponse) {
  option (google.api.http) = {
    post: "/v1/my-new-api"
    body: "*"
  };
}

// Add messages
message MyRequest {
  string field1 = 1;
}

message MyResponse {
  string result = 1;
}
```

**Step 2:** Regenerate proto files
```powershell
.\scripts\generate-proto-simple.ps1
```

**Step 3:** Implement handler in `internal/handler/grpc_handler.go`
```go
func (h *GRPCHandler) MyNewAPI(ctx context.Context, req *pb.MyRequest) (*pb.MyResponse, error) {
    // Implementation
    return &pb.MyResponse{Result: "success"}, nil
}
```

**Step 4:** Restart service và test trong Swagger UI

### 2. Update existing API

**Step 1:** Modify proto definition
**Step 2:** Regenerate: `.\scripts\generate-proto-simple.ps1`
**Step 3:** Update implementation
**Step 4:** Test

### 3. Change Swagger config

**Step 1:** Update `.env` file
**Step 2:** Restart service
**Step 3:** Clear browser cache
**Step 4:** Reload Swagger UI

## Best Practices

### 1. Security

✅ **DO:**
- Change default credentials trong production
- Use HTTPS trong production
- Disable Swagger trong production nếu không cần
- Use strong passwords

❌ **DON'T:**
- Commit credentials vào git
- Expose Swagger UI publicly without auth
- Use HTTP trong production

### 2. Development

✅ **DO:**
- Regenerate proto files sau khi update proto
- Test APIs trong Swagger UI trước khi commit
- Document request/response examples trong proto comments
- Keep proto files organized và documented

❌ **DON'T:**
- Manually edit generated files (*.pb.go, *.pb.gw.go)
- Skip proto regeneration
- Ignore proto lint warnings

### 3. Documentation

✅ **DO:**
- Add comments trong proto files
- Use descriptive field names
- Include example values
- Document error cases

❌ **DON'T:**
- Leave APIs undocumented
- Use cryptic field names
- Skip error documentation

## Additional Resources

- [gRPC Gateway Documentation](https://grpc-ecosystem.github.io/grpc-gateway/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)

## Support

Nếu gặp vấn đề:
1. Check `fix_error_ci_cd.md` section "21) Swagger UI Integration"
2. Review logs trong console
3. Verify proto files generated correctly
4. Check environment variables

---

**Happy API Testing! 🚀**

