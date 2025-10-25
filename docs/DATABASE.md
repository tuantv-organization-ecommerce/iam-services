# Database Schema Documentation

## Tổng quan

IAM Service sử dụng PostgreSQL để lưu trữ thông tin người dùng, vai trò và quyền hạn. Database được thiết kế theo mô hình nhiều-nhiều (many-to-many) giữa Users, Roles và Permissions.

## Entity Relationship Diagram

```
┌─────────────────────────────────┐
│           users                  │
├─────────────────────────────────┤
│ id              VARCHAR(36) PK   │
│ username        VARCHAR(100) UK  │
│ email           VARCHAR(255) UK  │
│ password_hash   VARCHAR(255)     │
│ full_name       VARCHAR(255)     │
│ is_active       BOOLEAN          │
│ created_at      TIMESTAMP        │
│ updated_at      TIMESTAMP        │
└────────────┬────────────────────┘
             │
             │ 1:N
             │
┌────────────▼────────────────────┐
│       user_roles                 │
├─────────────────────────────────┤
│ user_id         VARCHAR(36) PK,FK│
│ role_id         VARCHAR(36) PK,FK│
│ created_at      TIMESTAMP        │
└────────────┬────────────────────┘
             │
             │ N:1
             │
┌────────────▼────────────────────┐
│           roles                  │
├─────────────────────────────────┤
│ id              VARCHAR(36) PK   │
│ name            VARCHAR(100) UK  │
│ description     TEXT             │
│ created_at      TIMESTAMP        │
│ updated_at      TIMESTAMP        │
└────────────┬────────────────────┘
             │
             │ 1:N
             │
┌────────────▼────────────────────┐
│    role_permissions              │
├─────────────────────────────────┤
│ role_id         VARCHAR(36) PK,FK│
│ permission_id   VARCHAR(36) PK,FK│
│ created_at      TIMESTAMP        │
└────────────┬────────────────────┘
             │
             │ N:1
             │
┌────────────▼────────────────────┐
│       permissions                │
├─────────────────────────────────┤
│ id              VARCHAR(36) PK   │
│ name            VARCHAR(100) UK  │
│ resource        VARCHAR(100)     │
│ action          VARCHAR(50)      │
│ description     TEXT             │
│ created_at      TIMESTAMP        │
│ updated_at      TIMESTAMP        │
└─────────────────────────────────┘
```

## Tables

### 1. users

Lưu trữ thông tin người dùng.

```sql
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Columns**:
- `id`: UUID - Primary key
- `username`: Tên đăng nhập (unique)
- `email`: Email (unique)
- `password_hash`: Password đã được hash bằng bcrypt
- `full_name`: Tên đầy đủ của user
- `is_active`: User có đang active không
- `created_at`: Thời gian tạo
- `updated_at`: Thời gian cập nhật lần cuối

**Indexes**:
- `idx_users_username ON users(username)`
- `idx_users_email ON users(email)`
- `idx_users_is_active ON users(is_active)`

---

### 2. roles

Lưu trữ thông tin vai trò.

```sql
CREATE TABLE roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Columns**:
- `id`: UUID - Primary key
- `name`: Tên vai trò (unique)
- `description`: Mô tả vai trò
- `created_at`: Thời gian tạo
- `updated_at`: Thời gian cập nhật lần cuối

**Indexes**:
- `idx_roles_name ON roles(name)`

**Default Roles**:
- `admin`: Quản trị viên - có tất cả quyền
- `user`: Người dùng thông thường - quyền hạn chế
- `moderator`: Người kiểm duyệt - quyền trung gian

---

### 3. permissions

Lưu trữ thông tin quyền hạn.

```sql
CREATE TABLE permissions (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(resource, action)
);
```

**Columns**:
- `id`: UUID - Primary key
- `name`: Tên quyền
- `resource`: Tài nguyên (vd: user, role, product)
- `action`: Hành động (vd: read, create, update, delete)
- `description`: Mô tả quyền
- `created_at`: Thời gian tạo
- `updated_at`: Thời gian cập nhật lần cuối

**Indexes**:
- `idx_permissions_resource_action ON permissions(resource, action)`

**Permission Format**:
```
Resource: Action
Examples:
- user:read
- user:create
- user:update
- user:delete
- role:read
- role:create
- permission:read
```

---

### 4. user_roles

Junction table giữa users và roles (many-to-many).

