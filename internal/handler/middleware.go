package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// AppError represents a custom application error
type AppError struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	ErrorMessage error  `json:"error"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new custom application error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:         code,
		Message:      message,
		ErrorMessage: err,
	}
}

func ErrorMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Check if it's our custom AppError
			var appErr *AppError
			if errors.As(err, &appErr) {
				logger.Error("error",
					slog.Int("status", appErr.Code),
					slog.String("error", appErr.ErrorMessage.Error()),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
				)
				c.JSON(appErr.Code, ErrorResponse{
					Code:  appErr.Code,
					Error: appErr.ErrorMessage.Error(),
				})
				c.Abort()
				return
			}

			// Handle generic errors
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "internal server error",
				Code:  http.StatusInternalServerError,
			})
			c.Abort()
		}
	}
}
