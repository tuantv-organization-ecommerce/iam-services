package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/internal/repository"
)

// RoleService handles role management business logic
type RoleService interface {
	CreateRole(ctx context.Context, name, description string, permissionIDs []string) (*domain.Role, error)
	UpdateRole(ctx context.Context, roleID, name, description string, permissionIDs []string) error
	DeleteRole(ctx context.Context, roleID string) error
	GetRole(ctx context.Context, roleID string) (*domain.Role, error)
	ListRoles(ctx context.Context, page, pageSize int) ([]*domain.Role, int, error)
}

type roleService struct {
	roleRepo  repository.RoleRepository
	authzRepo repository.AuthorizationRepository
}

// NewRoleService creates a new instance of RoleService
func NewRoleService(roleRepo repository.RoleRepository, authzRepo repository.AuthorizationRepository) RoleService {
	return &roleService{
		roleRepo:  roleRepo,
		authzRepo: authzRepo,
	}
}

func (s *roleService) CreateRole(ctx context.Context, name, description string, permissionIDs []string) (*domain.Role, error) {
	if name == "" {
		return nil, fmt.Errorf("role name is required")
	}

	role := &domain.Role{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.roleRepo.CreateRole(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Assign permissions to role
	if len(permissionIDs) > 0 {
		if err := s.authzRepo.UpdateRolePermissions(ctx, role.ID, permissionIDs); err != nil {
			return nil, fmt.Errorf("failed to assign permissions to role: %w", err)
		}
	}

	return role, nil
}

func (s *roleService) UpdateRole(ctx context.Context, roleID, name, description string, permissionIDs []string) error {
	if roleID == "" {
		return fmt.Errorf("role ID is required")
	}

	// Get existing role
	role, err := s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Update role fields
	if name != "" {
		role.Name = name
	}
	role.Description = description
	role.UpdatedAt = time.Now()

	if err := s.roleRepo.UpdateRole(ctx, role); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	// Update permissions
	if permissionIDs != nil {
		if err := s.authzRepo.UpdateRolePermissions(ctx, roleID, permissionIDs); err != nil {
			return fmt.Errorf("failed to update role permissions: %w", err)
		}
	}

	return nil
}

func (s *roleService) DeleteRole(ctx context.Context, roleID string) error {
	if roleID == "" {
		return fmt.Errorf("role ID is required")
	}

	return s.roleRepo.DeleteRole(ctx, roleID)
}

func (s *roleService) GetRole(ctx context.Context, roleID string) (*domain.Role, error) {
	if roleID == "" {
		return nil, fmt.Errorf("role ID is required")
	}

	role, err := s.roleRepo.GetRoleWithPermissions(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

func (s *roleService) ListRoles(ctx context.Context, page, pageSize int) ([]*domain.Role, int, error) {
	return s.roleRepo.ListRoles(ctx, page, pageSize)
}
