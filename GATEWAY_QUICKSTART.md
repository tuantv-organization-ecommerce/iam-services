# 🚀 gRPC Gateway Quick Start

## 📋 Tóm Tắt

IAM Service hiện hỗ trợ **cả gRPC và REST API** thông qua gRPC Gateway!

- **gRPC Server**: `localhost:50051`
- **REST API**: `http://localhost:8080`

## ⚡ Quick Start

### 1. Setup & Generate Code

```powershell
# Chạy script tự động (khuyến nghị)
powershell -ExecutionPolicy Bypass -File .\scripts\setup-proto.ps1
```

Script sẽ:
- ✅ Cài đặt protoc plugins (go, grpc, gateway, openapi)
- ✅ Tải Google API proto files
- ✅ Generate gRPC code
- ✅ Generate Gateway code (REST)
- ✅ Generate OpenAPI/Swagger docs

### 2. Run Service

```bash
go run cmd/server/main.go
```

Output:
```
INFO  gRPC server is running  address=0.0.0.0:50051
INFO  HTTP gateway server is running  address=0.0.0.0:8080
```

### 3. Test REST API

```bash
# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "Pass123!",
    "full_name": "John Doe"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "Pass123!"
  }'
```

## 📚 REST API Endpoints

### Authentication
```
POST   /api/v1/auth/register     - Đăng ký user
POST   /api/v1/auth/login        - Đăng nhập
POST   /api/v1/auth/refresh      - Refresh token
POST   /api/v1/auth/logout       - Đăng xuất
POST   /api/v1/auth/verify       - Verify token
```

### Roles
```
POST   /api/v1/roles             - Tạo role
GET    /api/v1/roles             - List roles
GET    /api/v1/roles/{id}        - Get role
PUT    /api/v1/roles/{id}        - Update role
DELETE /api/v1/roles/{id}        - Delete role
```

### Permissions
```
POST   /api/v1/permissions       - Tạo permission
GET    /api/v1/permissions       - List permissions
DELETE /api/v1/permissions/{id}  - Delete permission
```

### Casbin (RBAC)
```
POST   /api/v1/casbin/check-api-access  - Check API access
POST   /api/v1/casbin/check-cms-access  - Check CMS access
POST   /api/v1/casbin/enforce           - Enforce policy
```

### CMS Roles
```
POST   /api/v1/cms/roles                   - Tạo CMS role
GET    /api/v1/cms/roles                   - List CMS roles
POST   /api/v1/cms/assign-role             - Assign role
GET    /api/v1/cms/users/{id}/tabs         - Get user tabs
DELETE /api/v1/cms/users/{id}/roles/{rid}  - Remove role
```

### API Resources
```
POST   /api/v1/resources         - Tạo API resource
GET    /api/v1/resources         - List resources
```

## 🔧 Configuration

### Environment Variables

```env
# gRPC Server
SERVER_HOST=0.0.0.0
SERVER_PORT=50051

# HTTP Gateway
HTTP_HOST=0.0.0.0
HTTP_PORT=8080
```

### File Structure

```
pkg/proto/
├── iam.proto                    # gRPC definitions
├── iam_gateway.proto            # Gateway annotations
├── iam.pb.go                    # Generated gRPC code
├── iam_grpc.pb.go              # Generated gRPC server
├── iam_gateway.pb.go           # Generated gateway code
├── iam_gateway.pb.gw.go        # Generated gateway handler
└── iam_gateway.swagger.json   # OpenAPI spec
```

## 📖 Documentation

Chi tiết đầy đủ xem tại:
- [GRPC_GATEWAY.md](docs/GRPC_GATEWAY.md) - REST API documentation đầy đủ
- [pkg/proto/iam_gateway.swagger.json](pkg/proto/iam_gateway.swagger.json) - OpenAPI spec

## 🎯 Features

✅ **Dual Protocol**: gRPC + REST API cùng lúc
✅ **Auto-generated**: REST API tự động sinh từ proto
✅ **CORS Support**: Cross-origin requests ready
✅ **OpenAPI/Swagger**: API documentation chuẩn
✅ **Type Safety**: Share definitions giữa gRPC và REST
✅ **Graceful Shutdown**: Cả 2 servers đều graceful

## 🔍 Testing Tools

### cURL
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"pass"}'
```

### Postman
1. Import OpenAPI: `pkg/proto/iam_gateway.swagger.json`
2. Set base URL: `http://localhost:8080`

### Browser
```
http://localhost:8080/api/v1/roles
```

## 🆚 gRPC vs REST

| Feature | gRPC (port 50051) | REST (port 8080) |
|---------|-------------------|------------------|
| **Protocol** | HTTP/2 + Protobuf | HTTP/1.1 + JSON |
| **Performance** | ⚡ Nhanh hơn | 🐌 Chậm hơn |
| **Browser** | ⚠️ Giới hạn | ✅ Full support |
| **Client** | 🔧 Cần gRPC lib | ✅ Bất kỳ HTTP client |
| **Debugging** | 🔍 Khó | 🔎 Dễ |
| **Use Case** | Service-to-service | Browser, testing |

## 💡 Tips

1. **Development**: Dùng REST API cho dễ test
2. **Production**: Dùng gRPC cho performance
3. **Frontend**: Dùng REST API cho browser
4. **Microservices**: Dùng gRPC cho inter-service

## 🐛 Troubleshooting

### Proto generation fails
```bash
# Re-run setup script
powershell -ExecutionPolicy Bypass -File .\scripts\setup-proto.ps1
```

### Port already in use
```bash
# Change ports in .env
HTTP_PORT=8081
SERVER_PORT=50052
```

### Import errors after generation
```bash
go mod tidy
```

## 📞 Support

Gặp vấn đề? Check:
1. [GRPC_GATEWAY.md](docs/GRPC_GATEWAY.md) - Full documentation
2. [ARCHITECTURE_NEW.md](ARCHITECTURE_NEW.md) - System architecture
3. [README.md](README.md) - Main documentation

Happy coding! 🎉

