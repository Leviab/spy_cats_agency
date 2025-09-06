package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a structured error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorMiddleware creates a new error handling middleware.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// For now, we'll just take the last error.
			// In a real app, you might want to log all of them.
			err := c.Errors.Last()
			// Check the error type and set the status code accordingly.
			// This is a simple example; you could have custom error types.
			status := http.StatusInternalServerError
			msg := err.Error()

			if strings.Contains(msg, "not found") {
				status = http.StatusNotFound
			} else if strings.Contains(msg, "invalid") || strings.Contains(msg, "must have between") || strings.Contains(msg, "cannot have more than") {
				status = http.StatusBadRequest
			} else if strings.Contains(msg, "cannot delete") || strings.Contains(msg, "not available") {
				status = http.StatusConflict
			}

			c.JSON(status, ErrorResponse{Error: msg})
		}
	}
}
