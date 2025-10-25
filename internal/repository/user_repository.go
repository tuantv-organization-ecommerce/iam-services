package repository

import (
	"context"
	"fmt"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/domain"
)

// UserRepository provides higher-level operations on User entity
type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
	UserExists(ctx context.Context, username, email string) (bool, error)
}

type userRepository struct {
	userDAO dao.UserDAO
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(userDAO dao.UserDAO) UserRepository {
	return &userRepository{
		userDAO: userDAO,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	// Check if user already exists
	exists, err := r.UserExists(ctx, user.Username, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("user with username '%s' or email '%s' already exists", user.Username, user.Email)
	}

	return r.userDAO.Create(ctx, user)
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := r.userDAO.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := r.userDAO.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.userDAO.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return r.userDAO.Update(ctx, user)
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	return r.userDAO.Delete(ctx, id)
}

func (r *userRepository) UserExists(ctx context.Context, username, email string) (bool, error) {
	userByUsername, err := r.userDAO.FindByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	if userByUsername != nil {
		return true, nil
	}

	userByEmail, err := r.userDAO.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	if userByEmail != nil {
		return true, nil
	}

	return false, nil
}
