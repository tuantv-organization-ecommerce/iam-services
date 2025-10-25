# Authorization Architecture - Separated User/App and CMS

## 🎯 Overview

IAM Service sử dụng **2 hệ thống authorization độc lập** với 2 Casbin models riêng biệt:

### 1. User/App Authorization (`roles` table)
- **Purpose**: Phân quyền cho end users truy cập web/app
- **Model**: `rbac_user_model.conf`
- **Database**: `casbin_rule_user`
- **Domains**: `user`, `api`
- **Use cases**: 
  - User browsing products
  - User creating orders
  - User managing profile
  - API access control

### 2. CMS Authorization (`cms_roles` table)
- **Purpose**: Phân quyền cho admin/staff trên CMS
- **Model**: `rbac_cms_model.conf`
- **Database**: `casbin_rule_cms`
- **Structure**: Tab-based authorization
- **Use cases**:
  - Admin managing products
  - Staff managing inventory
  - Manager viewing reports
  - Tab-level access control

---

## 🏗️ Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                         IAM SERVICE                                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌────────────────────────┐      ┌──────────────────────────┐      │
│  │  User/App Authorization│      │   CMS Authorization      │      │
│  │  (End Users)           │      │   (Admin/Staff)          │      │
│  └────────────────────────┘      └──────────────────────────┘      │
│           │                                   │                      │
│           ▼                                   ▼                      │
│  ┌────────────────────────┐      ┌──────────────────────────┐      │
│  │  rbac_user_model.conf  │      │  rbac_cms_model.conf     │      │
│  │  - Domains: user, api  │      │  - Tab-based             │      │
│  │  - Resource paths      │      │  - Tab → APIs            │      │
│  │  - HTTP methods        │      │  - HTTP methods          │      │
│  └────────────────────────┘      └──────────────────────────┘      │
│           │                                   │                      │
│           ▼                                   ▼                      │
│  ┌────────────────────────┐      ┌──────────────────────────┐      │
│  │ casbin_rule_user       │      │  casbin_rule_cms         │      │
│  │ (PostgreSQL)           │      │  (PostgreSQL)            │      │
│  └────────────────────────┘      └──────────────────────────┘      │
│           │                                   │                      │
│           ▼                                   ▼                      │
│  ┌────────────────────────┐      ┌──────────────────────────┐      │
│  │  roles                 │      │  cms_roles               │      │
│  │  - user                │      │  - cms_admin             │      │
│  │  - premium_user        │      │  - cms_product_manager   │      │
│  │  - api_admin           │      │  - cms_order_manager     │      │
│  └────────────────────────┘      │  - cms_viewer            │      │
│                                   └──────────────────────────┘      │
│                                            │                         │
│                                            ▼                         │
│                                   ┌──────────────────────────┐      │
│                                   │  cms_tab_apis            │      │
│                                   │  Maps: Tab → APIs        │      │
│                                   │  (Many-to-Many)          │      │
│                                   └──────────────────────────┘      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 📊 Database Schema

### User/App Authorization Tables

