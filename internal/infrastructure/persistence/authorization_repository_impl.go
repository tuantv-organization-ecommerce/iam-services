package persistence

import (
	"context"
	"time"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/internal/domain/model"
	domainRepo "github.com/tvttt/iam-services/internal/domain/repository"
)

// authorizationRepositoryImpl implements domain.repository.AuthorizationRepository using DAO
type authorizationRepositoryImpl struct {
	userRoleDAO       dao.UserRoleDAO
	rolePermissionDAO dao.RolePermissionDAO
}

// NewAuthorizationRepository creates a new authorization repository implementation
func NewAuthorizationRepository(
	userRoleDAO dao.UserRoleDAO,
	rolePermissionDAO dao.RolePermissionDAO,
) domainRepo.AuthorizationRepository {
	return &authorizationRepositoryImpl{
		userRoleDAO:       userRoleDAO,
		rolePermissionDAO: rolePermissionDAO,
	}
}

func (r *authorizationRepositoryImpl) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	userRole := &domain.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: time.Now(),
	}
	return r.userRoleDAO.AssignRole(ctx, userRole)
}

func (r *authorizationRepositoryImpl) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	return r.userRoleDAO.RemoveRole(ctx, userID, roleID)
}

func (r *authorizationRepositoryImpl) GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error) {
	daoRoles, err := r.userRoleDAO.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	roles := make([]*model.Role, len(daoRoles))
	for i, daoRole := range daoRoles {
		roles[i] = r.daoToModelRole(daoRole)
	}

	return roles, nil
}

func (r *authorizationRepositoryImpl) GetUserPermissions(ctx context.Context, userID string) ([]*model.Permission, error) {
	// Get user roles
	daoRoles, err := r.userRoleDAO.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Collect all permissions from all roles (deduplicate)
	permMap := make(map[string]*model.Permission)
	for _, role := range daoRoles {
		rolePerms, err := r.rolePermissionDAO.GetRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, err
		}

		for _, daoPerm := range rolePerms {
			if _, exists := permMap[daoPerm.ID]; !exists {
				permMap[daoPerm.ID] = r.daoToModelPerm(daoPerm)
			}
		}
	}

	// Convert map to slice
	perms := make([]*model.Permission, 0, len(permMap))
	for _, perm := range permMap {
		perms = append(perms, perm)
	}

	return perms, nil
}

func (r *authorizationRepositoryImpl) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	rolePermission := &domain.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
		CreatedAt:    time.Now(),
	}
	return r.rolePermissionDAO.AssignPermission(ctx, rolePermission)
}

func (r *authorizationRepositoryImpl) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	return r.rolePermissionDAO.RemovePermission(ctx, roleID, permissionID)
}

func (r *authorizationRepositoryImpl) GetRolePermissions(ctx context.Context, roleID string) ([]*model.Permission, error) {
	daoPerms, err := r.rolePermissionDAO.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	perms := make([]*model.Permission, len(daoPerms))
	for i, daoPerm := range daoPerms {
		perms[i] = r.daoToModelPerm(daoPerm)
	}

	return perms, nil
}

func (r *authorizationRepositoryImpl) UpdateRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	// Remove all existing permissions
	if err := r.rolePermissionDAO.RemoveAllPermissionsFromRole(ctx, roleID); err != nil {
		return err
	}

	// Add new permissions
	for _, permID := range permissionIDs {
		rolePermission := &domain.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
			CreatedAt:    time.Now(),
		}
		if err := r.rolePermissionDAO.AssignPermission(ctx, rolePermission); err != nil {
			return err
		}
	}

	return nil
}

func (r *authorizationRepositoryImpl) UserHasPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	permissions, err := r.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, perm := range permissions {
		if perm.Matches(resource, action) {
			return true, nil
		}
	}

	return false, nil
}

// Converters

func (r *authorizationRepositoryImpl) daoToModelRole(daoRole *domain.Role) *model.Role {
	if daoRole == nil {
		return nil
	}

	return model.ReconstructRole(
		daoRole.ID,
		daoRole.Name,
		daoRole.Description,
		"user", // default domain
		daoRole.CreatedAt,
		daoRole.UpdatedAt,
	)
}

func (r *authorizationRepositoryImpl) daoToModelPerm(daoPerm *domain.Permission) *model.Permission {
	if daoPerm == nil {
		return nil
	}

	return model.ReconstructPermission(
		daoPerm.ID,
		daoPerm.Name,
		daoPerm.Resource,
		daoPerm.Action,
		daoPerm.Description,
		daoPerm.CreatedAt,
		daoPerm.UpdatedAt,
	)
}
