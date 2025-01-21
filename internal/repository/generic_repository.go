package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	apperrors "github.com/yourusername/projectname/pkg/errors"
	"gorm.io/gorm"
)

// QueryModifier allows customizing the query before execution
type QueryModifier func(*gorm.DB) *gorm.DB

// GenericRepository provides a generic implementation of Repository interface
type GenericRepository[T any] struct {
	*BaseRepository
	model      T
	softDelete bool
}

// NewGenericRepository creates a new generic repository instance
func NewGenericRepository[T any](db *gorm.DB, model T, opts ...func(*GenericRepository[T])) *GenericRepository[T] {
	repo := &GenericRepository[T]{
		BaseRepository: NewBaseRepository(db),
		model:          model,
		softDelete:     true, // Enable soft delete by default
	}

	for _, opt := range opts {
		opt(repo)
	}

	return repo
}

// WithHardDelete disables soft delete for the repository
func WithHardDelete[T any](r *GenericRepository[T]) {
	r.softDelete = false
}

// applyMultitenancyFilters enhances the existing method with better filtering
func (r *GenericRepository[T]) applyMultitenancyFilters(query *gorm.DB, opts ListOptions) *gorm.DB {
	// ... existing multitenancy code ...

	// Add support for additional tenant-specific filters
	if opts.TenantFilters != nil {
		for field, value := range opts.TenantFilters {
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
	}

	return query
}

// Create with transaction support
func (r *GenericRepository[T]) CreateTx(ctx context.Context, entity *T, tx *gorm.DB) error {
	db := r.getDBContext(ctx, tx)
	if err := db.Create(entity).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return apperrors.ErrDuplicateEntry
		}
		return fmt.Errorf("%w: %v", apperrors.ErrInternalServer, err)
	}
	return nil
}

// BulkCreate implements batch insertion
func (r *GenericRepository[T]) BulkCreate(ctx context.Context, entities []T) error {
	if len(entities) == 0 {
		return nil
	}

	return r.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		const batchSize = 1000
		for i := 0; i < len(entities); i += batchSize {
			end := i + batchSize
			if end > len(entities) {
				end = len(entities)
			}

			if err := tx.Create(entities[i:end]).Error; err != nil {
				return fmt.Errorf("%w: %v", apperrors.ErrInternalServer, err)
			}
		}
		return nil
	})
}

// GetByIDWithQuery allows custom query modifications
func (r *GenericRepository[T]) GetByIDWithQuery(ctx context.Context, id interface{}, opts ListOptions, modifiers ...QueryModifier) (*T, error) {
	var entity T
	query := r.GetDB().WithContext(ctx)
	query = r.applyMultitenancyFilters(query, opts)

	for _, modifier := range modifiers {
		query = modifier(query)
	}

	if err := query.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInternalServer, err)
	}
	return &entity, nil
}

// Delete removes an entity by its ID
func (r *GenericRepository[T]) Delete(ctx context.Context, id string) error {
	var entity T
	if err := r.GetDB().WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound
		}
		return apperrors.ErrInternalServer
	}

	result := r.GetDB().WithContext(ctx).Delete(&entity)
	if result.Error != nil {
		return apperrors.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// Search implements a generic search functionality
func (r *GenericRepository[T]) Search(ctx context.Context, opts ListOptions, searchFields []string) ([]T, error) {
	var entities []T
	query := r.GetDB().WithContext(ctx)
	query = r.applyMultitenancyFilters(query, opts)

	if opts.Search != "" && len(searchFields) > 0 {
		searchQuery := "%" + opts.Search + "%"
		conditions := make([]string, len(searchFields))
		values := make([]interface{}, len(searchFields))

		for i, field := range searchFields {
			conditions[i] = fmt.Sprintf("%s LIKE ?", field)
			values[i] = searchQuery
		}

		query = query.Where(strings.Join(conditions, " OR "), values...)
	}

	// Apply standard options
	query = r.ApplyIncludes(query, opts.Includes)
	query = r.ApplySort(query, opts.Sort)
	query = r.ApplyPagination(query, opts.Pagination)

	if err := query.Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrInternalServer, err)
	}

	return entities, nil
}

// Create inserts a new entity into the database
func (r *GenericRepository[T]) Create(ctx context.Context, entity *T) error {
	if err := r.GetDB().WithContext(ctx).Create(entity).Error; err != nil {
		return apperrors.ErrInternalServer
	}
	return nil
}

// GetByID retrieves an entity by ID
func (r *GenericRepository[T]) GetByID(ctx context.Context, id string) (*T, error) {
	var entity T
	if err := r.GetDB().WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrInternalServer
	}
	return &entity, nil
}

// GetAll retrieves all entities
func (r *GenericRepository[T]) GetAll(ctx context.Context) ([]T, error) {
	var entities []T
	if err := r.GetDB().WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, apperrors.ErrInternalServer
	}
	return entities, nil
}

// GetPaged retrieves a paginated list of entities
func (r *GenericRepository[T]) GetPaged(ctx context.Context, opts ListOptions) ([]T, error) {
	var entities []T
	query := r.GetDB().WithContext(ctx)

	// Apply base query modifiers
	query = r.ApplyIncludes(query, opts.Includes)
	query = r.ApplySort(query, opts.Sort)
	query = r.ApplyPagination(query, opts.Pagination)

	if err := query.Find(&entities).Error; err != nil {
		return nil, apperrors.ErrInternalServer
	}

	return entities, nil
}

// Update updates an entity in the database
func (r *GenericRepository[T]) Update(ctx context.Context, entity *T) error {
	result := r.GetDB().WithContext(ctx).Save(entity)
	if result.Error != nil {
		return apperrors.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// Count returns the total number of entities matching the criteria
func (r *GenericRepository[T]) Count(ctx context.Context, opts ListOptions) (int64, error) {
	var count int64
	query := r.GetDB().WithContext(ctx).Model(new(T))

	if err := query.Count(&count).Error; err != nil {
		return 0, apperrors.ErrInternalServer
	}

	return count, nil
}

// getDBContext returns either the transaction or the regular db connection
func (r *GenericRepository[T]) getDBContext(ctx context.Context, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.GetDB().WithContext(ctx)
}
