package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/internal/repository"
)

// PermissionService handles permission management business logic
type PermissionService interface {
	CreatePermission(ctx context.Context, name, resource, action, description string) (*domain.Permission, error)
	DeletePermission(ctx context.Context, permissionID string) error
	ListPermissions(ctx context.Context, page, pageSize int) ([]*domain.Permission, int, error)
}

type permissionService struct {
	permissionRepo repository.PermissionRepository
}

// NewPermissionService creates a new instance of PermissionService
func NewPermissionService(permissionRepo repository.PermissionRepository) PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
	}
}

func (s *permissionService) CreatePermission(ctx context.Context, name, resource, action, description string) (*domain.Permission, error) {
	if name == "" || resource == "" || action == "" {
		return nil, fmt.Errorf("name, resource, and action are required")
	}

	permission := &domain.Permission{
		ID:          uuid.New().String(),
		Name:        name,
		Resource:    resource,
		Action:      action,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.permissionRepo.CreatePermission(ctx, permission); err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return permission, nil
}

func (s *permissionService) DeletePermission(ctx context.Context, permissionID string) error {
	if permissionID == "" {
		return fmt.Errorf("permission ID is required")
	}

	return s.permissionRepo.DeletePermission(ctx, permissionID)
}

func (s *permissionService) ListPermissions(ctx context.Context, page, pageSize int) ([]*domain.Permission, int, error) {
	return s.permissionRepo.ListPermissions(ctx, page, pageSize)
}
