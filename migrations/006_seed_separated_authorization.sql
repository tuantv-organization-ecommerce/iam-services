-- Seed Data for Separated User/App and CMS Authorization

-- ============================================
-- 1. Seed CMS Tab-API Mappings
-- ============================================
-- Product Tab APIs
INSERT INTO cms_tab_apis (id, tab_name, api_path, api_method, description) VALUES
('tab-api-001', 'product', '/api/v1/products', 'GET', 'List all products'),
('tab-api-002', 'product', '/api/v1/products', 'POST', 'Create new product'),
('tab-api-003', 'product', '/api/v1/products/*', 'GET', 'Get product by ID'),
('tab-api-004', 'product', '/api/v1/products/*', 'PUT', 'Update product'),
('tab-api-005', 'product', '/api/v1/products/*', 'DELETE', 'Delete product'),
('tab-api-006', 'product', '/api/v1/products/*/variants', 'GET', 'List product variants'),
('tab-api-007', 'product', '/api/v1/products/*/variants', 'POST', 'Add product variant'),
('tab-api-008', 'product', '/api/v1/categories', 'GET', 'List categories'),
('tab-api-009', 'product', '/api/v1/categories', 'POST', 'Create category');

-- Inventory Tab APIs
INSERT INTO cms_tab_apis (id, tab_name, api_path, api_method, description) VALUES
('tab-api-101', 'inventory', '/api/v1/inventory', 'GET', 'List inventory items'),
('tab-api-102', 'inventory', '/api/v1/inventory/*', 'GET', 'Get inventory by product'),
('tab-api-103', 'inventory', '/api/v1/inventory/*/adjust', 'POST', 'Adjust inventory'),
('tab-api-104', 'inventory', '/api/v1/inventory/*/history', 'GET', 'Get inventory history'),
('tab-api-105', 'inventory', '/api/v1/warehouses', 'GET', 'List warehouses'),
('tab-api-106', 'inventory', '/api/v1/warehouses', 'POST', 'Create warehouse');

-- Note: Some APIs can belong to multiple tabs
-- Example: Product APIs also used in Inventory tab
INSERT INTO cms_tab_apis (id, tab_name, api_path, api_method, description) VALUES
('tab-api-107', 'inventory', '/api/v1/products', 'GET', 'List products for inventory'),
('tab-api-108', 'inventory', '/api/v1/products/*', 'GET', 'Get product details');

-- Order Tab APIs
INSERT INTO cms_tab_apis (id, tab_name, api_path, api_method, description) VALUES
('tab-api-201', 'order', '/api/v1/orders', 'GET', 'List all orders'),
('tab-api-202', 'order', '/api/v1/orders/*', 'GET', 'Get order by ID'),
('tab-api-203', 'order', '/api/v1/orders/*', 'PUT', 'Update order status'),
('tab-api-204', 'order', '/api/v1/orders/*/cancel', 'POST', 'Cancel order'),
('tab-api-205', 'order', '/api/v1/orders/*/ship', 'POST', 'Ship order'),
('tab-api-206', 'order', '/api/v1/orders/*/refund', 'POST', 'Refund order');

-- User Management Tab APIs
INSERT INTO cms_tab_apis (id, tab_name, api_path, api_method, description) VALUES
('tab-api-301', 'user', '/api/v1/users', 'GET', 'List all users'),
('tab-api-302', 'user', '/api/v1/users/*', 'GET', 'Get user by ID'),
('tab-api-303', 'user', '/api/v1/users/*', 'PUT', 'Update user'),
('tab-api-304', 'user', '/api/v1/users/*/activate', 'POST', 'Activate user'),
('tab-api-305', 'user', '/api/v1/users/*/deactivate', 'POST', 'Deactivate user'),
('tab-api-306', 'user', '/api/v1/users/*/roles', 'GET', 'Get user roles'),
('tab-api-307', 'user', '/api/v1/users/*/roles', 'POST', 'Assign role to user');

-- Report Tab APIs
INSERT INTO cms_tab_apis (id, tab_name, api_path, api_method, description) VALUES
('tab-api-401', 'report', '/api/v1/reports/sales', 'GET', 'Get sales report'),
('tab-api-402', 'report', '/api/v1/reports/revenue', 'GET', 'Get revenue report'),
('tab-api-403', 'report', '/api/v1/reports/inventory', 'GET', 'Get inventory report'),
('tab-api-404', 'report', '/api/v1/reports/users', 'GET', 'Get user statistics'),
('tab-api-405', 'report', '/api/v1/reports/orders', 'GET', 'Get order statistics'),
('tab-api-406', 'report', '/api/v1/reports/export', 'POST', 'Export report');

