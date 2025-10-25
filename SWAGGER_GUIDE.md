# Swagger UI - IAM Service API Documentation

## 🎯 Overview

IAM Service includes **Swagger UI** for interactive API documentation and testing. Swagger is automatically generated from gRPC proto files and served alongside the REST API.

### Features

✅ **Interactive Documentation** - Browse all endpoints  
✅ **Try it Out** - Test APIs directly from browser  
✅ **Request/Response Examples** - See sample data  
✅ **Code Generation** - Generate client code snippets  
✅ **Schema Definitions** - View all data models  
✅ **Auto-generated** - Always in sync with proto files  

---

## 🚀 Quick Start

### 1. Start the Service

```bash
cd ecommerce/back_end/iam-services
go run cmd/server/main.go
```

### 2. Access Swagger UI

Open in browser:
```
http://localhost:8080/swagger/
```

### 3. View OpenAPI Spec

Direct access to spec:
```
http://localhost:8080/swagger.json
```

---

## 📍 URLs

| Service | URL | Description |
|---------|-----|-------------|
| **gRPC Server** | `localhost:50051` | gRPC endpoint |
| **REST API** | `http://localhost:8080` | HTTP Gateway |
| **Swagger UI** | `http://localhost:8080/swagger/` | Interactive docs |
| **OpenAPI Spec** | `http://localhost:8080/swagger.json` | OpenAPI/Swagger specification |

---

## ⚙️ Configuration

### Environment Variables

Add to `.env` file:

```env
# Swagger Configuration
SWAGGER_ENABLED=true
SWAGGER_BASE_PATH=/swagger/
SWAGGER_SPEC_PATH=/swagger.json
SWAGGER_TITLE=IAM Service API Documentation
```

### Disable Swagger (Production)

```env
SWAGGER_ENABLED=false
```

---

## 🔧 Generate OpenAPI Spec

The OpenAPI/Swagger spec is auto-generated from proto files.

### Prerequisites

```bash
# Install protoc-gen-openapiv2
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

### Generate Spec

Using the setup script:

**Windows:**
```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\setup-proto.ps1
```

**Linux/MacOS:**
```bash
./scripts/setup.sh
```

Or manually:
```bash
protoc -I. -I./third_party \
  --openapiv2_out=. \
  --openapiv2_opt=logtostderr=true \
  --openapiv2_opt=allow_merge=true \
  --openapiv2_opt=merge_file_name=iam_gateway \
  pkg/proto/iam_gateway.proto
