package service

import (
	"context"

	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/pkg/errors"
)

// BuildingService handles building-related business logic
type BuildingService interface {
	Create(ctx context.Context, building *domain.Building) error
	GetByID(ctx context.Context, id string) (*domain.Building, error)
	GetPaged(ctx context.Context, opts repository.ListOptions) ([]domain.Building, error)
	Update(ctx context.Context, building *domain.Building) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, opts repository.ListOptions) (int64, error)
	Search(ctx context.Context, query string, opts repository.ListOptions) ([]*domain.Building, error)
}

type buildingService struct {
	repo repository.BuildingRepository
}

// NewBuildingService creates a new building service instance
func NewBuildingService(repo repository.BuildingRepository) BuildingService {
	return &buildingService{repo: repo}
}

func (s *buildingService) Create(ctx context.Context, building *domain.Building) error {
	// Validate building data
	if err := building.Validate(); err != nil {
		return err
	}

	// Check if a building with the same name exists
	existingBuilding, err := s.repo.GetByName(ctx, building.Name)
	if err != nil && err != errors.ErrNotFound {
		return err
	}
	if existingBuilding != nil {
		return errors.ErrDuplicateEntry
	}

	return s.repo.Create(ctx, building)
}

func (s *buildingService) GetByID(ctx context.Context, id string) (*domain.Building, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *buildingService) GetPaged(ctx context.Context, opts repository.ListOptions) ([]domain.Building, error) {
	// Ensure tenant filters are initialized
	if opts.TenantFilters == nil {
		opts.TenantFilters = make(map[string]interface{})
	}

	return s.repo.GetPaged(ctx, opts)
}

func (s *buildingService) Update(ctx context.Context, building *domain.Building) error {
	// Validate building data
	if err := building.Validate(); err != nil {
		return err
	}

	// Check if the building exists
	existingBuilding, err := s.repo.GetByID(ctx, building.ID)
	if err != nil {
		return err
	}

	// Check if the name is being changed and if it's already taken
	if existingBuilding.Name != building.Name {
		nameBuilding, err := s.repo.GetByName(ctx, building.Name)
		if err != nil && err != errors.ErrNotFound {
			return err
		}
		if nameBuilding != nil {
			return errors.ErrDuplicateEntry
		}
	}

	return s.repo.Update(ctx, building)
}

func (s *buildingService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *buildingService) Count(ctx context.Context, opts repository.ListOptions) (int64, error) {
	// Ensure tenant filters are initialized
	if opts.TenantFilters == nil {
		opts.TenantFilters = make(map[string]interface{})
	}

	return s.repo.Count(ctx, opts)
}

func (s *buildingService) Search(ctx context.Context, query string, opts repository.ListOptions) ([]*domain.Building, error) {
	// Ensure tenant filters are initialized
	if opts.TenantFilters == nil {
		opts.TenantFilters = make(map[string]interface{})
	}

	return s.repo.Search(ctx, query, opts)
}
