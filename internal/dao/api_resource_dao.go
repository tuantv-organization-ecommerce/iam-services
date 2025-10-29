package dao

import (
	"context"
	"database/sql"

	"github.com/tvttt/iam-services/internal/domain"
)

// APIResourceDAO defines the data access operations for API resources
type APIResourceDAO interface {
	Create(ctx context.Context, resource *domain.APIResource) error
	FindByID(ctx context.Context, id string) (*domain.APIResource, error)
	FindByPathAndMethod(ctx context.Context, path, method string) (*domain.APIResource, error)
	ListByService(ctx context.Context, service string) ([]*domain.APIResource, error)
	List(ctx context.Context, limit, offset int) ([]*domain.APIResource, error)
	Update(ctx context.Context, resource *domain.APIResource) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type apiResourceDAO struct {
	db *sql.DB
}

// NewAPIResourceDAO creates a new instance of APIResourceDAO
func NewAPIResourceDAO(db *sql.DB) APIResourceDAO {
	return &apiResourceDAO{db: db}
}

func (d *apiResourceDAO) Create(ctx context.Context, resource *domain.APIResource) error {
	query := `
		INSERT INTO api_resources (id, path, method, service, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := d.db.ExecContext(ctx, query,
		resource.ID,
		resource.Path,
		resource.Method,
		resource.Service,
		resource.Description,
		resource.CreatedAt,
		resource.UpdatedAt,
	)
	return err
}

func (d *apiResourceDAO) FindByID(ctx context.Context, id string) (*domain.APIResource, error) {
	query := `
		SELECT id, path, method, service, description, created_at, updated_at
		FROM api_resources
		WHERE id = $1
	`
	resource := &domain.APIResource{}
	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&resource.ID,
		&resource.Path,
		&resource.Method,
		&resource.Service,
		&resource.Description,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return resource, err
}

func (d *apiResourceDAO) FindByPathAndMethod(ctx context.Context, path, method string) (*domain.APIResource, error) {
	query := `
		SELECT id, path, method, service, description, created_at, updated_at
		FROM api_resources
		WHERE path = $1 AND method = $2
	`
	resource := &domain.APIResource{}
	err := d.db.QueryRowContext(ctx, query, path, method).Scan(
		&resource.ID,
		&resource.Path,
		&resource.Method,
		&resource.Service,
		&resource.Description,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return resource, err
}

func (d *apiResourceDAO) ListByService(ctx context.Context, service string) ([]*domain.APIResource, error) {
	query := `
		SELECT id, path, method, service, description, created_at, updated_at
		FROM api_resources
		WHERE service = $1
		ORDER BY path, method
	`
	rows, err := d.db.QueryContext(ctx, query, service)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []*domain.APIResource
	for rows.Next() {
		resource := &domain.APIResource{}
		err := rows.Scan(
			&resource.ID,
			&resource.Path,
			&resource.Method,
			&resource.Service,
			&resource.Description,
			&resource.CreatedAt,
			&resource.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, rows.Err()
}

func (d *apiResourceDAO) List(ctx context.Context, limit, offset int) ([]*domain.APIResource, error) {
	query := `
		SELECT id, path, method, service, description, created_at, updated_at
		FROM api_resources
		ORDER BY service, path, method
		LIMIT $1 OFFSET $2
	`
	rows, err := d.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []*domain.APIResource
	for rows.Next() {
		resource := &domain.APIResource{}
		err := rows.Scan(
			&resource.ID,
			&resource.Path,
			&resource.Method,
			&resource.Service,
			&resource.Description,
			&resource.CreatedAt,
			&resource.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, rows.Err()
}

func (d *apiResourceDAO) Update(ctx context.Context, resource *domain.APIResource) error {
	query := `
		UPDATE api_resources
		SET path = $2, method = $3, service = $4, description = $5, updated_at = $6
		WHERE id = $1
	`
	_, err := d.db.ExecContext(ctx, query,
		resource.ID,
		resource.Path,
		resource.Method,
		resource.Service,
		resource.Description,
		resource.UpdatedAt,
	)
	return err
}

func (d *apiResourceDAO) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM api_resources WHERE id = $1`
	_, err := d.db.ExecContext(ctx, query, id)
	return err
}

func (d *apiResourceDAO) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM api_resources`
	var count int
	err := d.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
