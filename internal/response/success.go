package response

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type successResponse[T any] struct {
	Success   bool   `json:"success"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Data      T      `json:"data"`
}

func respondSuccess[T any](c echo.Context, status int, message string, data T) error {
	if c.Response().Committed {
		return nil
	}

	if status == http.StatusNoContent {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(status, successResponse[T]{
		Success:   true,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   message,
		RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
		Data:      data,
	})
}

func OK[T any](c echo.Context, message string, data T) error {
	return respondSuccess(c, http.StatusOK, message, data)
}

func Created[T any](c echo.Context, message string, data T) error {
	return respondSuccess(c, http.StatusCreated, message, data)
}

func Accepted[T any](c echo.Context, message string, data T) error {
	return respondSuccess(c, http.StatusAccepted, message, data)
}

func NoContent(c echo.Context) error {
	return respondSuccess[any](c, http.StatusNoContent, "", nil)
}
