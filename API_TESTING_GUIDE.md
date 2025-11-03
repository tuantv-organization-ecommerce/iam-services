# API Testing Guide - IAM Services

Hướng dẫn đầy đủ về cách test API của IAM Services sau khi deploy.

**Date**: November 2024  
**Status**: ✅ Ready for Testing

---

## 📋 Table of Contents

1. [Quick Start](#quick-start)
2. [Test Locally](#test-locally)
3. [Test with Swagger UI](#test-with-swagger-ui)
4. [Test with curl](#test-with-curl)
5. [Test with Postman](#test-with-postman)
6. [Test with grpcurl](#test-with-grpcurl)
7. [Test After CI/CD Deploy](#test-after-cicd-deploy)
8. [Common Test Scenarios](#common-test-scenarios)

---

## Quick Start

### Prerequisites

**Service Must Be Running:**
- Local: `go run cmd/server/main.go`
- Docker: `docker-compose up`
- Remote: Check deployment status on server

**Service Endpoints:**
- HTTP API: `http://localhost:8080`
- gRPC API: `localhost:50051`
- Swagger UI: `http://localhost:8080/swagger/`
- Health Check: `http://localhost:8080/health`

---

## Test Locally

### Step 1: Start Service

```bash
# Option 1: Run directly
cd iam-services
go run cmd/server/main.go

# Option 2: Build and run
go build -o bin/iam-service cmd/server/main.go
./bin/iam-service

# Option 3: Docker Compose
docker-compose up -d
```

### Step 2: Verify Service Running

```bash
# Check health endpoint
curl http://localhost:8080/health

# Expected response:
{"status":"healthy"}
```

### Step 3: Check Logs

```bash
# Service logs should show:
{"level":"info","msg":"HTTP server starting","address":"0.0.0.0:8080"}
{"level":"info","msg":"gRPC server starting","address":"0.0.0.0:50051"}
```

---

## Test with Swagger UI

### Access Swagger UI

**URL**: `http://localhost:8080/swagger/`

**Authentication**:
- Username: `admin` (default, từ `SWAGGER_AUTH_USERNAME`)
- Password: `changeme` (default, từ `SWAGGER_AUTH_PASSWORD`)

### Step-by-Step Guide

#### 1. Open Swagger UI
```bash
# Windows
start http://localhost:8080/swagger/

# macOS
open http://localhost:8080/swagger/

# Linux
xdg-open http://localhost:8080/swagger/
```

#### 2. Login với Basic Auth
- Nhập username: `admin`
- Nhập password: `changeme`

#### 3. Test Register Endpoint

**Endpoint**: `POST /v1/auth/register`

1. Click vào endpoint để expand
2. Click "Try it out"
3. Fill request body:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "Test123!@#",
  "full_name": "Test User"
}
```
4. Click "Execute"
5. Check response (should be 200 OK)

#### 4. Test Login Endpoint

**Endpoint**: `POST /v1/auth/login`

1. Click "Try it out"
2. Fill request body:
```json
{
  "username": "testuser",
  "password": "Test123!@#"
}
```
3. Click "Execute"
4. **Copy access_token** from response

#### 5. Test Protected Endpoint

**Endpoint**: `GET /v1/users/{user_id}/roles`

1. Click "Authorize" button (top right)
2. Enter: `Bearer <your-access-token>`
3. Click "Authorize"
4. Now test any protected endpoint

---

## Test with curl

### Authentication Flow

#### 1. Register New User

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'
```

**Expected Response**:
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "is_active": true,
    "created_at": "2024-11-02T10:00:00Z"
  }
}
```

#### 2. Login

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecurePass123!"
  }'
```

**Expected Response**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 86400,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe"
  }
}
```

#### 3. Save Token to Variable

```bash
# Linux/macOS
TOKEN=$(curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"johndoe","password":"SecurePass123!"}' \
  | jq -r '.access_token')

echo "Token: $TOKEN"

# Windows PowerShell
$response = Invoke-RestMethod -Uri "http://localhost:8080/v1/auth/login" -Method POST -ContentType "application/json" -Body '{"username":"johndoe","password":"SecurePass123!"}'
$TOKEN = $response.access_token
Write-Host "Token: $TOKEN"
```

#### 4. Call Protected Endpoint

```bash
# Get user roles
curl -X GET "http://localhost:8080/v1/users/550e8400-e29b-41d4-a716-446655440000/roles" \
  -H "Authorization: Bearer $TOKEN"
```

### Role Management

#### Create Role

```bash
curl -X POST http://localhost:8080/v1/roles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "editor",
    "description": "Content editor role",
    "permission_ids": ["perm-001", "perm-002"]
  }'
```

#### List Roles

```bash
# With pagination
curl "http://localhost:8080/v1/roles?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"
```

#### Assign Role to User

```bash
curl -X POST http://localhost:8080/v1/roles/assign \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "user_id": "user-123",
    "role_id": "role-456"
  }'
```

### CMS Authorization

#### Check CMS Access

```bash
curl -X POST http://localhost:8080/v1/access/cms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "user_id": "user-123",
    "cms_tab": "product",
    "action": "POST"
  }'
