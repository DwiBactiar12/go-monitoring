package utils

import "github.com/go-playground/validator/v10"

type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationErrors converts validator.ValidationErrors to a slice of ValidationErrorResponse
func FormatValidationErrors(err error) []ValidationErrorResponse {
	var errors []ValidationErrorResponse
	if err == nil {
		return errors
	}
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Not a validation error, return generic message
		errors = append(errors, ValidationErrorResponse{
			Field:   "",
			Message: err.Error(),
		})
		return errors
	}
	for _, ve := range validationErrors {
		errors = append(errors, ValidationErrorResponse{
			Field:   ve.Field(),
			Message: ve.Field() + " failed on the " + ve.Tag() + " validation",
		})
	}
	return errors
}
