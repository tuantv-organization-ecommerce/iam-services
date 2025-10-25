# Casbin Authorization - RBAC Model

## Tổng quan

IAM Service sử dụng **Casbin** để implement mô hình phân quyền RBAC (Role-Based Access Control) với hỗ trợ **multi-domain** authorization.

### Tính năng chính

✅ **Phân quyền theo Domain**:
- `user` domain: Phân quyền cho end users
- `cms` domain: Phân quyền truy cập CMS admin panel
- `api` domain: Phân quyền truy cập API endpoints

✅ **CMS Role Management**:
- Quản lý roles riêng cho CMS
- Phân quyền theo tabs (product, inventory, order, report, ...)
- Flexible permission assignment

✅ **API Resource Management**:
- Track API endpoints (path + method)
- Phân quyền chi tiết cho từng API
- Support wildcard patterns

## Kiến trúc

### Casbin RBAC Model

```
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
```

### Các thành phần:

- **sub** (subject): User ID hoặc Role name
- **dom** (domain): `user`, `cms`, hoặc `api`
- **obj** (object): Resource path (API path hoặc CMS tab)
- **act** (action): HTTP method hoặc operation

### Database Schema

#### 1. casbin_rule table
Lưu trữ tất cả policies và role assignments.

```sql
CREATE TABLE casbin_rule (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100),      -- 'p' (policy) hoặc 'g' (grouping/role)
    v0 VARCHAR(100),         -- subject/user/role
    v1 VARCHAR(100),         -- domain hoặc parent role
    v2 VARCHAR(100),         -- object/resource
    v3 VARCHAR(100),         -- action
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);
```

#### 2. cms_roles table
Quản lý CMS roles với tabs.

