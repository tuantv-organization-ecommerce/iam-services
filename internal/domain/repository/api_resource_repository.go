package repository

import (
	"context"

	"github.com/tvttt/iam-services/internal/domain/model"
)

// APIResourceRepository handles API resource data access
type APIResourceRepository interface {
	// Save persists a new API resource
	Save(ctx context.Context, resource *model.APIResource) error

	// FindByID retrieves an API resource by ID
	FindByID(ctx context.Context, id string) (*model.APIResource, error)

	// FindByPathAndMethod retrieves an API resource by path and method
	FindByPathAndMethod(ctx context.Context, path string, method model.HTTPMethod) (*model.APIResource, error)

	// FindByService retrieves all API resources for a service
	FindByService(ctx context.Context, service string) ([]*model.APIResource, error)

	// FindAll retrieves all API resources with pagination
	FindAll(ctx context.Context, page, pageSize int) ([]*model.APIResource, int, error)

	// Update updates an existing API resource
	Update(ctx context.Context, resource *model.APIResource) error

	// Delete removes an API resource
	Delete(ctx context.Context, id string) error

	// ExistsByPathAndMethod checks if API resource already exists
	ExistsByPathAndMethod(ctx context.Context, path string, method model.HTTPMethod) (bool, error)
}
