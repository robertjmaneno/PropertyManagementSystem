package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/pkg/errors"
)

type TestClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewTestClient(baseURL string) *TestClient {
	return &TestClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *TestClient) CreateUser(user *domain.User) (*domain.User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	resp, err := c.httpClient.Post(fmt.Sprintf("%s/api/users", c.baseURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Success bool        `json:"success"`
		Data    domain.User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}

func (c *TestClient) GetUser(id string) (*domain.User, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/users/%s", c.baseURL, id))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Success bool        `json:"success"`
		Data    domain.User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}