```

This generates:
- `pkg/proto/iam_gateway.swagger.json` - OpenAPI specification

---

## 📖 Using Swagger UI

### 1. Browse Endpoints

- Swagger UI lists all available endpoints grouped by tags
- Click on an endpoint to expand details

### 2. View Request/Response

Each endpoint shows:
- HTTP method (GET, POST, PUT, DELETE)
- Request parameters
- Request body schema
- Response schema
- Response codes

### 3. Try It Out

**Step 1**: Click "Try it out" button

**Step 2**: Fill in parameters
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Step 3**: Click "Execute"

**Step 4**: View response
```json
{
  "access_token": "eyJhbGciOiJI...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

### 4. Authentication

For protected endpoints:

**Step 1**: Login to get token

**Step 2**: Click "Authorize" button (top right)

**Step 3**: Enter token:
```
Bearer eyJhbGciOiJIUzI1NiIs...
```

**Step 4**: Click "Authorize"

**Step 5**: Now you can call protected endpoints

---

## 🔒 Security

### Production Deployment

#### 1. Disable Swagger in Production

```env
SWAGGER_ENABLED=false
```

Or in code:
```go
cfg.Swagger.Enabled = os.Getenv("ENV") != "production"
```

#### 2. Add Authentication

Protect Swagger behind auth:

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check basic auth or other auth mechanism
        username, password, ok := r.BasicAuth()
        if !ok || username != "admin" || password != "secret" {
            w.Header().Set("WWW-Authenticate", `Basic realm="Swagger"`)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Apply to Swagger
mux.Handle("/swagger/", authMiddleware(swagger.Handler(cfg, logger)))
```

#### 3. IP Whitelist

```go
func ipWhitelistMiddleware(allowedIPs []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr
            if !contains(allowedIPs, ip) {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 📚 Example Workflows

### Workflow 1: User Registration & Login

**1. View Registration Endpoint**
- Go to Swagger UI
- Find `POST /api/v1/auth/register`
- Click to expand

**2. Try Registration**
- Click "Try it out"
- Fill in request body:
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe"
}
```
- Click "Execute"

**3. Login**
- Find `POST /api/v1/auth/login`
- Fill in credentials:
```json
{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```
- Execute
- Copy `access_token` from response

**4. Use Token**
- Click "Authorize" button
- Enter: `Bearer <access_token>`
- Click "Authorize"

**5. Test Protected Endpoint**
- Find `GET /api/v1/authorization/users/{user_id}/roles`
- Enter user_id
- Execute
- Should work now with token

### Workflow 2: Role Management

**1. List All Roles**
- `GET /api/v1/roles`
- Set page=1, page_size=10
- Execute

**2. Create New Role**
- `POST /api/v1/roles`
- Request body:
```json
{
  "name": "content_editor",
  "description": "Content editor role",
  "permission_ids": ["perm-001", "perm-002"]
}
```
- Execute

**3. Assign Role to User**
- `POST /api/v1/authorization/assign-role`
- Request body:
```json
{
  "user_id": "user-123",
  "role_id": "<role_id_from_step_2>"
}
```

### Workflow 3: CMS Authorization

**1. Check CMS Access**
- `POST /api/v1/casbin/check-cms-access`
- Request:
```json
{
  "user_id": "user-123",
  "cms_tab": "product",
  "action": "POST"
}
```
- Execute
- See if access is granted

**2. Get User CMS Tabs**
- `GET /api/v1/cms/users/{user_id}/tabs`
- Enter user_id
- Execute
- See all accessible tabs

**3. Create CMS Role**
- `POST /api/v1/cms/roles`
- Request:
```json
{
  "name": "cms_product_manager",
  "description": "Product management role",
  "tabs": ["product", "inventory"]
}
```

---

## 🛠️ Troubleshooting

### Issue 1: Swagger UI Not Loading

**Symptoms**: Blank page or 404

**Solutions**:
1. Check if Swagger is enabled:
   ```env
   SWAGGER_ENABLED=true
   ```

2. Verify service is running:
   ```bash
   curl http://localhost:8080/swagger.json
   ```

3. Check logs for errors

### Issue 2: OpenAPI Spec Not Found

**Symptoms**: "Failed to load spec"

**Solutions**:
1. Verify spec file exists:
   ```bash
   ls -la pkg/proto/iam_gateway.swagger.json
   ```

2. Regenerate spec:
   ```bash
   powershell -ExecutionPolicy Bypass -File .\scripts\setup-proto.ps1
   ```

3. Check file path in config:
   ```go
   swagger.ServeSpec("./pkg/proto/iam_gateway.swagger.json", logger)
   ```

### Issue 3: "Try it out" Returns CORS Error

**Symptoms**: CORS error in browser console

**Solutions**:
1. CORS is already enabled in `app.go`:
   ```go
   w.Header().Set("Access-Control-Allow-Origin", "*")
   ```

2. If still failing, check browser network tab

3. Try from Postman instead (bypasses CORS)

### Issue 4: Authentication Not Working

**Symptoms**: 401 Unauthorized for protected endpoints

**Solutions**:
1. Get fresh token via Login endpoint

2. Click "Authorize" button in Swagger

3. Enter token with "Bearer " prefix:
   ```
   Bearer eyJhbGciOiJIUzI1NiIs...
   ```

4. Verify token is not expired

---

## 🎨 Customization

### Change Swagger Theme

Edit `gokits/swagger/ui/index.html`:

```html
<style>
    .swagger-ui .topbar {
        background-color: #2c3e50 !important; /* Custom color */
    }
</style>
```

### Custom Title

In `.env`:
```env
SWAGGER_TITLE=My Custom API Documentation
```

### Change Base Path

```env
SWAGGER_BASE_PATH=/api-docs/
```

Then access at:
```
http://localhost:8080/api-docs/
```

---

## 📝 Best Practices

### 1. Keep Proto Files Updated

Always regenerate Swagger when proto changes:
```bash
powershell -ExecutionPolicy Bypass -File .\scripts\setup-proto.ps1
```

### 2. Add Descriptions to Proto

```protobuf
message LoginRequest {
  // Username for authentication
  string username = 1;
  
  // User's password (min 8 characters)
  string password = 2;
}
```

### 3. Document Response Codes

```protobuf
rpc Login(LoginRequest) returns (LoginResponse) {
  option (google.api.http) = {
    post: "/api/v1/auth/login"
    body: "*"
  };
  // Success: 200
  // Bad Request: 400
  // Unauthorized: 401
  // Internal Error: 500
}
```

### 4. Test in Development, Disable in Production

```go
cfg.Swagger.Enabled = os.Getenv("ENV") == "development"
```

### 5. Version Your API

```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  info: {
    title: "IAM Service API";
    version: "1.0";
  };
};
```

---

## 🔗 Integration with Other Tools

### Postman

**Import OpenAPI spec:**
1. Open Postman
2. File → Import
3. URL: `http://localhost:8080/swagger.json`
4. Click Import

### Insomnia

**Import OpenAPI spec:**
1. Open Insomnia
2. Application → Import
3. URL: `http://localhost:8080/swagger.json`

### Code Generation

Generate client SDKs:

**Go Client:**
```bash
swagger-codegen generate -i http://localhost:8080/swagger.json -l go -o ./client/go
```

**Python Client:**
```bash
swagger-codegen generate -i http://localhost:8080/swagger.json -l python -o ./client/python
```

**JavaScript Client:**
```bash
swagger-codegen generate -i http://localhost:8080/swagger.json -l javascript -o ./client/js
```

---

## 📚 References

- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [gRPC-Gateway OpenAPI](https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_openapi_output/)
- [GoKits Swagger Package](../gokits/swagger/README.md)

---

**Last Updated**: 2024-01-XX  
**Version**: 1.0.0

