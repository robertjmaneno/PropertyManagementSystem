package dto

import (
	"github.com/yourusername/projectname/pkg/pagination"
	"github.com/yourusername/projectname/pkg/sorting"
)

// ListRequest represents the request parameters for list operations
type ListRequest struct {
	Page          int    `form:"page" json:"page"`
	PageSize      int    `form:"page_size" json:"page_size"`
	Sort          string `form:"sort" json:"sort"` // Format: "field1:asc:cs,field2:desc"
	Search        string `form:"search" json:"search"`
	CaseSensitive bool   `form:"case_sensitive" json:"case_sensitive"`
}

// ToListOptions converts ListRequest to repository.ListOptions
func (r *ListRequest) ToListOptions() (*pagination.Params, *sorting.Options, error) {
	// Create pagination params
	paginationParams := &pagination.Params{
		Page:     r.Page,
		PageSize: r.PageSize,
	}

	// Create sorting options
	sortOptions := &sorting.Options{}
	if err := sortOptions.ParseFromString(r.Sort); err != nil {
		return nil, nil, err
	}

	return paginationParams, sortOptions, nil
}

// ListResponse represents a paginated response with metadata
type ListResponse struct {
	Data          interface{}       `json:"data"`
	Pagination    pagination.Params `json:"pagination"`
	Sort          []sorting.Field   `json:"sort,omitempty"`
	CaseSensitive bool              `json:"case_sensitive,omitempty"`
	Total         int64             `json:"total"`
}

// NewListResponse creates a new ListResponse
func NewListResponse(data interface{}, params *pagination.Params, sort *sorting.Options, total int64) *ListResponse {
	return &ListResponse{
		Data:       data,
		Pagination: *params,
		Sort:       sort.Fields,
		Total:      total,
	}
}
