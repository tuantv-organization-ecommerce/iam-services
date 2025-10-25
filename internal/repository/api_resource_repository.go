package repository

import (
	"context"
	"fmt"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/domain"
)

// APIResourceRepository provides operations for API resource management
type APIResourceRepository interface {
	CreateAPIResource(ctx context.Context, resource *domain.APIResource) error
	GetAPIResourceByID(ctx context.Context, id string) (*domain.APIResource, error)
	GetAPIResourceByPathAndMethod(ctx context.Context, path, method string) (*domain.APIResource, error)
	ListAPIResourcesByService(ctx context.Context, service string) ([]*domain.APIResource, error)
	ListAPIResources(ctx context.Context, page, pageSize int) ([]*domain.APIResource, int, error)
	UpdateAPIResource(ctx context.Context, resource *domain.APIResource) error
	DeleteAPIResource(ctx context.Context, id string) error
}

type apiResourceRepository struct {
	apiResourceDAO dao.APIResourceDAO
}

// NewAPIResourceRepository creates a new instance of APIResourceRepository
func NewAPIResourceRepository(apiResourceDAO dao.APIResourceDAO) APIResourceRepository {
	return &apiResourceRepository{
		apiResourceDAO: apiResourceDAO,
	}
}

func (r *apiResourceRepository) CreateAPIResource(ctx context.Context, resource *domain.APIResource) error {
	// Check if resource already exists
	existing, err := r.apiResourceDAO.FindByPathAndMethod(ctx, resource.Path, resource.Method)
	if err != nil {
		return fmt.Errorf("failed to check API resource existence: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("API resource with path '%s' and method '%s' already exists", resource.Path, resource.Method)
	}

	return r.apiResourceDAO.Create(ctx, resource)
}

func (r *apiResourceRepository) GetAPIResourceByID(ctx context.Context, id string) (*domain.APIResource, error) {
	resource, err := r.apiResourceDAO.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get API resource by ID: %w", err)
	}
	if resource == nil {
		return nil, fmt.Errorf("API resource not found")
	}
	return resource, nil
}

func (r *apiResourceRepository) GetAPIResourceByPathAndMethod(ctx context.Context, path, method string) (*domain.APIResource, error) {
	resource, err := r.apiResourceDAO.FindByPathAndMethod(ctx, path, method)
	if err != nil {
		return nil, fmt.Errorf("failed to get API resource: %w", err)
	}
	if resource == nil {
		return nil, fmt.Errorf("API resource not found")
	}
	return resource, nil
}

func (r *apiResourceRepository) ListAPIResourcesByService(ctx context.Context, service string) ([]*domain.APIResource, error) {
	return r.apiResourceDAO.ListByService(ctx, service)
}

func (r *apiResourceRepository) ListAPIResources(ctx context.Context, page, pageSize int) ([]*domain.APIResource, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	resources, err := r.apiResourceDAO.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list API resources: %w", err)
	}

	total, err := r.apiResourceDAO.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count API resources: %w", err)
	}

	return resources, total, nil
}

func (r *apiResourceRepository) UpdateAPIResource(ctx context.Context, resource *domain.APIResource) error {
	return r.apiResourceDAO.Update(ctx, resource)
}

func (r *apiResourceRepository) DeleteAPIResource(ctx context.Context, id string) error {
	return r.apiResourceDAO.Delete(ctx, id)
}
