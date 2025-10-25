package repository

import (
	"context"
	"fmt"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/domain"
)

// RoleRepository provides higher-level operations on Role entity
type RoleRepository interface {
	CreateRole(ctx context.Context, role *domain.Role) error
	GetRoleByID(ctx context.Context, id string) (*domain.Role, error)
	GetRoleByName(ctx context.Context, name string) (*domain.Role, error)
	GetRoleWithPermissions(ctx context.Context, id string) (*domain.Role, error)
	ListRoles(ctx context.Context, page, pageSize int) ([]*domain.Role, int, error)
	UpdateRole(ctx context.Context, role *domain.Role) error
	DeleteRole(ctx context.Context, id string) error
}

type roleRepository struct {
	roleDAO           dao.RoleDAO
	rolePermissionDAO dao.RolePermissionDAO
}

// NewRoleRepository creates a new instance of RoleRepository
func NewRoleRepository(roleDAO dao.RoleDAO, rolePermissionDAO dao.RolePermissionDAO) RoleRepository {
	return &roleRepository{
		roleDAO:           roleDAO,
		rolePermissionDAO: rolePermissionDAO,
	}
}

func (r *roleRepository) CreateRole(ctx context.Context, role *domain.Role) error {
	// Check if role already exists
	existingRole, err := r.roleDAO.FindByName(ctx, role.Name)
	if err != nil {
		return fmt.Errorf("failed to check role existence: %w", err)
	}
	if existingRole != nil {
		return fmt.Errorf("role with name '%s' already exists", role.Name)
	}

	return r.roleDAO.Create(ctx, role)
}

func (r *roleRepository) GetRoleByID(ctx context.Context, id string) (*domain.Role, error) {
	role, err := r.roleDAO.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by id: %w", err)
	}
	if role == nil {
		return nil, fmt.Errorf("role not found")
	}
	return role, nil
}

func (r *roleRepository) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	role, err := r.roleDAO.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}
	if role == nil {
		return nil, fmt.Errorf("role not found")
	}
	return role, nil
}

func (r *roleRepository) GetRoleWithPermissions(ctx context.Context, id string) (*domain.Role, error) {
	role, err := r.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	permissions, err := r.rolePermissionDAO.GetRolePermissions(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	role.Permissions = permissions
	return role, nil
}

func (r *roleRepository) ListRoles(ctx context.Context, page, pageSize int) ([]*domain.Role, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	roles, err := r.roleDAO.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list roles: %w", err)
	}

	total, err := r.roleDAO.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count roles: %w", err)
	}

	return roles, total, nil
}

func (r *roleRepository) UpdateRole(ctx context.Context, role *domain.Role) error {
	return r.roleDAO.Update(ctx, role)
}

func (r *roleRepository) DeleteRole(ctx context.Context, id string) error {
	return r.roleDAO.Delete(ctx, id)
}
