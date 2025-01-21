package repository

import (
	"context"
	"fmt"

	apperrors "github.com/yourusername/projectname/pkg/errors"
	"github.com/yourusername/projectname/pkg/pagination"
	"github.com/yourusername/projectname/pkg/sorting"
	"gorm.io/gorm"
)

// BaseRepository provides a base implementation of Repository interface
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository instance
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// GetDB returns the GORM DB instance
func (r *BaseRepository) GetDB() *gorm.DB {
	return r.db
}

// WithTx executes function within transaction with better error handling
func (r *BaseRepository) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := fn(tx); err != nil {
			return fmt.Errorf("%w: %v", apperrors.ErrInternalServer, err)
		}
		return nil
	})
}

// AllowedSortFields returns a map of allowed sort fields for a model
func (r *BaseRepository) AllowedSortFields() map[string]bool {
	return map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
		"deleted_at": true,
	}
}

// DefaultSortOptions returns default sort options for a model
func (r *BaseRepository) DefaultSortOptions() *sorting.Options {
	return sorting.NewOptions("created_at", sorting.DESC)
}

// ApplyPagination applies pagination to the query
func (r *BaseRepository) ApplyPagination(query *gorm.DB, params *pagination.Params) *gorm.DB {
	if params == nil {
		return query
	}
	offset := (params.Page - 1) * params.PageSize
	return query.Offset(int(offset)).Limit(int(params.PageSize))
}

// ApplySort with improved field validation and error handling
func (r *BaseRepository) ApplySort(query *gorm.DB, opts *sorting.Options) *gorm.DB {
	if opts == nil || len(opts.Fields) == 0 {
		// Apply default sorting if no options provided
		return query.Order("created_at DESC")
	}

	allowedFields := r.AllowedSortFields()
	for _, field := range opts.Fields {
		if allowed := allowedFields[field.Name]; allowed {
			query = query.Order(fmt.Sprintf("%s %s", field.Name, field.Direction))
		}
	}

	return query
}

// ApplyIncludes applies includes (preloads) to the query
func (r *BaseRepository) ApplyIncludes(query *gorm.DB, includes []string) *gorm.DB {
	for _, include := range includes {
		query = query.Preload(include)
	}
	return query
}

// ApplyFilters adds a new method for applying generic filters
func (r *BaseRepository) ApplyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for field, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}
	return query
}
