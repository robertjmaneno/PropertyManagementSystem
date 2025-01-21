package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/internal/dto"
	"github.com/yourusername/projectname/internal/middleware"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/internal/service"
	"github.com/yourusername/projectname/pkg/errors"
	"github.com/yourusername/projectname/pkg/response"
)

type CommunityHandler struct {
	communityService service.CommunityService
}

func (h *CommunityHandler) RegisterPublicRoutes(public *gin.RouterGroup) {
	panic("unimplemented")
}

func NewCommunityHandler(communityService service.CommunityService) *CommunityHandler {
	return &CommunityHandler{
		communityService: communityService,
	}
}

// getTenancyContext extracts tenancy information from gin context
func (h *CommunityHandler) getTenancyContext(c *gin.Context) *middleware.TenancyContext {
	tenancy, exists := c.Get(middleware.TenancyKey)
	if !exists {
		return nil
	}
	return tenancy.(*middleware.TenancyContext)
}

// getListOptions gets list options with tenant filters from context
func (h *CommunityHandler) getListOptions(c *gin.Context) repository.ListOptions {
	opts, exists := c.Get("list_options")
	if !exists {
		return repository.ListOptions{
			TenantFilters: make(map[string]interface{}),
		}
	}
	return opts.(repository.ListOptions)
}

// RegisterRoutes registers community routes
func (h *CommunityHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Change the base route to /communities
	communities := router.Group("/communities")
	{
		communities.POST("", h.Create)
		communities.GET("", h.List)
		communities.GET("/search", h.Search)
		communities.GET("/:id", h.GetByID)
		communities.PUT("/:id", h.Update)
		communities.DELETE("/:id", h.Delete)
	}
}

// Create handles POST /community
// @Summary Create community
// @Description Create a new community with the provided details.
// @Tags community
// @Accept json
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param request body dto.CreateCommunityRequest true "Community information"
// @Success 201 {object} dto.CommunityResponse "Successfully created community"
// @Failure 400 {object} map[string]string "Invalid input, bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /communities [post]
func (h *CommunityHandler) Create(c *gin.Context) {
	var req dto.CreateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get tenancy context
	tenancy := h.getTenancyContext(c)
	if tenancy == nil {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	community := req.ToCommunity()
	community.OrganizationID = tenancy.OrganizationID
	community.BranchID = tenancy.BranchID

	if err := h.communityService.Create(c.Request.Context(), community); err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusCreated, dto.FromCommunity(community))
}

// List handles GET /community
// @Summary List communities
// @Description Retrieve a list of communities with pagination and sorting support.
// @Tags community
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param page_size query int false "Number of items per page (default: 10)" minimum(1) maximum(100)
// @Param sort query string false "Sort fields (format: field1:asc,field2:desc) - Available fields: first_name, last_name, email, created_at"
// @Success 200 {object} dto.ListCommunitiesResponse "List of communities with pagination information"
// @Failure 400 {object} map[string]string "Invalid input parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /communities [get]
func (h *CommunityHandler) List(c *gin.Context) {
	var req dto.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get list options with tenant filters
	opts := h.getListOptions(c)

	// Get total count
	total, err := h.communityService.Count(c.Request.Context(), opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	// Get communities with pagination and sorting
	communities, err := h.communityService.GetPaged(c.Request.Context(), opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromCommunities(communities, total))
}

// Search handles GET /community/search
// @Summary Search community
// @Description Search for a community by name or email with optional case sensitivity, pagination, and sorting.
// @Tags community
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param search query string true "Search query (e.g., community name )"
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param case_sensitive query bool false "Case sensitive search (default: false)"
// @Success 200 {object} dto.ListCommunitiesResponse "List of communities matching search criteria with pagination info"
// @Failure 400 {object} map[string]string "Invalid input parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /communities/search [get]
func (h *CommunityHandler) Search(c *gin.Context) {
	var req dto.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	query := c.Query("search")
	if query == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get list options with tenant filters
	opts := h.getListOptions(c)

	// Get communities with search, pagination, and sorting
	communities, err := h.communityService.Search(c.Request.Context(), query, opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	// Convert pointer slice to value slice
	communitySlice := make([]domain.Community, len(communities))
	for i, community := range communities {
		communitySlice[i] = *community
	}

	// Get total count for search
	total, err := h.communityService.Count(c.Request.Context(), opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromCommunities(communitySlice, total))
}

// GetByID handles GET /communities/:id
// @Summary Get community by ID
// @Description Retrieve a community's details using its ID
// @Tags community
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param id path string true "Community ID"
// @Success 200 {object} dto.CommunityResponse "Community details"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 404 {object} map[string]string "Community not found"
// @Router /communities/{id} [get]
func (h *CommunityHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	community, err := h.communityService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == errors.ErrNotFound {
			response.SendError(c, http.StatusNotFound, err)
			return
		}
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	// Verify tenant access
	tenancy := h.getTenancyContext(c)
	if tenancy == nil || community.OrganizationID != tenancy.OrganizationID || community.BranchID != tenancy.BranchID {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromCommunity(community))
}

// Update handles PUT /community/:id
// @Summary Update community
// @Description Update a community's details by its ID
// @Tags community
// @Accept json
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param id path string true "Community ID"
// @Param request body dto.UpdateCommunityRequest true "Community information to be updated"
// @Success 200 {object} dto.CommunityResponse "Updated community details"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 404 {object} map[string]string "Community not found"
// @Router /communities/{id} [put]
func (h *CommunityHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	var req dto.UpdateCommunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get existing community
	existingCommunity, err := h.communityService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == errors.ErrNotFound {
			response.SendError(c, http.StatusNotFound, err)
			return
		}
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	// Verify tenant access
	tenancy := h.getTenancyContext(c)
	if tenancy == nil || existingCommunity.OrganizationID != tenancy.OrganizationID || existingCommunity.BranchID != tenancy.BranchID {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	// Update fields
	if req.Name != "" {
		existingCommunity.Name = req.Name
	}
	if req.Address != "" {
		existingCommunity.Address = req.Address
	}
	if req.Description != "" {
		existingCommunity.Description = req.Description
	}

	if err := h.communityService.Update(c.Request.Context(), existingCommunity); err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromCommunity(existingCommunity))
}

// Delete handles DELETE /community/:id
// @Summary Delete community
// @Description Delete a community by its ID
// @Tags community
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param id path string true "Community ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 404 {object} map[string]string "Community not found"
// @Router /communities/{id} [delete]
func (h *CommunityHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get existing community
	existingCommunity, err := h.communityService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == errors.ErrNotFound {
			response.SendError(c, http.StatusNotFound, err)
			return
		}
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	// Verify tenant access
	tenancy := h.getTenancyContext(c)
	if tenancy == nil || existingCommunity.OrganizationID != tenancy.OrganizationID || existingCommunity.BranchID != tenancy.BranchID {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	// Perform deletion
	err = h.communityService.Delete(c.Request.Context(), id)
	if err != nil {
		if err == errors.ErrNotFound {
			response.SendError(c, http.StatusNotFound, err)
			return
		}
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	// Success response for deletion
	response.SendSuccess(c, http.StatusNoContent, nil)
}
