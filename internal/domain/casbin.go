package domain

import (
	"time"
)

// CasbinDomain represents different authorization domains
type CasbinDomain string

const (
	DomainUser CasbinDomain = "user" // Domain for end users
	DomainCMS  CasbinDomain = "cms"  // Domain for CMS/admin panel
	DomainAPI  CasbinDomain = "api"  // Domain for API access
)

// CMSTab represents different tabs/sections in CMS
type CMSTab string

const (
	CMSTabProduct   CMSTab = "product"
	CMSTabInventory CMSTab = "inventory"
	CMSTabOrder     CMSTab = "order"
	CMSTabUser      CMSTab = "user"
	CMSTabReport    CMSTab = "report"
	CMSTabSetting   CMSTab = "setting"
)

// APIResource represents API resource paths and methods
type APIResource struct {
	ID          string    `json:"id" db:"id"`
	Path        string    `json:"path" db:"path"`       // e.g., /api/v1/products
	Method      string    `json:"method" db:"method"`   // e.g., GET, POST, PUT, DELETE
	Service     string    `json:"service" db:"service"` // e.g., product-service, user-service
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CMSRole represents roles for CMS access control
type CMSRole struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Tabs        []CMSTab  `json:"tabs"` // CMS tabs this role can access
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UserCMSRole represents the relationship between users and CMS roles
type UserCMSRole struct {
	UserID    string    `json:"user_id" db:"user_id"`
	CMSRoleID string    `json:"cms_role_id" db:"cms_role_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// PolicyRule represents a Casbin policy rule
type PolicyRule struct {
	PType string `json:"p_type"` // p or g
	V0    string `json:"v0"`     // subject or role
	V1    string `json:"v1"`     // domain or parent role
	V2    string `json:"v2"`     // object or domain
	V3    string `json:"v3"`     // action
	V4    string `json:"v4"`     // extra
	V5    string `json:"v5"`     // extra
}

// AuthorizationRequest represents an authorization check request
type AuthorizationRequest struct {
	UserID   string       `json:"user_id"`
	Domain   CasbinDomain `json:"domain"`
	Resource string       `json:"resource"` // API path or CMS tab
	Action   string       `json:"action"`   // HTTP method or operation
}

// AuthorizationResponse represents an authorization check response
type AuthorizationResponse struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
}
