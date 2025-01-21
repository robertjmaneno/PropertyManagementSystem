package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/projectname/internal/domain"
)

// CreateCommunityRequest represents the request body for creating a community
type CreateCommunityRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

// UpdateCommunityRequest represents the request body for updating a community
type UpdateCommunityRequest struct {
	Name           string `json:"name,omitempty"`
	Address        string `json:"address,omitempty"`
	Description    string `json:"description,omitempty"`
	OrganizationID uint   `json:"organization_id,omitempty"`
	BranchID       string `json:"branch_id,omitempty"`
}

// CommunityResponse represents the response body for community operations
type CommunityResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	Description    string    `json:"description,omitempty"`
	OrganizationID uint      `json:"organization_id"`
	BranchID       string    `json:"branch_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ListCommunitiesResponse represents the response body for listing communities
type ListCommunitiesResponse struct {
	Communities []CommunityResponse `json:"communities"`
	Total       int64               `json:"total"`
}

// ToCommunity converts a CreateCommunityRequest to a domain.Community
func (r *CreateCommunityRequest) ToCommunity() *domain.Community {
	return &domain.Community{
		ID:          uuid.New().String(),
		Name:        r.Name,
		Address:     r.Address,
		Description: r.Description,
	}
}

// FromCommunity creates a CommunityResponse from a domain.Community
func FromCommunity(community *domain.Community) *CommunityResponse {
	return &CommunityResponse{
		ID:             community.ID,
		Name:           community.Name,
		Address:        community.Address,
		Description:    community.Description,
		OrganizationID: community.OrganizationID,
		BranchID:       community.BranchID,
		CreatedAt:      community.CreatedAt,
		UpdatedAt:      community.UpdatedAt,
	}
}

// FromCommunities creates a ListCommunitiesResponse from a slice of domain.Community
func FromCommunities(communities []domain.Community, total int64) *ListCommunitiesResponse {
	responses := make([]CommunityResponse, len(communities))
	for i, community := range communities {
		response := FromCommunity(&community)
		responses[i] = *response
	}

	return &ListCommunitiesResponse{
		Communities: responses,
		Total:       total,
	}
}
