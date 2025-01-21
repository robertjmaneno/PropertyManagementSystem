package service

import (
	"context"

	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/pkg/errors"
)

// UnitService handles unit-related business logic
type UnitService interface {
	Create(ctx context.Context, unit *domain.Unit) error
	GetByID(ctx context.Context, id string) (*domain.Unit, error)
	GetPaged(ctx context.Context, opts repository.ListOptions) ([]domain.Unit, error)
	Update(ctx context.Context, unit *domain.Unit) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, opts repository.ListOptions) (int64, error)
	Search(ctx context.Context, query string, opts repository.ListOptions) ([]*domain.Unit, error)
}

type unitService struct {
	repo repository.UnitRepository
}

// NewUnitService creates a new unit service instance
func NewUnitService(repo repository.UnitRepository) UnitService {
	return &unitService{repo: repo}
}

func (s *unitService) Create(ctx context.Context, unit *domain.Unit) error {
	// Validate unit data
	if err := unit.Validate(); err != nil {
		return err
	}

	// Check if a unit with the same number exists
	existingUnit, err := s.repo.GetByUnitNumber(ctx, unit.UnitNumber)
	if err != nil && err != errors.ErrNotFound {
		return err
	}
	if existingUnit != nil {
		return errors.ErrDuplicateEntry
	}

	return s.repo.Create(ctx, unit)
}

func (s *unitService) GetByID(ctx context.Context, id string) (*domain.Unit, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *unitService) GetPaged(ctx context.Context, opts repository.ListOptions) ([]domain.Unit, error) {
	// Ensure tenant filters are initialized
	if opts.TenantFilters == nil {
		opts.TenantFilters = make(map[string]interface{})
	}

	return s.repo.GetPaged(ctx, opts)
}

func (s *unitService) Update(ctx context.Context, unit *domain.Unit) error {
	// Validate unit data
	if err := unit.Validate(); err != nil {
		return err
	}

	// Check if the unit exists
	existingUnit, err := s.repo.GetByID(ctx, unit.ID)
	if err != nil {
		return err
	}

	// Check if the unit number is being changed and if it's already taken
	if existingUnit.UnitNumber != unit.UnitNumber {
		unitWithNumber, err := s.repo.GetByUnitNumber(ctx, unit.UnitNumber)
		if err != nil && err != errors.ErrNotFound {
			return err
		}
		if unitWithNumber != nil {
			return errors.ErrDuplicateEntry
		}
	}

	return s.repo.Update(ctx, unit)
}

func (s *unitService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *unitService) Count(ctx context.Context, opts repository.ListOptions) (int64, error) {
	// Ensure tenant filters are initialized
	if opts.TenantFilters == nil {
		opts.TenantFilters = make(map[string]interface{})
	}

	return s.repo.Count(ctx, opts)
}

func (s *unitService) Search(ctx context.Context, query string, opts repository.ListOptions) ([]*domain.Unit, error) {
	// Ensure tenant filters are initialized
	if opts.TenantFilters == nil {
		opts.TenantFilters = make(map[string]interface{})
	}

	return s.repo.Search(ctx, query, opts)
}
