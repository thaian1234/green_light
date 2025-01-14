package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func ParseError(err error) map[string]string {
	errorMessages := make(map[string]string)

	var unmarshalTypeError *json.UnmarshalTypeError
	var syntaxError *json.SyntaxError

	switch {
	case errors.As(err, &syntaxError):
		fieldName := strings.ToLower(unmarshalTypeError.Field)
		errorMessages[fieldName] = fmt.Sprintf("malformed JSON at position %d", syntaxError.Offset)
		return errorMessages
	case errors.As(err, &unmarshalTypeError):
		fieldName := strings.ToLower(unmarshalTypeError.Field)
		errorMessages[fieldName] = fmt.Sprintf("must be a %s value", unmarshalTypeError.Type)
		return errorMessages
	case errors.Is(err, io.EOF):
		errorMessages["error"] = "request body is empty"
		return errorMessages
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		errorMessages[fieldName] = fmt.Sprintf("unknown field %s", fieldName)
		return errorMessages
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errorMessages[field] = fmt.Sprintf("%s is required", field)
			case "year_range":
				errorMessages[field] = fmt.Sprintf("%s must be between 1888 and %v", field, time.Now().Year())
			case "min":
				errorMessages[field] = fmt.Sprintf("%s must have minimum length of %s", field, e.Param())
			case "max":
				errorMessages[field] = fmt.Sprintf("%s must have maximum length of %s", field, e.Param())
			case "unique":
				errorMessages[field] = fmt.Sprintf("%s must contain unique values", field)
			default:
				errorMessages[field] = fmt.Sprintf("invalid value for %s", field)
			}
		}
		return errorMessages
	}

	// Handle non-validation errors
	errorMessages["error"] = err.Error()
	return errorMessages
}
