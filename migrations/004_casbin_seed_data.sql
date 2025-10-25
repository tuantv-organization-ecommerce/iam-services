-- Insert API resources for common endpoints
INSERT INTO api_resources (id, path, method, service, description) VALUES
    -- User service APIs
    ('api-001', '/api/v1/users', 'GET', 'user-service', 'List all users'),
    ('api-002', '/api/v1/users/:id', 'GET', 'user-service', 'Get user by ID'),
    ('api-003', '/api/v1/users', 'POST', 'user-service', 'Create new user'),
    ('api-004', '/api/v1/users/:id', 'PUT', 'user-service', 'Update user'),
    ('api-005', '/api/v1/users/:id', 'DELETE', 'user-service', 'Delete user'),
    
    -- Product service APIs
    ('api-006', '/api/v1/products', 'GET', 'product-service', 'List all products'),
    ('api-007', '/api/v1/products/:id', 'GET', 'product-service', 'Get product by ID'),
    ('api-008', '/api/v1/products', 'POST', 'product-service', 'Create new product'),
    ('api-009', '/api/v1/products/:id', 'PUT', 'product-service', 'Update product'),
    ('api-010', '/api/v1/products/:id', 'DELETE', 'product-service', 'Delete product'),
    
    -- Inventory service APIs
    ('api-011', '/api/v1/inventory', 'GET', 'inventory-service', 'List inventory'),
    ('api-012', '/api/v1/inventory/:id', 'GET', 'inventory-service', 'Get inventory by ID'),
    ('api-013', '/api/v1/inventory', 'POST', 'inventory-service', 'Create inventory record'),
    ('api-014', '/api/v1/inventory/:id', 'PUT', 'inventory-service', 'Update inventory'),
    
    -- Order service APIs
    ('api-015', '/api/v1/orders', 'GET', 'order-service', 'List orders'),
    ('api-016', '/api/v1/orders/:id', 'GET', 'order-service', 'Get order by ID'),
    ('api-017', '/api/v1/orders', 'POST', 'order-service', 'Create new order'),
    ('api-018', '/api/v1/orders/:id', 'PUT', 'order-service', 'Update order'),
    ('api-019', '/api/v1/orders/:id', 'DELETE', 'order-service', 'Cancel order')
ON CONFLICT (path, method) DO NOTHING;

-- Insert CMS roles with tab permissions
INSERT INTO cms_roles (id, name, description, tabs) VALUES
    ('cms-role-001', 'cms_admin', 'Full CMS access', ARRAY['product', 'inventory', 'order', 'user', 'report', 'setting']),
    ('cms-role-002', 'cms_product_manager', 'Manage products', ARRAY['product', 'inventory']),
    ('cms-role-003', 'cms_order_manager', 'Manage orders', ARRAY['order', 'report']),
    ('cms-role-004', 'cms_content_editor', 'Edit content', ARRAY['product']),
    ('cms-role-005', 'cms_viewer', 'View only access', ARRAY['product', 'inventory', 'order', 'report'])
ON CONFLICT (id) DO NOTHING;

-- Update existing roles with domain
UPDATE roles SET domain = 'user' WHERE domain IS NULL;

-- Insert Casbin policies for CMS roles
-- Format: p, role, domain, resource, action

-- CMS Admin - Full access to all CMS tabs
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
    ('p', 'cms_admin', 'cms', '/cms/product/*', '(GET|POST|PUT|DELETE)'),
    ('p', 'cms_admin', 'cms', '/cms/inventory/*', '(GET|POST|PUT|DELETE)'),
    ('p', 'cms_admin', 'cms', '/cms/order/*', '(GET|POST|PUT|DELETE)'),
    ('p', 'cms_admin', 'cms', '/cms/user/*', '(GET|POST|PUT|DELETE)'),
    ('p', 'cms_admin', 'cms', '/cms/report/*', '(GET)'),
    ('p', 'cms_admin', 'cms', '/cms/setting/*', '(GET|POST|PUT|DELETE)')
ON CONFLICT DO NOTHING;

-- CMS Product Manager
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
    ('p', 'cms_product_manager', 'cms', '/cms/product/*', '(GET|POST|PUT|DELETE)'),
    ('p', 'cms_product_manager', 'cms', '/cms/inventory/*', '(GET|POST|PUT)')
ON CONFLICT DO NOTHING;

-- CMS Order Manager
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
    ('p', 'cms_order_manager', 'cms', '/cms/order/*', '(GET|POST|PUT)'),
    ('p', 'cms_order_manager', 'cms', '/cms/report/*', '(GET)')
ON CONFLICT DO NOTHING;

-- CMS Content Editor
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
    ('p', 'cms_content_editor', 'cms', '/cms/product/*', '(GET|POST|PUT)')
ON CONFLICT DO NOTHING;

-- CMS Viewer - Read only
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
    ('p', 'cms_viewer', 'cms', '/cms/product/*', 'GET'),
    ('p', 'cms_viewer', 'cms', '/cms/inventory/*', 'GET'),
    ('p', 'cms_viewer', 'cms', '/cms/order/*', 'GET'),
    ('p', 'cms_viewer', 'cms', '/cms/report/*', 'GET')
ON CONFLICT DO NOTHING;

-- Insert Casbin policies for API access based on user roles
-- Admin role - Full API access
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
    ('p', 'admin', 'api', '/api/v1/**', '(GET|POST|PUT|DELETE)')
ON CONFLICT DO NOTHING;

-- User role - Limited API access
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
    ('p', 'user', 'api', '/api/v1/products', 'GET'),
    ('p', 'user', 'api', '/api/v1/products/*', 'GET'),
    ('p', 'user', 'api', '/api/v1/orders', '(GET|POST)'),
    ('p', 'user', 'api', '/api/v1/orders/*', 'GET')
ON CONFLICT DO NOTHING;

-- Moderator role - Intermediate access
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES
    ('p', 'moderator', 'api', '/api/v1/products', '(GET|POST|PUT)'),
    ('p', 'moderator', 'api', '/api/v1/products/*', '(GET|POST|PUT)'),
    ('p', 'moderator', 'api', '/api/v1/orders', 'GET'),
    ('p', 'moderator', 'api', '/api/v1/orders/*', 'GET')
ON CONFLICT DO NOTHING;

