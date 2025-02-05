package util

import (
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/thaian1234/green_light/internal/core/domain"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v, _ := binding.Validator.Engine().(*validator.Validate)
	return &Validator{
		validate: v,
	}
}

func (v *Validator) SetupValidator() {
	v.validate.RegisterValidation("year_range", validateYear)
	v.validate.RegisterValidation("page", validatePage)
	v.validate.RegisterValidation("size", validateSize)
	v.validate.RegisterValidation("sort", validateSort)
}

func validateYear(fl validator.FieldLevel) bool {
	year := fl.Field().Int()
	return year >= 1888 && year <= int64(time.Now().Year())
}

func validatePage(fl validator.FieldLevel) bool {
	page := fl.Field().Int()
	return page >= 1 && page <= 100
}

func validateSize(fl validator.FieldLevel) bool {
	size := fl.Field().Int()
	return size >= 1 && size <= 10_000_000
}

func validateSort(fl validator.FieldLevel) bool {
	sort := fl.Field().String()
	parent := fl.Parent().Interface()

	if filter, ok := parent.(domain.Filter); ok {
		for _, safeValue := range filter.SortSafeList {
			if sort == safeValue {
				return true
			}
		}
	}
	return false
}
