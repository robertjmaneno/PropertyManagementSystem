package handler

import (
	"fmt"
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

type BuildingHandler struct {
	buildingService service.BuildingService
}

func NewBuildingHandler(buildingService service.BuildingService) *BuildingHandler {
	return &BuildingHandler{
		buildingService: buildingService,
	}
}

// getTenancyContext extracts tenancy information from gin context
func (h *BuildingHandler) getTenancyContext(c *gin.Context) *middleware.TenancyContext {
	tenancy, exists := c.Get(middleware.TenancyKey)
	if !exists {
		return nil
	}
	return tenancy.(*middleware.TenancyContext)
}

// getListOptions gets list options with tenant filters from context
func (h *BuildingHandler) getListOptions(c *gin.Context) repository.ListOptions {
	opts, exists := c.Get("list_options")
	if !exists {
		return repository.ListOptions{
			TenantFilters: make(map[string]interface{}),
		}
	}
	return opts.(repository.ListOptions)
}

// RegisterRoutes registers building routes
func (h *BuildingHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Change the base route to /buildings
	buildings := router.Group("/buildings")
	{
		buildings.POST("", h.Create)
		buildings.GET("", h.List)
		buildings.GET("/search", h.Search)
		buildings.GET("/:id", h.GetByID)
		buildings.PUT("/:id", h.Update)
		buildings.DELETE("/:id", h.Delete)
	}
}

// Create handles POST /building
// @Summary Create building
// @Description Create a new building with the provided details.
// @Tags building
// @Accept json
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param request body dto.CreateBuildingRequest true "Building information"
// @Success 201 {object} dto.BuildingResponse "Successfully created building"
// @Failure 400 {object} map[string]string "Invalid input, bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /buildings [post]
func (h *BuildingHandler) Create(c *gin.Context) {
	var req dto.CreateBuildingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get tenancy context
	tenancy := h.getTenancyContext(c)
	fmt.Println(tenancy)
	if tenancy == nil {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	building := req.ToBuilding()
	building.OrganizationID = tenancy.OrganizationID
	building.BranchID = tenancy.BranchID

	if err := h.buildingService.Create(c.Request.Context(), building); err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusCreated, dto.FromBuilding(building))
}

// List handles GET /building
// @Summary List buildings
// @Description Retrieve a list of buildings with pagination and sorting support.
// @Tags building
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param page_size query int false "Number of items per page (default: 10)" minimum(1) maximum(100)
// @Param sort query string false "Sort fields (format: field1:asc,field2:desc) - Available fields: first_name, last_name, email, created_at"
// @Success 200 {object} dto.ListBuildingsResponse "List of buildings with pagination information"
// @Failure 400 {object} map[string]string "Invalid input parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /buildings [get]
func (h *BuildingHandler) List(c *gin.Context) {
	var req dto.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get list options with tenant filters
	opts := h.getListOptions(c)

	// Get total count
	total, err := h.buildingService.Count(c.Request.Context(), opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	// Get buildings with pagination and sorting
	buildings, err := h.buildingService.GetPaged(c.Request.Context(), opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromBuildings(buildings, total))
}

// Search handles GET /building/search
// @Summary Search building
// @Description Search for a building by name or email with optional case sensitivity, pagination, and sorting.
// @Tags building
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param search query string true "Search query (e.g., building name )"
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param case_sensitive query bool false "Case sensitive search (default: false)"
// @Success 200 {object} dto.ListBuildingsResponse "List of buildings matching search criteria with pagination info"
// @Failure 400 {object} map[string]string "Invalid input parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /buildings/search [get]
func (h *BuildingHandler) Search(c *gin.Context) {
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

	// Get buildings with search, pagination, and sorting
	buildings, err := h.buildingService.Search(c.Request.Context(), query, opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	// Convert pointer slice to value slice
	buildingSlice := make([]domain.Building, len(buildings))
	for i, building := range buildings {
		buildingSlice[i] = *building
	}

	// Get total count for search
	total, err := h.buildingService.Count(c.Request.Context(), opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromBuildings(buildingSlice, total))
}

// GetByID handles GET /building/:id
// @Summary Get building by ID
// @Description Retrieve a building's details using its ID
// @Tags building
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param id path string true "building ID"
// @Success 200 {object} dto.BuildingResponse "Building details"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 404 {object} map[string]string "Building not found"
// @Router /buildings/{id} [get]
func (h *BuildingHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	building, err := h.buildingService.GetByID(c.Request.Context(), id)
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
	if tenancy == nil || building.OrganizationID != tenancy.OrganizationID || building.BranchID != tenancy.BranchID {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromBuilding(building))
}

// Update handles PUT /building/:id
// @Summary Update building
// @Description Update a building's details by its ID
// @Tags building
// @Accept json
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param id path string true "Building ID"
// @Param request body dto.UpdateBuildingRequest true "Building information to be updated"
// @Success 200 {object} dto.BuildingResponse "Updated Building details"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 404 {object} map[string]string "Building not found"
// @Router /buildings/{id} [put]
func (h *BuildingHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	var req dto.UpdateBuildingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get existing building
	existingBuilding, err := h.buildingService.GetByID(c.Request.Context(), id)
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
	if tenancy == nil || existingBuilding.OrganizationID != tenancy.OrganizationID || existingBuilding.BranchID != tenancy.BranchID {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	// Update fields
	if req.Name != "" {
		existingBuilding.Name = req.Name
	}
	if req.Address != "" {
		existingBuilding.Address = req.Address
	}
	if req.Description != "" {
		existingBuilding.Description = req.Description
	}

	if err := h.buildingService.Update(c.Request.Context(), existingBuilding); err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromBuilding(existingBuilding))
}

// Delete handles DELETE /building/:id
// @Summary Delete building
// @Description Delete a building by its ID
// @Tags building
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param id path string true "Building ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 404 {object} map[string]string "Building not found"
// @Router /buildings/{id} [delete]
func (h *BuildingHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get existing building
	existingBuilding, err := h.buildingService.GetByID(c.Request.Context(), id)
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
	if tenancy == nil || existingBuilding.OrganizationID != tenancy.OrganizationID || existingBuilding.BranchID != tenancy.BranchID {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	// Perform deletion
	err = h.buildingService.Delete(c.Request.Context(), id)
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
