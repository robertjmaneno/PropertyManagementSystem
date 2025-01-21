package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/internal/service"
	"github.com/yourusername/projectname/pkg/errors"
	"github.com/yourusername/projectname/test/testutil"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) GetPaged(ctx context.Context, opts repository.ListOptions) ([]domain.User, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Count(ctx context.Context, opts repository.ListOptions) (int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) Search(ctx context.Context, query string, opts repository.ListOptions) ([]*domain.User, error) {
	args := m.Called(ctx, query, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func TestUserService(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	t.Run("Create", func(t *testing.T) {
		user := testutil.NewTestUser()

		// Setup expectations
		mockRepo.On("Create", ctx, user).Return(nil)

		// Test
		err := userService.Create(ctx, user)
		require.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByID", func(t *testing.T) {
		user := testutil.NewTestUser()

		// Setup expectations
		mockRepo.On("GetByID", ctx, user.ID).Return(user, nil)

		// Test success
		found, err := userService.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)

		// Test not found
		mockRepo.On("GetByID", ctx, "non-existent").Return(nil, errors.ErrNotFound)
		_, err = userService.GetByID(ctx, "non-existent")
		assert.ErrorIs(t, err, errors.ErrNotFound)

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByEmail", func(t *testing.T) {
		user := testutil.NewTestUser()

		// Setup expectations
		mockRepo.On("GetByEmail", ctx, user.Email).Return(user, nil)

		// Test success
		found, err := userService.GetByEmail(ctx, user.Email)
		require.NoError(t, err)
		assert.Equal(t, user.Email, found.Email)

		// Test not found
		mockRepo.On("GetByEmail", ctx, "nonexistent@example.com").Return(nil, errors.ErrNotFound)
		_, err = userService.GetByEmail(ctx, "nonexistent@example.com")
		assert.ErrorIs(t, err, errors.ErrNotFound)

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetPaged", func(t *testing.T) {
		users := []domain.User{*testutil.NewTestUser(), *testutil.NewTestUser()}
		opts := repository.ListOptions{}

		// Setup expectations
		mockRepo.On("GetPaged", ctx, opts).Return(users, nil)
		mockRepo.On("Count", ctx, opts).Return(int64(len(users)), nil)

		// Test
		found, err := userService.GetPaged(ctx, opts)
		require.NoError(t, err)
		assert.Len(t, found, len(users))

		mockRepo.AssertExpectations(t)
	})

	t.Run("Update", func(t *testing.T) {
		user := testutil.NewTestUser()

		// Setup expectations
		mockRepo.On("Update", ctx, user).Return(nil)

		// Test
		err := userService.Update(ctx, user)
		require.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		user := testutil.NewTestUser()

		// Setup expectations
		mockRepo.On("Delete", ctx, user.ID).Return(nil)

		// Test
		err := userService.Delete(ctx, user.ID)
		require.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Search", func(t *testing.T) {
		users := []*domain.User{testutil.NewTestUser(), testutil.NewTestUser()}
		query := "test"
		opts := repository.ListOptions{}

		// Setup expectations
		mockRepo.On("Search", ctx, query, opts).Return(users, nil)

		// Test
		found, err := userService.Search(ctx, query, opts)
		require.NoError(t, err)
		assert.Len(t, found, len(users))

		mockRepo.AssertExpectations(t)
	})
}
