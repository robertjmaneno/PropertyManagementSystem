package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/projectname/internal/repository"
)

// TenancyContext holds tenant information
type TenancyContext struct {
	OrganizationID uint
	BranchID       string
}

// TenancyKey is the key used to store tenancy information in the context
const TenancyKey = "tenancy"

// TenancyMiddleware extracts tenant information from request headers and adds it to the context
func TenancyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant information from headers
		orgIDStr := c.GetHeader("Organization-ID")
		branchID := c.GetHeader("Branch-ID")

		if orgIDStr == "" || branchID == "" {
			c.Next()
			return
		}

		// Convert organization ID to uint
		orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
		if err != nil {
			c.Next()
			return
		}

		tenancy := &TenancyContext{
			OrganizationID: uint(orgID),
			BranchID:       branchID,
		}

		// Add tenant filters to the request context
		opts := repository.ListOptions{
			TenantFilters: map[string]interface{}{
				"organization_id": tenancy.OrganizationID,
				"branch_id":       tenancy.BranchID,
			},
		}

		// Store tenancy information in context
		c.Set(TenancyKey, tenancy)
		c.Set("list_options", opts)

		c.Next()
	}
}
