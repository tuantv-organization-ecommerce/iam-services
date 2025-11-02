# Quick Proto Generation Script
# Run: powershell -ExecutionPolicy Bypass -File .\scripts\generate-proto.ps1

Write-Host "`n========================================" -ForegroundColor Green
Write-Host "Generating Proto Files for IAM Service" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Green

# Step 1: Generate gRPC files
Write-Host "[1/3] Generating gRPC proto files..." -ForegroundColor Yellow
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/iam.proto

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ gRPC files generated" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to generate gRPC files" -ForegroundColor Red
    Write-Host "`nMake sure you have installed:" -ForegroundColor Yellow
    Write-Host "  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1" -ForegroundColor White
    Write-Host "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0" -ForegroundColor White
    exit 1
}

# Step 2: Create third_party directory for Google API protos
Write-Host "`n[2/3] Setting up Google API protos..." -ForegroundColor Yellow
$thirdPartyDir = "third_party/google/api"
if (-not (Test-Path $thirdPartyDir)) {
    New-Item -ItemType Directory -Force -Path $thirdPartyDir | Out-Null
}

# Download required Google API proto files
$annotationsProto = Join-Path $thirdPartyDir "annotations.proto"
$httpProto = Join-Path $thirdPartyDir "http.proto"

if (-not (Test-Path $annotationsProto)) {
    Write-Host "  Downloading annotations.proto..." -ForegroundColor Cyan
    try {
        Invoke-WebRequest -Uri "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto" -OutFile $annotationsProto -ErrorAction Stop
        Write-Host "  ✓ annotations.proto downloaded" -ForegroundColor Green
    } catch {
        Write-Host "  ✗ Failed to download annotations.proto" -ForegroundColor Red
        Write-Host "  Creating manual file..." -ForegroundColor Yellow
        
        $annotationsContent = @"
syntax = "proto3";

package google.api;

import "google/api/http.proto";
import "google/protobuf/descriptor.proto";

option go_package = "google.golang.org/genproto/googleapis/api/annotations;annotations";

extend google.protobuf.MethodOptions {
  HttpRule http = 72295728;
}
"@
        Set-Content -Path $annotationsProto -Value $annotationsContent
    }
}

if (-not (Test-Path $httpProto)) {
    Write-Host "  Downloading http.proto..." -ForegroundColor Cyan
    try {
        Invoke-WebRequest -Uri "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto" -OutFile $httpProto -ErrorAction Stop
        Write-Host "  ✓ http.proto downloaded" -ForegroundColor Green
    } catch {
        Write-Host "  ✗ Failed to download http.proto" -ForegroundColor Red
        Write-Host "  Creating manual file..." -ForegroundColor Yellow
        
        $httpContent = @"
syntax = "proto3";

package google.api;

option go_package = "google.golang.org/genproto/googleapis/api/annotations;annotations";

message Http {
  repeated HttpRule rules = 1;
  bool fully_decode_reserved_expansion = 2;
}

message HttpRule {
  string selector = 1;
  oneof pattern {
    string get = 2;
    string put = 3;
    string post = 4;
    string delete = 5;
    string patch = 6;
    CustomHttpPattern custom = 8;
  }
  string body = 7;
  string response_body = 12;
  repeated HttpRule additional_bindings = 11;
}

message CustomHttpPattern {
  string kind = 1;
  string path = 2;
}
"@
        Set-Content -Path $httpProto -Value $httpContent
    }
}

Write-Host "Google API protos ready" -ForegroundColor Green

# Step 3: Generate Gateway and OpenAPI files
Write-Host "`n[3/3] Generating Gateway and OpenAPI files..." -ForegroundColor Yellow
protoc -I. -I./third_party `
  --go_out=. --go_opt=paths=source_relative `
  --go-grpc_out=. --go-grpc_opt=paths=source_relative `
  --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative `
  --grpc-gateway_opt=generate_unbound_methods=true `
  --openapiv2_out=. --openapiv2_opt=logtostderr=true `
  pkg/proto/iam_gateway.proto

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Gateway and OpenAPI files generated" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to generate Gateway files" -ForegroundColor Red
    Write-Host "`nMake sure you have installed:" -ForegroundColor Yellow
    Write-Host "  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1" -ForegroundColor White
    Write-Host "  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.18.1" -ForegroundColor White
    exit 1
}

Write-Host "`n========================================" -ForegroundColor Green
Write-Host "✓ Proto generation completed!" -ForegroundColor Green
Write-Host "========================================`n" -ForegroundColor Green

Write-Host "Generated files:" -ForegroundColor Yellow
Write-Host "  ✓ pkg/proto/iam.pb.go" -ForegroundColor White
Write-Host "  ✓ pkg/proto/iam_grpc.pb.go" -ForegroundColor White
Write-Host "  ✓ pkg/proto/iam_gateway.pb.go" -ForegroundColor White
Write-Host "  ✓ pkg/proto/iam_gateway.pb.gw.go" -ForegroundColor White
Write-Host "  ✓ pkg/proto/iam_gateway.swagger.json`n" -ForegroundColor White

Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Run: go mod tidy" -ForegroundColor White
Write-Host "  2. Run: go run cmd/server/main.go" -ForegroundColor White
Write-Host "  3. Access Swagger UI: http://localhost:8080/swagger/`n" -ForegroundColor White

