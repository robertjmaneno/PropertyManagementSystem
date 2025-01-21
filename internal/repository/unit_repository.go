package repository

import (
	"context"
	"errors"

	"github.com/yourusername/projectname/internal/domain"
	apperrors "github.com/yourusername/projectname/pkg/errors"
	"gorm.io/gorm"
)

type unitRepository struct {
	*GenericRepository[domain.Unit]
}

// NewUnitRepository creates a new unit repository
func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepository{
		GenericRepository: NewGenericRepository(db, domain.Unit{}, WithHardDelete[domain.Unit]),
	}
}

// GetByUnit retrieves a unit by its unitnumber
func (r *unitRepository) GetByUnitNumber(ctx context.Context, unit_number string) (*domain.Unit, error) {
	var unit domain.Unit
	query := r.GetDB().WithContext(ctx)

	// Apply default tenant filters if needed
	opts := ListOptions{
		TenantFilters: make(map[string]interface{}),
	}
	query = r.applyMultitenancyFilters(query, opts)

	if err := query.Where("unit_number = ?", unit_number).First(&unit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	return &unit, nil
}

// Search searches for units based on a query string
func (r *unitRepository) Search(ctx context.Context, query string, opts ListOptions) ([]*domain.Unit, error) {
	searchFields := []string{"unit_number"}

	// Use the generic search with our search fields
	entities, err := r.GenericRepository.Search(ctx, opts, searchFields)
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice
	units := make([]*domain.Unit, len(entities))
	for i := range entities {
		units[i] = &entities[i]
	}

	return units, nil
}

// GetByID, GetPaged, Create, Update, Delete, and Count are implemented by GenericRepository
// We only need to override them if we want to add unit-specific functionality
