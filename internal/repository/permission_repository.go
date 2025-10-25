package repository

import (
	"context"
	"fmt"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/domain"
)

// PermissionRepository provides higher-level operations on Permission entity
type PermissionRepository interface {
	CreatePermission(ctx context.Context, permission *domain.Permission) error
	GetPermissionByID(ctx context.Context, id string) (*domain.Permission, error)
	GetPermissionByResourceAndAction(ctx context.Context, resource, action string) (*domain.Permission, error)
	ListPermissions(ctx context.Context, page, pageSize int) ([]*domain.Permission, int, error)
	DeletePermission(ctx context.Context, id string) error
}

type permissionRepository struct {
	permissionDAO dao.PermissionDAO
}

// NewPermissionRepository creates a new instance of PermissionRepository
func NewPermissionRepository(permissionDAO dao.PermissionDAO) PermissionRepository {
	return &permissionRepository{
		permissionDAO: permissionDAO,
	}
}

func (r *permissionRepository) CreatePermission(ctx context.Context, permission *domain.Permission) error {
	// Check if permission already exists
	existingPerm, err := r.permissionDAO.FindByResourceAndAction(ctx, permission.Resource, permission.Action)
	if err != nil {
		return fmt.Errorf("failed to check permission existence: %w", err)
	}
	if existingPerm != nil {
		return fmt.Errorf("permission for resource '%s' and action '%s' already exists", permission.Resource, permission.Action)
	}

	return r.permissionDAO.Create(ctx, permission)
}

func (r *permissionRepository) GetPermissionByID(ctx context.Context, id string) (*domain.Permission, error) {
	permission, err := r.permissionDAO.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission by id: %w", err)
	}
	if permission == nil {
		return nil, fmt.Errorf("permission not found")
	}
	return permission, nil
}

func (r *permissionRepository) GetPermissionByResourceAndAction(ctx context.Context, resource, action string) (*domain.Permission, error) {
	permission, err := r.permissionDAO.FindByResourceAndAction(ctx, resource, action)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	if permission == nil {
		return nil, fmt.Errorf("permission not found")
	}
	return permission, nil
}

func (r *permissionRepository) ListPermissions(ctx context.Context, page, pageSize int) ([]*domain.Permission, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	permissions, err := r.permissionDAO.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list permissions: %w", err)
	}

	total, err := r.permissionDAO.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count permissions: %w", err)
	}

	return permissions, total, nil
}

func (r *permissionRepository) DeletePermission(ctx context.Context, id string) error {
	return r.permissionDAO.Delete(ctx, id)
}
