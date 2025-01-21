package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/projectname/internal/client"
	"github.com/yourusername/projectname/internal/domain"
	pkgerrors "github.com/yourusername/projectname/pkg/errors"
	"github.com/yourusername/projectname/pkg/logger"
)

// PermittedPaths contains paths that don't require authentication
var PermittedPaths = []string{
	"/swagger",
	"/api-docs",
	"/health",
	"/metrics",
	"/error",
	"/webjars",
}

// AuthMiddleware validates the authentication token and sets user information in the context
func AuthMiddleware(authClient client.AuthClient, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "MISSING_TOKEN: No authentication token was provided"})
			c.Abort()
			return
		}

		// Extract token, handling both "Bearer" prefix and raw token cases
		var token string
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// If no Bearer prefix, use the entire header as the token
			token = authHeader
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "MISSING_TOKEN: No authentication token was provided"})
			c.Abort()
			return
		}

		// Validate token with auth service
		userInfo, err := authClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			logger.Error("Failed to validate token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "INVALID_TOKEN: The provided authentication token is invalid"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user", userInfo)
		c.Next()
	}
}

// GetAuthenticatedUser retrieves the authenticated user from the context
func GetAuthenticatedUser(c *gin.Context) *domain.AuthenticatedUser {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user.(*domain.AuthenticatedUser)
}

// RequireRoles creates middleware to check if user has required roles
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetAuthenticatedUser(c)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": pkgerrors.ErrUnauthorized.Error(),
			})
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, requiredRole := range roles {
			for _, userRole := range user.Roles {
				if userRole.Name == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": pkgerrors.ErrInsufficientPermissions.Error(),
			})
			return
		}

		c.Next()
	}
} 