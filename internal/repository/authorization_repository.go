package repository

import (
	"context"
	"fmt"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/domain"
)

// AuthorizationRepository provides operations for authorization (user-role-permission relationships)
type AuthorizationRepository interface {
	AssignRoleToUser(ctx context.Context, userID, roleID string) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*domain.Role, error)
	GetUserPermissions(ctx context.Context, userID string) ([]*domain.Permission, error)
	UserHasPermission(ctx context.Context, userID, resource, action string) (bool, error)
	AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]*domain.Permission, error)
	UpdateRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error
}

type authorizationRepository struct {
	userRoleDAO       dao.UserRoleDAO
	rolePermissionDAO dao.RolePermissionDAO
	permissionDAO     dao.PermissionDAO
}

// NewAuthorizationRepository creates a new instance of AuthorizationRepository
func NewAuthorizationRepository(
	userRoleDAO dao.UserRoleDAO,
	rolePermissionDAO dao.RolePermissionDAO,
	permissionDAO dao.PermissionDAO,
) AuthorizationRepository {
	return &authorizationRepository{
		userRoleDAO:       userRoleDAO,
		rolePermissionDAO: rolePermissionDAO,
		permissionDAO:     permissionDAO,
	}
}

func (r *authorizationRepository) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	userRole := &domain.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.userRoleDAO.AssignRole(ctx, userRole)
}

func (r *authorizationRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	return r.userRoleDAO.RemoveRole(ctx, userID, roleID)
}

func (r *authorizationRepository) GetUserRoles(ctx context.Context, userID string) ([]*domain.Role, error) {
	return r.userRoleDAO.GetUserRoles(ctx, userID)
}

func (r *authorizationRepository) GetUserPermissions(ctx context.Context, userID string) ([]*domain.Permission, error) {
	// Get all roles for the user
	roles, err := r.userRoleDAO.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Collect all unique permissions from all roles
	permissionMap := make(map[string]*domain.Permission)
	for _, role := range roles {
		permissions, err := r.rolePermissionDAO.GetRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get role permissions: %w", err)
		}

		for _, perm := range permissions {
			permissionMap[perm.ID] = perm
		}
	}

	// Convert map to slice
	permissions := make([]*domain.Permission, 0, len(permissionMap))
	for _, perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (r *authorizationRepository) UserHasPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	// Get all user permissions
	permissions, err := r.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Check if user has the required permission
	for _, perm := range permissions {
		if perm.Resource == resource && perm.Action == action {
			return true, nil
		}
	}

	return false, nil
}

func (r *authorizationRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	rolePermission := &domain.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return r.rolePermissionDAO.AssignPermission(ctx, rolePermission)
}

func (r *authorizationRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	return r.rolePermissionDAO.RemovePermission(ctx, roleID, permissionID)
}

func (r *authorizationRepository) GetRolePermissions(ctx context.Context, roleID string) ([]*domain.Permission, error) {
	return r.rolePermissionDAO.GetRolePermissions(ctx, roleID)
}

func (r *authorizationRepository) UpdateRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	// Remove all existing permissions
	if err := r.rolePermissionDAO.RemoveAllPermissionsFromRole(ctx, roleID); err != nil {
		return fmt.Errorf("failed to remove existing permissions: %w", err)
	}

	// Add new permissions
	for _, permID := range permissionIDs {
		rolePermission := &domain.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
		}
		if err := r.rolePermissionDAO.AssignPermission(ctx, rolePermission); err != nil {
			return fmt.Errorf("failed to assign permission: %w", err)
		}
	}

	return nil
}
