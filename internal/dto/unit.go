package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/projectname/internal/domain"
)

// CreateUnitRequest represents the request body for creating a unit
type CreateUnitRequest struct {
	UnitNumber     string  `json:"unit_number" binding:"required"`
	Rooms          int     `json:"rooms" binding:"required"`
	Bathrooms      float64 `json:"bathrooms" binding:"required"`
	Description    string  `json:"description,omitempty"`
	SquareFeet     int     `json:"square_feet" binding:"required"`
	OrganizationID uint    `json:"organization_id" binding:"required"`
	BranchID       string  `json:"branch_id" binding:"required"`
}

// UpdateUnitRequest represents the request body for updating a unit
type UpdateUnitRequest struct {
	UnitNumber     string  `json:"unit_number,omitempty"`
	Rooms          int     `json:"rooms,omitempty"`
	Bathrooms      float64 `json:"bathrooms,omitempty"`
	Description    string  `json:"description,omitempty"`
	SquareFeet     int     `json:"square_feet,omitempty"`
	OrganizationID uint    `json:"organization_id,omitempty"`
	BranchID       string  `json:"branch_id,omitempty"`
}

// UnitResponse represents the response body for unit operations
type UnitResponse struct {
	ID             string    `json:"id"`
	UnitNumber     string    `json:"unit_number"`
	Rooms          int       `json:"rooms"`
	Bathrooms      float64   `json:"bathrooms"`
	Description    string    `json:"description,omitempty"`
	SquareFeet     int       `json:"square_feet"`
	OrganizationID uint      `json:"organization_id"`
	BranchID       string    `json:"branch_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ListUnitsResponse represents the response body for listing units
type ListUnitsResponse struct {
	Units []UnitResponse `json:"units"`
	Total int64          `json:"total"`
}

// ToUnit converts a CreateUnitRequest to a domain.Unit
func (r *CreateUnitRequest) ToUnit() *domain.Unit {
	return &domain.Unit{
		ID:             uuid.New().String(),
		UnitNumber:     r.UnitNumber,
		Rooms:          r.Rooms,
		Bathrooms:      r.Bathrooms,
		Description:    r.Description,
		SquareFeet:     r.SquareFeet,
		OrganizationID: r.OrganizationID,
		BranchID:       r.BranchID,
	}
}

// FromUnit creates a UnitResponse from a domain.Unit
func FromUnit(unit *domain.Unit) *UnitResponse {
	return &UnitResponse{
		ID:             unit.ID,
		UnitNumber:     unit.UnitNumber,
		Rooms:          unit.Rooms,
		Bathrooms:      unit.Bathrooms,
		Description:    unit.Description,
		SquareFeet:     unit.SquareFeet,
		OrganizationID: unit.OrganizationID,
		BranchID:       unit.BranchID,
		CreatedAt:      unit.CreatedAt,
		UpdatedAt:      unit.UpdatedAt,
	}
}

// FromUnits creates a ListUnitsResponse from a slice of domain.Unit
func FromUnits(units []domain.Unit, total int64) *ListUnitsResponse {
	responses := make([]UnitResponse, len(units))
	for i, unit := range units {
		response := FromUnit(&unit)
		responses[i] = *response
	}

	return &ListUnitsResponse{
		Units: responses,
		Total: total,
	}
}
