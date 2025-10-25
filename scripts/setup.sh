#!/bin/bash

echo "========================================="
echo "IAM Service Setup Script"
echo "========================================="

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go 1.21 or higher"
    exit 1
fi

echo -e "${GREEN}✓ Go is installed: $(go version)${NC}"

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo -e "${RED}Error: PostgreSQL is not installed${NC}"
    echo "Please install PostgreSQL 12 or higher"
    exit 1
fi

echo -e "${GREEN}✓ PostgreSQL is installed: $(psql --version)${NC}"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo -e "${YELLOW}Warning: protoc is not installed${NC}"
    echo "You'll need it to regenerate Protocol Buffer files"
    echo "Install from: https://github.com/protocolbuffers/protobuf/releases"
else
    echo -e "${GREEN}✓ protoc is installed: $(protoc --version)${NC}"
fi

# Install Go dependencies
echo ""
echo "Installing Go dependencies..."
go mod download
go mod tidy
echo -e "${GREEN}✓ Dependencies installed${NC}"

# Create .env file if not exists
if [ ! -f .env ]; then
    echo ""
    echo "Creating .env file from template..."
    cp .env.example .env
    echo -e "${GREEN}✓ .env file created${NC}"
    echo -e "${YELLOW}⚠ Please update .env with your configuration${NC}"
else
    echo -e "${GREEN}✓ .env file already exists${NC}"
fi

# Ask to create database
echo ""
read -p "Do you want to create the database now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Creating database..."
    psql -U postgres -c "CREATE DATABASE iam_db;" 2>/dev/null
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Database created${NC}"
    else
        echo -e "${YELLOW}⚠ Database might already exist or PostgreSQL is not running${NC}"
    fi
    
    echo "Running migrations..."
    psql -U postgres -d iam_db -f migrations/001_init_schema.sql
    psql -U postgres -d iam_db -f migrations/002_seed_data.sql
    echo -e "${GREEN}✓ Migrations completed${NC}"
fi

# Generate protobuf files
echo ""
read -p "Do you want to generate protobuf files? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if command -v protoc &> /dev/null; then
        echo "Generating protobuf files..."
        protoc --go_out=. --go_opt=paths=source_relative \
               --go-grpc_out=. --go-grpc_opt=paths=source_relative \
               pkg/proto/iam.proto
        echo -e "${GREEN}✓ Protobuf files generated${NC}"
    else
        echo -e "${RED}Error: protoc is not installed${NC}"
    fi
fi

# Build the application
echo ""
read -p "Do you want to build the application? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Building application..."
    mkdir -p bin
    go build -o bin/iam-service cmd/server/main.go
    echo -e "${GREEN}✓ Build completed: bin/iam-service${NC}"
fi

echo ""
echo "========================================="
echo -e "${GREEN}Setup completed!${NC}"
echo "========================================="
echo ""
echo "Next steps:"
echo "1. Update .env file with your configuration"
echo "2. Run the service: go run cmd/server/main.go"
echo "   or use the binary: ./bin/iam-service"
echo ""
echo "For more information, see README.md"

