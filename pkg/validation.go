package util

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(s interface{}) []ErrorResponse {
	var errors []ErrorResponse
	err := Validate.Struct(s)
	
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var message string
			
			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("Field %s wajib diisi", err.Field())
			case "email":
				message = fmt.Sprintf("Format %s tidak valid", err.Field())
			case "min":
				message = fmt.Sprintf("%s minimal mengandung %s karakter", err.Field(), err.Param())
			case "max":
				message = fmt.Sprintf("%s maksimal mengandung %s karakter", err.Field(), err.Param())
			case "oneof":
				message = fmt.Sprintf("%s harus salah satu dari: %s", err.Field(), err.Param())
			default:
				message = fmt.Sprintf("Field %s tidak valid", err.Field())
			}

			element := ErrorResponse{
				Field:   err.Field(),
				Message: message,
			}
			errors = append(errors, element)
		}
	}
	return errors
}