package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/projectname/internal/dto"
	"github.com/yourusername/projectname/internal/middleware"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/internal/service"
	"github.com/yourusername/projectname/pkg/errors"
	"github.com/yourusername/projectname/pkg/response"
)

type UnitHandler struct {
	unitService service.UnitService
}

func NewUnitHandler(unitService service.UnitService) *UnitHandler {
	return &UnitHandler{
		unitService: unitService,
	}
}

// getTenancyContext extracts tenancy information from gin context
func (h *UnitHandler) getTenancyContext(c *gin.Context) *middleware.TenancyContext {
	tenancy, exists := c.Get(middleware.TenancyKey)
	if !exists {
		return nil
	}
	return tenancy.(*middleware.TenancyContext)
}

// getListOptions gets list options with tenant filters from context
func (h *UnitHandler) getListOptions(c *gin.Context) repository.ListOptions {
	opts, exists := c.Get("list_options")
	if !exists {
		return repository.ListOptions{
			TenantFilters: make(map[string]interface{}),
		}
	}
	return opts.(repository.ListOptions)
}

// RegisterRoutes registers building routes
func (h *UnitHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Change the base route to /buildings
	units := router.Group("/units")
	{
		units.POST("", h.Create)
		units.GET("", h.List)
		units.GET("/:id", h.GetByID)
		units.PUT("/:id", h.Update)
		units.DELETE("/:id", h.Delete)
	}
}

// Create handles POST /units
// @Summary Create unit
// @Description Create a new unit with the provided details.
// @Tags unit
// @Accept json
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param request body dto.CreateUnitRequest true "Unit information"
// @Success 201 {object} dto.UnitResponse "Successfully created unit"
// @Failure 400 {object} map[string]string "Invalid input, bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /units [post]
func (h *UnitHandler) Create(c *gin.Context) {
	var req dto.CreateUnitRequest
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

	unit := req.ToUnit()
	unit.OrganizationID = tenancy.OrganizationID
	unit.BranchID = tenancy.BranchID

	if err := h.unitService.Create(c.Request.Context(), unit); err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusCreated, dto.FromUnit(unit))
}

// List handles GET /units
// @Summary List units
// @Description Retrieve a list of units with pagination and sorting support.
// @Tags unit
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param page_size query int false "Number of items per page (default: 10)" minimum(1) maximum(100)
// @Param sort query string false "Sort fields (format: field1:asc,field2:desc) - Available fields: unit_number, name, created_at"
// @Success 200 {object} dto.ListUnitsResponse "List of units with pagination information"
// @Failure 400 {object} map[string]string "Invalid input parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /units [get]
func (h *UnitHandler) List(c *gin.Context) {
	var req dto.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get list options with tenant filters
	opts := h.getListOptions(c)

	// Get total count
	total, err := h.unitService.Count(c.Request.Context(), opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	// Get units with pagination and sorting
	units, err := h.unitService.GetPaged(c.Request.Context(), opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromUnits(units, total))
}

// GetByID handles GET /unit/:id
// @Summary Get unit by ID
// @Description Retrieve a unit's details using its ID
// @Tags unit
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param id path string true "unit ID"
// @Success 200 {object} dto.UnitResponse "Unit details"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 404 {object} map[string]string "Unit not found"
// @Router /units/{id} [get]
func (h *UnitHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Fetch unit details by ID
	unit, err := h.unitService.GetByID(c.Request.Context(), id)
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
	if tenancy == nil || unit.OrganizationID != tenancy.OrganizationID || unit.BranchID != tenancy.BranchID {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromUnit(unit))
}

// Update handles PUT /unit/:id
// @Summary Update unit
// @Description Update a unit's details by its ID
// @Tags unit
// @Accept json
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param id path string true "Unit ID"
// @Param request body dto.UpdateUnitRequest true "Unit information to be updated"
// @Success 200 {object} dto.UnitResponse "Updated Unit details"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 404 {object} map[string]string "Unit not found"
// @Router /units/{id} [put]
func (h *UnitHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	var req dto.UpdateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get existing unit
	existingUnit, err := h.unitService.GetByID(c.Request.Context(), id)
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
	if tenancy == nil || existingUnit.OrganizationID != tenancy.OrganizationID || existingUnit.BranchID != tenancy.BranchID {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	// Update fields
	if req.UnitNumber != "" {
		existingUnit.UnitNumber = req.UnitNumber
	}
	if req.Description != "" {
		existingUnit.Description = req.Description
	}

	// Update unit in the service
	if err := h.unitService.Update(c.Request.Context(), existingUnit); err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	// Respond with updated unit
	response.SendSuccess(c, http.StatusOK, dto.FromUnit(existingUnit))
}

// Delete handles DELETE /unit/:id
// @Summary Delete unit
// @Description Delete a unit by its ID
// @Tags unit
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param id path string true "Unit ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Forbidden access"
// @Failure 404 {object} map[string]string "Unit not found"
// @Router /units/{id} [delete]
func (h *UnitHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get existing unit
	existingUnit, err := h.unitService.GetByID(c.Request.Context(), id)
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
	if tenancy == nil || existingUnit.OrganizationID != tenancy.OrganizationID || existingUnit.BranchID != tenancy.BranchID {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	// Perform deletion
	err = h.unitService.Delete(c.Request.Context(), id)
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
