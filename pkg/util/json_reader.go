package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ReadJSON(ctx *gin.Context, dst interface{}) error {
	validate = validator.New()
	maxBytes := 1_048_576
	ctx.Request.Body = http.MaxBytesReader(nil, ctx.Request.Body, int64(maxBytes))

	dec := json.NewDecoder(ctx.Request.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("malformed JSON at position %d", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("malformed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("invalid type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("invalid type at position %d", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("request body is empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("unknown field %s", fieldName)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("request body exceeds %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("request body must contain single JSON value")
	}

	if err := validate.Struct(dst); err != nil {
		return err
	}

	return nil
}