-- Setting Tab APIs
INSERT INTO cms_tab_apis (id, tab_name, api_path, api_method, description) VALUES
('tab-api-501', 'setting', '/api/v1/settings', 'GET', 'Get all settings'),
('tab-api-502', 'setting', '/api/v1/settings/*', 'GET', 'Get setting by key'),
('tab-api-503', 'setting', '/api/v1/settings/*', 'PUT', 'Update setting'),
('tab-api-504', 'setting', '/api/v1/roles', 'GET', 'List roles'),
('tab-api-505', 'setting', '/api/v1/roles', 'POST', 'Create role'),
('tab-api-506', 'setting', '/api/v1/roles/*', 'PUT', 'Update role'),
('tab-api-507', 'setting', '/api/v1/roles/*', 'DELETE', 'Delete role');

-- ============================================
-- 2. Seed User/App Authorization Policies
-- ============================================
-- Basic user role policies (for regular app users)
INSERT INTO casbin_rule_user (ptype, v0, v1, v2, v3) VALUES
-- Regular user can browse products
('p', 'user', 'user', '/api/v1/products', 'GET'),
('p', 'user', 'user', '/api/v1/products/*', 'GET'),
('p', 'user', 'user', '/api/v1/categories', 'GET'),

-- Regular user can manage their own orders
('p', 'user', 'user', '/api/v1/orders', '(GET|POST)'),
('p', 'user', 'user', '/api/v1/users/self', 'GET'),
('p', 'user', 'user', '/api/v1/users/self', 'PUT');

-- Premium user with more access
INSERT INTO casbin_rule_user (ptype, v0, v1, v2, v3) VALUES
('p', 'premium_user', 'api', '/api/v1/products/**', '(GET|POST)'),
('p', 'premium_user', 'api', '/api/v1/orders/**', '(GET|POST|PUT)');

-- Admin with full API access (for app/web)
INSERT INTO casbin_rule_user (ptype, v0, v1, v2, v3) VALUES
('p', 'api_admin', 'api', '/api/v1/**', '(GET|POST|PUT|DELETE)');

-- ============================================
-- 3. Seed CMS Authorization Policies
-- ============================================
-- CMS Admin - Full access to all tabs
INSERT INTO casbin_rule_cms (ptype, v0, v1, v2, v3) VALUES
-- Product tab
('p', 'cms_admin', 'product', '/api/v1/products', '(GET|POST)'),
('p', 'cms_admin', 'product', '/api/v1/products/*', '(GET|POST|PUT|DELETE)'),
('p', 'cms_admin', 'product', '/api/v1/products/*/variants', '(GET|POST)'),
('p', 'cms_admin', 'product', '/api/v1/categories', '(GET|POST)'),

-- Inventory tab
('p', 'cms_admin', 'inventory', '/api/v1/inventory', 'GET'),
('p', 'cms_admin', 'inventory', '/api/v1/inventory/*', 'GET'),
('p', 'cms_admin', 'inventory', '/api/v1/inventory/*/adjust', 'POST'),
('p', 'cms_admin', 'inventory', '/api/v1/inventory/*/history', 'GET'),
('p', 'cms_admin', 'inventory', '/api/v1/warehouses', '(GET|POST)'),
('p', 'cms_admin', 'inventory', '/api/v1/products', 'GET'),
('p', 'cms_admin', 'inventory', '/api/v1/products/*', 'GET'),

-- Order tab
('p', 'cms_admin', 'order', '/api/v1/orders', 'GET'),
('p', 'cms_admin', 'order', '/api/v1/orders/*', '(GET|PUT)'),
('p', 'cms_admin', 'order', '/api/v1/orders/*/cancel', 'POST'),
('p', 'cms_admin', 'order', '/api/v1/orders/*/ship', 'POST'),
('p', 'cms_admin', 'order', '/api/v1/orders/*/refund', 'POST'),

-- User tab
('p', 'cms_admin', 'user', '/api/v1/users', 'GET'),
('p', 'cms_admin', 'user', '/api/v1/users/*', '(GET|PUT)'),
('p', 'cms_admin', 'user', '/api/v1/users/*/activate', 'POST'),
('p', 'cms_admin', 'user', '/api/v1/users/*/deactivate', 'POST'),
('p', 'cms_admin', 'user', '/api/v1/users/*/roles', '(GET|POST)'),

