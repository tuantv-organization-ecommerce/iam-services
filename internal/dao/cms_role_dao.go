package dao

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/tvttt/iam-services/internal/domain"
)

// CMSRoleDAO defines the data access operations for CMS roles
type CMSRoleDAO interface {
	Create(ctx context.Context, role *domain.CMSRole) error
	FindByID(ctx context.Context, id string) (*domain.CMSRole, error)
	FindByName(ctx context.Context, name string) (*domain.CMSRole, error)
	List(ctx context.Context, limit, offset int) ([]*domain.CMSRole, error)
	Update(ctx context.Context, role *domain.CMSRole) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type cmsRoleDAO struct {
	db *sql.DB
}

// NewCMSRoleDAO creates a new instance of CMSRoleDAO
func NewCMSRoleDAO(db *sql.DB) CMSRoleDAO {
	return &cmsRoleDAO{db: db}
}

func (d *cmsRoleDAO) Create(ctx context.Context, role *domain.CMSRole) error {
	query := `
		INSERT INTO cms_roles (id, name, description, tabs, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	tabs := make([]string, len(role.Tabs))
	for i, tab := range role.Tabs {
		tabs[i] = string(tab)
	}

	_, err := d.db.ExecContext(ctx, query,
		role.ID,
		role.Name,
		role.Description,
		pq.Array(tabs),
		role.CreatedAt,
		role.UpdatedAt,
	)
	return err
}

func (d *cmsRoleDAO) FindByID(ctx context.Context, id string) (*domain.CMSRole, error) {
	query := `
		SELECT id, name, description, tabs, created_at, updated_at
		FROM cms_roles
		WHERE id = $1
	`
	role := &domain.CMSRole{}
	var tabs []string

	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		pq.Array(&tabs),
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	role.Tabs = make([]domain.CMSTab, len(tabs))
	for i, tab := range tabs {
		role.Tabs[i] = domain.CMSTab(tab)
	}

	return role, nil
}

func (d *cmsRoleDAO) FindByName(ctx context.Context, name string) (*domain.CMSRole, error) {
	query := `
		SELECT id, name, description, tabs, created_at, updated_at
		FROM cms_roles
		WHERE name = $1
	`
	role := &domain.CMSRole{}
	var tabs []string

	err := d.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		pq.Array(&tabs),
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	role.Tabs = make([]domain.CMSTab, len(tabs))
	for i, tab := range tabs {
		role.Tabs[i] = domain.CMSTab(tab)
	}

	return role, nil
}

func (d *cmsRoleDAO) List(ctx context.Context, limit, offset int) ([]*domain.CMSRole, error) {
	query := `
		SELECT id, name, description, tabs, created_at, updated_at
		FROM cms_roles
		ORDER BY name
		LIMIT $1 OFFSET $2
	`
	rows, err := d.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var roles []*domain.CMSRole
	for rows.Next() {
		role := &domain.CMSRole{}
		var tabs []string

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			pq.Array(&tabs),
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		role.Tabs = make([]domain.CMSTab, len(tabs))
		for i, tab := range tabs {
			role.Tabs[i] = domain.CMSTab(tab)
		}

		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (d *cmsRoleDAO) Update(ctx context.Context, role *domain.CMSRole) error {
	query := `
		UPDATE cms_roles
		SET name = $2, description = $3, tabs = $4, updated_at = $5
		WHERE id = $1
	`
	tabs := make([]string, len(role.Tabs))
	for i, tab := range role.Tabs {
		tabs[i] = string(tab)
	}

	_, err := d.db.ExecContext(ctx, query,
		role.ID,
		role.Name,
		role.Description,
		pq.Array(tabs),
		role.UpdatedAt,
	)
	return err
}

func (d *cmsRoleDAO) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM cms_roles WHERE id = $1`
	_, err := d.db.ExecContext(ctx, query, id)
	return err
}

func (d *cmsRoleDAO) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM cms_roles`
	var count int
	err := d.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