```

#### Get User's CMS Tabs

```bash
curl "http://localhost:8080/v1/cms/users/user-123/tabs" \
  -H "Authorization: Bearer $TOKEN"
```

---

## Test with Postman

### Import OpenAPI Spec

#### Option 1: Import from URL

1. Open Postman
2. Click **Import**
3. Select **Link** tab
4. Enter URL: `http://localhost:8080/swagger.json`
5. Click **Continue** → **Import**

#### Option 2: Import File

1. Download spec:
```bash
curl http://localhost:8080/swagger.json -o iam-api.json
```
2. Postman → **Import** → **Upload Files**
3. Select `iam-api.json`

### Setup Environment

1. Click **Environments** (left sidebar)
2. Click **+** to create new environment
3. Name: `IAM Service Local`
4. Add variables:

| Variable | Initial Value | Current Value |
|----------|---------------|---------------|
| baseUrl | http://localhost:8080 | http://localhost:8080 |
| token | (empty) | (will be set by script) |

5. **Save**

### Create Authentication Flow

#### 1. Create Collection: "IAM Service"

1. Right-click Collections → **New Collection**
2. Name: `IAM Service`

#### 2. Add Pre-request Script (Collection Level)

Collection → **Scripts** → **Pre-request**:

```javascript
// Auto-refresh token if expired
const token = pm.environment.get("token");
const tokenExpiry = pm.environment.get("token_expiry");

if (!token || Date.now() >= tokenExpiry) {
    // Token expired, get new one
    pm.sendRequest({
        url: pm.environment.get("baseUrl") + "/v1/auth/login",
        method: "POST",
        header: {
            "Content-Type": "application/json"
        },
        body: {
            mode: "raw",
            raw: JSON.stringify({
                username: "johndoe",
                password: "SecurePass123!"
            })
        }
    }, function(err, response) {
        if (err) {
            console.log(err);
        } else {
            const jsonData = response.json();
            pm.environment.set("token", jsonData.access_token);
            pm.environment.set("token_expiry", Date.now() + (jsonData.expires_in * 1000));
        }
    });
}
```

#### 3. Add Requests

**Register Request**:
- Method: `POST`
- URL: `{{baseUrl}}/v1/auth/register`
- Body (JSON):
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "Test123!",
  "full_name": "Test User"
}
```

**Login Request**:
- Method: `POST`
- URL: `{{baseUrl}}/v1/auth/login`
- Body (JSON):
```json
{
  "username": "testuser",
  "password": "Test123!"
}
```
- Tests tab:
```javascript
pm.test("Status is 200", function() {
    pm.response.to.have.status(200);
});

