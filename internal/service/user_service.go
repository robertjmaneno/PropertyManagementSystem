package service

import (
	"context"

	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/pkg/errors"
)

// UserService handles user-related business logic
type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetPaged(ctx context.Context, opts repository.ListOptions) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, opts repository.ListOptions) (int64, error)
	Search(ctx context.Context, query string, opts repository.ListOptions) ([]*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, user *domain.User) error {
	// Validate user data
	if err := user.Validate(); err != nil {
		return err
	}

	// Check if email already exists
	existingUser, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil && err != errors.ErrNotFound {
		return err
	}
	if existingUser != nil {
		return errors.ErrDuplicateEntry
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) GetPaged(ctx context.Context, opts repository.ListOptions) ([]domain.User, error) {
	// Ensure tenant filters are initialized
	if opts.TenantFilters == nil {
		opts.TenantFilters = make(map[string]interface{})
	}

	return s.repo.GetPaged(ctx, opts)
}

func (s *userService) Update(ctx context.Context, user *domain.User) error {
	// Validate user data
	if err := user.Validate(); err != nil {
		return err
	}

	// Check if user exists
	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}

	// Check if email is being changed and if it's already taken
	if existingUser.Email != user.Email {
		emailUser, err := s.repo.GetByEmail(ctx, user.Email)
		if err != nil && err != errors.ErrNotFound {
			return err
		}
		if emailUser != nil {
			return errors.ErrDuplicateEntry
		}
	}

	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) Count(ctx context.Context, opts repository.ListOptions) (int64, error) {
	// Ensure tenant filters are initialized
	if opts.TenantFilters == nil {
		opts.TenantFilters = make(map[string]interface{})
	}

	return s.repo.Count(ctx, opts)
}

func (s *userService) Search(ctx context.Context, query string, opts repository.ListOptions) ([]*domain.User, error) {
	// Ensure tenant filters are initialized
	if opts.TenantFilters == nil {
		opts.TenantFilters = make(map[string]interface{})
	}

	return s.repo.Search(ctx, query, opts)
}
