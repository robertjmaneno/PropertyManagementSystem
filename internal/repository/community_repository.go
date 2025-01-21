package repository

import (
	"context"
	"errors"

	"github.com/yourusername/projectname/internal/domain"
	apperrors "github.com/yourusername/projectname/pkg/errors"
	"gorm.io/gorm"
)

type communityRepository struct {
	*GenericRepository[domain.Community]
}

// NewCommunityRepository creates a new community repository
func NewCommunityRepository(db *gorm.DB) CommunityRepository {
	return &communityRepository{
		GenericRepository: NewGenericRepository(db, domain.Community{}, WithHardDelete[domain.Community]),
	}
}

// GetByName retrieves a community by its name
func (r *communityRepository) GetByName(ctx context.Context, name string) (*domain.Community, error) {
	var community domain.Community
	query := r.GetDB().WithContext(ctx)

	// Apply default tenant filters if needed
	opts := ListOptions{
		TenantFilters: make(map[string]interface{}),
	}
	query = r.applyMultitenancyFilters(query, opts)

	if err := query.Where("name = ?", name).First(&community).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	return &community, nil
}

// Search searches for communities based on a query string
func (r *communityRepository) Search(ctx context.Context, query string, opts ListOptions) ([]*domain.Community, error) {
	searchFields := []string{"name", "address", "description"}

	// Use the generic search with our search fields
	entities, err := r.GenericRepository.Search(ctx, opts, searchFields)
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice
	communities := make([]*domain.Community, len(entities))
	for i := range entities {
		communities[i] = &entities[i]
	}

	return communities, nil
}

// GetByOrganizationID retrieves communities by organization ID
func (r *communityRepository) GetByOrganizationID(ctx context.Context, organizationID string, opts ListOptions) ([]*domain.Community, error) {
	var communities []domain.Community
	query := r.GetDB().WithContext(ctx)

	// Apply multitenancy filters if necessary
	query = r.applyMultitenancyFilters(query, opts)

	if err := query.Where("organization_id = ?", organizationID).Find(&communities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrInternalServer
	}

	// Convert to pointer slice
	communityPtrs := make([]*domain.Community, len(communities))
	for i := range communities {
		communityPtrs[i] = &communities[i]
	}

	return communityPtrs, nil
}

// GetByID, GetPaged, Create, Update, Delete, and Count are implemented by GenericRepository
// Additional methods for domain-specific logic can be added as needed
