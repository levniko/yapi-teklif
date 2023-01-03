package utils

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type ValidatorErrorResponse struct {
	FailedField string
	Text        string
}

type CustomValidator struct {
	validator *validator.Validate
}

func (v *CustomValidator) Validate(value interface{}) error {
	return v.validator.Struct(value)
}

func GetCustomValidator() *CustomValidator {
	v := validator.New()
	return &CustomValidator{validator: v}
}

func CustomValidatorErr(err error) interface{} {
	var errors []ValidatorErrorResponse
	for _, field := range err.(validator.ValidationErrors) {
		switch field.Tag() {
		case "required":
			errors = append(errors, ValidatorErrorResponse{FailedField: field.Field(), Text: required})
		case "required_without":
			errors = append(errors, ValidatorErrorResponse{FailedField: field.Field(), Text: fmt.Sprintf(requiredwithout, field.Param())})
		case "number":
			errors = append(errors, ValidatorErrorResponse{FailedField: field.Field(), Text: fmt.Sprintf(number, field.Value())})
		case "numeric":
			errors = append(errors, ValidatorErrorResponse{FailedField: field.Field(), Text: fmt.Sprintf(numeric, field.Value())})
		case "alphanum":
			errors = append(errors, ValidatorErrorResponse{FailedField: field.Field(), Text: fmt.Sprintf(alphanum, field.Value())})
		case "email":
			errors = append(errors, ValidatorErrorResponse{FailedField: field.Field(), Text: fmt.Sprintf(email, field.Value())})
		case "max":
			errors = append(errors, ValidatorErrorResponse{FailedField: field.Field(), Text: fmt.Sprintf(max, field.Param())})
		case "min":
			errors = append(errors, ValidatorErrorResponse{FailedField: field.Field(), Text: fmt.Sprintf(max, field.Param())})
		case "eqfield":
			errors = append(errors, ValidatorErrorResponse{FailedField: field.Field(), Text: fmt.Sprintf(eqfield, field.Param(), field.Value())})
		default:
			errors = append(errors, ValidatorErrorResponse{FailedField: field.Field(), Text: defaultErrorMessage})
		}
	}
	return errors
}

func StartEndValidator(fl validator.FieldLevel) bool {
	startEnd := fl.Field().String()

	stageRegexString := "^[0-9]{4} - [1-4]{1}.Çeyrek$"
	stageRegex := regexp.MustCompile(stageRegexString)

	return stageRegex.MatchString(startEnd)
}
