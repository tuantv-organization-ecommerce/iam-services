package persistence

import (
	"context"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/internal/domain/model"
	domainRepo "github.com/tvttt/iam-services/internal/domain/repository"
)

// roleRepositoryImpl implements domain.repository.RoleRepository using DAO
type roleRepositoryImpl struct {
	roleDAO           dao.RoleDAO
	rolePermissionDAO dao.RolePermissionDAO
}

// NewRoleRepository creates a new role repository implementation
func NewRoleRepository(
	roleDAO dao.RoleDAO,
	rolePermissionDAO dao.RolePermissionDAO,
) domainRepo.RoleRepository {
	return &roleRepositoryImpl{
		roleDAO:           roleDAO,
		rolePermissionDAO: rolePermissionDAO,
	}
}

func (r *roleRepositoryImpl) Save(ctx context.Context, role *model.Role) error {
	daoRole := r.modelToDAO(role)
	return r.roleDAO.Create(ctx, daoRole)
}

func (r *roleRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Role, error) {
	daoRole, err := r.roleDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if daoRole == nil {
		return nil, nil
	}

	// Load permissions
	role := r.daoToModel(daoRole)
	permissions, err := r.rolePermissionDAO.GetRolePermissions(ctx, id)
	if err != nil {
		return nil, err
	}

	modelPerms := make([]model.Permission, 0, len(permissions))
	for _, p := range permissions {
		modelPerms = append(modelPerms, *r.permToModel(p))
	}
	role.SetPermissions(modelPerms)

	return role, nil
}

func (r *roleRepositoryImpl) FindByName(ctx context.Context, name string) (*model.Role, error) {
	daoRole, err := r.roleDAO.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if daoRole == nil {
		return nil, nil
	}

	return r.daoToModel(daoRole), nil
}

func (r *roleRepositoryImpl) FindByDomain(ctx context.Context, domainStr string) ([]*model.Role, error) {
	// Note: Current DAO doesn't have FindByDomain, would need to add it
	// For now, return empty list
	return []*model.Role{}, nil
}

func (r *roleRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]*model.Role, int, error) {
	offset := (page - 1) * pageSize
	daoRoles, err := r.roleDAO.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	roles := make([]*model.Role, len(daoRoles))
	for i, daoRole := range daoRoles {
		roles[i] = r.daoToModel(daoRole)
	}

	// Get total count
	count, err := r.roleDAO.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return roles, count, nil
}

func (r *roleRepositoryImpl) Update(ctx context.Context, role *model.Role) error {
	daoRole := r.modelToDAO(role)
	return r.roleDAO.Update(ctx, daoRole)
}

func (r *roleRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.roleDAO.Delete(ctx, id)
}

func (r *roleRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	role, err := r.roleDAO.FindByName(ctx, name)
	if err != nil {
		return false, err
	}
	return role != nil, nil
}

// Converters

func (r *roleRepositoryImpl) modelToDAO(role *model.Role) *domain.Role {
	return &domain.Role{
		ID:          role.ID(),
		Name:        role.Name(),
		Description: role.Description(),
		CreatedAt:   role.CreatedAt(),
		UpdatedAt:   role.UpdatedAt(),
	}
}

func (r *roleRepositoryImpl) daoToModel(daoRole *domain.Role) *model.Role {
	if daoRole == nil {
		return nil
	}

	// Note: domain.Role doesn't have Domain field, defaulting to "user"
	return model.ReconstructRole(
		daoRole.ID,
		daoRole.Name,
		daoRole.Description,
		"user", // default domain
		daoRole.CreatedAt,
		daoRole.UpdatedAt,
	)
}

func (r *roleRepositoryImpl) permToModel(daoPerm *domain.Permission) *model.Permission {
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
