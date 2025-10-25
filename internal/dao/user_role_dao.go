package dao

import (
	"context"
	"database/sql"

	"github.com/tvttt/iam-services/internal/domain"
)

// UserRoleDAO defines the data access operations for UserRole relationship
type UserRoleDAO interface {
	AssignRole(ctx context.Context, userRole *domain.UserRole) error
	RemoveRole(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*domain.Role, error)
	GetRoleUsers(ctx context.Context, roleID string) ([]*domain.User, error)
	HasRole(ctx context.Context, userID, roleID string) (bool, error)
}

type userRoleDAO struct {
	db *sql.DB
}

// NewUserRoleDAO creates a new instance of UserRoleDAO
func NewUserRoleDAO(db *sql.DB) UserRoleDAO {
	return &userRoleDAO{db: db}
}

func (d *userRoleDAO) AssignRole(ctx context.Context, userRole *domain.UserRole) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err := d.db.ExecContext(ctx, query,
		userRole.UserID,
		userRole.RoleID,
		userRole.CreatedAt,
	)
	return err
}

func (d *userRoleDAO) RemoveRole(ctx context.Context, userID, roleID string) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	_, err := d.db.ExecContext(ctx, query, userID, roleID)
	return err
}

func (d *userRoleDAO) GetUserRoles(ctx context.Context, userID string) ([]*domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.name
	`
	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		role := &domain.Role{}
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (d *userRoleDAO) GetRoleUsers(ctx context.Context, roleID string) ([]*domain.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, 
		       u.is_active, u.created_at, u.updated_at
		FROM users u
		INNER JOIN user_roles ur ON u.id = ur.user_id
		WHERE ur.role_id = $1
		ORDER BY u.username
	`
	rows, err := d.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.FullName,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (d *userRoleDAO) HasRole(ctx context.Context, userID, roleID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_roles WHERE user_id = $1 AND role_id = $2)`
	var exists bool
	err := d.db.QueryRowContext(ctx, query, userID, roleID).Scan(&exists)
	return exists, err
}
