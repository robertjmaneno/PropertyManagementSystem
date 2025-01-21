package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/internal/dto"
	"github.com/yourusername/projectname/internal/handler"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/pkg/errors"
	"github.com/yourusername/projectname/test/testutil"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetAll(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserService) GetPaged(ctx context.Context, opts repository.ListOptions) ([]domain.User, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) Count(ctx context.Context, opts repository.ListOptions) (int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserService) Search(ctx context.Context, query string, opts repository.ListOptions) ([]*domain.User, error) {
	args := m.Called(ctx, query, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func setupTest(t *testing.T) (*gin.Engine, *MockUserService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(MockUserService)
	userHandler := handler.NewUserHandler(mockService)

	api := router.Group("/api/v1")
	userHandler.RegisterRoutes(api)

	return router, mockService
}

func TestUserHandler(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		router, mockService := setupTest(t)
		user := testutil.NewTestUser()

		// Setup expectations
		mockService.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

		// Create request
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/users", createJSONBody(t, dto.CreateUserRequest{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}))

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("GetByID", func(t *testing.T) {
		router, mockService := setupTest(t)
		user := testutil.NewTestUser()

		// Setup expectations
		mockService.On("GetByID", mock.Anything, user.ID).Return(user, nil)

		// Create request
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/users/"+user.ID, nil)

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.UserResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, user.ID, response.ID)

		// Test not found
		mockService.On("GetByID", mock.Anything, "non-existent").Return(nil, errors.ErrNotFound)
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "/api/v1/users/non-existent", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("List", func(t *testing.T) {
		router, mockService := setupTest(t)
		users := []domain.User{*testutil.NewTestUser(), *testutil.NewTestUser()}

		// Setup expectations
		mockService.On("GetPaged", mock.Anything, mock.AnythingOfType("repository.ListOptions")).Return(users, nil)
		mockService.On("Count", mock.Anything, mock.AnythingOfType("repository.ListOptions")).Return(int64(len(users)), nil)

		// Create request
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/users?page=1&page_size=10", nil)

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.ListResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Len(t, response.Data.([]interface{}), len(users))

		mockService.AssertExpectations(t)
	})

	t.Run("Search", func(t *testing.T) {
		router, mockService := setupTest(t)
		users := []*domain.User{testutil.NewTestUser(), testutil.NewTestUser()}

		// Setup expectations
		mockService.On("Search", mock.Anything, "test", mock.AnythingOfType("repository.ListOptions")).Return(users, nil)

		// Create request
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/users/search?search=test", nil)

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		var response []*dto.UserResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Len(t, response, len(users))

		mockService.AssertExpectations(t)
	})

	t.Run("Update", func(t *testing.T) {
		router, mockService := setupTest(t)
		user := testutil.NewTestUser()

		// Setup expectations
		mockService.On("GetByID", mock.Anything, user.ID).Return(user, nil)
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

		// Create request
		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/api/v1/users/"+user.ID, createJSONBody(t, dto.UpdateUserRequest{
			Email:     "updated@example.com",
			FirstName: "Updated",
			LastName:  "User",
		}))

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		router, mockService := setupTest(t)
		user := testutil.NewTestUser()

		// Setup expectations
		mockService.On("Delete", mock.Anything, user.ID).Return(nil)

		// Create request
		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/api/v1/users/"+user.ID, nil)

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})
}

func createJSONBody(t *testing.T, v interface{}) *bytes.Buffer {
	var buf bytes.Buffer
	require.NoError(t, json.NewEncoder(&buf).Encode(v))
	return &buf
}
