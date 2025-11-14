package response

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Success   bool        `json:"success"`
	Error     ErrorDetail `json:"error"`
	Message   string      `json:"message"`
	Timestamp int64       `json:"timestamp"`
}

type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

const (
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
	ErrCodeConflict     = "CONFLICT"
	ErrCodeBadRequest   = "BAD_REQUEST"
	ErrCodeInternal     = "INTERNAL_ERROR"
	ErrCodeRateLimit    = "RATE_LIMIT_EXCEEDED"
)

func BadRequest(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
		Success:   false,
		Message:   "Unable to process your request",
		Error:     ErrorDetail{Code: ErrCodeBadRequest, Message: message},
		Timestamp: time.Now().Unix(),
	})
}

func ValidationError(c *gin.Context, err error) {
	details := make(map[string]string)
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			details[e.Field()] = formatValidationError(e)
		}
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
		Success:   false,
		Message:   "The submitted data is invalid",
		Error:     ErrorDetail{Code: ErrCodeValidation, Message: "Validation failed", Details: details},
		Timestamp: time.Now().Unix(),
	})
}

func NotFound(c *gin.Context, resource string) {
	c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{
		Success:   false,
		Message:   "The requested resource could not be found",
		Error:     ErrorDetail{Code: ErrCodeNotFound, Message: resource + " not found"},
		Timestamp: time.Now().Unix(),
	})
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Authentication required"
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
		Success:   false,
		Message:   "Please log in to continue",
		Error:     ErrorDetail{Code: ErrCodeUnauthorized, Message: message},
		Timestamp: time.Now().Unix(),
	})
}

func Forbidden(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
		Success:   false,
		Message:   "You don't have permission to access this resource",
		Error:     ErrorDetail{Code: ErrCodeForbidden, Message: "Access denied"},
		Timestamp: time.Now().Unix(),
	})
}

func Conflict(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusConflict, ErrorResponse{
		Success:   false,
		Message:   "This action conflicts with existing data",
		Error:     ErrorDetail{Code: ErrCodeConflict, Message: message},
		Timestamp: time.Now().Unix(),
	})
}

func InternalError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
		Success:   false,
		Message:   "Something went wrong on our end",
		Error:     ErrorDetail{Code: ErrCodeInternal, Message: "An internal error occurred"},
		Timestamp: time.Now().Unix(),
	})
}

func RateLimitExceeded(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusTooManyRequests, ErrorResponse{
		Success:   false,
		Message:   "Too many requests, please try again later",
		Error:     ErrorDetail{Code: ErrCodeRateLimit, Message: "Rate limit exceeded"},
		Timestamp: time.Now().Unix(),
	})
}

func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + e.Param() + " characters"
	case "max":
		return "Must be at most " + e.Param() + " characters"
	case "oneof":
		return "Must be one of: " + e.Param()
	case "uuid":
		return "Must be a valid UUID"
	case "indian_phone":
		return "Must be a valid 10-digit Indian phone number"
	case "court_id":
		return "Must be a valid court ID"
	default:
		return "Invalid value"
	}
}
