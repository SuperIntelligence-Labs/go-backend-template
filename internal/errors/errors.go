package errors

import "errors"

// Common application errors
var (
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrBadRequest     = errors.New("bad request")
	ErrConflict       = errors.New("resource conflict")
	ErrInternalServer = errors.New("internal server error")
	ErrValidation     = errors.New("validation failed")
)
