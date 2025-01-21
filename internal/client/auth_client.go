package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"bytes"
	"strings"

	"github.com/yourusername/projectname/internal/domain"
)

const (
	userValidationPath = "/users/uservalidation"
	loginPath         = "/authenticate/login"
)

// AuthClient defines the interface for authentication operations
type AuthClient interface {
	ValidateToken(ctx context.Context, token string) (*domain.AuthenticatedUser, error)
	Login(ctx context.Context, request domain.LoginRequest) (*domain.AuthResponse, error)
}

// authClientImpl implements the AuthClient interface
type authClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

// NewAuthClient creates a new instance of AuthClient
func NewAuthClient(baseURL string) AuthClient {
	return &authClientImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ValidateToken validates the provided token with the auth service
func (c *authClientImpl) ValidateToken(ctx context.Context, token string) (*domain.AuthenticatedUser, error) {
	// Create request to validate token
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s", c.baseURL, userValidationPath), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add token to request header with Bearer prefix if not present
	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}
	req.Header.Set("Authorization", token)

	// Make request to auth service
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	// Parse response
	var user domain.AuthenticatedUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// Login authenticates user credentials
func (c *authClientImpl) Login(ctx context.Context, request domain.LoginRequest) (*domain.AuthResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.baseURL, loginPath), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var authResponse domain.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &authResponse, nil
}

// handleErrorResponse handles error responses from the auth service
func (c *authClientImpl) handleErrorResponse(resp *http.Response) error {
	var errorResp struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	return fmt.Errorf("%s", errorResp.Error)
} 