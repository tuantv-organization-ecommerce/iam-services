package dao

import (
	"context"
	"database/sql"
	"time"

	"github.com/tvttt/iam-services/internal/domain"
)

// UserDAO defines the data access operations for User entity
type UserDAO interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

type userDAO struct {
	db *sql.DB
}

// NewUserDAO creates a new instance of UserDAO
func NewUserDAO(db *sql.DB) UserDAO {
	return &userDAO{db: db}
}

func (d *userDAO) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := d.db.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (d *userDAO) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &domain.User{}
	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *userDAO) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, is_active, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	user := &domain.User{}
	err := d.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *userDAO) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &domain.User{}
	err := d.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *userDAO) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET username = $2, email = $3, password_hash = $4, full_name = $5, 
		    is_active = $6, updated_at = $7
		WHERE id = $1
	`
	user.UpdatedAt = time.Now()
	_, err := d.db.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.IsActive,
		user.UpdatedAt,
	)
	return err
}

func (d *userDAO) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := d.db.ExecContext(ctx, query, id)
	return err
}
