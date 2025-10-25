package security

import (
	"github.com/tvttt/iam-services/internal/domain/service"
	"github.com/tvttt/iam-services/pkg/password"
)

// passwordServiceImpl implements domain.service.PasswordService using password manager
type passwordServiceImpl struct {
	passwordManager *password.PasswordManager
}

// NewPasswordService creates a new password service implementation
func NewPasswordService(passwordManager *password.PasswordManager) service.PasswordService {
	return &passwordServiceImpl{
		passwordManager: passwordManager,
	}
}

func (s *passwordServiceImpl) Hash(password string) (string, error) {
	return s.passwordManager.HashPassword(password)
}

func (s *passwordServiceImpl) Verify(password, hash string) bool {
	return s.passwordManager.CheckPassword(password, hash)
}
