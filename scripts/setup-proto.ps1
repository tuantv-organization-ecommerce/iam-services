# Setup Proto Generation for Go 1.19
# Run this script in PowerShell

Write-Host "========================================" -ForegroundColor Green
Write-Host "Proto Setup Script for IAM Service" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""

# Check Go version
Write-Host "Checking Go version..." -ForegroundColor Yellow
go version
Write-Host ""

# Install protoc-gen-go (compatible with Go 1.19)
Write-Host "Installing protoc-gen-go v1.28.1..." -ForegroundColor Yellow
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
if ($LASTEXITCODE -eq 0) {
    Write-Host "Success: protoc-gen-go installed successfully" -ForegroundColor Green
} else {
    Write-Host "Error: Failed to install protoc-gen-go" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Install protoc-gen-go-grpc (compatible with Go 1.19)
Write-Host "Installing protoc-gen-go-grpc v1.2.0..." -ForegroundColor Yellow
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
if ($LASTEXITCODE -eq 0) {
    Write-Host "Success: protoc-gen-go-grpc installed successfully" -ForegroundColor Green
} else {
    Write-Host "Error: Failed to install protoc-gen-go-grpc" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Install protoc-gen-grpc-gateway (for REST API Gateway)
Write-Host "Installing protoc-gen-grpc-gateway v2.18.1..." -ForegroundColor Yellow
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1
if ($LASTEXITCODE -eq 0) {
    Write-Host "Success: protoc-gen-grpc-gateway installed successfully" -ForegroundColor Green
} else {
    Write-Host "Error: Failed to install protoc-gen-grpc-gateway" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Install protoc-gen-openapiv2 (for OpenAPI/Swagger docs)
Write-Host "Installing protoc-gen-openapiv2 v2.18.1..." -ForegroundColor Yellow
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.18.1
if ($LASTEXITCODE -eq 0) {
    Write-Host "Success: protoc-gen-openapiv2 installed successfully" -ForegroundColor Green
} else {
    Write-Host "Error: Failed to install protoc-gen-openapiv2" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Check if GOPATH/bin is in PATH
$gopath = go env GOPATH
$gobinPath = Join-Path $gopath "bin"
Write-Host "Checking if GOPATH/bin is in PATH..." -ForegroundColor Yellow
if ($env:PATH -like "*$gobinPath*") {
    Write-Host "Success: GOPATH/bin is in PATH" -ForegroundColor Green
} else {
    Write-Host "Warning: GOPATH/bin is not in PATH" -ForegroundColor Red
    Write-Host "Please add it to your PATH environment variable" -ForegroundColor Yellow
}
Write-Host ""

# Create third_party directory for Google API protos
Write-Host "Setting up Google API protos..." -ForegroundColor Yellow
$thirdPartyDir = "third_party/google/api"
if (-not (Test-Path $thirdPartyDir)) {
    New-Item -ItemType Directory -Force -Path $thirdPartyDir | Out-Null
}

# Download Google API protos if not exist
$annotationsProto = Join-Path $thirdPartyDir "annotations.proto"
$httpProto = Join-Path $thirdPartyDir "http.proto"

if (-not (Test-Path $annotationsProto)) {
    Write-Host "Downloading annotations.proto..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto" -OutFile $annotationsProto
}

if (-not (Test-Path $httpProto)) {
    Write-Host "Downloading http.proto..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto" -OutFile $httpProto
}
Write-Host "Success: Google API protos ready" -ForegroundColor Green
Write-Host ""

# Generate proto files for gRPC
Write-Host "Generating gRPC proto files..." -ForegroundColor Yellow
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/iam.proto
if ($LASTEXITCODE -eq 0) {
    Write-Host "Success: gRPC proto files generated" -ForegroundColor Green
} else {
    Write-Host "Error: Failed to generate gRPC proto files" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Generate gateway proto files for REST API
Write-Host "Generating Gateway proto files..." -ForegroundColor Yellow
protoc -I. -I./third_party --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=generate_unbound_methods=true pkg/proto/iam_gateway.proto
if ($LASTEXITCODE -eq 0) {
    Write-Host "Success: Gateway proto files generated" -ForegroundColor Green
} else {
    Write-Host "Error: Failed to generate gateway proto files" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Generate OpenAPI/Swagger documentation
Write-Host "Generating OpenAPI/Swagger documentation..." -ForegroundColor Yellow
protoc -I. -I./third_party --openapiv2_out=. --openapiv2_opt=logtostderr=true pkg/proto/iam_gateway.proto
if ($LASTEXITCODE -eq 0) {
    Write-Host "Success: OpenAPI documentation generated" -ForegroundColor Green
} else {
    Write-Host "Warning: Failed to generate OpenAPI documentation (optional)" -ForegroundColor Yellow
}
Write-Host ""

# Run go mod tidy
Write-Host "Running go mod tidy..." -ForegroundColor Yellow
go mod tidy
if ($LASTEXITCODE -eq 0) {
    Write-Host "Success: go mod tidy completed successfully" -ForegroundColor Green
} else {
    Write-Host "Error: Failed to run go mod tidy" -ForegroundColor Red
    exit 1
}
Write-Host ""

Write-Host "========================================" -ForegroundColor Green
Write-Host "Setup completed successfully!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Generated files:" -ForegroundColor Yellow
Write-Host "  gRPC files:" -ForegroundColor Cyan
Write-Host "    - pkg/proto/iam.pb.go" -ForegroundColor White
Write-Host "    - pkg/proto/iam_grpc.pb.go" -ForegroundColor White
Write-Host "  Gateway files:" -ForegroundColor Cyan
Write-Host "    - pkg/proto/iam_gateway.pb.go" -ForegroundColor White
Write-Host "    - pkg/proto/iam_gateway.pb.gw.go" -ForegroundColor White
Write-Host "    - pkg/proto/iam_gateway.swagger.json (OpenAPI)" -ForegroundColor White
Write-Host ""
Write-Host "You can now run the service:" -ForegroundColor Yellow
Write-Host "  go run cmd/server/main.go" -ForegroundColor White
Write-Host ""
Write-Host "The service will expose:" -ForegroundColor Yellow
Write-Host "  - gRPC API on port 50051" -ForegroundColor White
Write-Host "  - REST API on port 8080" -ForegroundColor White
