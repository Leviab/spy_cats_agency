package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			status := http.StatusInternalServerError
			msg := "Internal server error"
			code := "INTERNAL_ERROR"

			errMsg := strings.ToLower(err.Error())

			switch {
			case strings.Contains(errMsg, "not found"):
				status = http.StatusNotFound
				msg = "Resource not found"
				code = "NOT_FOUND"

			case strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "unauthenticated"):
				status = http.StatusUnauthorized
				msg = "Authentication required"
				code = "UNAUTHORIZED"

			case strings.Contains(errMsg, "forbidden") || strings.Contains(errMsg, "access denied"):
				status = http.StatusForbidden
				msg = "Access denied"
				code = "FORBIDDEN"

			case strings.Contains(errMsg, "invalid") ||
				strings.Contains(errMsg, "must have between") ||
				strings.Contains(errMsg, "cannot have more than") ||
				strings.Contains(errMsg, "validation failed") ||
				strings.Contains(errMsg, "bad request"):
				status = http.StatusBadRequest
				msg = "Invalid request"
				code = "BAD_REQUEST"

			case strings.Contains(errMsg, "cannot delete") ||
				strings.Contains(errMsg, "not available") ||
				strings.Contains(errMsg, "conflict") ||
				strings.Contains(errMsg, "already exists"):
				status = http.StatusConflict
				msg = "Resource conflict"
				code = "CONFLICT"

			case strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline exceeded"):
				status = http.StatusRequestTimeout
				msg = "Request timeout"
				code = "TIMEOUT"

			case strings.Contains(errMsg, "too many"):
				status = http.StatusTooManyRequests
				msg = "Rate limit exceeded"
				code = "RATE_LIMIT"
			}

			response := ErrorResponse{
				Error: msg,
				Code:  code,
			}

			// In development mode, include original error message
			if gin.Mode() == gin.DebugMode {
				response.Error = err.Error()
			}

			c.JSON(status, response)
			c.Abort() // Prevent further processing
		}
	}
}
