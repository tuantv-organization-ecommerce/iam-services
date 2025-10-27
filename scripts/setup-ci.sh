#!/bin/bash

# CI/CD Setup Helper Script
# This script prepares the project for CI/CD pipeline

set -e

echo "🚀 IAM Service CI/CD Setup"
echo "======================================"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Error: go.mod not found. Please run this script from iam-services directory.${NC}"
    exit 1
fi

echo "📍 Current directory: $(pwd)"
echo ""

# Step 1: Create .env.example if not exists
echo "Step 1: Creating .env.example..."
if [ ! -f ".env.example" ]; then
    cp .env.template .env.example 2>/dev/null || {
        cat > .env.example << 'EOF'
# IAM Service Environment Variables
SERVER_HOST=0.0.0.0
SERVER_PORT=50051
HTTP_HOST=0.0.0.0
HTTP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=iam_db
DB_SSL_MODE=disable
JWT_SECRET=CHANGE-THIS-TO-64-CHAR-RANDOM-STRING
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_EXPIRATION_HOURS=168
CASBIN_MODEL_PATH=./configs/rbac_model.conf
LOG_LEVEL=info
LOG_ENCODING=json
EOF
    }
    echo -e "${GREEN}✅ .env.example created${NC}"
else
    echo -e "${YELLOW}⚠️  .env.example already exists${NC}"
fi
echo ""

# Step 2: Download Go dependencies
echo "Step 2: Downloading Go dependencies..."
go mod download
echo -e "${GREEN}✅ Dependencies downloaded${NC}"
echo ""

echo "Step 3: Tidying Go modules..."
go mod tidy
echo -e "${GREEN}✅ Go modules tidied${NC}"
echo ""

# Step 4: Run local tests (optional)
echo "Step 4: Running local tests..."
read -p "Do you want to run tests locally? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Running tests that don't require database..."
    go test -v ./pkg/jwt/ || echo -e "${YELLOW}⚠️  JWT tests failed (expected if dependencies missing)${NC}"
    go test -v ./pkg/password/ || echo -e "${YELLOW}⚠️  Password tests failed (expected if dependencies missing)${NC}"
    echo -e "${GREEN}✅ Local tests completed${NC}"
else
    echo -e "${YELLOW}⏭️  Skipping local tests${NC}"
fi
echo ""

# Step 5: Check Git status
echo "Step 5: Checking Git status..."
echo ""
git status --short
echo ""

# Step 6: Prepare commit
echo "======================================"
echo -e "${GREEN}✅ Setup Complete!${NC}"
echo ""
echo "📝 Next steps:"
echo "   1. Review the changes: git status"
echo "   2. Stage changes: git add ."
echo "   3. Commit: git commit -m \"ci: setup CI/CD pipeline with basic tests\""
echo "   4. Create branch: git checkout -b feature/setup-cicd"
echo "   5. Push: git push origin feature/setup-cicd"
echo ""
echo "Or run the quick commit script:"
echo "   ./scripts/quick-commit.sh"
echo ""
echo -e "${YELLOW}⚠️  Remember to check GitHub Actions after pushing!${NC}"
echo ""

