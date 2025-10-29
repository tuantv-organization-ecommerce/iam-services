package dao

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/tvttt/iam-services/internal/domain"
	"log"
)

// UserCMSRoleDAO defines the data access operations for UserCMSRole relationship
type UserCMSRoleDAO interface {
	AssignCMSRole(ctx context.Context, userCMSRole *domain.UserCMSRole) error
	RemoveCMSRole(ctx context.Context, userID, cmsRoleID string) error
	GetUserCMSRoles(ctx context.Context, userID string) ([]*domain.CMSRole, error)
	GetCMSRoleUsers(ctx context.Context, cmsRoleID string) ([]string, error)
	HasCMSRole(ctx context.Context, userID, cmsRoleID string) (bool, error)
	RemoveAllCMSRolesFromUser(ctx context.Context, userID string) error
}

type userCMSRoleDAO struct {
	db *sql.DB
}

// NewUserCMSRoleDAO creates a new instance of UserCMSRoleDAO
func NewUserCMSRoleDAO(db *sql.DB) UserCMSRoleDAO {
	return &userCMSRoleDAO{db: db}
}

func (d *userCMSRoleDAO) AssignCMSRole(ctx context.Context, userCMSRole *domain.UserCMSRole) error {
	query := `
		INSERT INTO user_cms_roles (user_id, cms_role_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, cms_role_id) DO NOTHING
	`
	_, err := d.db.ExecContext(ctx, query,
		userCMSRole.UserID,
		userCMSRole.CMSRoleID,
		userCMSRole.CreatedAt,
	)
	return err
}

func (d *userCMSRoleDAO) RemoveCMSRole(ctx context.Context, userID, cmsRoleID string) error {
	query := `DELETE FROM user_cms_roles WHERE user_id = $1 AND cms_role_id = $2`
	_, err := d.db.ExecContext(ctx, query, userID, cmsRoleID)
	return err
}

func (d *userCMSRoleDAO) GetUserCMSRoles(ctx context.Context, userID string) ([]*domain.CMSRole, error) {
	query := `
		SELECT r.id, r.name, r.description, r.tabs, r.created_at, r.updated_at
		FROM cms_roles r
		INNER JOIN user_cms_roles ucr ON r.id = ucr.cms_role_id
		WHERE ucr.user_id = $1
		ORDER BY r.name
	`
	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows:", err)
		}
	}()

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

func (d *userCMSRoleDAO) GetCMSRoleUsers(ctx context.Context, cmsRoleID string) ([]string, error) {
	query := `
		SELECT user_id
		FROM user_cms_roles
		WHERE cms_role_id = $1
		ORDER BY created_at
	`
	rows, err := d.db.QueryContext(ctx, query, cmsRoleID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println("Error closing rows:", err)
		}
	}()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, rows.Err()
}

func (d *userCMSRoleDAO) HasCMSRole(ctx context.Context, userID, cmsRoleID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_cms_roles WHERE user_id = $1 AND cms_role_id = $2)`
	var exists bool
	err := d.db.QueryRowContext(ctx, query, userID, cmsRoleID).Scan(&exists)
	return exists, err
}

func (d *userCMSRoleDAO) RemoveAllCMSRolesFromUser(ctx context.Context, userID string) error {
	query := `DELETE FROM user_cms_roles WHERE user_id = $1`
	_, err := d.db.ExecContext(ctx, query, userID)
	return err
}
