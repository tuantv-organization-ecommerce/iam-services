# Casbin Integration - Quick Start

IAM Service đã được tích hợp **Casbin RBAC** với multi-domain authorization support.

## 🚀 Quick Start

### 1. Chạy migrations

```bash
# Di chuyển đến thư mục project
cd ecommerce/back_end/iam-services

# Chạy migration cho Casbin tables
psql -U postgres -d iam_db -f migrations/003_casbin_tables.sql

# Seed default data (roles, policies)
psql -U postgres -d iam_db -f migrations/004_casbin_seed_data.sql
```

### 2. Download dependencies

```bash
go mod download
go mod tidy
```

### 3. Generate proto files

```bash
make proto
# Hoặc
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       pkg/proto/iam.proto
```

### 4. Chạy service

```bash
go run cmd/server/main.go
```

## 📋 Các tính năng mới

### 1. Multi-Domain Authorization

- **User Domain**: Phân quyền cho end users
- **CMS Domain**: Phân quyền CMS admin panel theo tabs
- **API Domain**: Phân quyền API endpoints

### 2. CMS Role Management

Quản lý roles cho CMS với phân quyền theo tabs:
- `product`: Quản lý sản phẩm
- `inventory`: Quản lý tồn kho
- `order`: Quản lý đơn hàng
- `user`: Quản lý người dùng
- `report`: Xem báo cáo
- `setting`: Cấu hình

### 3. API Resource Tracking

Track và phân quyền cho từng API endpoint (path + method).

## 🔧 New APIs

### Authorization APIs

```bash
# Check API access
CheckAPIAccess(user_id, api_path, method)

# Check CMS access
CheckCMSAccess(user_id, cms_tab, action)

# General policy enforcement
EnforcePolicy(user_id, domain, resource, action)
```

### CMS Management APIs

```bash
# Create CMS role
CreateCMSRole(name, description, tabs[])

# Assign CMS role to user
AssignCMSRole(user_id, cms_role_id)

# Remove CMS role from user
RemoveCMSRole(user_id, cms_role_id)

# Get user's accessible CMS tabs
GetUserCMSTabs(user_id)

# List all CMS roles
ListCMSRoles(page, page_size)
```

### API Resource APIs

```bash
# Create API resource
CreateAPIResource(path, method, service, description)

# List API resources
ListAPIResources(service, page, page_size)
```

## 📝 Example Usage

### Example 1: Check CMS Access

```bash
grpcurl -plaintext -d '{
  "user_id": "user-123",
  "cms_tab": "product",
  "action": "POST"
}' localhost:50051 iam.IAMService/CheckCMSAccess
```

Response:
```json
{
  "allowed": true,
  "message": "Access granted to CMS tab",
  "accessible_tabs": ["product", "inventory", "order"]
}
```

### Example 2: Check API Access

```bash
grpcurl -plaintext -d '{
  "user_id": "user-456",
  "api_path": "/api/v1/products",
  "method": "POST"
}' localhost:50051 iam.IAMService/CheckAPIAccess
```

Response:
```json
{
  "allowed": false,
  "message": "Access denied"
}
```

### Example 3: Create CMS Role

```bash
grpcurl -plaintext -d '{
  "name": "cms_product_editor",
  "description": "Product content editor",
  "tabs": ["product", "inventory"]
}' localhost:50051 iam.IAMService/CreateCMSRole
```

## 🗂️ Database Schema Changes

### New Tables:

1. **casbin_rule**: Stores Casbin policies and role assignments
2. **cms_roles**: CMS roles with tabs
3. **user_cms_roles**: User-CMS role relationships
4. **api_resources**: API endpoint definitions

### Updated Tables:

- **roles**: Added `domain` column

## 📁 New Files Structure

```
iam-services/
├── configs/
│   └── rbac_model.conf          # Casbin RBAC model
├── internal/
│   ├── dao/
│   │   ├── api_resource_dao.go  # API resource DAO
│   │   ├── cms_role_dao.go      # CMS role DAO
│   │   └── user_cms_role_dao.go # User-CMS role DAO
│   ├── domain/
│   │   └── casbin.go            # Casbin domain models
│   ├── handler/
│   │   └── casbin_handler.go    # Casbin gRPC handlers
│   ├── repository/
│   │   ├── api_resource_repository.go
│   │   └── cms_repository.go
│   └── service/
│       └── casbin_service.go    # Casbin business logic
├── pkg/
│   └── casbin/
│       └── enforcer.go          # Casbin enforcer wrapper
├── migrations/
│   ├── 003_casbin_tables.sql    # Schema migration
│   └── 004_casbin_seed_data.sql # Seed data
└── docs/
    └── CASBIN.md                # Detailed documentation
```

## 🎯 Use Cases

### Use Case 1: CMS Admin Setup

1. Tạo CMS admin role với full access
2. Assign role cho admin user
3. Admin có thể truy cập tất cả CMS tabs

### Use Case 2: Product Manager

1. Tạo CMS product manager role
2. Chỉ có quyền với `product` và `inventory` tabs
3. Không thể truy cập `order`, `user`, hoặc `setting`

### Use Case 3: API Gateway Integration

1. Đăng ký tất cả API endpoints
2. Mỗi request check authorization
3. Allow/deny based on user roles

## 📖 Documentation

Xem thêm chi tiết tại:
- [docs/CASBIN.md](docs/CASBIN.md) - Hướng dẫn chi tiết
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) - Kiến trúc tổng quan
- [docs/API.md](docs/API.md) - API documentation

## 🔐 Default Roles & Policies

### User Roles (user domain)
- **admin**: Full access
- **user**: Read products, create orders
- **moderator**: Moderate content

### CMS Roles (cms domain)
- **cms_admin**: Full CMS access
- **cms_product_manager**: Product & Inventory
- **cms_order_manager**: Order & Report
- **cms_content_editor**: Product editing
- **cms_viewer**: Read-only

### API Policies (api domain)
- Admin: All APIs
- User: Limited read access
- Moderator: Moderate access

## ⚙️ Configuration

Casbin model file: `configs/rbac_model.conf`

```conf
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
```

## 🐛 Troubleshooting

### Problem: "casbin_rule table not found"
**Solution**: Chạy migration `003_casbin_tables.sql`

### Problem: "No policies loaded"
**Solution**: Chạy seed data `004_casbin_seed_data.sql`

### Problem: "Access always denied"
**Solution**: 
1. Check user has role assignment
2. Check role has correct policies
3. Verify domain matching

## 🤝 Contributing

Khi thêm tính năng mới:
1. Đăng ký API resources trong `api_resources` table
2. Tạo policies trong `casbin_rule` table
3. Test authorization với các roles khác nhau

## 📞 Support

Có vấn đề? Xem:
- [docs/CASBIN.md](docs/CASBIN.md) - Detailed guide
- [docs/SETUP.md](docs/SETUP.md) - Setup instructions

