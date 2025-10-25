package persistence

import (
	"context"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/internal/domain/model"
	domainRepo "github.com/tvttt/iam-services/internal/domain/repository"
)

// permissionRepositoryImpl implements domain.repository.PermissionRepository using DAO
type permissionRepositoryImpl struct {
	permissionDAO dao.PermissionDAO
}

// NewPermissionRepository creates a new permission repository implementation
func NewPermissionRepository(permissionDAO dao.PermissionDAO) domainRepo.PermissionRepository {
	return &permissionRepositoryImpl{
		permissionDAO: permissionDAO,
	}
}

func (r *permissionRepositoryImpl) Save(ctx context.Context, perm *model.Permission) error {
	daoPerm := r.modelToDAO(perm)
	return r.permissionDAO.Create(ctx, daoPerm)
}

func (r *permissionRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Permission, error) {
	daoPerm, err := r.permissionDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if daoPerm == nil {
		return nil, nil
	}

	return r.daoToModel(daoPerm), nil
}

func (r *permissionRepositoryImpl) FindByResourceAndAction(ctx context.Context, resource, action string) (*model.Permission, error) {
	daoPerm, err := r.permissionDAO.FindByResourceAndAction(ctx, resource, action)
	if err != nil {
		return nil, err
	}

	if daoPerm == nil {
		return nil, nil
	}

	return r.daoToModel(daoPerm), nil
}

func (r *permissionRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]*model.Permission, int, error) {
	offset := (page - 1) * pageSize
	daoPerms, err := r.permissionDAO.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	perms := make([]*model.Permission, len(daoPerms))
	for i, daoPerm := range daoPerms {
		perms[i] = r.daoToModel(daoPerm)
	}

	// Get total count
	count, err := r.permissionDAO.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return perms, count, nil
}

func (r *permissionRepositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]*model.Permission, error) {
	perms := make([]*model.Permission, 0, len(ids))
	for _, id := range ids {
		perm, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if perm != nil {
			perms = append(perms, perm)
		}
	}
	return perms, nil
}

func (r *permissionRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.permissionDAO.Delete(ctx, id)
}

func (r *permissionRepositoryImpl) ExistsByResourceAndAction(ctx context.Context, resource, action string) (bool, error) {
	perm, err := r.permissionDAO.FindByResourceAndAction(ctx, resource, action)
	if err != nil {
		return false, err
	}
	return perm != nil, nil
}

// Converters

func (r *permissionRepositoryImpl) modelToDAO(perm *model.Permission) *domain.Permission {
	return &domain.Permission{
		ID:          perm.ID(),
		Name:        perm.Name(),
		Resource:    perm.Resource(),
		Action:      perm.Action(),
		Description: perm.Description(),
		CreatedAt:   perm.CreatedAt(),
		UpdatedAt:   perm.UpdatedAt(),
	}
}

func (r *permissionRepositoryImpl) daoToModel(daoPerm *domain.Permission) *model.Permission {
	if daoPerm == nil {
		return nil
	}

	return model.ReconstructPermission(
		daoPerm.ID,
		daoPerm.Name,
		daoPerm.Resource,
		daoPerm.Action,
		daoPerm.Description,
		daoPerm.CreatedAt,
		daoPerm.UpdatedAt,
	)
}
