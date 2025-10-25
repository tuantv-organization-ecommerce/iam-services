package repository

import (
	"context"

	"github.com/tvttt/iam-services/internal/domain/model"
)

// RoleRepository defines the contract for role data access
type RoleRepository interface {
	// Save persists a new role
	Save(ctx context.Context, role *model.Role) error

	// FindByID retrieves a role by ID
	FindByID(ctx context.Context, id string) (*model.Role, error)

	// FindByName retrieves a role by name
	FindByName(ctx context.Context, name string) (*model.Role, error)

	// FindByDomain retrieves all roles in a domain
	FindByDomain(ctx context.Context, domain string) ([]*model.Role, error)

	// FindAll retrieves all roles with pagination
	FindAll(ctx context.Context, page, pageSize int) ([]*model.Role, int, error)

	// Update updates an existing role
	Update(ctx context.Context, role *model.Role) error

	// Delete removes a role
	Delete(ctx context.Context, id string) error

	// ExistsByName checks if role name already exists
	ExistsByName(ctx context.Context, name string) (bool, error)
}