```sql
CREATE TABLE cms_roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    tabs TEXT[],             -- Array of CMS tabs
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

#### 3. api_resources table
Track API endpoints.

```sql
CREATE TABLE api_resources (
    id VARCHAR(36) PRIMARY KEY,
    path VARCHAR(500) NOT NULL,
    method VARCHAR(20) NOT NULL,
    service VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    UNIQUE(path, method)
);
```

## Các Domain Authorization

### 1. User Domain - End User Authorization

Phân quyền cho end users truy cập các resources thông thường.

**Example policies**:

```
# Role: user
p, user, user, /api/v1/products, GET
p, user, user, /api/v1/products/*, GET
p, user, user, /api/v1/orders, (GET|POST)

# Assign role to user
g, user-123, user, user
```

**Use case**:
- Customer browsing products
- Creating orders
- View own profile

### 2. CMS Domain - Admin Panel Authorization

Phân quyền cho admin/staff truy cập CMS panel theo tabs.

**CMS Tabs**:
- `product`: Quản lý sản phẩm
- `inventory`: Quản lý tồn kho
- `order`: Quản lý đơn hàng
- `user`: Quản lý người dùng
- `report`: Xem báo cáo
- `setting`: Cấu hình hệ thống

**Example policies**:

```
# CMS Admin - Full access
p, cms_admin, cms, /cms/product/*, (GET|POST|PUT|DELETE)
p, cms_admin, cms, /cms/inventory/*, (GET|POST|PUT|DELETE)
p, cms_admin, cms, /cms/order/*, (GET|POST|PUT|DELETE)

# CMS Product Manager - Only product & inventory
p, cms_product_manager, cms, /cms/product/*, (GET|POST|PUT|DELETE)
p, cms_product_manager, cms, /cms/inventory/*, (GET|POST|PUT)

# CMS Viewer - Read only
p, cms_viewer, cms, /cms/product/*, GET
p, cms_viewer, cms, /cms/inventory/*, GET
p, cms_viewer, cms, /cms/order/*, GET

# Assign CMS role to user
g, user-456, cms_admin, cms
```

**Use case**:
- Admin quản lý toàn bộ hệ thống
- Product manager chỉ quản lý products
- Staff xem reports nhưng không thể sửa

### 3. API Domain - API Access Control

Phân quyền chi tiết cho từng API endpoint.

**Example policies**:

```
# Admin - Full API access
p, admin, api, /api/v1/**, (GET|POST|PUT|DELETE)

# Moderator - Limited access
p, moderator, api, /api/v1/products, (GET|POST|PUT)
p, moderator, api, /api/v1/products/*, (GET|POST|PUT)

# Assign API role
g, user-789, admin, api
```

## API Usage

### 1. Check API Access

Kiểm tra user có quyền truy cập API không.

```bash
grpcurl -plaintext -d '{
  "user_id": "user-123",
  "api_path": "/api/v1/products",
  "method": "POST"
}' localhost:50051 iam.IAMService/CheckAPIAccess
```

**Response**:
```json
{
  "allowed": true,
  "message": "Access granted"
}
```

### 2. Check CMS Access

Kiểm tra user có quyền truy cập CMS tab không.

```bash
grpcurl -plaintext -d '{
  "user_id": "user-456",
  "cms_tab": "product",
  "action": "POST"
}' localhost:50051 iam.IAMService/CheckCMSAccess
```

**Response**:
```json
{
  "allowed": true,
  "message": "Access granted to CMS tab",
  "accessible_tabs": ["product", "inventory", "order"]
}
```

### 3. Enforce Policy (General)

General policy enforcement cho bất kỳ domain nào.

```bash
grpcurl -plaintext -d '{
  "user_id": "user-789",
  "domain": "api",
  "resource": "/api/v1/users/123",
  "action": "DELETE"
}' localhost:50051 iam.IAMService/EnforcePolicy
```

### 4. Create CMS Role

Tạo CMS role mới với tabs.

```bash
grpcurl -plaintext -d '{
  "name": "cms_content_editor",
  "description": "Content editor role",
  "tabs": ["product", "inventory"]
}' localhost:50051 iam.IAMService/CreateCMSRole
```

### 5. Assign CMS Role

Gán CMS role cho user.

```bash
grpcurl -plaintext -d '{
  "user_id": "user-123",
  "cms_role_id": "cms-role-001"
}' localhost:50051 iam.IAMService/AssignCMSRole
```

### 6. Get User CMS Tabs

Lấy danh sách tabs mà user có quyền truy cập.

```bash
grpcurl -plaintext -d '{
  "user_id": "user-123"
}' localhost:50051 iam.IAMService/GetUserCMSTabs
```

**Response**:
```json
{
  "tabs": ["product", "inventory", "report"]
}
```

### 7. Create API Resource

Đăng ký API endpoint mới.

```bash
grpcurl -plaintext -d '{
  "path": "/api/v1/products",
  "method": "POST",
  "service": "product-service",
  "description": "Create new product"
}' localhost:50051 iam.IAMService/CreateAPIResource
```

### 8. List API Resources

Liệt kê các API resources.

```bash
grpcurl -plaintext -d '{
  "service": "product-service"
}' localhost:50051 iam.IAMService/ListAPIResources
```

## Pattern Matching

Casbin hỗ trợ các pattern matching:

### KeyMatch2

Sử dụng cho resource paths:

- `/api/v1/products` matches `/api/v1/products`
- `/api/v1/products/*` matches `/api/v1/products/123`
- `/api/v1/**` matches tất cả sub-paths

### RegexMatch

Sử dụng cho actions:

- `GET` matches exact "GET"
- `(GET|POST)` matches "GET" OR "POST"
- `(GET|POST|PUT|DELETE)` matches tất cả CRUD operations

## Workflow Examples

### Workflow 1: Setup CMS Admin

```go
// 1. Create CMS admin role
CreateCMSRole("cms_admin", "Full CMS access", 
    []string{"product", "inventory", "order", "user", "report", "setting"})

// 2. Add policies for each tab
AddPolicy("cms_admin", "cms", "/cms/product/*", "(GET|POST|PUT|DELETE)")
AddPolicy("cms_admin", "cms", "/cms/inventory/*", "(GET|POST|PUT|DELETE)")
// ... more policies

// 3. Assign to user
AssignCMSRole("user-123", "cms-role-admin-001")

// 4. Check access
allowed := CheckCMSAccess("user-123", "product", "POST")  // true
```

### Workflow 2: Setup API Access

```go
// 1. Register API resources
CreateAPIResource("/api/v1/products", "POST", "product-service", "Create product")

// 2. Add policy for role
AddPolicy("product_manager", "api", "/api/v1/products", "POST")

// 3. Assign role to user in API domain
AssignUserRole("user-456", "product_manager", "api")

// 4. Check access
allowed := CheckAPIAccess("user-456", "/api/v1/products", "POST")  // true
```

### Workflow 3: Multi-Domain User

User có thể có roles ở nhiều domains:

```go
// User là end user
AssignUserRole("user-789", "user", "user")

// Đồng thời là CMS staff
AssignCMSRole("user-789", "cms-product-manager")

// Và có API access
AssignUserRole("user-789", "moderator", "api")

// Check different domains
CheckAPIAccess("user-789", "/api/v1/products", "GET")     // Check in API domain
CheckCMSAccess("user-789", "product", "POST")              // Check in CMS domain
```

## Best Practices

### 1. Least Privilege Principle
- Chỉ gán quyền tối thiểu cần thiết
- Sử dụng specific paths thay vì wildcards khi có thể

### 2. Role Hierarchy
- Tạo roles từ general → specific
- VD: `viewer` → `editor` → `admin`

### 3. Domain Separation
- Tách biệt rõ ràng giữa user/cms/api domains
- Không mix permissions giữa domains

### 4. Regular Audit
- Review policies định kỳ
- Remove unused roles và permissions

### 5. Testing
- Test authorization cho tất cả endpoints
- Verify negative cases (should deny)

## Performance Considerations

### 1. Policy Caching
Casbin tự động cache policies trong memory.

### 2. Database Indexes
Đã tạo indexes cho:
- `casbin_rule.ptype`
- `casbin_rule.v0` (subject)
- `casbin_rule.v1` (domain)

### 3. Policy Reload
```go
enforcer.LoadPolicy()  // Reload from database
```

## Troubleshooting

### Issue: Authorization always denies

**Giải pháp**:
1. Check user có role assignment không:
   ```sql
   SELECT * FROM casbin_rule WHERE ptype = 'g' AND v0 = 'user-id';
   ```

2. Check policies cho role:
   ```sql
   SELECT * FROM casbin_rule WHERE ptype = 'p' AND v0 = 'role-name';
   ```

3. Verify domain matching

### Issue: Wildcard không hoạt động

**Giải pháp**:
- Sử dụng `keyMatch2` trong matcher
- Đảm bảo pattern đúng: `/api/v1/*` hoặc `/api/v1/**`

## Migration Guide

### Chạy migrations:

```bash
# Schema
psql -U postgres -d iam_db -f migrations/003_casbin_tables.sql

# Seed data (default policies)
psql -U postgres -d iam_db -f migrations/004_casbin_seed_data.sql
```

### Verify:

```sql
-- Check policies
SELECT COUNT(*) FROM casbin_rule;

-- Check CMS roles
SELECT * FROM cms_roles;

-- Check API resources
SELECT * FROM api_resources;
```

## Tài liệu tham khảo

- [Casbin Documentation](https://casbin.org/)
- [RBAC Model](https://casbin.org/docs/rbac)
- [Pattern Matching](https://casbin.org/docs/function)

