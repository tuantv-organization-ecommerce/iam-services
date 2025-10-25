package repository

import (
	"context"
	"fmt"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/domain"
)

// CMSRepository provides operations for CMS role management
type CMSRepository interface {
	CreateCMSRole(ctx context.Context, role *domain.CMSRole) error
	GetCMSRoleByID(ctx context.Context, id string) (*domain.CMSRole, error)
	GetCMSRoleByName(ctx context.Context, name string) (*domain.CMSRole, error)
	ListCMSRoles(ctx context.Context, page, pageSize int) ([]*domain.CMSRole, int, error)
	UpdateCMSRole(ctx context.Context, role *domain.CMSRole) error
	DeleteCMSRole(ctx context.Context, id string) error

	AssignCMSRoleToUser(ctx context.Context, userID, cmsRoleID string) error
	RemoveCMSRoleFromUser(ctx context.Context, userID, cmsRoleID string) error
	GetUserCMSRoles(ctx context.Context, userID string) ([]*domain.CMSRole, error)
	GetUserCMSTabs(ctx context.Context, userID string) ([]domain.CMSTab, error)
}

type cmsRepository struct {
	cmsRoleDAO     dao.CMSRoleDAO
	userCMSRoleDAO dao.UserCMSRoleDAO
}

// NewCMSRepository creates a new instance of CMSRepository
func NewCMSRepository(cmsRoleDAO dao.CMSRoleDAO, userCMSRoleDAO dao.UserCMSRoleDAO) CMSRepository {
	return &cmsRepository{
		cmsRoleDAO:     cmsRoleDAO,
		userCMSRoleDAO: userCMSRoleDAO,
	}
}

func (r *cmsRepository) CreateCMSRole(ctx context.Context, role *domain.CMSRole) error {
	// Check if role already exists
	existingRole, err := r.cmsRoleDAO.FindByName(ctx, role.Name)
	if err != nil {
		return fmt.Errorf("failed to check CMS role existence: %w", err)
	}
	if existingRole != nil {
		return fmt.Errorf("CMS role with name '%s' already exists", role.Name)
	}

	return r.cmsRoleDAO.Create(ctx, role)
}

func (r *cmsRepository) GetCMSRoleByID(ctx context.Context, id string) (*domain.CMSRole, error) {
	role, err := r.cmsRoleDAO.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get CMS role by ID: %w", err)
	}
	if role == nil {
		return nil, fmt.Errorf("CMS role not found")
	}
	return role, nil
}

func (r *cmsRepository) GetCMSRoleByName(ctx context.Context, name string) (*domain.CMSRole, error) {
	role, err := r.cmsRoleDAO.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get CMS role by name: %w", err)
	}
	if role == nil {
		return nil, fmt.Errorf("CMS role not found")
	}
	return role, nil
}

func (r *cmsRepository) ListCMSRoles(ctx context.Context, page, pageSize int) ([]*domain.CMSRole, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	roles, err := r.cmsRoleDAO.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list CMS roles: %w", err)
	}

	total, err := r.cmsRoleDAO.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count CMS roles: %w", err)
	}

	return roles, total, nil
}

func (r *cmsRepository) UpdateCMSRole(ctx context.Context, role *domain.CMSRole) error {
	return r.cmsRoleDAO.Update(ctx, role)
}

func (r *cmsRepository) DeleteCMSRole(ctx context.Context, id string) error {
	return r.cmsRoleDAO.Delete(ctx, id)
}

func (r *cmsRepository) AssignCMSRoleToUser(ctx context.Context, userID, cmsRoleID string) error {
	userCMSRole := &domain.UserCMSRole{
		UserID:    userID,
		CMSRoleID: cmsRoleID,
	}
	return r.userCMSRoleDAO.AssignCMSRole(ctx, userCMSRole)
}

func (r *cmsRepository) RemoveCMSRoleFromUser(ctx context.Context, userID, cmsRoleID string) error {
	return r.userCMSRoleDAO.RemoveCMSRole(ctx, userID, cmsRoleID)
}

func (r *cmsRepository) GetUserCMSRoles(ctx context.Context, userID string) ([]*domain.CMSRole, error) {
	return r.userCMSRoleDAO.GetUserCMSRoles(ctx, userID)
}

func (r *cmsRepository) GetUserCMSTabs(ctx context.Context, userID string) ([]domain.CMSTab, error) {
	roles, err := r.userCMSRoleDAO.GetUserCMSRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user CMS roles: %w", err)
	}

	// Collect all unique tabs from all CMS roles
	tabMap := make(map[domain.CMSTab]bool)
	for _, role := range roles {
		for _, tab := range role.Tabs {
			tabMap[tab] = true
		}
	}

	// Convert map to slice
	tabs := make([]domain.CMSTab, 0, len(tabMap))
	for tab := range tabMap {
		tabs = append(tabs, tab)
	}

	return tabs, nil
}