#### `roles` table
```sql
CREATE TABLE roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

**Default Roles**:
- `user`: Regular user (browse products, create orders)
- `premium_user`: Premium user (more API access)
- `api_admin`: Full API access

#### `casbin_rule_user` table
```sql
CREATE TABLE casbin_rule_user (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100),      -- 'p' or 'g'
    v0 VARCHAR(100),         -- subject (user_id or role_name)
    v1 VARCHAR(100),         -- domain ('user' or 'api')
    v2 VARCHAR(100),         -- resource path
    v3 VARCHAR(100),         -- action (GET, POST, etc)
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);
```

**Example Policies**:
```
p, user, user, /api/v1/products, GET
p, user, user, /api/v1/orders, (GET|POST)
p, premium_user, api, /api/v1/products/**, (GET|POST)
```

**Example Role Assignments**:
```
g, user-123, user, user
g, user-456, premium_user, api
```

---

### CMS Authorization Tables

#### `cms_roles` table
```sql
CREATE TABLE cms_roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    tabs TEXT[],             -- Array of tab names
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

**Default CMS Roles**:
- `cms_admin`: Full access to all tabs
- `cms_product_manager`: Product + Inventory tabs
- `cms_order_manager`: Order tab only
- `cms_viewer`: Read-only access to selected tabs

#### `cms_tab_apis` table (NEW)
```sql
CREATE TABLE cms_tab_apis (
    id VARCHAR(36) PRIMARY KEY,
    tab_name VARCHAR(100) NOT NULL,
    api_path VARCHAR(500) NOT NULL,
    api_method VARCHAR(20) NOT NULL,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    UNIQUE(tab_name, api_path, api_method)
);
```

**Purpose**: Maps which APIs belong to which CMS tabs.  
**Relationship**: Many-to-Many (An API can belong to multiple tabs)

**Example Mappings**:
```
tab_name='product'    | api_path='/api/v1/products'     | method='GET'
tab_name='product'    | api_path='/api/v1/products'     | method='POST'
tab_name='product'    | api_path='/api/v1/products/*'   | method='PUT'
tab_name='inventory'  | api_path='/api/v1/inventory/*'  | method='GET'
tab_name='inventory'  | api_path='/api/v1/products'     | method='GET'  ← Shared API
```

#### `casbin_rule_cms` table
```sql
CREATE TABLE casbin_rule_cms (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100),      -- 'p' or 'g'
    v0 VARCHAR(100),         -- subject (user_id or cms_role_name)
    v1 VARCHAR(100),         -- tab name
    v2 VARCHAR(100),         -- api path
    v3 VARCHAR(100),         -- action (GET, POST, etc)
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);
```

**Example Policies**:
```
p, cms_admin, product, /api/v1/products/*, (GET|POST|PUT|DELETE)
p, cms_product_manager, product, /api/v1/products/*, (GET|POST|PUT)
p, cms_viewer, product, /api/v1/products/*, GET
```

**Example Role Assignments**:
```
g, user-789, cms_admin
g, user-456, cms_product_manager
```

---

## 🔄 Authorization Flow

### User/App Authorization Flow

```
1. User Request → API Endpoint
         ↓
2. Extract: user_id, api_path, method
         ↓
3. Check Casbin (User Enforcer)
   - Model: rbac_user_model.conf
   - Request: (user_id, domain, api_path, method)
         ↓
4. Casbin evaluates:
   - Get user roles from casbin_rule_user
   - Check policies for roles
   - Match resource path (keyMatch2)
   - Match action (regexMatch)
         ↓
5. Return: Allow/Deny
```

**Example**:
```go
// User-123 trying to GET /api/v1/products
allowed := userEnforcer.Enforce("user-123", "user", "/api/v1/products", "GET")
// Returns true if user-123 has 'user' role with permission
```

---

### CMS Authorization Flow

```
1. Admin Request → CMS Endpoint
         ↓
2. Extract: user_id, tab, api_path, method
         ↓
3. Verify Tab-API Mapping (cms_tab_apis)
   - Check if api_path+method belongs to tab
         ↓
4. Check Casbin (CMS Enforcer)
   - Model: rbac_cms_model.conf
   - Request: (user_id, tab, api_path, method)
         ↓
5. Casbin evaluates:
   - Get user CMS roles from casbin_rule_cms
   - Check policies for roles in that tab
   - Match API path (keyMatch2)
   - Match action (regexMatch)
         ↓
6. Return: Allow/Deny
```

**Example**:
```go
// Admin-789 trying to POST /api/v1/products in 'product' tab
allowed := cmsEnforcer.Enforce("admin-789", "product", "/api/v1/products", "POST")
// Returns true if admin-789 has cms_admin or cms_product_manager role
```

---

## 🎨 CMS Tab Structure

### Available Tabs

| Tab | Description | Example APIs |
|-----|-------------|--------------|
| `product` | Product management | `/api/v1/products`, `/api/v1/categories` |
| `inventory` | Inventory management | `/api/v1/inventory/*`, `/api/v1/warehouses` |
| `order` | Order management | `/api/v1/orders/*`, `/api/v1/orders/*/ship` |
| `user` | User management | `/api/v1/users/*`, `/api/v1/users/*/roles` |
| `report` | Reports & analytics | `/api/v1/reports/sales`, `/api/v1/reports/revenue` |
| `setting` | System settings | `/api/v1/settings/*`, `/api/v1/roles` |

### Tab-API Sharing

**Scenario**: Product API được sử dụng ở cả `product` và `inventory` tabs.

**Solution**: Multiple mappings in `cms_tab_apis`
```
| id  | tab_name  | api_path            | api_method |
|-----|-----------|---------------------|------------|
| 001 | product   | /api/v1/products    | GET        |
| 002 | product   | /api/v1/products/*  | GET        |
| 107 | inventory | /api/v1/products    | GET        |  ← Same API
| 108 | inventory | /api/v1/products/*  | GET        |  ← Same API
```

**Authorization**:
- User với `cms_product_manager` role: Can access via `product` tab
- User với `cms_inventory_manager` role: Can access via `inventory` tab
- Same API, different tab context

---

## 🔧 Implementation

### 1. Initialize Two Casbin Enforcers

```go
// User/App Enforcer
userEnforcer, err := casbin.NewEnforcer(
    "configs/rbac_user_model.conf",
    gormadapter.NewAdapterByDBWithCustomTable(db, nil, "casbin_rule_user"),
)

// CMS Enforcer
cmsEnforcer, err := casbin.NewEnforcer(
    "configs/rbac_cms_model.conf",
    gormadapter.NewAdapterByDBWithCustomTable(db, nil, "casbin_rule_cms"),
)
```

### 2. Check User/App Permission

```go
func CheckUserPermission(userID, domain, resource, action string) (bool, error) {
    // domain: "user" or "api"
    allowed, err := userEnforcer.Enforce(userID, domain, resource, action)
    return allowed, err
}

// Example usage:
allowed := CheckUserPermission("user-123", "user", "/api/v1/products", "GET")
```

### 3. Check CMS Permission

```go
func CheckCMSPermission(userID, tab, apiPath, method string) (bool, error) {
    // 1. Verify tab-API mapping exists
    exists := cmsTabAPIDAO.Exists(tab, apiPath, method)
    if !exists {
        return false, fmt.Errorf("API not registered for this tab")
    }
    
    // 2. Check Casbin permission
    allowed, err := cmsEnforcer.Enforce(userID, tab, apiPath, method)
    return allowed, err
}

// Example usage:
allowed := CheckCMSPermission("admin-789", "product", "/api/v1/products", "POST")
```

### 4. Get User's Accessible Tabs

```go
func GetUserCMSTabs(userID string) ([]string, error) {
    // Get all CMS roles for user
    roles := cmsEnforcer.GetRolesForUser(userID)
    
    // Get tabs for each role
    var allTabs []string
    for _, role := range roles {
        cmsRole := cmsRoleDAO.FindByName(role)
        allTabs = append(allTabs, cmsRole.Tabs...)
    }
    
    // Deduplicate
    tabs := unique(allTabs)
    return tabs, nil
}
```

### 5. Get APIs for Tab

```go
func GetTabAPIs(tabName string) ([]*CMSTabAPI, error) {
    apis := cmsTabAPIDAO.FindByTab(tabName)
    return apis, nil
}
```

---

## 📝 Usage Examples

### Example 1: Regular User (Web/App)

```go
// Register user
user := CreateUser("john_doe", "john@example.com", "password")

// Assign 'user' role for web access
userEnforcer.AddGroupingPolicy("user-123", "user", "user")

// Check permissions
canBrowse := userEnforcer.Enforce("user-123", "user", "/api/v1/products", "GET")
// → true (user role can browse products)

canCreate := userEnforcer.Enforce("user-123", "user", "/api/v1/products", "POST")
// → false (user role cannot create products)
```

### Example 2: CMS Admin

```go
// Create CMS admin user
admin := CreateUser("admin", "admin@company.com", "securepass")

// Assign 'cms_admin' role
cmsEnforcer.AddGroupingPolicy("admin-789", "cms_admin")

// Check CMS permissions
canManageProducts := cmsEnforcer.Enforce(
    "admin-789", 
    "product", 
    "/api/v1/products/*", 
    "DELETE",
)
// → true (cms_admin has full access)

canViewReports := cmsEnforcer.Enforce(
    "admin-789", 
    "report", 
    "/api/v1/reports/sales", 
    "GET",
)
// → true (cms_admin can access all tabs)
```

### Example 3: CMS Product Manager (Limited Access)

```go
// Create product manager
manager := CreateUser("product_manager", "pm@company.com", "password")

// Assign 'cms_product_manager' role
cmsEnforcer.AddGroupingPolicy("manager-456", "cms_product_manager")

// Get accessible tabs
tabs := GetUserCMSTabs("manager-456")
// → ["product", "inventory"]

// Check permissions
canManageProducts := cmsEnforcer.Enforce(
    "manager-456", 
    "product", 
    "/api/v1/products/*", 
    "PUT",
)
// → true (product_manager can edit products)

canDeleteProducts := cmsEnforcer.Enforce(
    "manager-456", 
    "product", 
    "/api/v1/products/*", 
    "DELETE",
)
// → false (product_manager cannot delete)

canManageOrders := cmsEnforcer.Enforce(
    "manager-456", 
    "order", 
    "/api/v1/orders/*", 
    "GET",
)
// → false (product_manager has no access to order tab)
```

---

## 🔄 Migration from Old System

### Migration Steps

1. **Run Migration Script**
   ```bash
   psql -U postgres -d iam_db -f migrations/005_separate_user_cms_authorization.sql
   ```

2. **Seed Data**
   ```bash
   psql -U postgres -d iam_db -f migrations/006_seed_separated_authorization.sql
   ```

3. **Verify Migration**
   ```sql
   -- Check User policies
   SELECT COUNT(*) FROM casbin_rule_user;
   
   -- Check CMS policies
   SELECT COUNT(*) FROM casbin_rule_cms;
   
   -- Check tab-API mappings
   SELECT COUNT(*) FROM cms_tab_apis;
   
   -- Old table backed up
   SELECT COUNT(*) FROM casbin_rule_old_backup;
   ```

4. **Update Application Code**
   - Initialize 2 Casbin enforcers
   - Update authorization checks
   - Use appropriate enforcer based on context

5. **Test Both Systems**
   - Test user/app authorization
   - Test CMS authorization
   - Verify tab-level access control

---

## ✅ Benefits of Separated Architecture

### 1. Clear Separation of Concerns
- User authorization logic separate from CMS authorization
- Different models for different use cases
- Easier to understand and maintain

### 2. Flexibility
- Can modify user authorization without affecting CMS
- Can add new CMS tabs without changing user policies
- Independent scaling and optimization

### 3. Security
- Principle of least privilege
- Tab-level granular control for CMS
- Isolated policy enforcement

### 4. Scalability
- Can use different Casbin adapters for each system
- Can cache policies separately
- Can deploy enforcers on different services

### 5. Maintainability
- Clearer code structure
- Easier debugging (know which enforcer to check)
- Better testing (test each system independently)

---

## 🚀 Best Practices

### 1. Role Naming Convention
- User roles: `user`, `premium_user`, `api_admin`
- CMS roles: `cms_admin`, `cms_product_manager`, `cms_viewer`

### 2. Policy Management
- Use wildcard patterns carefully: `/api/v1/products/*`
- Use regex for actions: `(GET|POST)` or `(GET|POST|PUT|DELETE)`
- Always specify domain/tab explicitly

### 3. Tab-API Mapping
- Register all APIs in `cms_tab_apis` table
- Document which APIs belong to which tabs
- Update mappings when adding new APIs

### 4. Testing
```go
// Test user authorization
func TestUserAuthorization(t *testing.T) {
    // Test regular user
    // Test premium user
    // Test API admin
}

// Test CMS authorization
func TestCMSAuthorization(t *testing.T) {
    // Test admin access
    // Test manager access
    // Test viewer access
    // Test tab restrictions
}
```

### 5. Monitoring
- Log authorization decisions
- Monitor policy load time
- Track denied access attempts
- Audit role assignments

---

## 📚 References

- [Casbin Documentation](https://casbin.org/)
- [RBAC Model](https://casbin.org/docs/rbac)
- [Multi-Tenancy Pattern](https://casbin.org/docs/multi-tenancy)
- [Pattern Matching](https://casbin.org/docs/function)

---

**Last Updated**: 2024-01-XX  
**Version**: 2.0 (Separated Architecture)

