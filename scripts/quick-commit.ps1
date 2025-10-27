# Quick Commit & Push Script for Windows PowerShell
# This script commits and pushes all CI/CD setup changes

$ErrorActionPreference = "Stop"

Write-Host "🚀 Quick Commit & Push for CI/CD Setup" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# Check if we're in a git repository
if (-not (Test-Path ".git")) {
    Write-Host "❌ Error: Not a git repository" -ForegroundColor Red
    exit 1
}

# Show current status
Write-Host "📊 Current Git Status:" -ForegroundColor Yellow
git status --short
Write-Host ""

# Confirm with user
$confirm = Read-Host "Do you want to commit and push these changes? (y/n)"
if ($confirm -ne "y" -and $confirm -ne "Y") {
    Write-Host "⏹️  Aborted by user" -ForegroundColor Yellow
    exit 0
}

# Stage all changes
Write-Host ""
Write-Host "📦 Staging changes..." -ForegroundColor Yellow
git add .
Write-Host "✅ Changes staged" -ForegroundColor Green

# Create commit
$commitMsg = @"
ci: setup CI/CD pipeline with basic tests

Features:
- Add test files for jwt, password, dao, and service layers (25 tests)
- Comment out deploy jobs temporarily (no servers yet)
- Fix Dockerfile Go version to 1.19
- Add migrations 005, 006 to CI workflow
- Add testify test dependency to go.mod
- Create CI documentation and helper scripts

Changes:
- pkg/jwt/jwt_manager_test.go: JWT token generation and verification tests
- pkg/password/password_manager_test.go: Password hashing tests
- internal/dao/user_dao_test.go: User DAO CRUD tests
- internal/service/auth_service_test.go: Auth service business logic tests
- .github/workflows/ci-cd.yml: Disabled deploy jobs, added migrations
- .github/workflows/test.yml: Added migrations 005, 006
- Dockerfile: Fixed Go version from 1.21 to 1.19
- go.mod: Added github.com/stretchr/testify v1.8.4

Docs:
- CI_QUICK_START.md: Quick start guide for CI/CD
- scripts/setup-ci.sh: Automated setup helper
- scripts/quick-commit.sh: Quick commit script
- .env.template: Environment variable template

Status: Ready for CI testing ✅
"@

Write-Host ""
Write-Host "📝 Creating commit..." -ForegroundColor Yellow
git commit -m $commitMsg
Write-Host "✅ Commit created" -ForegroundColor Green

# Get current branch
$currentBranch = git branch --show-current
Write-Host ""
Write-Host "📍 Current branch: $currentBranch" -ForegroundColor Cyan

# Ask if user wants to create a new branch
if ($currentBranch -eq "main" -or $currentBranch -eq "master" -or $currentBranch -eq "develop") {
    Write-Host ""
    Write-Host "⚠️  You're on $currentBranch branch" -ForegroundColor Yellow
    $createBranch = Read-Host "Do you want to create a new feature branch? (y/n)"
    if ($createBranch -eq "y" -or $createBranch -eq "Y") {
        $branchName = "feature/setup-cicd"
        Write-Host "🌿 Creating branch: $branchName" -ForegroundColor Cyan
        git checkout -b $branchName
        Write-Host "✅ Branch created and switched" -ForegroundColor Green
        $currentBranch = $branchName
    }
}

# Push to remote
Write-Host ""
Write-Host "🚀 Pushing to origin/$currentBranch..." -ForegroundColor Yellow
git push origin $currentBranch

Write-Host ""
Write-Host "======================================" -ForegroundColor Green
Write-Host "✅ Successfully Pushed!" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green
Write-Host ""
Write-Host "📊 Next Steps:" -ForegroundColor Cyan
Write-Host ""
Write-Host "1. 🌐 Go to GitHub Actions:"
Write-Host "   https://github.com/YOUR_USERNAME/YOUR_REPO/actions"
Write-Host ""
Write-Host "2. ⏱️  Wait for CI pipeline (10-15 minutes):" -ForegroundColor Yellow
Write-Host "   ✅ Lint      - Code quality checks"
Write-Host "   ✅ Test      - Run 25 tests with PostgreSQL"
Write-Host "   ✅ Build     - Build binary"
Write-Host "   ✅ Security  - Vulnerability scanning"
Write-Host ""
Write-Host "3. 📝 If all jobs pass:"
Write-Host "   - Create Pull Request to develop/main"
Write-Host "   - Review and merge"
Write-Host ""
Write-Host "4. ⚠️  If any job fails:"
Write-Host "   - Check logs in GitHub Actions"
Write-Host "   - Fix issues locally"
Write-Host "   - Commit and push again"
Write-Host ""
Write-Host "🔗 Quick Links:" -ForegroundColor Cyan
Write-Host "   Actions: https://github.com/YOUR_USERNAME/YOUR_REPO/actions"
Write-Host "   Branch:  https://github.com/YOUR_USERNAME/YOUR_REPO/tree/$currentBranch"
Write-Host ""
Write-Host "Happy coding! 🎉" -ForegroundColor Green
Write-Host ""

