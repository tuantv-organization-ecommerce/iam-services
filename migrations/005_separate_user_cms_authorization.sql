-- Migration: Separate User/App and CMS Authorization
-- Purpose: Split authorization into two independent systems:
--   1. User/App Authorization (roles table) - for end users
--   2. CMS Authorization (cms_roles table) - for admin/staff

-- ============================================
-- 1. Create CMS Tab-API Mapping Table
-- ============================================
-- This table maps which APIs belong to which CMS tabs
-- An API can belong to multiple tabs
CREATE TABLE IF NOT EXISTS cms_tab_apis (
    id VARCHAR(36) PRIMARY KEY,
    tab_name VARCHAR(100) NOT NULL,
    api_path VARCHAR(500) NOT NULL,
    api_method VARCHAR(20) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tab_name, api_path, api_method)
);

-- Index for faster lookups
CREATE INDEX idx_cms_tab_apis_tab_name ON cms_tab_apis(tab_name);
CREATE INDEX idx_cms_tab_apis_api_path ON cms_tab_apis(api_path);
CREATE INDEX idx_cms_tab_apis_tab_api ON cms_tab_apis(tab_name, api_path);

-- ============================================
-- 2. Create Separate Casbin Tables
-- ============================================
-- User/App Casbin Rules (uses rbac_user_model.conf)
CREATE TABLE IF NOT EXISTS casbin_rule_user (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100),      -- 'p' (policy) or 'g' (grouping/role)
    v0 VARCHAR(100),         -- subject/user/role
    v1 VARCHAR(100),         -- domain (user, api)
    v2 VARCHAR(100),         -- object/resource path
    v3 VARCHAR(100),         -- action (GET, POST, etc)
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);

-- CMS Casbin Rules (uses rbac_cms_model.conf)
CREATE TABLE IF NOT EXISTS casbin_rule_cms (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100),      -- 'p' (policy) or 'g' (grouping/role)
    v0 VARCHAR(100),         -- subject/user/cms_role
    v1 VARCHAR(100),         -- tab name
    v2 VARCHAR(100),         -- api path
    v3 VARCHAR(100),         -- action (GET, POST, etc)
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);

-- Indexes for Casbin tables
CREATE INDEX idx_casbin_rule_user_ptype ON casbin_rule_user(ptype);
CREATE INDEX idx_casbin_rule_user_v0 ON casbin_rule_user(v0);
CREATE INDEX idx_casbin_rule_user_v1 ON casbin_rule_user(v1);

CREATE INDEX idx_casbin_rule_cms_ptype ON casbin_rule_cms(ptype);
CREATE INDEX idx_casbin_rule_cms_v0 ON casbin_rule_cms(v0);
CREATE INDEX idx_casbin_rule_cms_v1 ON casbin_rule_cms(v1);

-- ============================================
-- 3. Migrate Data from Old casbin_rule
-- ============================================
-- Migrate User/App policies (domain = 'user' or 'api')
INSERT INTO casbin_rule_user (ptype, v0, v1, v2, v3, v4, v5)
SELECT ptype, v0, v1, v2, v3, v4, v5
FROM casbin_rule
WHERE v1 IN ('user', 'api') OR (ptype = 'g' AND v2 IN ('user', 'api'));

-- Migrate CMS policies (domain = 'cms')
-- Note: This converts old format to new format
-- Old: (ptype, subject, 'cms', resource, action)
-- New: (ptype, subject, tab, api_path, action)
INSERT INTO casbin_rule_cms (ptype, v0, v1, v2, v3, v4, v5)
SELECT 
    ptype,
    v0,                                    -- subject (user or role)
    CASE 
        WHEN v2 LIKE '/cms/product%' THEN 'product'
        WHEN v2 LIKE '/cms/inventory%' THEN 'inventory'
        WHEN v2 LIKE '/cms/order%' THEN 'order'
        WHEN v2 LIKE '/cms/user%' THEN 'user'
        WHEN v2 LIKE '/cms/report%' THEN 'report'
        WHEN v2 LIKE '/cms/setting%' THEN 'setting'
        ELSE 'general'
    END as tab,                            -- extract tab from path
    REPLACE(v2, '/cms/', '/api/v1/'),     -- convert /cms/* to /api/v1/*
    v3,                                    -- action
    v4,
    v5
FROM casbin_rule
WHERE v1 = 'cms' OR (ptype = 'g' AND v2 = 'cms');

-- ============================================
-- 4. Add Comments
-- ============================================
COMMENT ON TABLE cms_tab_apis IS 'Maps CMS tabs to API endpoints. An API can belong to multiple tabs.';
COMMENT ON TABLE casbin_rule_user IS 'Casbin policies for User/App authorization (end users on web/app)';
COMMENT ON TABLE casbin_rule_cms IS 'Casbin policies for CMS authorization (admin/staff on CMS)';

COMMENT ON COLUMN cms_tab_apis.tab_name IS 'CMS tab name (e.g., product, inventory, order)';
COMMENT ON COLUMN cms_tab_apis.api_path IS 'API endpoint path (e.g., /api/v1/products/*)';
COMMENT ON COLUMN cms_tab_apis.api_method IS 'HTTP method (GET, POST, PUT, DELETE)';

-- ============================================
-- 5. Drop Old Table (Optional - keep for now)
-- ============================================
-- To keep backwards compatibility, we'll keep the old table
-- You can drop it later after verifying the migration:
-- DROP TABLE IF EXISTS casbin_rule;

-- Or rename it for backup:
ALTER TABLE casbin_rule RENAME TO casbin_rule_old_backup;

COMMENT ON TABLE casbin_rule_old_backup IS 'Backup of old unified casbin_rule table. Can be dropped after verification.';

