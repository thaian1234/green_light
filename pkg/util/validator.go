package util

import (
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/thaian1234/green_light/pkg/logger"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		if err := v.RegisterValidation("year_range", validateYear); err != nil {
			logger.Fatal(err.Error())
		}
	}
	return &Validator{
		validate: v,
	}
}

func (v *Validator) SetupValidator() {
	v.validate.RegisterValidation("year_range", validateYear)
}

func validateYear(fl validator.FieldLevel) bool {
	year := fl.Field().Int()
	return year >= 1888 && year <= int64(time.Now().Year())
}
