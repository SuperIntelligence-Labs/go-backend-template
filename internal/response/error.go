package response

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/SuperIntelligence-Labs/go-backend-template/internal/config"
)

type errorResponse struct {
	Success    bool        `json:"success"`
	Timestamp  string      `json:"timestamp"`
	Message    string      `json:"message"`
	RequestID  string      `json:"request_id"`
	ErrorCode  string      `json:"error_code"`
	Details    interface{} `json:"details,omitempty"`
	DebugStack string      `json:"debug_stack,omitempty"`
}

type AppError struct {
	StatusCode int
	Message    string
	Code       string
	Details    interface{}
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func ErrorHandler(err error, c echo.Context) {
	if err == nil || c.Response().Committed {
		return
	}

	var appErr *AppError
	var httpErr *echo.HTTPError

	if errors.As(err, &appErr) && appErr != nil {
		if appErr.StatusCode == 0 {
			appErr.StatusCode = http.StatusInternalServerError
		}
		if appErr.Code == "" {
			appErr.Code = "ERR_UNKNOWN"
		}

	} else if errors.As(err, &httpErr) && httpErr != nil {
		msg := "Unknown HTTP error"

		switch m := httpErr.Message.(type) {
		case string:
			msg = m
		case error:
			msg = m.Error()
		default:
			msg = fmt.Sprint(httpErr.Message)
		}

		appErr = &AppError{
			StatusCode: httpErr.Code,
			Message:    msg,
			Code:       "ERR_ECHO",
			Err:        httpErr,
		}

	} else {
		appErr = ErrInternalError(err)
	}

	resp := errorResponse{
		Success:   false,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   appErr.Message,
		ErrorCode: appErr.Code,
		RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
		Details:   appErr.Details,
	}

	if config.IsDev() && appErr.StatusCode >= 500 && appErr.Err != nil {
		resp.DebugStack = string(debug.Stack())
	}

	_ = c.JSON(appErr.StatusCode, resp)
}

func NewAppError(status int, code, message string, details interface{}, err error) *AppError {
	return &AppError{
		StatusCode: status,
		Message:    message,
		Code:       code,
		Details:    details,
		Err:        err,
	}
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ErrBadRequest(message string, details interface{}) *AppError {
	return NewAppError(http.StatusBadRequest, "ERR_BAD_REQUEST", message, details, nil)
}

func ErrValidationFailed(details []ValidationError) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, "ERR_VALIDATION", "Validation failed", details, nil)
}

func ErrUnauthorized(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", message, nil, nil)
}

func ErrForbidden(message string) *AppError {
	return NewAppError(http.StatusForbidden, "ERR_FORBIDDEN", message, nil, nil)
}

func ErrNotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, "ERR_NOT_FOUND", message, nil, nil)
}

func ErrConflict(message string) *AppError {
	return NewAppError(http.StatusConflict, "ERR_CONFLICT", message, nil, nil)
}

func ErrTooManyRequests(message string) *AppError {
	return NewAppError(http.StatusTooManyRequests, "ERR_TOO_MANY_REQUESTS", message, nil, nil)
}

func ErrUnsupportedMediaType(message string) *AppError {
	return NewAppError(http.StatusUnsupportedMediaType, "ERR_UNSUPPORTED_MEDIA_TYPE", message, nil, nil)
}

func ErrServiceUnavailable(message string) *AppError {
	return NewAppError(http.StatusServiceUnavailable, "ERR_SERVICE_UNAVAILABLE", message, nil, nil)
}

func ErrInternalError(err error) *AppError {
	return NewAppError(http.StatusInternalServerError, "ERR_INTERNAL", "Something went wrong", nil, err)
}

func ErrInternalErrorMsg(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, "ERR_INTERNAL", message, nil, err)
}
