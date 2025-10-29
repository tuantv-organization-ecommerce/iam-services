package dao

import (
	"context"
	"database/sql"

	"github.com/tvttt/iam-services/internal/domain"
)

// RoleDAO defines the data access operations for Role entity
type RoleDAO interface {
	Create(ctx context.Context, role *domain.Role) error
	FindByID(ctx context.Context, id string) (*domain.Role, error)
	FindByName(ctx context.Context, name string) (*domain.Role, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Role, error)
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type roleDAO struct {
	db *sql.DB
}

// NewRoleDAO creates a new instance of RoleDAO
func NewRoleDAO(db *sql.DB) RoleDAO {
	return &roleDAO{db: db}
}

func (d *roleDAO) Create(ctx context.Context, role *domain.Role) error {
	query := `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := d.db.ExecContext(ctx, query,
		role.ID,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.UpdatedAt,
	)
	return err
}

func (d *roleDAO) FindByID(ctx context.Context, id string) (*domain.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = $1
	`
	role := &domain.Role{}
	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (d *roleDAO) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE name = $1
	`
	role := &domain.Role{}
	err := d.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (d *roleDAO) List(ctx context.Context, limit, offset int) ([]*domain.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := d.db.QueryContext(ctx, query, limit, offset)
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

func (d *roleDAO) Update(ctx context.Context, role *domain.Role) error {
	query := `
		UPDATE roles
		SET name = $2, description = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := d.db.ExecContext(ctx, query,
		role.ID,
		role.Name,
		role.Description,
		role.UpdatedAt,
	)
	return err
}

func (d *roleDAO) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM roles WHERE id = $1`
	_, err := d.db.ExecContext(ctx, query, id)
	return err
}

func (d *roleDAO) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM roles`
	var count int
	err := d.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
