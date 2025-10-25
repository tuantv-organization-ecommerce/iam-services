package service

import "context"

// CasbinDomain represents authorization domains
type CasbinDomain string

const (
	DomainUser CasbinDomain = "user"
	DomainCMS  CasbinDomain = "cms"
	DomainAPI  CasbinDomain = "api"
)

// AuthorizationService defines the contract for Casbin authorization
// This is a domain service interface
type AuthorizationService interface {
	// Enforce checks if subject can perform action on object in domain
	Enforce(ctx context.Context, subject string, domain CasbinDomain, object, action string) (bool, error)

	// AddPolicy adds a policy rule
	AddPolicy(ctx context.Context, subject string, domain CasbinDomain, object, action string) error

	// RemovePolicy removes a policy rule
	RemovePolicy(ctx context.Context, subject string, domain CasbinDomain, object, action string) error

	// AddRoleForUser assigns a role to user in a domain
	AddRoleForUser(ctx context.Context, userID, role string, domain CasbinDomain) error

	// RemoveRoleForUser removes a role from user in a domain
	RemoveRoleForUser(ctx context.Context, userID, role string, domain CasbinDomain) error

	// GetRolesForUser gets all roles for user in a domain
	GetRolesForUser(ctx context.Context, userID string, domain CasbinDomain) ([]string, error)

	// GetPermissionsForUser gets all permissions for user in a domain
	GetPermissionsForUser(ctx context.Context, userID string, domain CasbinDomain) ([][]string, error)
}
