package repository

import (
	"context"

	"github.com/tvttt/iam-services/internal/domain/model"
)

// AuthorizationRepository handles user-role-permission relationships
type AuthorizationRepository interface {
	// AssignRoleToUser assigns a role to a user
	AssignRoleToUser(ctx context.Context, userID, roleID string) error

	// RemoveRoleFromUser removes a role from a user
	RemoveRoleFromUser(ctx context.Context, userID, roleID string) error

	// GetUserRoles retrieves all roles for a user
	GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error)

	// GetUserPermissions retrieves all permissions for a user (through roles)
	GetUserPermissions(ctx context.Context, userID string) ([]*model.Permission, error)

	// AssignPermissionToRole assigns a permission to a role
	AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error

	// RemovePermissionFromRole removes a permission from a role
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error

	// GetRolePermissions retrieves all permissions for a role
	GetRolePermissions(ctx context.Context, roleID string) ([]*model.Permission, error)

	// UpdateRolePermissions replaces all permissions for a role
	UpdateRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error

	// UserHasPermission checks if user has a specific permission
	UserHasPermission(ctx context.Context, userID, resource, action string) (bool, error)
}
