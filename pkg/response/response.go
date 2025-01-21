package response

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/projectname/pkg/errors"
)

// Response represents a standard API response
// @Description Standard API response format
type Response struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents an error response
// @Description Error details
type Error struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}

// SendSuccess sends a success response
func SendSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

// SendError sends an error response
func SendError(c *gin.Context, statusCode int, err error) {
	var appErr *errors.AppError
	if e, ok := err.(*errors.AppError); ok {
		appErr = e
	} else {
		appErr = errors.WrapError(err)
	}

	c.JSON(statusCode, Response{
		Success: false,
		Error: &Error{
			Code:    statusCode,
			Message: appErr.Message,
		},
	})
}
