package persistence

import (
	"context"
	"time"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/internal/domain/model"
	domainRepo "github.com/tvttt/iam-services/internal/domain/repository"
)

// cmsRepositoryImpl implements domain.repository.CMSRepository using DAO
type cmsRepositoryImpl struct {
	cmsRoleDAO     dao.CMSRoleDAO
	userCMSRoleDAO dao.UserCMSRoleDAO
}

// NewCMSRepository creates a new CMS repository implementation
func NewCMSRepository(
	cmsRoleDAO dao.CMSRoleDAO,
	userCMSRoleDAO dao.UserCMSRoleDAO,
) domainRepo.CMSRepository {
	return &cmsRepositoryImpl{
		cmsRoleDAO:     cmsRoleDAO,
		userCMSRoleDAO: userCMSRoleDAO,
	}
}

func (r *cmsRepositoryImpl) Save(ctx context.Context, role *model.CMSRole) error {
	daoRole := r.modelToDAO(role)
	return r.cmsRoleDAO.Create(ctx, daoRole)
}

func (r *cmsRepositoryImpl) FindByID(ctx context.Context, id string) (*model.CMSRole, error) {
	daoRole, err := r.cmsRoleDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if daoRole == nil {
		return nil, nil
	}

	return r.daoToModel(daoRole), nil
}

func (r *cmsRepositoryImpl) FindByName(ctx context.Context, name string) (*model.CMSRole, error) {
	daoRole, err := r.cmsRoleDAO.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if daoRole == nil {
		return nil, nil
	}

	return r.daoToModel(daoRole), nil
}

func (r *cmsRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]*model.CMSRole, int, error) {
	offset := (page - 1) * pageSize
	daoRoles, err := r.cmsRoleDAO.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	roles := make([]*model.CMSRole, len(daoRoles))
	for i, daoRole := range daoRoles {
		roles[i] = r.daoToModel(daoRole)
	}

	// Get total count
	count, err := r.cmsRoleDAO.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return roles, count, nil
}

func (r *cmsRepositoryImpl) Update(ctx context.Context, role *model.CMSRole) error {
	daoRole := r.modelToDAO(role)
	return r.cmsRoleDAO.Update(ctx, daoRole)
}

func (r *cmsRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.cmsRoleDAO.Delete(ctx, id)
}

func (r *cmsRepositoryImpl) AssignToUser(ctx context.Context, userID, roleID string) error {
	userCMSRole := &domain.UserCMSRole{
		UserID:    userID,
		CMSRoleID: roleID,
		CreatedAt: time.Now(),
	}
	return r.userCMSRoleDAO.AssignCMSRole(ctx, userCMSRole)
}

func (r *cmsRepositoryImpl) RemoveFromUser(ctx context.Context, userID, roleID string) error {
	return r.userCMSRoleDAO.RemoveCMSRole(ctx, userID, roleID)
}

func (r *cmsRepositoryImpl) GetUserCMSRoles(ctx context.Context, userID string) ([]*model.CMSRole, error) {
	daoRoles, err := r.userCMSRoleDAO.GetUserCMSRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	roles := make([]*model.CMSRole, len(daoRoles))
	for i, daoRole := range daoRoles {
		roles[i] = r.daoToModel(daoRole)
	}

	return roles, nil
}

func (r *cmsRepositoryImpl) GetUserAccessibleTabs(ctx context.Context, userID string) ([]model.CMSTab, error) {
	daoRoles, err := r.userCMSRoleDAO.GetUserCMSRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Deduplicate tabs
	tabMap := make(map[model.CMSTab]bool)
	for _, role := range daoRoles {
		for _, tab := range role.Tabs {
			tabMap[model.CMSTab(tab)] = true
		}
	}

	// Convert to slice
	tabs := make([]model.CMSTab, 0, len(tabMap))
	for tab := range tabMap {
		tabs = append(tabs, tab)
	}

	return tabs, nil
}

// Converters

func (r *cmsRepositoryImpl) modelToDAO(role *model.CMSRole) *domain.CMSRole {
	tabs := make([]domain.CMSTab, len(role.Tabs()))
	for i, tab := range role.Tabs() {
		tabs[i] = domain.CMSTab(tab)
	}

	return &domain.CMSRole{
		ID:          role.ID(),
		Name:        role.Name(),
		Description: role.Description(),
		Tabs:        tabs,
		CreatedAt:   role.CreatedAt(),
		UpdatedAt:   role.UpdatedAt(),
	}
}

func (r *cmsRepositoryImpl) daoToModel(daoRole *domain.CMSRole) *model.CMSRole {
	if daoRole == nil {
		return nil
	}

	tabs := make([]model.CMSTab, len(daoRole.Tabs))
	for i, tab := range daoRole.Tabs {
		tabs[i] = model.CMSTab(tab)
	}

	return model.ReconstructCMSRole(
		daoRole.ID,
		daoRole.Name,
		daoRole.Description,
		tabs,
		daoRole.CreatedAt,
		daoRole.UpdatedAt,
	)
}
