package repository

import (
	"context"
	"errors"

	"github.com/yourusername/projectname/internal/domain"
	apperrors "github.com/yourusername/projectname/pkg/errors"
	"gorm.io/gorm"
)

type buildingRepository struct {
	*GenericRepository[domain.Building]
}

// NewBuildingRepository creates a new building repository
func NewBuildingRepository(db *gorm.DB) BuildingRepository {
	return &buildingRepository{
		GenericRepository: NewGenericRepository(db, domain.Building{}, WithHardDelete[domain.Building]),
	}
}

// GetByName retrieves a building by its name
func (r *buildingRepository) GetByName(ctx context.Context, name string) (*domain.Building, error) {
	var building domain.Building
	query := r.GetDB().WithContext(ctx)

	// Apply default tenant filters if needed
	opts := ListOptions{
		TenantFilters: make(map[string]interface{}),
	}
	query = r.applyMultitenancyFilters(query, opts)

	if err := query.Where("name = ?", name).First(&building).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	return &building, nil
}

// Search searches for buildings based on a query string
func (r *buildingRepository) Search(ctx context.Context, query string, opts ListOptions) ([]*domain.Building, error) {
	searchFields := []string{"name", "address", "description"}

	// Use the generic search with our search fields
	entities, err := r.GenericRepository.Search(ctx, opts, searchFields)
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice
	buildings := make([]*domain.Building, len(entities))
	for i := range entities {
		buildings[i] = &entities[i]
	}

	return buildings, nil
}

// GetByOrganizationID retrieves buildings by organization ID
func (r *buildingRepository) GetByOrganizationID(ctx context.Context, organizationID string, opts ListOptions) ([]*domain.Building, error) {
	var buildings []domain.Building
	query := r.GetDB().WithContext(ctx)

	// Apply multitenancy filters if necessary
	query = r.applyMultitenancyFilters(query, opts)

	if err := query.Where("organization_id = ?", organizationID).Find(&buildings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrInternalServer
	}

	// Convert to pointer slice
	buildingPtrs := make([]*domain.Building, len(buildings))
	for i := range buildings {
		buildingPtrs[i] = &buildings[i]
	}

	return buildingPtrs, nil
}

// GetByID, GetPaged, Create, Update, Delete, and Count are implemented by GenericRepository
// Additional methods for domain-specific logic can be added as needed
