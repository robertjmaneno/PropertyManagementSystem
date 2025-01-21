package service

import (
	"context"

	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/pkg/errors"
)

// CommunityService handles community-related business logic
type CommunityService interface {
	Create(ctx context.Context, community *domain.Community) error
	GetByID(ctx context.Context, id string) (*domain.Community, error)
	GetPaged(ctx context.Context, opts repository.ListOptions) ([]domain.Community, error)
	Update(ctx context.Context, community *domain.Community) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, opts repository.ListOptions) (int64, error)
	Search(ctx context.Context, query string, opts repository.ListOptions) ([]*domain.Community, error)
}

type communityService struct {
	repo repository.CommunityRepository
}

// NewCommunityService creates a new community service instance
func NewCommunityService(repo repository.CommunityRepository) CommunityService {
	return &communityService{repo: repo}
}

func (s *communityService) Create(ctx context.Context, community *domain.Community) error {
	// Validate community data
	if err := community.Validate(); err != nil {
		return err
	}

	// Check if a community with the same name exists
	existingCommunity, err := s.repo.GetByName(ctx, community.Name)
	if err != nil && err != errors.ErrNotFound {
		return err
	}
	if existingCommunity != nil {
		return errors.ErrDuplicateEntry
	}

	return s.repo.Create(ctx, community)
}

func (s *communityService) GetByID(ctx context.Context, id string) (*domain.Community, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *communityService) GetPaged(ctx context.Context, opts repository.ListOptions) ([]domain.Community, error) {
	// Ensure tenant filters are initialized
	if opts.TenantFilters == nil {
		opts.TenantFilters = make(map[string]interface{})
	}

	return s.repo.GetPaged(ctx, opts)
}

func (s *communityService) Update(ctx context.Context, community *domain.Community) error {
	// Validate community data
	if err := community.Validate(); err != nil {
		return err
	}

	// Check if the community exists
	existingCommunity, err := s.repo.GetByID(ctx, community.ID)
	if err != nil {
		return err
	}

	// Check if the name is being changed and if it's already taken
	if existingCommunity.Name != community.Name {
		nameCommunity, err := s.repo.GetByName(ctx, community.Name)
		if err != nil && err != errors.ErrNotFound {
			return err
		}
		if nameCommunity != nil {
			return errors.ErrDuplicateEntry
		}
	}

	return s.repo.Update(ctx, community)
}

func (s *communityService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *communityService) Count(ctx context.Context, opts repository.ListOptions) (int64, error) {
	// Ensure tenant filters are initialized
	if opts.TenantFilters == nil {
		opts.TenantFilters = make(map[string]interface{})
	}

	return s.repo.Count(ctx, opts)
}

func (s *communityService) Search(ctx context.Context, query string, opts repository.ListOptions) ([]*domain.Community, error) {
	// Ensure tenant filters are initialized
	if opts.TenantFilters == nil {
		opts.TenantFilters = make(map[string]interface{})
	}

	return s.repo.Search(ctx, query, opts)
}
