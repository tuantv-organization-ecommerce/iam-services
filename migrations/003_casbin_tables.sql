-- Create Casbin policy table
CREATE TABLE IF NOT EXISTS casbin_rule (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100),
    v0 VARCHAR(100),
    v1 VARCHAR(100),
    v2 VARCHAR(100),
    v3 VARCHAR(100),
    v4 VARCHAR(100),
    v5 VARCHAR(100),
    CONSTRAINT unique_key_casbin_rule UNIQUE(ptype, v0, v1, v2, v3, v4, v5)
);

-- Create API resources table
CREATE TABLE IF NOT EXISTS api_resources (
    id VARCHAR(36) PRIMARY KEY,
    path VARCHAR(500) NOT NULL,
    method VARCHAR(20) NOT NULL,
    service VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(path, method)
);

-- Create CMS roles table
CREATE TABLE IF NOT EXISTS cms_roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    tabs TEXT[], -- Array of CMS tabs
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create user_cms_roles junction table
CREATE TABLE IF NOT EXISTS user_cms_roles (
    user_id VARCHAR(36) NOT NULL,
    cms_role_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, cms_role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (cms_role_id) REFERENCES cms_roles(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX idx_casbin_rule_ptype ON casbin_rule(ptype);
CREATE INDEX idx_casbin_rule_v0 ON casbin_rule(v0);
CREATE INDEX idx_casbin_rule_v1 ON casbin_rule(v1);
CREATE INDEX idx_api_resources_path_method ON api_resources(path, method);
CREATE INDEX idx_api_resources_service ON api_resources(service);
CREATE INDEX idx_cms_roles_name ON cms_roles(name);
CREATE INDEX idx_user_cms_roles_user_id ON user_cms_roles(user_id);
CREATE INDEX idx_user_cms_roles_cms_role_id ON user_cms_roles(cms_role_id);

-- Add domain column to roles table to distinguish between user roles and system roles
ALTER TABLE roles ADD COLUMN IF NOT EXISTS domain VARCHAR(50) DEFAULT 'user';
CREATE INDEX IF NOT EXISTS idx_roles_domain ON roles(domain);

