package persistence

import (
	"context"

	"github.com/tvttt/iam-services/internal/dao"
	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/internal/domain/model"
	domainRepo "github.com/tvttt/iam-services/internal/domain/repository"
)

// userRepositoryImpl implements domain.repository.UserRepository using DAO
type userRepositoryImpl struct {
	userDAO dao.UserDAO
}

// NewUserRepository creates a new user repository implementation
func NewUserRepository(userDAO dao.UserDAO) domainRepo.UserRepository {
	return &userRepositoryImpl{
		userDAO: userDAO,
	}
}

func (r *userRepositoryImpl) Save(ctx context.Context, user *model.User) error {
	// Convert domain model to DAO entity
	daoUser := r.modelToDAO(user)

	if err := r.userDAO.Create(ctx, daoUser); err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id string) (*model.User, error) {
	daoUser, err := r.userDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if daoUser == nil {
		return nil, nil
	}

	return r.daoToModel(daoUser), nil
}

func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	daoUser, err := r.userDAO.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if daoUser == nil {
		return nil, nil
	}

	return r.daoToModel(daoUser), nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	daoUser, err := r.userDAO.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if daoUser == nil {
		return nil, nil
	}

	return r.daoToModel(daoUser), nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	daoUser := r.modelToDAO(user)

	if err := r.userDAO.Update(ctx, daoUser); err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.userDAO.Delete(ctx, id)
}

func (r *userRepositoryImpl) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	user, err := r.userDAO.FindByUsername(ctx, username)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func (r *userRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	user, err := r.userDAO.FindByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

// Converters between domain.User (DAO) and model.User (domain model)

func (r *userRepositoryImpl) modelToDAO(user *model.User) *domain.User {
	return &domain.User{
		ID:           user.ID(),
		Username:     user.Username(),
		Email:        user.Email(),
		PasswordHash: user.PasswordHash(),
		FullName:     user.FullName(),
		IsActive:     user.IsActive(),
		CreatedAt:    user.CreatedAt(),
		UpdatedAt:    user.UpdatedAt(),
	}
}

func (r *userRepositoryImpl) daoToModel(daoUser *domain.User) *model.User {
	if daoUser == nil {
		return nil
	}

	return model.ReconstructUser(
		daoUser.ID,
		daoUser.Username,
		daoUser.Email,
		daoUser.PasswordHash,
		daoUser.FullName,
		daoUser.IsActive,
		daoUser.CreatedAt,
		daoUser.UpdatedAt,
	)
}
