# Quick Proto Generation Script
# Run: powershell -ExecutionPolicy Bypass -File .\scripts\generate-proto-simple.ps1

Write-Host "`n========================================"
Write-Host "Generating Proto Files for IAM Service"
Write-Host "========================================`n"

# Step 1: Setup third_party first
Write-Host "[1/4] Setting up Google API protos..."
$thirdPartyDir = "third_party/google/api"
if (-not (Test-Path $thirdPartyDir)) {
    New-Item -ItemType Directory -Force -Path $thirdPartyDir | Out-Null
}

# Download required Google API proto files
$annotationsProto = Join-Path $thirdPartyDir "annotations.proto"
$httpProto = Join-Path $thirdPartyDir "http.proto"

if (-not (Test-Path $annotationsProto)) {
    Write-Host "  Creating annotations.proto..."
    
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
    Write-Host "  [OK] annotations.proto created"
}

if (-not (Test-Path $httpProto)) {
    Write-Host "  Creating http.proto..."
    
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
    Write-Host "  [OK] http.proto created"
}

Write-Host "[OK] Google API protos ready"

# Step 2: Generate gRPC files
Write-Host "`n[2/4] Generating gRPC proto files..."

# Try to find protoc include directory
$protocPath = (Get-Command protoc -ErrorAction SilentlyContinue).Source
$includeDir = $null

if ($protocPath) {
    $protocDir = Split-Path $protocPath -Parent
    # Try different common locations
    $possibleIncludes = @(
        (Join-Path (Split-Path $protocDir -Parent) "include"),
        "C:\ProgramData\chocolatey\lib\protoc\tools\include",
        (Join-Path $protocDir "include")
    )
    
    foreach ($dir in $possibleIncludes) {
        if (Test-Path $dir) {
            $includeDir = $dir
            Write-Host "  Using protoc include: $includeDir" -ForegroundColor Cyan
            break
        }
    }
}

if ($includeDir) {
    protoc -I. -I./third_party -I"$includeDir" --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/iam.proto
} else {
    Write-Host "  Warning: Could not find protoc include directory" -ForegroundColor Yellow
    protoc -I. -I./third_party --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/iam.proto
}

if ($LASTEXITCODE -eq 0) {
    Write-Host "[OK] gRPC files generated"
} else {
    Write-Host "[ERROR] Failed to generate gRPC files"
    Write-Host "`nMake sure you have installed:"
    Write-Host "  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1"
    Write-Host "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0"
    exit 1
}

# Step 3: Generate Gateway and OpenAPI files
Write-Host "`n[3/4] Generating Gateway and OpenAPI files..."

# Use the same include directory
if ($includeDir -and (Test-Path $includeDir)) {
    protoc -I. -I./third_party -I"$includeDir" `
      --go_out=. --go_opt=paths=source_relative `
      --go-grpc_out=. --go-grpc_opt=paths=source_relative `
      --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative `
      --grpc-gateway_opt=generate_unbound_methods=true `
      --openapiv2_out=./pkg/proto --openapiv2_opt=logtostderr=true `
      pkg/proto/iam.proto
} else {
    protoc -I. -I./third_party `
      --go_out=. --go_opt=paths=source_relative `
      --go-grpc_out=. --go-grpc_opt=paths=source_relative `
      --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative `
      --grpc-gateway_opt=generate_unbound_methods=true `
      --openapiv2_out=./pkg/proto --openapiv2_opt=logtostderr=true `
      pkg/proto/iam.proto
}

if ($LASTEXITCODE -eq 0) {
    Write-Host "[OK] Gateway and OpenAPI files generated"
} else {
    Write-Host "[ERROR] Failed to generate Gateway files"
    Write-Host "`nMake sure you have installed:"
    Write-Host "  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1"
    Write-Host "  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.18.1"
    exit 1
}

Write-Host "`n========================================"
Write-Host "[OK] Proto generation completed!"
Write-Host "========================================`n"

Write-Host "Generated files:"
Write-Host "  - pkg/proto/iam.pb.go"
Write-Host "  - pkg/proto/iam_grpc.pb.go"
Write-Host "  - pkg/proto/iam.pb.gw.go"
Write-Host "  - pkg/proto/iam.swagger.json`n"

Write-Host "Next steps:"
Write-Host "  1. Run: go mod tidy"
Write-Host "  2. Run: go run cmd/server/main.go"
Write-Host "  3. Access Swagger UI: http://localhost:8080/swagger/"
Write-Host "     Username: admin"
Write-Host "     Password: changeme`n"

