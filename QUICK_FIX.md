# Quick Fix - Generate Proto Files

## ❌ Current Error

```
RegisterIAMServiceHandlerFromEndpoint not declared by package proto
```

## ✅ Solution

Bạn cần generate proto files với gRPC Gateway support.

### Step 1: Cài đặt tools (nếu chưa có)

```powershell
# Trong PowerShell, chạy:
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.18.1
```

**Note**: Đảm bảo `GOPATH/bin` đã có trong PATH:
```powershell
# Check GOPATH
go env GOPATH

# Add to PATH nếu chưa có:
# Go to: System Properties → Environment Variables → PATH
# Add: C:\Users\<your-user>\go\bin
```

### Step 2: Generate Proto Files

```powershell
# Từ thư mục iam-services
cd ecommerce\back_end\iam-services

# Chạy script
powershell -ExecutionPolicy Bypass -File .\scripts\generate-proto.ps1
```

### Step 3: Uncomment HTTP Gateway Code

Sau khi generate xong, mở file `internal/app/app.go` và:

**1. Uncomment imports (line 15-17):**
```go
// Uncomment after generating proto files
"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
"google.golang.org/grpc/credentials/insecure"
"github.com/tvttt/gokits/swagger"
```

Thành:
```go
"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
"google.golang.org/grpc/credentials/insecure"
"github.com/tvttt/gokits/swagger"
```

**2. Uncomment setupHTTPGateway call (line 134-140):**
```go
// TODO: Uncomment after generating proto files
/*
if err := a.setupHTTPGateway(); err != nil {
    return fmt.Errorf("failed to setup HTTP gateway: %w", err)
}
*/
```

Thành:
```go
if err := a.setupHTTPGateway(); err != nil {
    return fmt.Errorf("failed to setup HTTP gateway: %w", err)
}
```

**3. Uncomment functions (line 152-234):**
Remove `/*` at line 152 and `*/` at line 234.

### Step 4: Run go mod tidy

```bash
go mod tidy
```

### Step 5: Run Service

```bash
go run cmd/server/main.go
```

---

## 🎯 Expected Result

Service sẽ chạy với:
- **gRPC Server**: `localhost:50051`
- **HTTP Gateway**: `http://localhost:8080`
- **Swagger UI**: `http://localhost:8080/swagger/`

---

## 🔧 Alternative: Run gRPC Only (Skip HTTP Gateway)

Nếu bạn chỉ muốn chạy gRPC server (không cần REST API và Swagger UI), service hiện tại đã OK!

Just run:
```bash
go run cmd/server/main.go
```

Service sẽ chỉ expose gRPC trên port 50051.

Test với grpcurl:
```bash
grpcurl -plaintext -d '{
  "username": "admin",
  "password": "admin123"
}' localhost:50051 iam.IAMService/Login
```

---

## 📁 Generated Files

Sau khi chạy script, bạn sẽ có:

```
pkg/proto/
├── iam.pb.go                      # gRPC message types
├── iam_grpc.pb.go                # gRPC service
├── iam_gateway.pb.go             # Gateway message types
├── iam_gateway.pb.gw.go          # Gateway handlers ✅
└── iam_gateway.swagger.json      # OpenAPI spec for Swagger UI
```

File `iam_gateway.pb.gw.go` chứa function `RegisterIAMServiceHandlerFromEndpoint` mà service đang cần.

---

## ❓ Troubleshooting

### Issue 1: protoc not found

```
'protoc' is not recognized as an internal or external command
```

**Solution**: Download and install protoc:
- Go to: https://github.com/protocolbuffers/protobuf/releases
- Download: `protoc-<version>-win64.zip`
- Extract to `C:\protoc`
- Add to PATH: `C:\protoc\bin`

### Issue 2: protoc-gen-go not found

```
'protoc-gen-go' is not recognized as an internal or external command
```

**Solution**: 
```powershell
# Check GOPATH/bin
go env GOPATH
# Output: C:\Users\<you>\go

# Add to PATH if not there
# System Properties → Environment Variables → PATH
# Add: C:\Users\<you>\go\bin

# Re-install tools
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
```

### Issue 3: Import errors after generation

```
cannot find module providing package google.golang.org/genproto/googleapis/api
```

**Solution**:
```bash
go mod tidy
```

---

## ✅ Current Status

Service **CAN RUN** với gRPC only:
- ✅ gRPC Server works on port 50051  
- ⏳ HTTP Gateway (pending proto generation)  
- ⏳ Swagger UI (pending proto generation)

Sau khi generate proto files → Uncomment code → Full features! 🚀

