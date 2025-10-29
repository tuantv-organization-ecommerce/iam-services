package model

import (
	"errors"
	"time"
)

var (
	// ErrInvalidRoleName indicates role name validation failed
	ErrInvalidRoleName = errors.New("invalid role name")
	// ErrEmptyRoleName indicates role name is empty
	ErrEmptyRoleName = errors.New("role name cannot be empty")
)

// Role represents a role entity in the domain
type Role struct {
	id          string
	name        string
	description string
	domain      string // user, cms, api
	permissions []Permission
	createdAt   time.Time
	updatedAt   time.Time
}

// NewRole creates a new Role entity
func NewRole(id, name, description, domain string) *Role {
	now := time.Now()
	return &Role{
		id:          id,
		name:        name,
		description: description,
		domain:      domain,
		permissions: make([]Permission, 0),
		createdAt:   now,
		updatedAt:   now,
	}
}

// ReconstructRole reconstructs a Role from persistence
func ReconstructRole(id, name, description, domain string, createdAt, updatedAt time.Time) *Role {
	return &Role{
		id:          id,
		name:        name,
		description: description,
		domain:      domain,
		permissions: make([]Permission, 0),
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// Getters

// ID returns the role's unique identifier
func (r *Role) ID() string { return r.id }

// Name returns the role's name
func (r *Role) Name() string { return r.name }

// Description returns the role's description
func (r *Role) Description() string { return r.description }

// Domain returns the role's authorization domain
func (r *Role) Domain() string { return r.domain }

// Permissions returns the permissions assigned to this role
func (r *Role) Permissions() []Permission { return r.permissions }

// CreatedAt returns when the role was created
func (r *Role) CreatedAt() time.Time { return r.createdAt }

// UpdatedAt returns when the role was last updated
func (r *Role) UpdatedAt() time.Time { return r.updatedAt }

// UpdateDetails updates role details
func (r *Role) UpdateDetails(name, description string) {
	if name != "" {
		r.name = name
	}
	r.description = description
	r.updatedAt = time.Now()
}

// AddPermission adds a permission to the role
func (r *Role) AddPermission(permission Permission) {
	for _, p := range r.permissions {
		if p.ID() == permission.ID() {
			return // Already has permission
		}
	}
	r.permissions = append(r.permissions, permission)
}

// RemovePermission removes a permission from the role
func (r *Role) RemovePermission(permissionID string) {
	for i, p := range r.permissions {
		if p.ID() == permissionID {
			r.permissions = append(r.permissions[:i], r.permissions[i+1:]...)
			return
		}
	}
}

// HasPermission checks if role has a specific permission
func (r *Role) HasPermission(permissionID string) bool {
	for _, p := range r.permissions {
		if p.ID() == permissionID {
			return true
		}
	}
	return false
}

// SetPermissions sets all permissions for the role
func (r *Role) SetPermissions(permissions []Permission) {
	r.permissions = permissions
	r.updatedAt = time.Now()
}

// Validate validates the role entity
func (r *Role) Validate() error {
	if r.name == "" {
		return ErrEmptyRoleName
	}
	if len(r.name) < 2 {
		return ErrInvalidRoleName
	}
	return nil
}
