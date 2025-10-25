package repository

import (
	"context"

	"github.com/tvttt/iam-services/internal/domain/model"
)

// PermissionRepository defines the contract for permission data access
type PermissionRepository interface {
	// Save persists a new permission
	Save(ctx context.Context, permission *model.Permission) error

	// FindByID retrieves a permission by ID
	FindByID(ctx context.Context, id string) (*model.Permission, error)

	// FindByResourceAndAction retrieves a permission by resource and action
	FindByResourceAndAction(ctx context.Context, resource, action string) (*model.Permission, error)

	// FindAll retrieves all permissions with pagination
	FindAll(ctx context.Context, page, pageSize int) ([]*model.Permission, int, error)

	// FindByIDs retrieves multiple permissions by their IDs
	FindByIDs(ctx context.Context, ids []string) ([]*model.Permission, error)

	// Delete removes a permission
	Delete(ctx context.Context, id string) error

	// ExistsByResourceAndAction checks if permission already exists
	ExistsByResourceAndAction(ctx context.Context, resource, action string) (bool, error)
}
