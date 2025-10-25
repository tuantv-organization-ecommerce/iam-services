package repository

import (
	"context"

	"github.com/tvttt/iam-services/internal/domain/model"
)

// UserRepository defines the contract for user data access
// This is a PORT in hexagonal architecture
type UserRepository interface {
	// Save persists a new user or updates existing
	Save(ctx context.Context, user *model.User) error

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id string) (*model.User, error)

	// FindByUsername retrieves a user by username
	FindByUsername(ctx context.Context, username string) (*model.User, error)

	// FindByEmail retrieves a user by email
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *model.User) error

	// Delete removes a user
	Delete(ctx context.Context, id string) error

	// ExistsByUsername checks if username already exists
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ExistsByEmail checks if email already exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