-- Report tab
('p', 'cms_admin', 'report', '/api/v1/reports/sales', 'GET'),
('p', 'cms_admin', 'report', '/api/v1/reports/revenue', 'GET'),
('p', 'cms_admin', 'report', '/api/v1/reports/inventory', 'GET'),
('p', 'cms_admin', 'report', '/api/v1/reports/users', 'GET'),
('p', 'cms_admin', 'report', '/api/v1/reports/orders', 'GET'),
('p', 'cms_admin', 'report', '/api/v1/reports/export', 'POST'),

-- Setting tab
('p', 'cms_admin', 'setting', '/api/v1/settings', 'GET'),
('p', 'cms_admin', 'setting', '/api/v1/settings/*', '(GET|PUT)'),
('p', 'cms_admin', 'setting', '/api/v1/roles', '(GET|POST)'),
('p', 'cms_admin', 'setting', '/api/v1/roles/*', '(PUT|DELETE)');

-- CMS Product Manager - Product and Inventory tabs only
INSERT INTO casbin_rule_cms (ptype, v0, v1, v2, v3) VALUES
-- Product tab - full CRUD
('p', 'cms_product_manager', 'product', '/api/v1/products', '(GET|POST)'),
('p', 'cms_product_manager', 'product', '/api/v1/products/*', '(GET|POST|PUT|DELETE)'),
('p', 'cms_product_manager', 'product', '/api/v1/products/*/variants', '(GET|POST)'),
('p', 'cms_product_manager', 'product', '/api/v1/categories', '(GET|POST)'),

-- Inventory tab - read and adjust only
('p', 'cms_product_manager', 'inventory', '/api/v1/inventory', 'GET'),
('p', 'cms_product_manager', 'inventory', '/api/v1/inventory/*', 'GET'),
('p', 'cms_product_manager', 'inventory', '/api/v1/inventory/*/adjust', 'POST'),
('p', 'cms_product_manager', 'inventory', '/api/v1/products', 'GET'),
('p', 'cms_product_manager', 'inventory', '/api/v1/products/*', 'GET');

-- CMS Order Manager - Order tab only
INSERT INTO casbin_rule_cms (ptype, v0, v1, v2, v3) VALUES
('p', 'cms_order_manager', 'order', '/api/v1/orders', 'GET'),
('p', 'cms_order_manager', 'order', '/api/v1/orders/*', '(GET|PUT)'),
('p', 'cms_order_manager', 'order', '/api/v1/orders/*/ship', 'POST');

-- CMS Viewer - Read-only access to specific tabs
INSERT INTO casbin_rule_cms (ptype, v0, v1, v2, v3) VALUES
-- Product tab - read only
('p', 'cms_viewer', 'product', '/api/v1/products', 'GET'),
('p', 'cms_viewer', 'product', '/api/v1/products/*', 'GET'),
('p', 'cms_viewer', 'product', '/api/v1/categories', 'GET'),

-- Inventory tab - read only
('p', 'cms_viewer', 'inventory', '/api/v1/inventory', 'GET'),
('p', 'cms_viewer', 'inventory', '/api/v1/inventory/*', 'GET'),

-- Report tab - read only
('p', 'cms_viewer', 'report', '/api/v1/reports/sales', 'GET'),
('p', 'cms_viewer', 'report', '/api/v1/reports/revenue', 'GET'),
('p', 'cms_viewer', 'report', '/api/v1/reports/inventory', 'GET');

-- ============================================
-- 4. Update CMS Roles Table with Tabs
-- ============================================
-- Make sure cms_roles have tabs assigned
UPDATE cms_roles SET tabs = ARRAY['product', 'inventory', 'order', 'user', 'report', 'setting']
WHERE name = 'cms_admin';

UPDATE cms_roles SET tabs = ARRAY['product', 'inventory']
WHERE name = 'cms_product_manager';

UPDATE cms_roles SET tabs = ARRAY['order']
WHERE name = 'cms_order_manager';

UPDATE cms_roles SET tabs = ARRAY['product', 'inventory', 'report']
WHERE name = 'cms_viewer';

-- ============================================
-- 5. Example: Assign Users to Roles
-- ============================================
-- These are example assignments
-- In production, use the API to assign roles

-- Assign user-789 as CMS Admin (has access to all tabs)
-- INSERT INTO casbin_rule_cms (ptype, v0, v1) VALUES
-- ('g', 'user-789', 'cms_admin');

-- Assign user-456 as Product Manager (has access to product + inventory tabs)
-- INSERT INTO casbin_rule_cms (ptype, v0, v1) VALUES
-- ('g', 'user-456', 'cms_product_manager');

-- Assign user-123 as regular user (app access)
-- INSERT INTO casbin_rule_user (ptype, v0, v1, v2) VALUES
-- ('g', 'user-123', 'user', 'user');

