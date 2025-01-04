package util

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// ParseError parses error messages from the error object and returns a slice of error messages
func ParseError(err error) []string {
	var errMsgs []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Error())
		}
	} else {
		errMsgs = append(errMsgs, err.Error())
	}

	return errMsgs
}
