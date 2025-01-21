package repository

import (
	"context"
	"errors"

	"github.com/yourusername/projectname/internal/domain"
	apperrors "github.com/yourusername/projectname/pkg/errors"
	"gorm.io/gorm"
)

type userRepository struct {
	*GenericRepository[domain.User]
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		GenericRepository: NewGenericRepository(db, domain.User{}, WithHardDelete[domain.User]),
	}
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := r.GetDB().WithContext(ctx)

	// Apply default tenant filters
	opts := ListOptions{
		TenantFilters: make(map[string]interface{}),
	}
	query = r.applyMultitenancyFilters(query, opts)

	if err := query.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	return &user, nil
}

// Search searches for users based on a query string
func (r *userRepository) Search(ctx context.Context, query string, opts ListOptions) ([]*domain.User, error) {
	searchFields := []string{"email", "first_name", "last_name"}

	// Use the generic search with our search fields
	entities, err := r.GenericRepository.Search(ctx, opts, searchFields)
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice
	users := make([]*domain.User, len(entities))
	for i := range entities {
		users[i] = &entities[i]
	}

	return users, nil
}

// GetByID, GetPaged, Create, Update, Delete, and Count are implemented by GenericRepository
// We only need to override them if we want to add user-specific functionality