const jsonData = pm.response.json();
pm.environment.set("token", jsonData.access_token);
```

**Get User Roles Request**:
- Method: `GET`
- URL: `{{baseUrl}}/v1/users/:user_id/roles`
- Headers:
  - Key: `Authorization`
  - Value: `Bearer {{token}}`

---

## Test with grpcurl

### Install grpcurl

```bash
# macOS
brew install grpcurl

# Windows
choco install grpcurl

# Or download from GitHub
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### List Services

```bash
grpcurl -plaintext localhost:50051 list
```

**Expected Output**:
```
iam.IAMService
grpc.health.v1.Health
```

### Describe Service

```bash
grpcurl -plaintext localhost:50051 describe iam.IAMService
```

### Test gRPC Endpoints

#### Register

```bash
grpcurl -plaintext -d '{
  "username": "grpcuser",
  "email": "grpc@example.com",
  "password": "GrpcPass123!",
  "full_name": "gRPC User"
}' localhost:50051 iam.IAMService/Register
```

#### Login

```bash
grpcurl -plaintext -d '{
  "username": "grpcuser",
  "password": "GrpcPass123!"
}' localhost:50051 iam.IAMService/Login
```

#### Verify Token

```bash
grpcurl -plaintext -d '{
  "token": "your-access-token-here"
}' localhost:50051 iam.IAMService/VerifyToken
```

---

## Test After CI/CD Deploy

### Test Staging Environment

Sau khi CI/CD deploy lên staging (khi push lên branch `develop`):

#### 1. Check Deployment Status

```bash
# Check GitHub Actions
# Go to: Repository → Actions → Latest workflow run
# Verify: deploy-staging job passed ✓
```

#### 2. Test Staging Endpoint

```bash
# Replace with your staging URL
STAGING_URL="https://iam-staging.example.com"

# Health check
curl $STAGING_URL/health

# Register
curl -X POST $STAGING_URL/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "staginguser",
    "email": "staging@example.com",
    "password": "Staging123!",
    "full_name": "Staging User"
  }'
```

#### 3. Swagger UI on Staging

```bash
# Open in browser
open https://iam-staging.example.com/swagger/
```

### Test Production Environment

Sau khi merge lên `main`:

#### 1. Verify Production Deploy

```bash
# Check GitHub Actions
# Verify: deploy-production job passed ✓
# Check: GitHub Release created
```

#### 2. Test Production Endpoint

```bash
PROD_URL="https://iam.example.com"

# Health check
curl $PROD_URL/health

# Test với production credentials
curl -X POST $PROD_URL/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "$PRODUCTION_PASSWORD"
  }'
```

**⚠️ Important**: Use strong credentials in production!

---

## Common Test Scenarios

### Scenario 1: Complete User Flow

```bash
# 1. Register
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testflow",
    "email": "testflow@example.com",
    "password": "Flow123!",
    "full_name": "Test Flow"
  }'

# 2. Login and save token
TOKEN=$(curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testflow","password":"Flow123!"}' \
  | jq -r '.access_token')

# 3. Verify token
curl -X POST http://localhost:8080/v1/auth/verify \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}"

# 4. Get user profile (protected endpoint)
USER_ID=$(curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testflow","password":"Flow123!"}' \
  | jq -r '.user.id')

curl "http://localhost:8080/v1/users/$USER_ID/roles" \
  -H "Authorization: Bearer $TOKEN"

# 5. Logout
curl -X POST http://localhost:8080/v1/auth/logout \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"user_id\":\"$USER_ID\",\"token\":\"$TOKEN\"}"
```

### Scenario 2: Role Management Flow

