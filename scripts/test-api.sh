#!/bin/bash

# IAM Service API Test Script
# This script tests the basic functionality of the IAM service

SERVER="localhost:50051"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "========================================="
echo "IAM Service API Test"
echo "========================================="

# Check if grpcurl is installed
if ! command -v grpcurl &> /dev/null; then
    echo -e "${RED}Error: grpcurl is not installed${NC}"
    echo "Install it with: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    exit 1
fi

# Check if service is running
echo ""
echo "Checking if service is running..."
grpcurl -plaintext $SERVER list > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Service is not running on $SERVER${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Service is running${NC}"

# Test 1: Register a user
echo ""
echo "Test 1: Register a new user"
REGISTER_RESPONSE=$(grpcurl -plaintext -d '{
  "username": "testuser_'$(date +%s)'",
  "email": "test_'$(date +%s)'@example.com",
  "password": "password123",
  "full_name": "Test User"
}' $SERVER iam.IAMService/Register 2>&1)

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Registration successful${NC}"
    echo "$REGISTER_RESPONSE"
else
    echo -e "${RED}✗ Registration failed${NC}"
    echo "$REGISTER_RESPONSE"
fi

# Test 2: Login
echo ""
echo "Test 2: Login with test user"
LOGIN_RESPONSE=$(grpcurl -plaintext -d '{
  "username": "testuser",
  "password": "password123"
}' $SERVER iam.IAMService/Login 2>&1)

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Login successful${NC}"
    # Extract token (basic extraction, may need adjustment)
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token": "[^"]*' | grep -o '[^"]*$')
    echo "Access Token: ${TOKEN:0:50}..."
else
    echo -e "${YELLOW}⚠ Login failed (user may not exist)${NC}"
    echo "$LOGIN_RESPONSE"
fi

# Test 3: Verify Token
if [ ! -z "$TOKEN" ]; then
    echo ""
    echo "Test 3: Verify Token"
    VERIFY_RESPONSE=$(grpcurl -plaintext -d "{
      \"token\": \"$TOKEN\"
    }" $SERVER iam.IAMService/VerifyToken 2>&1)
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Token verification successful${NC}"
        echo "$VERIFY_RESPONSE"
    else
        echo -e "${RED}✗ Token verification failed${NC}"
        echo "$VERIFY_RESPONSE"
    fi
fi

# Test 4: List Roles
echo ""
echo "Test 4: List Roles"
ROLES_RESPONSE=$(grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' $SERVER iam.IAMService/ListRoles 2>&1)

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ List roles successful${NC}"
    echo "$ROLES_RESPONSE"
else
    echo -e "${RED}✗ List roles failed${NC}"
    echo "$ROLES_RESPONSE"
fi

# Test 5: List Permissions
echo ""
echo "Test 5: List Permissions"
PERMS_RESPONSE=$(grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' $SERVER iam.IAMService/ListPermissions 2>&1)

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ List permissions successful${NC}"
    echo "$PERMS_RESPONSE"
else
    echo -e "${RED}✗ List permissions failed${NC}"
    echo "$PERMS_RESPONSE"
fi

echo ""
echo "========================================="
echo "API Test Completed"
echo "========================================="

