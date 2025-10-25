package repository

import (
	"context"

	"github.com/tvttt/iam-services/internal/domain/model"
)

// CMSRepository handles CMS role data access
type CMSRepository interface {
	// Save persists a new CMS role
	Save(ctx context.Context, role *model.CMSRole) error

	// FindByID retrieves a CMS role by ID
	FindByID(ctx context.Context, id string) (*model.CMSRole, error)

	// FindByName retrieves a CMS role by name
	FindByName(ctx context.Context, name string) (*model.CMSRole, error)

	// FindAll retrieves all CMS roles with pagination
	FindAll(ctx context.Context, page, pageSize int) ([]*model.CMSRole, int, error)

	// Update updates an existing CMS role
	Update(ctx context.Context, role *model.CMSRole) error

	// Delete removes a CMS role
	Delete(ctx context.Context, id string) error

	// AssignToUser assigns a CMS role to a user
	AssignToUser(ctx context.Context, userID, roleID string) error

	// RemoveFromUser removes a CMS role from a user
	RemoveFromUser(ctx context.Context, userID, roleID string) error

	// GetUserCMSRoles retrieves all CMS roles for a user
	GetUserCMSRoles(ctx context.Context, userID string) ([]*model.CMSRole, error)

	// GetUserAccessibleTabs retrieves all accessible tabs for a user
	GetUserAccessibleTabs(ctx context.Context, userID string) ([]model.CMSTab, error)
}
