package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func ToValidationErrors(err error) []ValidationError {
	var details []ValidationError

	var validationErrors validator.ValidationErrors
	ok := errors.As(err, &validationErrors)
	if !ok {
		return details
	}

	for _, fieldError := range validationErrors {
		details = append(details, ValidationError{
			Field:   fieldError.Field(),
			Message: validationMessage(fieldError),
		})
	}
	return details
}

func validationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + e.Param() + " characters/items"
	case "max":
		return "Must be at most " + e.Param() + " characters/items"
	case "gte":
		return "Must be greater than or equal to " + e.Param()
	case "gt":
		return "Must be greater than " + e.Param()
	case "lte":
		return "Must be less than or equal to " + e.Param()
	case "lt":
		return "Must be less than " + e.Param()
	case "len":
		return "Must be exactly " + e.Param() + " characters long"
	case "alpha":
		return "Must contain only letters"
	case "alphanum":
		return "Must contain only letters and numbers"
	case "numeric":
		return "Must be a valid numeric value"
	case "url":
		return "Must be a valid URL"
	case "uuid":
		return "Must be a valid UUID"
	case "ip":
		return "Must be a valid IP address"
	case "oneof":
		return "Must be one of: " + e.Param()
	}

	return "Invalid value"
}
