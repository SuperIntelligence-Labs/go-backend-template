package config

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func ParseValidationErrors(err error) error {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return err
	}

	var errMsgs []string

	for _, e := range verrs {
		field := strings.TrimPrefix(e.StructNamespace(), "Config.")

		parts := strings.Split(field, ".")

		var snakeParts []string
		for _, part := range parts {
			snakeParts = append(snakeParts, toSnakeCase(part))
		}

		envName := strings.ToUpper(strings.Join(snakeParts, "_"))

		reason := getValidationReason(e)
		errMsg := fmt.Sprintf("%s (%s)", envName, reason)
		errMsgs = append(errMsgs, errMsg)
	}

	return fmt.Errorf("configuration validation failed: %s", strings.Join(errMsgs, "; "))
}

func toSnakeCase(s string) string {
	var result strings.Builder
	runes := []rune(s)

	for i := range runes {
		r := runes[i]

		if i > 0 {
			prev := runes[i-1]

			if unicode.IsLower(prev) && unicode.IsUpper(r) {
				result.WriteRune('_')
			}

			if unicode.IsUpper(prev) && unicode.IsUpper(r) {
				if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					result.WriteRune('_')
				}
			}
		}

		result.WriteRune(unicode.ToLower(r))
	}

	return result.String()
}

func getValidationReason(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required but not set"
	case "numeric":
		return "must be a number"
	case "oneof":
		return fmt.Sprintf("must be one of: %s", e.Param())
	case "min":
		return fmt.Sprintf("value must be at least %s", e.Param())
	case "max":
		return fmt.Sprintf("value must be at most %s", e.Param())
	case "email":
		return "must be a valid email address"
	case "url":
		return "must be a valid URL"
	case "len":
		return fmt.Sprintf("must be exactly %s characters", e.Param())
	case "gt":
		return fmt.Sprintf("must be greater than %s", e.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", e.Param())
	case "lt":
		return fmt.Sprintf("must be less than %s", e.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", e.Param())
	default:
		return fmt.Sprintf("failed validation '%s'", e.Tag())
	}
}
