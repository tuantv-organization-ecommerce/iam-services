package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/internal/repository"
)

// AuthorizationService handles authorization business logic
type AuthorizationService interface {
	AssignRole(ctx context.Context, userID, roleID string) error
	RemoveRole(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*domain.Role, error)
	CheckPermission(ctx context.Context, userID, resource, action string) (bool, error)
}

type authorizationService struct {
	authzRepo repository.AuthorizationRepository
	userRepo  repository.UserRepository
	roleRepo  repository.RoleRepository
	logger    *zap.Logger
}

// NewAuthorizationService creates a new instance of AuthorizationService
func NewAuthorizationService(
	authzRepo repository.AuthorizationRepository,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	logger *zap.Logger,
) AuthorizationService {
	return &authorizationService{
		authzRepo: authzRepo,
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		logger:    logger,
	}
}

func (s *authorizationService) AssignRole(ctx context.Context, userID, roleID string) error {
	if userID == "" || roleID == "" {
		err := NewServiceError(ErrCodeInvalidParameters, "user ID and role ID are required", nil)
		LogError(s.logger, err, "AssignRole", zap.String("user_id", userID), zap.String("role_id", roleID))
		return err
	}

	// Verify user exists
	if _, err := s.userRepo.GetUserByID(ctx, userID); err != nil {
		serviceErr := NewServiceError(ErrCodeUserNotFound, "user not found", err)
		LogError(s.logger, serviceErr, "AssignRole", zap.String("user_id", userID), zap.String("role_id", roleID))
		return serviceErr
	}

	// Verify role exists
	if _, err := s.roleRepo.GetRoleByID(ctx, roleID); err != nil {
		serviceErr := NewServiceError(ErrCodeRoleNotFound, "role not found", err)
		LogError(s.logger, serviceErr, "AssignRole", zap.String("user_id", userID), zap.String("role_id", roleID))
		return serviceErr
	}

	if err := s.authzRepo.AssignRoleToUser(ctx, userID, roleID); err != nil {
		serviceErr := NewServiceError(ErrCodeRoleAssignmentFailed, "failed to assign role to user", err)
		LogError(s.logger, serviceErr, "AssignRole", zap.String("user_id", userID), zap.String("role_id", roleID))
		return serviceErr
	}

	s.logger.Info("Role assigned successfully",
		zap.String("user_id", userID),
		zap.String("role_id", roleID))

	return nil
}

func (s *authorizationService) RemoveRole(ctx context.Context, userID, roleID string) error {
	if userID == "" || roleID == "" {
		err := NewServiceError(ErrCodeInvalidParameters, "user ID and role ID are required", nil)
		LogError(s.logger, err, "RemoveRole", zap.String("user_id", userID), zap.String("role_id", roleID))
		return err
	}

	if err := s.authzRepo.RemoveRoleFromUser(ctx, userID, roleID); err != nil {
		serviceErr := NewServiceError(ErrCodeRoleRemovalFailed, "failed to remove role from user", err)
		LogError(s.logger, serviceErr, "RemoveRole", zap.String("user_id", userID), zap.String("role_id", roleID))
		return serviceErr
	}

	s.logger.Info("Role removed successfully",
		zap.String("user_id", userID),
		zap.String("role_id", roleID))

	return nil
}

func (s *authorizationService) GetUserRoles(ctx context.Context, userID string) ([]*domain.Role, error) {
	if userID == "" {
		err := NewServiceError(ErrCodeInvalidParameters, "user ID is required", nil)
		LogError(s.logger, err, "GetUserRoles", zap.String("user_id", userID))
		return nil, err
	}

	roles, err := s.authzRepo.GetUserRoles(ctx, userID)
	if err != nil {
		serviceErr := NewServiceError(ErrCodeGetRolesFailed, "failed to get user roles", err)
		LogError(s.logger, serviceErr, "GetUserRoles", zap.String("user_id", userID))
		return nil, serviceErr
	}

	// Get permissions for each role
	for i, role := range roles {
		permissions, err := s.authzRepo.GetRolePermissions(ctx, role.ID)
		if err != nil {
			serviceErr := NewServiceError(ErrCodeGetPermissionsFailed, "failed to get role permissions", err)
			LogError(s.logger, serviceErr, "GetUserRoles", zap.String("user_id", userID), zap.String("role_id", role.ID))
			return nil, serviceErr
		}
		// Convert []*Permission to []Permission
		roles[i].Permissions = make([]domain.Permission, len(permissions))
		for j, p := range permissions {
			if p != nil {
				roles[i].Permissions[j] = *p
			}
		}
	}

	s.logger.Debug("Retrieved user roles successfully",
		zap.String("user_id", userID),
		zap.Int("roles_count", len(roles)))

	return roles, nil
}

func (s *authorizationService) CheckPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	if userID == "" || resource == "" || action == "" {
		err := NewServiceError(ErrCodeInvalidParameters, "user ID, resource, and action are required", nil)
		LogError(s.logger, err, "CheckPermission", zap.String("user_id", userID), zap.String("resource", resource), zap.String("action", action))
		return false, err
	}

	hasPermission, err := s.authzRepo.UserHasPermission(ctx, userID, resource, action)
	if err != nil {
		serviceErr := NewServiceError(ErrCodePermissionCheckFailed, "failed to check user permission", err)
		LogError(s.logger, serviceErr, "CheckPermission", zap.String("user_id", userID), zap.String("resource", resource), zap.String("action", action))
		return false, serviceErr
	}

	s.logger.Debug("Permission check completed",
		zap.String("user_id", userID),
		zap.String("resource", resource),
		zap.String("action", action),
		zap.Bool("has_permission", hasPermission))

	return hasPermission, nil
}
