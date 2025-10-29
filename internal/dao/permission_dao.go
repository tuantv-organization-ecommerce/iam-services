package dao

import (
	"context"
	"database/sql"

	"github.com/tvttt/iam-services/internal/domain"
	"log"
)

// PermissionDAO defines the data access operations for Permission entity
type PermissionDAO interface {
	Create(ctx context.Context, permission *domain.Permission) error
	FindByID(ctx context.Context, id string) (*domain.Permission, error)
	FindByResourceAndAction(ctx context.Context, resource, action string) (*domain.Permission, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Permission, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type permissionDAO struct {
	db *sql.DB
}

// NewPermissionDAO creates a new instance of PermissionDAO
func NewPermissionDAO(db *sql.DB) PermissionDAO {
	return &permissionDAO{db: db}
}

func (d *permissionDAO) Create(ctx context.Context, permission *domain.Permission) error {
	query := `
		INSERT INTO permissions (id, name, resource, action, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := d.db.ExecContext(ctx, query,
		permission.ID,
		permission.Name,
		permission.Resource,
		permission.Action,
		permission.Description,
		permission.CreatedAt,
		permission.UpdatedAt,
	)
	return err
}

func (d *permissionDAO) FindByID(ctx context.Context, id string) (*domain.Permission, error) {
	query := `
		SELECT id, name, resource, action, description, created_at, updated_at
		FROM permissions
		WHERE id = $1
	`
	permission := &domain.Permission{}
	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&permission.ID,
		&permission.Name,
		&permission.Resource,
		&permission.Action,
		&permission.Description,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (d *permissionDAO) FindByResourceAndAction(ctx context.Context, resource, action string) (*domain.Permission, error) {
	query := `
		SELECT id, name, resource, action, description, created_at, updated_at
		FROM permissions
		WHERE resource = $1 AND action = $2
	`
	permission := &domain.Permission{}
	err := d.db.QueryRowContext(ctx, query, resource, action).Scan(
		&permission.ID,
		&permission.Name,
		&permission.Resource,
		&permission.Action,
		&permission.Description,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (d *permissionDAO) List(ctx context.Context, limit, offset int) ([]*domain.Permission, error) {
	query := `
		SELECT id, name, resource, action, description, created_at, updated_at
		FROM permissions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := d.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows:", err)
		}
	}()

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

func (d *permissionDAO) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM permissions WHERE id = $1`
	_, err := d.db.ExecContext(ctx, query, id)
	return err
}

func (d *permissionDAO) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM permissions`
	var count int
	err := d.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
