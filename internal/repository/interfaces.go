package repository

import (
	"context"

	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/pkg/pagination"
	"github.com/yourusername/projectname/pkg/sorting"
)

// ListOptions combines pagination and sorting options
type ListOptions struct {
	TenantFilters map[string]interface{}
	Pagination    *pagination.Params
	Sort          *sorting.Options
	Search        string
	Includes      []string
}

// Repository defines the generic repository interface
type Repository[T any] interface {
	GetByID(ctx context.Context, id string) (*T, error)
	GetAll(ctx context.Context) ([]T, error)
	GetPaged(ctx context.Context, opts ListOptions) ([]T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, opts ListOptions) (int64, error)
}

// UserRepository extends Repository for User-specific operations
type UserRepository interface {
	Repository[domain.User]
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Search(ctx context.Context, query string, opts ListOptions) ([]*domain.User, error)
}

// CommunityRepository extends Repository for Community-specific operations
type CommunityRepository interface {
	Repository[domain.Community]
	GetByName(ctx context.Context, name string) (*domain.Community, error)
	Search(ctx context.Context, query string, opts ListOptions) ([]*domain.Community, error)
}

// BuildingRepository extends Repository for Buiding-specific operations
type BuildingRepository interface {
	Repository[domain.Building]
	GetByName(ctx context.Context, name string) (*domain.Building, error)
	Search(ctx context.Context, query string, opts ListOptions) ([]*domain.Building, error)
}

// UnitRepository extends Repository for unit-specific operations
type UnitRepository interface {
	Repository[domain.Unit]
	GetByUnitNumber(ctx context.Context, unit_number string) (*domain.Unit, error)
	Search(ctx context.Context, query string, opts ListOptions) ([]*domain.Unit, error)
}

type UnitOfWork interface {
	Users() UserRepository
	Communities() CommunityRepository
	Buildings() BuildingRepository
	Units() UnitRepository
	Begin() error
	Commit() error
	Rollback() error
}
