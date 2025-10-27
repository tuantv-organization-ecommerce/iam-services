package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tvttt/iam-services/internal/domain"
	"github.com/tvttt/iam-services/pkg/jwt"
	"github.com/tvttt/iam-services/pkg/password"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) UserExists(ctx context.Context, username, email string) (bool, error) {
	args := m.Called(ctx, username, email)
	return args.Bool(0), args.Error(1)
}

// Mock AuthorizationRepository
type MockAuthorizationRepository struct {
	mock.Mock
}

func (m *MockAuthorizationRepository) GetUserRoles(ctx context.Context, userID string) ([]*domain.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Role), args.Error(1)
}

func (m *MockAuthorizationRepository) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockAuthorizationRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockAuthorizationRepository) GetUserPermissions(ctx context.Context, userID string) ([]*domain.Permission, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Permission), args.Error(1)
}

func (m *MockAuthorizationRepository) UserHasPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	args := m.Called(ctx, userID, resource, action)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthorizationRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

func (m *MockAuthorizationRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

func (m *MockAuthorizationRepository) GetRolePermissions(ctx context.Context, roleID string) ([]*domain.Permission, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Permission), args.Error(1)
}

func (m *MockAuthorizationRepository) UpdateRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	args := m.Called(ctx, roleID, permissionIDs)
	return args.Error(0)
}

func TestRegister_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockAuthzRepo := new(MockAuthorizationRepository)
	jwtManager := jwt.NewJWTManager("test-secret-min-32-chars-long", time.Hour, time.Hour*24)
	passwordManager := password.NewPasswordManager()

	service := NewAuthService(mockUserRepo, mockAuthzRepo, jwtManager, passwordManager)

	// Mock expectations
	mockUserRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	// Execute
	ctx := context.Background()
	user, err := service.Register(ctx, "testuser", "test@example.com", "password123", "Test User")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.FullName)
	assert.True(t, user.IsActive)
	assert.NotEmpty(t, user.ID)
	assert.NotEqual(t, "password123", user.PasswordHash) // Password should be hashed

	mockUserRepo.AssertExpectations(t)
}

func TestRegister_EmptyUsername(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockAuthzRepo := new(MockAuthorizationRepository)
	jwtManager := jwt.NewJWTManager("test-secret-min-32-chars-long", time.Hour, time.Hour*24)
	passwordManager := password.NewPasswordManager()

	service := NewAuthService(mockUserRepo, mockAuthzRepo, jwtManager, passwordManager)

	// Execute
	ctx := context.Background()
	user, err := service.Register(ctx, "", "test@example.com", "password123", "Test User")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "username")
}

func TestLogin_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockAuthzRepo := new(MockAuthorizationRepository)
	jwtManager := jwt.NewJWTManager("test-secret-min-32-chars-long", time.Hour, time.Hour*24)
	passwordManager := password.NewPasswordManager()

	service := NewAuthService(mockUserRepo, mockAuthzRepo, jwtManager, passwordManager)

	// Create hashed password
	hashedPassword, _ := passwordManager.HashPassword("password123")

	// Mock user
	mockUser := &domain.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		FullName:     "Test User",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRoles := []*domain.Role{
		{ID: "role-1", Name: "user"},
	}

	// Mock expectations
	mockUserRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(mockUser, nil)
	mockAuthzRepo.On("GetUserRoles", mock.Anything, "user-123").Return(mockRoles, nil)

	// Execute
	ctx := context.Background()
	user, tokenPair, err := service.Login(ctx, "testuser", "password123")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, tokenPair)
	assert.Equal(t, "testuser", user.Username)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.Equal(t, "Bearer", tokenPair.TokenType)

	mockUserRepo.AssertExpectations(t)
	mockAuthzRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockAuthzRepo := new(MockAuthorizationRepository)
	jwtManager := jwt.NewJWTManager("test-secret-min-32-chars-long", time.Hour, time.Hour*24)
	passwordManager := password.NewPasswordManager()

	service := NewAuthService(mockUserRepo, mockAuthzRepo, jwtManager, passwordManager)

	// Create hashed password
	hashedPassword, _ := passwordManager.HashPassword("password123")

	// Mock user
	mockUser := &domain.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		FullName:     "Test User",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Mock expectations
	mockUserRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(mockUser, nil)

	// Execute with wrong password
	ctx := context.Background()
	user, tokenPair, err := service.Login(ctx, "testuser", "wrongpassword")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, tokenPair)
	assert.Contains(t, err.Error(), "credentials")

	mockUserRepo.AssertExpectations(t)
}

func TestLogin_InactiveUser(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockAuthzRepo := new(MockAuthorizationRepository)
	jwtManager := jwt.NewJWTManager("test-secret-min-32-chars-long", time.Hour, time.Hour*24)
	passwordManager := password.NewPasswordManager()

	service := NewAuthService(mockUserRepo, mockAuthzRepo, jwtManager, passwordManager)

	// Create hashed password
	hashedPassword, _ := passwordManager.HashPassword("password123")

	// Mock inactive user
	mockUser := &domain.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		FullName:     "Test User",
		IsActive:     false, // Inactive
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Mock expectations
	mockUserRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(mockUser, nil)

	// Execute
	ctx := context.Background()
	user, tokenPair, err := service.Login(ctx, "testuser", "password123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, tokenPair)
	assert.Contains(t, err.Error(), "inactive")

	mockUserRepo.AssertExpectations(t)
}

func TestVerifyToken_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockAuthzRepo := new(MockAuthorizationRepository)
	jwtManager := jwt.NewJWTManager("test-secret-min-32-chars-long", time.Hour, time.Hour*24)
	passwordManager := password.NewPasswordManager()

	service := NewAuthService(mockUserRepo, mockAuthzRepo, jwtManager, passwordManager)

	// Generate a valid token
	token, err := jwtManager.GenerateAccessToken("user-123", "testuser", []string{"user", "admin"})
	require.NoError(t, err)

	// Execute
	ctx := context.Background()
	userID, roles, err := service.VerifyToken(ctx, token)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "user-123", userID)
	assert.Equal(t, []string{"user", "admin"}, roles)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockAuthzRepo := new(MockAuthorizationRepository)
	jwtManager := jwt.NewJWTManager("test-secret-min-32-chars-long", time.Hour, time.Hour*24)
	passwordManager := password.NewPasswordManager()

	service := NewAuthService(mockUserRepo, mockAuthzRepo, jwtManager, passwordManager)

	// Execute with invalid token
	ctx := context.Background()
	userID, roles, err := service.VerifyToken(ctx, "invalid.token.here")

	// Assert
	assert.Error(t, err)
	assert.Empty(t, userID)
	assert.Nil(t, roles)
}

func TestLogout_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockAuthzRepo := new(MockAuthorizationRepository)
	jwtManager := jwt.NewJWTManager("test-secret-min-32-chars-long", time.Hour, time.Hour*24)
	passwordManager := password.NewPasswordManager()

	service := NewAuthService(mockUserRepo, mockAuthzRepo, jwtManager, passwordManager)

	// Execute
	ctx := context.Background()
	err := service.Logout(ctx, "user-123")

	// Assert
	assert.NoError(t, err) // Currently logout just returns nil
}
