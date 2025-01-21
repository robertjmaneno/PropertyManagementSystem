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

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// getTenancyContext extracts tenancy information from gin context
func (h *UserHandler) getTenancyContext(c *gin.Context) *middleware.TenancyContext {
	tenancy, exists := c.Get(middleware.TenancyKey)
	if !exists {
		return nil
	}
	return tenancy.(*middleware.TenancyContext)
}

// getListOptions gets list options with tenant filters from context
func (h *UserHandler) getListOptions(c *gin.Context) repository.ListOptions {
	opts, exists := c.Get("list_options")
	if !exists {
		return repository.ListOptions{
			TenantFilters: make(map[string]interface{}),
		}
	}
	return opts.(repository.ListOptions)
}

// RegisterRoutes registers user routes
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("", h.Create)
		users.GET("", h.List)
		users.GET("/search", h.Search)
		users.GET("/:id", h.GetByID)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}

// Create handles POST /users
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param request body dto.CreateUserRequest true "User information"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
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

	user := req.ToUser()
	user.OrganizationID = tenancy.OrganizationID
	user.BranchID = tenancy.BranchID

	if err := h.userService.Create(c.Request.Context(), user); err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusCreated, dto.FromUser(user))
}

// List handles GET /users
// @Summary List users
// @Description Get a list of all users with pagination and sorting
// @Tags users
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param page_size query int false "Number of items per page (default: 10)" minimum(1) maximum(100)
// @Param sort query string false "Sort fields (format: field1:asc,field2:desc) - Available fields: first_name, last_name, email, created_at"
// @Success 200 {object} dto.ListUsersResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	var req dto.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get list options with tenant filters
	opts := h.getListOptions(c)

	// Get total count
	total, err := h.userService.Count(c.Request.Context(), opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	// Get users with pagination and sorting
	users, err := h.userService.GetPaged(c.Request.Context(), opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromUsers(users, total))
}

// Search handles GET /users/search
// @Summary Search users
// @Description Search users by name or email
// @Tags users
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param search query string true "Search query"
// @Param page query int false "Page number (default: 1)"
// @Param case_sensitive query bool false "Case sensitive search (default: false)"
// @Success 200 {object} dto.ListUsersResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /users/search [get]
func (h *UserHandler) Search(c *gin.Context) {
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

	// Get users with search, pagination and sorting
	users, err := h.userService.Search(c.Request.Context(), query, opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	// Convert pointer slice to value slice
	userSlice := make([]domain.User, len(users))
	for i, user := range users {
		userSlice[i] = *user
	}

	// Get total count for search
	total, err := h.userService.Count(c.Request.Context(), opts)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromUsers(userSlice, total))
}

// GetByID handles GET /users/:id
// @Summary Get user by ID
// @Description Get a user's information by their ID
// @Tags users
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
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
	if tenancy == nil || user.OrganizationID != tenancy.OrganizationID || user.BranchID != tenancy.BranchID {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromUser(user))
}

// Update handles PUT /users/:id
// @Summary Update user
// @Description Update a user's information
// @Tags users
// @Accept json
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRequest true "User information"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get existing user
	existingUser, err := h.userService.GetByID(c.Request.Context(), id)
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
	if tenancy == nil || existingUser.OrganizationID != tenancy.OrganizationID || existingUser.BranchID != tenancy.BranchID {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	// Update fields
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.FirstName != "" {
		existingUser.FirstName = req.FirstName
	}
	if req.LastName != "" {
		existingUser.LastName = req.LastName
	}

	if err := h.userService.Update(c.Request.Context(), existingUser); err != nil {
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusOK, dto.FromUser(existingUser))
}

// Delete handles DELETE /users/:id
// @Summary Delete user
// @Description Delete a user by their ID
// @Tags users
// @Produce json
// @Param Organization-ID header string true "Organization ID"
// @Param Branch-ID header string true "Branch ID"
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Get existing user
	existingUser, err := h.userService.GetByID(c.Request.Context(), id)
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
	if tenancy == nil || existingUser.OrganizationID != tenancy.OrganizationID || existingUser.BranchID != tenancy.BranchID {
		response.SendError(c, http.StatusForbidden, errors.ErrForbidden)
		return
	}

	err = h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		if err == errors.ErrNotFound {
			response.SendError(c, http.StatusNotFound, err)
			return
		}
		response.SendError(c, http.StatusInternalServerError, err)
		return
	}

	response.SendSuccess(c, http.StatusNoContent, nil)
}