```bash
# Login as admin
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.access_token')

# Create new role
ROLE_ID=$(curl -s -X POST http://localhost:8080/v1/roles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "content_editor",
    "description": "Can edit content"
  }' | jq -r '.id')

# Assign role to user
curl -X POST http://localhost:8080/v1/roles/assign \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{
    \"user_id\": \"$USER_ID\",
    \"role_id\": \"$ROLE_ID\"
  }"

# Verify assignment
curl "http://localhost:8080/v1/users/$USER_ID/roles" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### Scenario 3: CMS Authorization Flow

```bash
# Create CMS role
curl -X POST http://localhost:8080/v1/cms/roles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "cms_product_manager",
    "description": "Product manager",
    "tabs": ["product", "inventory"]
  }'

# Check user CMS access
curl -X POST http://localhost:8080/v1/access/cms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "user_id": "user-123",
    "cms_tab": "product",
    "action": "POST"
  }'

# Get user's accessible tabs
curl "http://localhost:8080/v1/cms/users/user-123/tabs" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

## Testing Scripts

### Comprehensive Test Script

Create `scripts/test-api.ps1`:

```powershell
# test-api.ps1 - Comprehensive API Testing Script

$BASE_URL = "http://localhost:8080"

Write-Host "=== IAM Service API Testing ===" -ForegroundColor Cyan

# 1. Health Check
Write-Host "`n1. Testing Health Endpoint..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$BASE_URL/health" -Method GET
    Write-Host "✓ Health: $($health.status)" -ForegroundColor Green
} catch {
    Write-Host "✗ Health check failed: $_" -ForegroundColor Red
    exit 1
}

# 2. Register
Write-Host "`n2. Testing Register..." -ForegroundColor Yellow
$registerBody = @{
    username = "testuser_$(Get-Random)"
    email = "test_$(Get-Random)@example.com"
    password = "Test123!@#"
    full_name = "Test User"
} | ConvertTo-Json

