package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yourusername/projectname/internal/client"
	"github.com/yourusername/projectname/internal/config"
	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/internal/security"
	"go.uber.org/zap"
)

// AccessControlService manages service authentication and token management
type AccessControlService struct {
	authClient client.AuthClient
	config     *config.Config
	logger     *zap.Logger

	// Cache for service token
	serviceToken string
	tokenExpiry  time.Time
	tokenMutex   sync.RWMutex
}

// NewAccessControlService creates a new AccessControlService
func NewAccessControlService(authClient client.AuthClient, config *config.Config, logger *zap.Logger) *AccessControlService {
	return &AccessControlService{
		authClient: authClient,
		config:     config,
		logger:     logger,
	}
}

// GetServiceToken gets a valid service token, refreshing if necessary
func (s *AccessControlService) GetServiceToken(ctx context.Context) (string, error) {
	s.tokenMutex.RLock()
	if s.serviceToken != "" && time.Now().Before(s.tokenExpiry) {
		token := s.serviceToken
		s.tokenMutex.RUnlock()
		return token, nil
	}
	s.tokenMutex.RUnlock()

	// Need to refresh token
	s.tokenMutex.Lock()
	defer s.tokenMutex.Unlock()

	// Double check if another goroutine refreshed the token
	if s.serviceToken != "" && time.Now().Before(s.tokenExpiry) {
		return s.serviceToken, nil
	}

	// Login to get new token
	loginReq := domain.LoginRequest{
		Username: s.config.Auth.Username,
		Password: s.config.Auth.Password,
	}

	authResp, err := s.authClient.Login(ctx, loginReq)
	if err != nil {
		return "", fmt.Errorf("failed to get service token: %w", err)
	}

	s.serviceToken = authResp.Token
	s.tokenExpiry = time.Unix(int64(authResp.EXP), 0).Add(-5 * time.Minute) // Refresh 5 minutes before expiry

	return s.serviceToken, nil
}

// GetBearerToken gets a bearer token for service authentication
func (s *AccessControlService) GetBearerToken(ctx context.Context) (string, error) {
	token, err := s.GetServiceToken(ctx)
	if err != nil {
		return "", err
	}
	return security.CreateBearerToken(token), nil
}
