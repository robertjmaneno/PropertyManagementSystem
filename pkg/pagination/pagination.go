package pagination

import (
	"math"

	"gorm.io/gorm"
)

// Params holds pagination parameters
type Params struct {
	Page     int   `json:"page" form:"page"`
	PageSize int   `json:"page_size" form:"page_size"`
	Total    int64 `json:"total"`
}

// Response represents a paginated response
type Response struct {
	Data       interface{} `json:"data"`
	Pagination Params      `json:"pagination"`
}

// Calculate returns the offset and limit for the query
func (p *Params) Calculate() (offset, limit int) {
	// Ensure valid page size
	if p.PageSize <= 0 {
		p.PageSize = 10 // Default page size
	}
	if p.PageSize > 100 {
		p.PageSize = 100 // Maximum page size
	}

	// Ensure valid page number
	if p.Page <= 0 {
		p.Page = 1
	}

	// Calculate offset and limit
	offset = (p.Page - 1) * p.PageSize
	limit = p.PageSize

	return offset, limit
}

// TotalPages calculates the total number of pages
func (p *Params) TotalPages() int {
	return int(math.Ceil(float64(p.Total) / float64(p.PageSize)))
}

// HasNext checks if there is a next page
func (p *Params) HasNext() bool {
	return p.Page < p.TotalPages()
}

// HasPrev checks if there is a previous page
func (p *Params) HasPrev() bool {
	return p.Page > 1
}

// Paginate adds pagination to a GORM query
func Paginate(params *Params) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset, limit := params.Calculate()

		// First, get total count
		var count int64
		db.Count(&count)
		params.Total = count

		// Then add pagination
		return db.Offset(offset).Limit(limit)
	}
} 