try {
    $registerResponse = Invoke-RestMethod -Uri "$BASE_URL/v1/auth/register" `
        -Method POST -ContentType "application/json" -Body $registerBody
    Write-Host "✓ User registered: $($registerResponse.user.username)" -ForegroundColor Green
    $username = $registerResponse.user.username
    $userId = $registerResponse.user.id
} catch {
    Write-Host "✗ Register failed: $_" -ForegroundColor Red
    exit 1
}

# 3. Login
Write-Host "`n3. Testing Login..." -ForegroundColor Yellow
$loginBody = @{
    username = $username
    password = "Test123!@#"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$BASE_URL/v1/auth/login" `
        -Method POST -ContentType "application/json" -Body $loginBody
    Write-Host "✓ Login successful" -ForegroundColor Green
    $token = $loginResponse.access_token
    Write-Host "  Token: $($token.Substring(0,20))..." -ForegroundColor Gray
} catch {
    Write-Host "✗ Login failed: $_" -ForegroundColor Red
    exit 1
}

# 4. Verify Token
Write-Host "`n4. Testing Token Verification..." -ForegroundColor Yellow
$verifyBody = @{
    token = $token
} | ConvertTo-Json

try {
    $verifyResponse = Invoke-RestMethod -Uri "$BASE_URL/v1/auth/verify" `
        -Method POST -ContentType "application/json" -Body $verifyBody
    Write-Host "✓ Token valid: $($verifyResponse.valid)" -ForegroundColor Green
} catch {
    Write-Host "✗ Verify failed: $_" -ForegroundColor Red
}

# 5. Get User Roles
Write-Host "`n5. Testing Get User Roles..." -ForegroundColor Yellow
try {
    $headers = @{
        "Authorization" = "Bearer $token"
    }
    $rolesResponse = Invoke-RestMethod -Uri "$BASE_URL/v1/users/$userId/roles" `
        -Method GET -Headers $headers
    Write-Host "✓ User has $($rolesResponse.roles.Count) role(s)" -ForegroundColor Green
} catch {
    Write-Host "✗ Get roles failed: $_" -ForegroundColor Red
}

Write-Host "`n=== All Tests Completed ===" -ForegroundColor Cyan
Write-Host "Summary:" -ForegroundColor Yellow
Write-Host "  Username: $username" -ForegroundColor Gray
Write-Host "  User ID: $userId" -ForegroundColor Gray
Write-Host "  Token: $($token.Substring(0,30))..." -ForegroundColor Gray
```

Run the script:
```powershell
.\scripts\test-api.ps1
```

### Quick Test Script (Bash)

Create `scripts/test-api.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== IAM Service API Testing ==="

# Health check
echo -e "\n1. Testing Health..."
curl -f $BASE_URL/health && echo " ✓" || echo " ✗"

# Register
echo -e "\n2. Testing Register..."
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser_'$RANDOM'",
    "email": "test_'$RANDOM'@example.com",
    "password": "Test123!",
    "full_name": "Test User"
  }')

USERNAME=$(echo $REGISTER_RESPONSE | jq -r '.user.username')
echo "✓ Registered: $USERNAME"

# Login
echo -e "\n3. Testing Login..."
TOKEN=$(curl -s -X POST $BASE_URL/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$USERNAME\",
    \"password\": \"Test123!\"
  }" | jq -r '.access_token')

echo "✓ Token: ${TOKEN:0:30}..."

echo -e "\n=== Tests Completed ==="
```

Make executable and run:
```bash
chmod +x scripts/test-api.sh
./scripts/test-api.sh
```

---

## Troubleshooting

### Issue 1: Connection Refused

**Symptoms**: `curl: (7) Failed to connect`

**Solutions**:
```bash
# Check if service is running
ps aux | grep iam-service

# Check port
netstat -an | grep 8080

# Restart service
go run cmd/server/main.go
```

### Issue 2: 401 Unauthorized

**Symptoms**: API returns 401

**Solutions**:
```bash
# Get fresh token
TOKEN=$(curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.access_token')

# Use Bearer prefix
curl -H "Authorization: Bearer $TOKEN" ...
```

### Issue 3: Swagger UI Not Loading

**Solutions**:
```bash
# Check Swagger enabled
cat .env | grep SWAGGER_ENABLED

# Regenerate proto files
.\scripts\generate-proto-simple.ps1

# Restart service
```

### Issue 4: Database Connection Error

**Solutions**:
```bash
# Check PostgreSQL
pg_isready

# Check migrations
psql -U postgres -d iam_db -c "\dt"

# Re-run migrations if needed
```

---

## Best Practices

### Security
1. ✅ **Never log tokens** in production
2. ✅ **Use HTTPS** in production
3. ✅ **Rotate credentials** regularly
4. ✅ **Use strong passwords**
5. ✅ **Test with test accounts** only

### Testing
1. ✅ **Test in order**: Register → Login → Protected endpoints
2. ✅ **Save tokens**: Store in variables for reuse
3. ✅ **Check responses**: Verify status codes and data
4. ✅ **Clean up**: Remove test users after testing
5. ✅ **Automate**: Use scripts for repeated tests

### CI/CD
1. ✅ **Test staging first** before production
2. ✅ **Monitor health** endpoints
3. ✅ **Check logs** after deployment
4. ✅ **Rollback plan** ready
5. ✅ **Verify migrations** applied

---

## Quick Reference

### Useful Commands

```bash
# Health check
curl http://localhost:8080/health

# Get token quickly
curl -s -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.access_token'

# Test with saved token
TOKEN="your-token-here"
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/v1/roles

# Check service logs
docker-compose logs -f iam-service

# Check database
psql -U postgres -d iam_db -c "SELECT * FROM users LIMIT 5;"
```

---

## Resources

- **Main Documentation**: [README.md](README.md)
- **Troubleshooting**: [fix_error_ci_cd.md](fix_error_ci_cd.md)
- **Swagger Guide**: [SWAGGER_GUIDE.md](SWAGGER_GUIDE.md)
- **CI/CD Setup**: [CI_CD_SETUP_GUIDE.md](CI_CD_SETUP_GUIDE.md)

---

**Last Updated**: November 2024  
**Status**: ✅ Ready for Testing

