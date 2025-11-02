# CI/CD Error Fixes & Troubleshooting Guide

Comprehensive guide for fixing common CI/CD errors in iam-services GitHub Actions workflows.

**Last Updated**: November 2024  
**Go Version**: 1.19  
**Status**: ✅ All CI checks passing

---

## 📋 Table of Contents

1. [Quick Reference](#quick-reference)
2. [Workflow Issues](#workflow-issues)
3. [Database & Migration Issues](#database--migration-issues)
4. [Linting Issues](#linting-issues)
5. [Security & Code Quality](#security--code-quality)
6. [Swagger & HTTP Gateway Issues](#swagger--http-gateway-issues)
7. [Deployment Issues](#deployment-issues)
8. [Best Practices](#best-practices)

---

## Quick Reference

### CI Pipeline Overview

```
Push/PR → Lint → Test → Build → Security → Docker → Deploy
           ↓       ↓       ↓        ↓         ↓        ↓
        Pass    Pass    Pass     Pass      Pass     Pass
```

### Common Quick Fixes

```bash
# Format code
go fmt ./...

# Fix imports
go mod tidy

# Run all checks locally
.\scripts\check-all.ps1     # Windows
./scripts/check-all.sh      # Linux/Mac

# Install linter (Go 1.19 compatible)
.\scripts\install-golangci-lint.ps1 -Version "v1.54.2"
```

---

## Workflow Issues

### 1. Deprecated Artifact Action (v3)

**Error**: "This request has been automatically failed because it uses a deprecated version of actions/upload-artifact: v3"

**Cause**: GitHub deprecated v3 of artifact actions

**Fix**: Upgrade to v4
```yaml
# In .github/workflows/ci-cd.yml and test.yml
- uses: actions/upload-artifact@v4  # Changed from v3
```

**Files to Update**:
- `.github/workflows/ci-cd.yml` (build job)
- `.github/workflows/test.yml` (unit-tests, benchmark-tests)

---

### 2. Working Directory Not Found

**Error**: "No such file or directory: ecommerce/back_end/iam-services"

**Cause**: Hard-coded `working-directory` doesn't match actual repo structure

**Fix**: Auto-detect service directory
```yaml
- name: Set service directory
  run: |
    if [ -d "ecommerce/back_end/iam-services" ]; then
      echo "SERVICE_DIR=ecommerce/back_end/iam-services" >> $GITHUB_ENV
    else
      echo "SERVICE_DIR=." >> $GITHUB_ENV
    fi

# Then use ${{ env.SERVICE_DIR }} in all steps
- name: Build
  working-directory: ${{ env.SERVICE_DIR }}
  run: go build ./...
```

**Apply to**: All jobs with `working-directory`

---

### 3. Environment File Missing

**Error**: ".env.example not found"

**Cause**: `.env.example` not created or gitignored

**Fix**: Create from template
```powershell
# Windows PowerShell
Copy-Item .env.template .env.example

# Linux/macOS
cp .env.template .env.example
```

**Or use script**:
```bash
.\scripts\setup-ci.ps1      # Windows
./scripts/setup-ci.sh       # Linux/Mac
```

---

## Database & Migration Issues

### 4. PostgreSQL Client Not Found

**Error**: "psql: command not found"

**Cause**: GitHub runner doesn't have PostgreSQL client pre-installed

**Fix**: Install client before migrations
```yaml
- name: Install PostgreSQL client
  run: |
    sudo apt-get update
    sudo apt-get install -y postgresql-client
```

**Apply to**: All jobs running `psql` commands

---

### 5. Missing Migrations

**Error**: "relation 'casbin_rule_user' does not exist"

**Cause**: New migration files not added to workflow

**Fix**: Add all migrations in order
```yaml
- name: Run database migrations
  run: |
    psql -h localhost -U postgres -d iam_db_test -f migrations/001_init_schema.sql
    psql -h localhost -U postgres -d iam_db_test -f migrations/002_seed_data.sql
    psql -h localhost -U postgres -d iam_db_test -f migrations/003_casbin_tables.sql
    psql -h localhost -U postgres -d iam_db_test -f migrations/004_casbin_seed_data.sql
    psql -h localhost -U postgres -d iam_db_test -f migrations/005_separate_user_cms_authorization.sql
    psql -h localhost -U postgres -d iam_db_test -f migrations/006_seed_separated_authorization.sql
```

**Current Migrations** (as of Nov 2024):
1. `001_init_schema.sql` - Core schema
2. `002_seed_data.sql` - Initial data
3. `003_casbin_tables.sql` - Casbin tables
4. `004_casbin_seed_data.sql` - Casbin seed
5. `005_separate_user_cms_authorization.sql` - Separated auth architecture
6. `006_seed_separated_authorization.sql` - Separated auth seed data

---

### 6. Database Connection Timeout

**Error**: "connection refused" or "timeout"

**Fix**: Wait for PostgreSQL to be ready
```yaml
services:
  postgres:
    options: >-
      --health-cmd pg_isready
      --health-interval 10s
      --health-timeout 5s
      --health-retries 5

- name: Wait for PostgreSQL
  run: |
    until pg_isready -h localhost -p 5432 -U postgres; do
      echo "Waiting for PostgreSQL..."; 
      sleep 2; 
    done
```

---

## Linting Issues

### 7. errcheck: "Error return value not checked"

**Error**: `errcheck` reports unchecked error returns

**Common Cases**:

#### A. rows.Close() in DAO
**Problem**: `defer func() { _ = rows.Close() }()` considered error suppression

**Fix**: Use idiomatic Go
```go
// Before (flagged by errcheck)
defer func() { _ = rows.Close() }()

// After (idiomatic Go)
defer rows.Close()
```

**Files Fixed**:
- `internal/dao/role_permission_dao.go`
- `internal/dao/role_dao.go`
- `internal/dao/permission_dao.go`
- `internal/dao/cms_tab_api_dao.go`
- `internal/dao/cms_role_dao.go`
- `internal/dao/api_resource_dao.go`

#### B. Logger.Sync() and Shutdown()
**Problem**: Not checking return values

**Fix**: Check and log errors
```go
// Sync() - can fail harmlessly on stdout/stderr
if syncErr := logger.Sync(); syncErr != nil {
    logger.Warn("Logger sync returned error (may be harmless)", zap.Error(syncErr))
}

// Shutdown() - should be checked
if err := a.Shutdown(context.Background()); err != nil {
    logger.Error("Error during shutdown", zap.Error(err))
}
```

**Files Fixed**:
- `internal/app/app.go`
- `internal/container/container.go`

---

### 8. revive: "Should have package comment"

**Error**: `revive` requires package-level comment

**Fix**: Add comment before `package` declaration
```go
// Package dto provides Data Transfer Objects for application layer communication.
// It contains request/response structures for authentication and authorization operations.
package dto
```

**Pattern**:
```go
// Package <name> <description of package purpose>.
package <name>
```

**Files Fixed**:
- `internal/application/dto/auth_dto.go`
- `internal/domain/casbin.go`
- All model files in `internal/domain/model/`

---

### 9. revive: "Exported symbol should have comment"

**Error**: Missing comments on exported vars/consts/methods

**Fix Examples**:

#### Exported Variables (Errors)
```go
var (
    // ErrInvalidUsername indicates username validation failed
    ErrInvalidUsername = errors.New("invalid username")
    
    // ErrInvalidEmail indicates email validation failed
    ErrInvalidEmail = errors.New("invalid email")
)
```

#### Exported Constants
```go
const (
    // DomainUser is the domain for end user authorization
    DomainUser CasbinDomain = "user"
    
    // DomainCMS is the domain for CMS/admin panel authorization
    DomainCMS CasbinDomain = "cms"
    
    // MethodGET represents HTTP GET method
    MethodGET HTTPMethod = "GET"
)
```

#### Exported Methods (Getters)
```go
// ID returns the user's unique identifier
func (u *User) ID() string { return u.id }

// Username returns the user's username
func (u *User) Username() string { return u.username }

// Email returns the user's email address
func (u *User) Email() string { return u.email }
```

**Note**: For `ID()` methods, use "ID" not "Id" to pass revive check

**Files Fixed**:
- `internal/domain/model/api_resource.go` - Error vars + HTTP method constants
- `internal/domain/model/user.go` - Error vars + all getters
- `internal/domain/model/permission.go` - Error vars + getters
- `internal/domain/model/role.go` - Error vars + getters
- `internal/domain/model/cms_role.go` - Error vars + getters
- `internal/domain/casbin.go` - Domain constants + CMS tab constants
- `internal/domain/service/authorization_service.go` - Domain constants

---

### 10. goconst: "String has N occurrences"

**Error**: Repeated string literals should be constants

**Fix**: Extract to constants
```go
// Before
func TestHashPassword(t *testing.T) {
    password := "MySecurePassword123!"  // Repeated
    // ...
}

// After
const (
    testPassword = "MySecurePassword123!"
    testUserID = "test-user-123"
)

func TestHashPassword(t *testing.T) {
    password := testPassword
    // ...
}
```

**Files Fixed**:
- `pkg/password/password_manager_test.go` - `testPassword` constant
- `pkg/jwt/jwt_manager_test.go` - `testUserID`, `testUsername` constants

---

### 11. gosec G115: "Integer overflow conversion"

**Error**: Converting `int` to `int32` may overflow

**Problem**: Direct conversion without overflow check
```go
// Flagged by gosec
totalInt32 := int32(total)  // Potential overflow
```

**Fix**: Helper function with overflow protection
```go
// safeIntToInt32 safely converts int to int32 with overflow protection
func safeIntToInt32(value int) int32 {
    const maxInt32 = 2147483647
    if value > maxInt32 {
        return maxInt32
    }
    if value < -2147483648 {
        return -2147483648
    }
    return int32(value) // #nosec G115 -- overflow checked above
}

// Usage
Total: safeIntToInt32(total),
```

**Files Fixed**:
- `internal/handler/grpc_handler.go` - Added helper, used in `ListRoles`, `ListPermissions`
- `internal/handler/casbin_handler.go` - Uses shared pattern in `ListAPIResources`

**Rationale**: `#nosec G115` is justified because we explicitly check for overflow before conversion

---

### 12. Go Version Mismatch

**Error**: "module requires Go 1.20" or "go.work file requires go >= 1.21"

**Cause**: Inconsistent Go version declarations

**Fix**: Sync all version declarations to 1.19
```bash
# Fix go.mod
(Get-Content go.mod) -replace 'go 1\.21', 'go 1.19' | Set-Content go.mod

# Fix Dockerfile
# Use: FROM golang:1.19-alpine AS builder

# Fix workflows
# Use: GO_VERSION: '1.19'
```

**Files to Check**:
- `go.mod` → `go 1.19`
- `Dockerfile` → `golang:1.19-alpine`
- `.github/workflows/ci-cd.yml` → `GO_VERSION: '1.19'`
- `.github/workflows/test.yml` → `GO_VERSION: '1.19'`

---

### 13. golangci-lint Installation (Go 1.19)

**Error**: "module requires Go 1.20" when using `go install`

**Cause**: Recent golangci-lint versions require Go 1.20+

**Fix**: Use binary download for Go 1.19
```powershell
# Recommended: Use script (auto-downloads binary)
.\scripts\install-golangci-lint.ps1 -Version "v1.54.2"

# Add to PATH
$env:PATH += ";$(go env GOPATH)\bin"

# Verify
golangci-lint version
```

**Compatibility Matrix**:
| Go Version | golangci-lint Version | Method |
|------------|----------------------|---------|
| Go 1.19    | v1.54.2              | Binary download (script) |
| Go 1.20+   | v1.55.2+             | `go install` or binary |

**See**: `INSTALLATION_GUIDE.md` for detailed instructions

---

## Security & Code Quality

### 14. Codecov Upload Failure

**Error**: Upload fails without `CODECOV_TOKEN`

**Cause**: Private repos need token for Codecov

**Fix**: Make upload optional
```yaml
- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v3
  with:
    file: ./${{ env.SERVICE_DIR }}/coverage.out
    token: ${{ secrets.CODECOV_TOKEN }}
    fail_ci_if_error: false  # Don't fail if token missing
```

**Note**: For public repos, token is optional

---

### 15. Trivy Security Scan Errors

**Error**: Trivy fails to download vulnerability database

**Fix**: Add retry logic
```yaml
- name: Run Trivy vulnerability scanner
  uses: aquasecurity/trivy-action@master
  continue-on-error: true  # Don't fail entire pipeline
  with:
    scan-type: 'fs'
    scan-ref: '.'
    format: 'sarif'
    output: 'trivy-results.sarif'
```

---

## Swagger & HTTP Gateway Issues

### 16. Swagger UI Not Loading

**Symptoms**: 404, blank page, or "Failed to load spec"

**Common Causes & Fixes**:

#### A. Swagger Disabled
```env
# Check .env file
SWAGGER_ENABLED=true
```

#### B. Missing OpenAPI Spec
```bash
# Verify file exists
ls pkg/proto/iam.swagger.json

# Regenerate if missing
.\scripts\generate-proto-simple.ps1  # Windows
./scripts/setup.sh                   # Linux/Mac
```

#### C. Proto Files Not Generated
```bash
# Install required tools
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.18.1

# Generate proto files
.\scripts\generate-proto-simple.ps1
```

#### D. HTTP Gateway Not Enabled
Check `internal/app/app.go`:
```go
// Should be uncommented
if a.config.Swagger.Enabled {
    setupHTTPGateway()  // or setupGinServer()
}
```

---

### 17. Swagger Basic Auth Loop

**Symptoms**: 401 Unauthorized, endless auth prompts

**Fix**: Check credentials and clear cache
```bash
# 1. Verify credentials in .env
cat .env | grep SWAGGER_AUTH

# 2. Clear browser cache/cookies

# 3. Try incognito mode

# 4. Test with curl
curl -u admin:changeme http://localhost:8080/swagger/
```

---

### 18. CORS Errors from Swagger

**Symptoms**: "CORS policy blocked" in browser console

**Fix**: Verify CORS middleware in Gin router
```go
// internal/router/gin_router.go
router.Use(middleware.GinCORS())  // Should be present

// internal/middleware/gin_middleware.go
func GinCORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}
```

---

## Deployment Issues

### 19. Health Check Fails (When Deploy Enabled)

**Error**: `curl -f https://.../health` returns error

**Common Causes**:

#### A. Endpoint Not Implemented
```go
// Add to internal/router/gin_router.go
router.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "healthy"})
})
```

#### B. Service Not Ready
```yaml
# Add delay in workflow
- name: Health check
  run: |
    sleep 10  # Wait for service to start
    curl -f http://localhost:8080/health || exit 1
```

#### C. Port Not Exposed
```yaml
# docker-compose.yml
ports:
  - "8080:8080"
  - "50051:50051"
```

---

### 20. Docker Build Fails

**Error**: Docker build fails in CI

**Common Issues**:

#### A. Go Version Mismatch
```dockerfile
# Dockerfile must match CI
FROM golang:1.19-alpine AS builder  # Not 1.21 or 1.24
```

#### B. Missing Dependencies
```dockerfile
# Add build tools
RUN apk add --no-cache git make
```

#### C. Context Issues
```bash
# Build from correct directory
docker build -t iam-service -f Dockerfile .
```

---

### 21. SSH Deployment Fails

**Error**: SSH connection refused or key rejected

**Troubleshooting**:

#### A. Test SSH Connection
```bash
# Test manually
ssh -i ~/.ssh/github_staging deploy@staging.example.com

# Check key format (should be OpenSSH format)
cat ~/.ssh/github_staging | head -1
# Should show: -----BEGIN OPENSSH PRIVATE KEY-----
```

#### B. Verify GitHub Secret
- Go to: Settings → Secrets → Actions
- Check `STAGING_SSH_KEY` includes full key with headers
- Should start with `-----BEGIN OPENSSH PRIVATE KEY-----`
- Should end with `-----END OPENSSH PRIVATE KEY-----`

#### C. Check Server SSH Config
```bash
# On server
sudo nano /etc/ssh/sshd_config

# Ensure these are set:
PubkeyAuthentication yes
PasswordAuthentication no
```

---

### 22. Docker Pull Fails on Server

**Error**: "pull access denied" or "manifest unknown"

**Fix**:

#### A. Login to Docker Hub
```bash
# On server
docker login -u your-username

# Enter Docker Hub access token (not password)
```

#### B. Verify Image Exists
```bash
# Check Docker Hub
docker pull your-username/iam-service:latest

# List local images
docker images | grep iam-service
```

#### C. Check docker-compose.yml
```yaml
services:
  iam-service:
    image: ${DOCKER_USERNAME}/iam-service:latest  # Correct format
```

---

## Best Practices

### 23. CI/CD Checklist

Before pushing to GitHub:

- [ ] Run `go fmt ./...` to format code
- [ ] Run `go mod tidy` to clean dependencies
- [ ] Run `.\scripts\lint.ps1` to check linting
- [ ] Run `go build ./...` to verify compilation
- [ ] Run `go test ./...` to run tests
- [ ] Create `.env.example` if missing
- [ ] Update migration steps in workflows if new migrations added
- [ ] Check `go.mod` has `go 1.19`
- [ ] Check `Dockerfile` uses `golang:1.19-alpine`

**All-in-one command**:
```bash
.\scripts\check-all.ps1    # Windows
./scripts/check-all.sh     # Linux/Mac
```

---

### 24. Local CI Testing

Test workflow locally before pushing:

```bash
# Install act (GitHub Actions locally)
# See: https://github.com/nektos/act

# Run specific job
act -j lint
act -j test
act -j build

# Run entire workflow
act push
```

**Docker Compose Test**:
```bash
# Test local build
docker-compose -f docker-compose.yml up --build

# Verify health
curl http://localhost:8080/health
```

---

### 25. Debugging Failed Workflows

**Step-by-step debugging**:

1. **View Workflow Logs**
   - Go to: GitHub → Actions → Failed workflow
   - Click on failed job
   - Expand failed step

2. **Check Working Directory**
   ```yaml
   - name: Debug
     run: |
       pwd
       ls -la
       echo $GITHUB_WORKSPACE
   ```

3. **Check Environment Variables**
   ```yaml
   - name: Debug env
     run: |
       echo "SERVICE_DIR: ${{ env.SERVICE_DIR }}"
       env | sort
   ```

4. **Re-run with Debug Logging**
   - Re-run workflow
   - Enable "Enable debug logging"

5. **Test Locally**
   - Clone fresh copy
   - Follow exact steps from workflow
   - Identify differences

---

### 26. Quick Fix Scripts

**Create .env.example**:
```powershell
# Windows
Copy-Item .env.template .env.example

# Linux/Mac
cp .env.template .env.example
```

**Install All Tools**:
```powershell
# Windows
.\scripts\setup-ci.ps1

# Linux/Mac
./scripts/setup-ci.sh
```

**Quick Commit & Push (for CI testing)**:
```powershell
# Windows
.\scripts\quick-commit.ps1

# Linux/Mac
./scripts/quick-commit.sh
```

---

## Summary of All Fixes

### ✅ Completed Fixes

#### Workflow Configuration
- [x] Upgraded artifact actions to v4
- [x] Added auto-detect for service directory
- [x] Added PostgreSQL client installation
- [x] Added all 6 migrations to workflows
- [x] Fixed Go version consistency (1.19)
- [x] Made Codecov upload optional

#### Code Quality (Linting)
- [x] Fixed `errcheck` - rows.Close() in DAOs (idiomatic `defer`)
- [x] Fixed `errcheck` - Logger.Sync() and Shutdown() checks
- [x] Added package comments to all packages
- [x] Added comments to all exported vars (errors)
- [x] Added comments to all exported constants
- [x] Added comments to all exported methods (getters)
- [x] Extracted test string constants (`goconst`)
- [x] Fixed integer overflow with safe conversion helpers (`gosec G115`)

#### Documentation
- [x] Created `.env.example` from template
- [x] Updated README.md with consolidated documentation
- [x] Updated fix_error_ci_cd.md (this file)
- [x] Added LINTING_SETUP.md
- [x] Added INSTALLATION_GUIDE.md

#### Features
- [x] Integrated Swagger UI with Basic Auth
- [x] Migrated to Gin web framework
- [x] Separated User/App and CMS authorization
- [x] Added health check endpoint

### 📊 Files Modified Summary

**Workflows**:
- `.github/workflows/ci-cd.yml`
- `.github/workflows/test.yml`

**Code Quality**:
- `internal/dao/*.go` (6 files - errcheck fixes)
- `internal/app/app.go` (error handling)
- `internal/container/container.go` (error handling)
- `internal/domain/**/*.go` (comments on all exports)
- `internal/handler/*.go` (overflow handling)
- `pkg/*/***_test.go` (string constants)

**Configuration**:
- `go.mod` (version 1.19)
- `Dockerfile` (golang:1.19-alpine)
- `.golangci.yml` (linter config)

**Scripts**:
- `scripts/install-golangci-lint.ps1`
- `scripts/lint.ps1`
- `scripts/check-all.ps1`
- `scripts/setup-ci.ps1`

---

## Quick Reference: Comment Patterns

### Package Comment
```go
// Package <name> provides <description>.
package <name>
```

### Exported Error Variables
```go
var (
    // ErrInvalidUsername indicates username validation failed
    ErrInvalidUsername = errors.New("invalid username")
)
```

### Exported Constants
```go
const (
    // DomainUser is the domain for end user authorization
    DomainUser CasbinDomain = "user"
)
```

### Exported Methods
```go
// ID returns the user's unique identifier
func (u *User) ID() string { return u.id }
```

**Important**: Use "ID" not "Id" for ID getter methods

---

## Contact & Support

If issues persist after trying these fixes:

1. **Check Documentation**:
   - [README.md](README.md) - Main documentation
   - [LINTING_SETUP.md](LINTING_SETUP.md) - Linting details
   - [CI_CD_SETUP_GUIDE.md](CI_CD_SETUP_GUIDE.md) - Full CI/CD guide

2. **Review Workflow Logs**:
   - GitHub → Actions → Click failed workflow
   - Review each failed step
   - Look for specific error messages

3. **Local Reproduction**:
   - Try reproducing issue locally
   - Run same commands as workflow
   - Compare environments

4. **Create Issue**:
   - Provide workflow run URL
   - Include relevant logs
   - Describe steps to reproduce

---

**Version**: 2.0  
**Last Updated**: November 2024  
**Status**: ✅ CI/CD Fully Operational
