package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/projectname/internal/domain"
)

// CreateBuildingRequest represents the request body for creating a building
type CreateBuildingRequest struct {
	Name           string `json:"name" binding:"required"`
	Address        string `json:"address" binding:"required"`
	Description    string `json:"description,omitempty"`
	OrganizationID uint   `json:"organization_id" binding:"required"`
	BranchID       string `json:"branch_id" binding:"required"`
}

// UpdateBuildingRequest represents the request body for updating a building
type UpdateBuildingRequest struct {
	Name           string `json:"name,omitempty"`
	Address        string `json:"address,omitempty"`
	Description    string `json:"description,omitempty"`
	OrganizationID uint   `json:"organization_id,omitempty"`
	BranchID       string `json:"branch_id,omitempty"`
}

// BuildingResponse represents the response body for building operations
type BuildingResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	Description    string    `json:"description,omitempty"`
	OrganizationID uint      `json:"organization_id"`
	BranchID       string    `json:"branch_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ListBuildingsResponse represents the response body for listing buildings
type ListBuildingsResponse struct {
	Buildings []BuildingResponse `json:"buildings"`
	Total     int64              `json:"total"`
}

// ToBuilding converts a CreateBuildingRequest to a domain.Building
func (r *CreateBuildingRequest) ToBuilding() *domain.Building {
	return &domain.Building{
		ID:             uuid.New().String(),
		Name:           r.Name,
		Address:        r.Address,
		Description:    r.Description,
		OrganizationID: r.OrganizationID,
		BranchID:       r.BranchID,
	}
}

// FromBuilding creates a BuildingResponse from a domain.Building
func FromBuilding(building *domain.Building) *BuildingResponse {
	return &BuildingResponse{
		ID:             building.ID,
		Name:           building.Name,
		Address:        building.Address,
		Description:    building.Description,
		OrganizationID: building.OrganizationID,
		BranchID:       building.BranchID,
		CreatedAt:      building.CreatedAt,
		UpdatedAt:      building.UpdatedAt,
	}
}

// FromBuildings creates a ListBuildingsResponse from a slice of domain.Building
func FromBuildings(buildings []domain.Building, total int64) *ListBuildingsResponse {
	responses := make([]BuildingResponse, len(buildings))
	for i, building := range buildings {
		response := FromBuilding(&building)
		responses[i] = *response
	}

	return &ListBuildingsResponse{
		Buildings: responses,
		Total:     total,
	}
}
