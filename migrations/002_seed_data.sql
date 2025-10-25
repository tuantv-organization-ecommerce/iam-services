-- Insert default permissions
INSERT INTO permissions (id, name, resource, action, description) VALUES
    ('perm-001', 'User Read', 'user', 'read', 'Permission to read user information'),
    ('perm-002', 'User Create', 'user', 'create', 'Permission to create new users'),
    ('perm-003', 'User Update', 'user', 'update', 'Permission to update user information'),
    ('perm-004', 'User Delete', 'user', 'delete', 'Permission to delete users'),
    ('perm-005', 'Role Read', 'role', 'read', 'Permission to read role information'),
    ('perm-006', 'Role Create', 'role', 'create', 'Permission to create new roles'),
    ('perm-007', 'Role Update', 'role', 'update', 'Permission to update role information'),
    ('perm-008', 'Role Delete', 'role', 'delete', 'Permission to delete roles'),
    ('perm-009', 'Permission Read', 'permission', 'read', 'Permission to read permission information'),
    ('perm-010', 'Permission Create', 'permission', 'create', 'Permission to create new permissions'),
    ('perm-011', 'Permission Delete', 'permission', 'delete', 'Permission to delete permissions')
ON CONFLICT (id) DO NOTHING;

-- Insert default roles
INSERT INTO roles (id, name, description) VALUES
    ('role-001', 'admin', 'Administrator with full access'),
    ('role-002', 'user', 'Regular user with limited access'),
    ('role-003', 'moderator', 'Moderator with intermediate access')
ON CONFLICT (id) DO NOTHING;

-- Assign all permissions to admin role
INSERT INTO role_permissions (role_id, permission_id) VALUES
    ('role-001', 'perm-001'),
    ('role-001', 'perm-002'),
    ('role-001', 'perm-003'),
    ('role-001', 'perm-004'),
    ('role-001', 'perm-005'),
    ('role-001', 'perm-006'),
    ('role-001', 'perm-007'),
    ('role-001', 'perm-008'),
    ('role-001', 'perm-009'),
    ('role-001', 'perm-010'),
    ('role-001', 'perm-011')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read permissions to user role
INSERT INTO role_permissions (role_id, permission_id) VALUES
    ('role-002', 'perm-001'),
    ('role-002', 'perm-005'),
    ('role-002', 'perm-009')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read and some write permissions to moderator role
INSERT INTO role_permissions (role_id, permission_id) VALUES
    ('role-003', 'perm-001'),
    ('role-003', 'perm-003'),
    ('role-003', 'perm-005'),
    ('role-003', 'perm-009')
ON CONFLICT (role_id, permission_id) DO NOTHING;

