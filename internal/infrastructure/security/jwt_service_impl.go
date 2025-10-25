package security

import (
	"github.com/tvttt/iam-services/internal/domain/service"
	"github.com/tvttt/iam-services/pkg/jwt"
)

// jwtServiceImpl implements domain.service.TokenService using JWT manager
type jwtServiceImpl struct {
	jwtManager *jwt.JWTManager
}

// NewJWTService creates a new JWT service implementation
func NewJWTService(jwtManager *jwt.JWTManager) service.TokenService {
	return &jwtServiceImpl{
		jwtManager: jwtManager,
	}
}

func (s *jwtServiceImpl) GenerateAccessToken(userID, username string, roles []string) (string, error) {
	token, err := s.jwtManager.GenerateAccessToken(userID, username, roles)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *jwtServiceImpl) GenerateRefreshToken(userID string) (string, error) {
	token, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *jwtServiceImpl) GenerateTokenPair(userID, username string, roles []string) (*service.TokenPair, error) {
	accessToken, err := s.GenerateAccessToken(userID, username, roles)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &service.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtManager.AccessTokenDuration.Seconds()),
	}, nil
}

func (s *jwtServiceImpl) VerifyToken(token string) (*service.TokenClaims, error) {
	claims, err := s.jwtManager.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	return &service.TokenClaims{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Roles:     claims.Roles,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

func (s *jwtServiceImpl) RefreshAccessToken(refreshToken string) (string, error) {
	// Verify refresh token
	claims, err := s.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Generate new access token
	// Note: We don't have username and roles from refresh token
	// In production, you might want to fetch user data from database
	token, err := s.jwtManager.GenerateAccessToken(claims.UserID, "", []string{})
	if err != nil {
		return "", err
	}

	return token, nil
}