```sql
CREATE TABLE user_roles (
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);
```

**Columns**:
- `user_id`: Foreign key đến users
- `role_id`: Foreign key đến roles
- `created_at`: Thời gian gán role

**Indexes**:
- `idx_user_roles_user_id ON user_roles(user_id)`
- `idx_user_roles_role_id ON user_roles(role_id)`

**Cascade Delete**: Khi xóa user hoặc role, các liên kết trong bảng này cũng bị xóa

---

### 5. role_permissions

Junction table giữa roles và permissions (many-to-many).

```sql
CREATE TABLE role_permissions (
    role_id VARCHAR(36) NOT NULL,
    permission_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);
```

**Columns**:
- `role_id`: Foreign key đến roles
- `permission_id`: Foreign key đến permissions
- `created_at`: Thời gian gán permission

**Indexes**:
- `idx_role_permissions_role_id ON role_permissions(role_id)`
- `idx_role_permissions_permission_id ON role_permissions(permission_id)`

**Cascade Delete**: Khi xóa role hoặc permission, các liên kết trong bảng này cũng bị xóa

---

## Common Queries

### Get user with roles and permissions

```sql
-- Get all roles of a user
SELECT r.*
FROM roles r
INNER JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = 'user-id';

-- Get all permissions of a user (through roles)
SELECT DISTINCT p.*
FROM permissions p
INNER JOIN role_permissions rp ON p.id = rp.permission_id
INNER JOIN user_roles ur ON rp.role_id = ur.role_id
WHERE ur.user_id = 'user-id';
```

### Check if user has specific permission

```sql
SELECT EXISTS(
    SELECT 1
    FROM permissions p
    INNER JOIN role_permissions rp ON p.id = rp.permission_id
    INNER JOIN user_roles ur ON rp.role_id = ur.role_id
    WHERE ur.user_id = 'user-id'
      AND p.resource = 'user'
      AND p.action = 'create'
) AS has_permission;
```

### Get role with all permissions

```sql
SELECT r.*, 
       json_agg(json_build_object(
           'id', p.id,
           'name', p.name,
           'resource', p.resource,
           'action', p.action
       )) AS permissions
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE r.id = 'role-id'
GROUP BY r.id;
```

### List users with their roles

```sql
SELECT u.id, u.username, u.email,
       json_agg(json_build_object(
           'id', r.id,
           'name', r.name
       )) AS roles
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.role_id
GROUP BY u.id, u.username, u.email;
```

## Migrations

### Initial Schema

File: `migrations/001_init_schema.sql`

Tạo tất cả tables, indexes và constraints.

### Seed Data

File: `migrations/002_seed_data.sql`

Insert dữ liệu mặc định:
- 3 roles: admin, user, moderator
- 11 permissions cơ bản
- Mapping role-permission mặc định

## Performance Considerations

### Indexes

Tất cả foreign keys đều có indexes để tối ưu:
- JOIN operations
- WHERE clauses trên foreign keys
- CASCADE deletes

### Query Optimization

1. **Use prepared statements** để tránh SQL injection và cache query plans
2. **Limit số lượng JOINs** khi có thể
3. **Use pagination** cho list operations
4. **Cache frequently accessed data** (roles, permissions)

### Connection Pooling

```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

## Backup và Recovery

### Backup Database

```bash
pg_dump -U postgres -d iam_db -F c -f iam_db_backup.dump
```

### Restore Database

```bash
pg_restore -U postgres -d iam_db -c iam_db_backup.dump
```

### Scheduled Backups

Sử dụng cron job:

```bash
# Daily backup at 2 AM
0 2 * * * pg_dump -U postgres -d iam_db -F c -f /backup/iam_db_$(date +\%Y\%m\%d).dump
```

## Security Best Practices

1. **Password Storage**: Luôn hash password với bcrypt (cost factor 10+)
2. **SQL Injection**: Sử dụng parameterized queries
3. **Connection Security**: Sử dụng SSL/TLS trong production
4. **Access Control**: Giới hạn database user permissions
5. **Audit Logging**: Log các thay đổi quan trọng (role assignments, permission changes)
6. **Regular Backups**: Backup database thường xuyên
7. **Encryption at Rest**: Encrypt database files trong production

## Monitoring

### Useful Queries

```sql
-- Check table sizes
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Check index usage
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- Check slow queries
SELECT 
    calls,
    total_time,
    mean_time,
    query
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;
```

