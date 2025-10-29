package dao

import (
	"context"
	"database/sql"

	"github.com/tvttt/iam-services/internal/domain"
)

// RolePermissionDAO defines the data access operations for RolePermission relationship
type RolePermissionDAO interface {
	AssignPermission(ctx context.Context, rolePermission *domain.RolePermission) error
	RemovePermission(ctx context.Context, roleID, permissionID string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]*domain.Permission, error)
	GetPermissionRoles(ctx context.Context, permissionID string) ([]*domain.Role, error)
	HasPermission(ctx context.Context, roleID, permissionID string) (bool, error)
	RemoveAllPermissionsFromRole(ctx context.Context, roleID string) error
}

type rolePermissionDAO struct {
	db *sql.DB
}

// NewRolePermissionDAO creates a new instance of RolePermissionDAO
func NewRolePermissionDAO(db *sql.DB) RolePermissionDAO {
	return &rolePermissionDAO{db: db}
}

func (d *rolePermissionDAO) AssignPermission(ctx context.Context, rolePermission *domain.RolePermission) error {
	query := `
		INSERT INTO role_permissions (role_id, permission_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (role_id, permission_id) DO NOTHING
	`
	_, err := d.db.ExecContext(ctx, query,
		rolePermission.RoleID,
		rolePermission.PermissionID,
		rolePermission.CreatedAt,
	)
	return err
}

func (d *rolePermissionDAO) RemovePermission(ctx context.Context, roleID, permissionID string) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`
	_, err := d.db.ExecContext(ctx, query, roleID, permissionID)
	return err
}

func (d *rolePermissionDAO) GetRolePermissions(ctx context.Context, roleID string) ([]*domain.Permission, error) {
	query := `
		SELECT p.id, p.name, p.resource, p.action, p.description, p.created_at, p.updated_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.resource, p.action
	`
	rows, err := d.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var permissions []*domain.Permission
	for rows.Next() {
		permission := &domain.Permission{}
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Resource,
			&permission.Action,
			&permission.Description,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, rows.Err()
}

func (d *rolePermissionDAO) GetPermissionRoles(ctx context.Context, permissionID string) ([]*domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN role_permissions rp ON r.id = rp.role_id
		WHERE rp.permission_id = $1
		ORDER BY r.name
	`
	rows, err := d.db.QueryContext(ctx, query, permissionID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var roles []*domain.Role
	for rows.Next() {
		role := &domain.Role{}
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (d *rolePermissionDAO) HasPermission(ctx context.Context, roleID, permissionID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM role_permissions WHERE role_id = $1 AND permission_id = $2)`
	var exists bool
	err := d.db.QueryRowContext(ctx, query, roleID, permissionID).Scan(&exists)
	return exists, err
}

func (d *rolePermissionDAO) RemoveAllPermissionsFromRole(ctx context.Context, roleID string) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1`
	_, err := d.db.ExecContext(ctx, query, roleID)
	return err
}
