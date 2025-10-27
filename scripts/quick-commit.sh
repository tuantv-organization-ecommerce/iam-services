#!/bin/bash

# Quick Commit & Push Script for CI/CD Setup
# This script commits and pushes all CI/CD setup changes

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 Quick Commit & Push for CI/CD Setup${NC}"
echo "======================================"
echo ""

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo -e "${RED}❌ Error: Not a git repository${NC}"
    exit 1
fi

# Show current status
echo "📊 Current Git Status:"
git status --short
echo ""

# Confirm with user
read -p "Do you want to commit and push these changes? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}⏹️  Aborted by user${NC}"
    exit 0
fi

# Stage all changes
echo ""
echo "📦 Staging changes..."
git add .
echo -e "${GREEN}✅ Changes staged${NC}"

# Create commit
COMMIT_MSG="ci: setup CI/CD pipeline with basic tests

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

Status: Ready for CI testing ✅"

echo ""
echo "📝 Creating commit..."
git commit -m "$COMMIT_MSG"
echo -e "${GREEN}✅ Commit created${NC}"

# Get current branch
CURRENT_BRANCH=$(git branch --show-current)
echo ""
echo "📍 Current branch: $CURRENT_BRANCH"

# Ask if user wants to create a new branch
if [ "$CURRENT_BRANCH" = "main" ] || [ "$CURRENT_BRANCH" = "master" ] || [ "$CURRENT_BRANCH" = "develop" ]; then
    echo ""
    echo -e "${YELLOW}⚠️  You're on $CURRENT_BRANCH branch${NC}"
    read -p "Do you want to create a new feature branch? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        BRANCH_NAME="feature/setup-cicd"
        echo "🌿 Creating branch: $BRANCH_NAME"
        git checkout -b $BRANCH_NAME
        echo -e "${GREEN}✅ Branch created and switched${NC}"
        CURRENT_BRANCH=$BRANCH_NAME
    fi
fi

# Push to remote
echo ""
echo "🚀 Pushing to origin/$CURRENT_BRANCH..."
git push origin $CURRENT_BRANCH

echo ""
echo -e "${GREEN}======================================"
echo "✅ Successfully Pushed!${NC}"
echo "======================================"
echo ""
echo "📊 Next Steps:"
echo ""
echo "1. 🌐 Go to GitHub Actions:"
echo "   https://github.com/YOUR_USERNAME/YOUR_REPO/actions"
echo ""
echo "2. ⏱️  Wait for CI pipeline (10-15 minutes):"
echo "   ✅ Lint      - Code quality checks"
echo "   ✅ Test      - Run 25 tests with PostgreSQL"
echo "   ✅ Build     - Build binary"
echo "   ✅ Security  - Vulnerability scanning"
echo ""
echo "3. 📝 If all jobs pass:"
echo "   - Create Pull Request to develop/main"
echo "   - Review and merge"
echo ""
echo "4. ⚠️  If any job fails:"
echo "   - Check logs in GitHub Actions"
echo "   - Fix issues locally"
echo "   - Commit and push again"
echo ""
echo -e "${BLUE}🔗 Quick Links:${NC}"
echo "   Actions: https://github.com/YOUR_USERNAME/YOUR_REPO/actions"
echo "   Branch:  https://github.com/YOUR_USERNAME/YOUR_REPO/tree/$CURRENT_BRANCH"
echo ""
echo -e "${GREEN}Happy coding! 🎉${NC}"
echo ""

