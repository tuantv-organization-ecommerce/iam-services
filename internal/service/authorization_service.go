package service

import (
	"context"
	"fmt"

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
}

// NewAuthorizationService creates a new instance of AuthorizationService
func NewAuthorizationService(
	authzRepo repository.AuthorizationRepository,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
) AuthorizationService {
	return &authorizationService{
		authzRepo: authzRepo,
		userRepo:  userRepo,
		roleRepo:  roleRepo,
	}
}

func (s *authorizationService) AssignRole(ctx context.Context, userID, roleID string) error {
	if userID == "" || roleID == "" {
		return fmt.Errorf("user ID and role ID are required")
	}

	// Verify user exists
	if _, err := s.userRepo.GetUserByID(ctx, userID); err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify role exists
	if _, err := s.roleRepo.GetRoleByID(ctx, roleID); err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	return s.authzRepo.AssignRoleToUser(ctx, userID, roleID)
}

func (s *authorizationService) RemoveRole(ctx context.Context, userID, roleID string) error {
	if userID == "" || roleID == "" {
		return fmt.Errorf("user ID and role ID are required")
	}

	return s.authzRepo.RemoveRoleFromUser(ctx, userID, roleID)
}

func (s *authorizationService) GetUserRoles(ctx context.Context, userID string) ([]*domain.Role, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	roles, err := s.authzRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Get permissions for each role
	for i, role := range roles {
		permissions, err := s.authzRepo.GetRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get role permissions: %w", err)
		}
		// Convert []*Permission to []Permission
		roles[i].Permissions = make([]domain.Permission, len(permissions))
		for j, p := range permissions {
			if p != nil {
				roles[i].Permissions[j] = *p
			}
		}
	}

	return roles, nil
}

func (s *authorizationService) CheckPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	if userID == "" || resource == "" || action == "" {
		return false, fmt.Errorf("user ID, resource, and action are required")
	}

	return s.authzRepo.UserHasPermission(ctx, userID, resource, action)
}
