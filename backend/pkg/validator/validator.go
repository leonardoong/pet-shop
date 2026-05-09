package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// indonesianPhone accepts formats: 08xx, +628xx, 628xx
var indonesianPhoneRegex = regexp.MustCompile(`^(\+62|62|0)8[1-9][0-9]{7,10}$`)

func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("indonesian_phone", validateIndonesianPhone)
	}
}

func validateIndonesianPhone(fl validator.FieldLevel) bool {
	return indonesianPhoneRegex.MatchString(fl.Field().String())
}
