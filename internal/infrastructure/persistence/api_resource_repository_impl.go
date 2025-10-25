package persistence

import (
	"context"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/internal/domain/model"
	domainRepo "github.com/tvttt/iam-services/internal/domain/repository"
)

// apiResourceRepositoryImpl implements domain.repository.APIResourceRepository using DAO
type apiResourceRepositoryImpl struct {
	apiResourceDAO dao.APIResourceDAO
}

// NewAPIResourceRepository creates a new API resource repository implementation
func NewAPIResourceRepository(apiResourceDAO dao.APIResourceDAO) domainRepo.APIResourceRepository {
	return &apiResourceRepositoryImpl{
		apiResourceDAO: apiResourceDAO,
	}
}

func (r *apiResourceRepositoryImpl) Save(ctx context.Context, resource *model.APIResource) error {
	daoResource := r.modelToDAO(resource)
	return r.apiResourceDAO.Create(ctx, daoResource)
}

func (r *apiResourceRepositoryImpl) FindByID(ctx context.Context, id string) (*model.APIResource, error) {
	daoResource, err := r.apiResourceDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if daoResource == nil {
		return nil, nil
	}

	return r.daoToModel(daoResource), nil
}

func (r *apiResourceRepositoryImpl) FindByPathAndMethod(ctx context.Context, path string, method model.HTTPMethod) (*model.APIResource, error) {
	daoResource, err := r.apiResourceDAO.FindByPathAndMethod(ctx, path, string(method))
	if err != nil {
		return nil, err
	}

	if daoResource == nil {
		return nil, nil
	}

	return r.daoToModel(daoResource), nil
}

func (r *apiResourceRepositoryImpl) FindByService(ctx context.Context, service string) ([]*model.APIResource, error) {
	daoResources, err := r.apiResourceDAO.ListByService(ctx, service)
	if err != nil {
		return nil, err
	}

	resources := make([]*model.APIResource, len(daoResources))
	for i, daoResource := range daoResources {
		resources[i] = r.daoToModel(daoResource)
	}

	return resources, nil
}

func (r *apiResourceRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]*model.APIResource, int, error) {
	offset := (page - 1) * pageSize
	daoResources, err := r.apiResourceDAO.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	resources := make([]*model.APIResource, len(daoResources))
	for i, daoResource := range daoResources {
		resources[i] = r.daoToModel(daoResource)
	}

	// Get total count
	count, err := r.apiResourceDAO.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return resources, count, nil
}

func (r *apiResourceRepositoryImpl) Update(ctx context.Context, resource *model.APIResource) error {
	daoResource := r.modelToDAO(resource)
	return r.apiResourceDAO.Update(ctx, daoResource)
}

func (r *apiResourceRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.apiResourceDAO.Delete(ctx, id)
}

func (r *apiResourceRepositoryImpl) ExistsByPathAndMethod(ctx context.Context, path string, method model.HTTPMethod) (bool, error) {
	resource, err := r.FindByPathAndMethod(ctx, path, method)
	if err != nil {
		return false, err
	}
	return resource != nil, nil
}

// Converters

func (r *apiResourceRepositoryImpl) modelToDAO(resource *model.APIResource) *domain.APIResource {
	return &domain.APIResource{
		ID:          resource.ID(),
		Path:        resource.Path(),
		Method:      string(resource.Method()),
		Service:     resource.Service(),
		Description: resource.Description(),
		CreatedAt:   resource.CreatedAt(),
		UpdatedAt:   resource.UpdatedAt(),
	}
}

func (r *apiResourceRepositoryImpl) daoToModel(daoResource *domain.APIResource) *model.APIResource {
	if daoResource == nil {
		return nil
	}

	return model.ReconstructAPIResource(
		daoResource.ID,
		daoResource.Path,
		model.HTTPMethod(daoResource.Method),
		daoResource.Service,
		daoResource.Description,
		daoResource.CreatedAt,
		daoResource.UpdatedAt,
	)
}
