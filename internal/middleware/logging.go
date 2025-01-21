package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/projectname/pkg/logger"
)

// LoggingMiddleware creates a middleware for request logging
func LoggingMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		if query != "" {
			path = path + "?" + query
		}

		switch {
		case status >= 500:
			log.Errorw("Server error",
				"status", status,
				"method", method,
				"path", path,
				"latency", latency,
				"client_ip", c.ClientIP(),
			)
		case status >= 400:
			log.Warnw("Client error",
				"status", status,
				"method", method,
				"path", path,
				"latency", latency,
				"client_ip", c.ClientIP(),
			)
		default:
			log.Infow("Request completed",
				"status", status,
				"method", method,
				"path", path,
				"latency", latency,
				"client_ip", c.ClientIP(),
			)
		}
	}
} 