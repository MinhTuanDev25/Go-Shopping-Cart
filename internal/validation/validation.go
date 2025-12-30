package validation

import (
	"fmt"
	"go-shopping-cart/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validator engine")
	}

	RegisterCustomValidation(v)
	return nil
}

func HandleValidationErrors(err error) gin.H {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		for _, e := range validationError {
			root := strings.Split(e.Namespace(), ".")[0]
			rawPath := strings.TrimPrefix(e.Namespace(), root+".")

			parts := strings.Split(rawPath, ".")

			for i, part := range parts {
				if strings.Contains("part", "[") {
					idx := strings.Index(part, "[")
					base := utils.CamelToSnake(part[:idx])
					index := part[idx:]
					parts[i] = base + index
				} else {
					parts[i] = utils.CamelToSnake(part)
				}
			}

			fieldPath := strings.Join(parts, ".")

			switch e.Tag() {
			case "gt":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than %s", fieldPath, e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("%s must be less than %s", fieldPath, e.Param())
			case "gte":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than or equal to %s", fieldPath, e.Param())
			case "lte":
				errors[fieldPath] = fmt.Sprintf("%s must be less than or equal to %s", fieldPath, e.Param())
			case "uuid":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid UUID", fieldPath)
			case "slug":
				errors[fieldPath] = fmt.Sprintf("%s may only contain lowercase letters, numbers, hyphens or dots", fieldPath)
			case "min":
				errors[fieldPath] = fmt.Sprintf("%s must be at least %s characters long", fieldPath, e.Param())
			case "max":
				errors[fieldPath] = fmt.Sprintf("%s must be at most %s characters long", fieldPath, e.Param())
			case "min_int":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than %s", fieldPath, e.Param())
			case "max_int":
				errors[fieldPath] = fmt.Sprintf("%s must be less than %s", fieldPath, e.Param())
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ", ")
				errors[fieldPath] = fmt.Sprintf("%s must be one of the following values: %s", fieldPath, allowedValues)
			case "required":
				errors[fieldPath] = fmt.Sprintf("%s is required", fieldPath)
			case "search":
				errors[fieldPath] = fmt.Sprintf("%s may only contain letters, numbers, and spaces", fieldPath)
			case "email":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid email address", fieldPath)
			case "datetime":
				errors[fieldPath] = fmt.Sprintf("%s must follow the format YYYY-MM-DD", fieldPath)
			case "email_advanced":
				errors[fieldPath] = fmt.Sprintf("%s is in the blocked email list", fieldPath)
			case "password_strong":
				errors[fieldPath] = fmt.Sprintf(
					"%s must be at least 8 characters and include lowercase, uppercase, number, and special character",
					fieldPath,
				)
			case "file_ext":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ", ")
				errors[fieldPath] = fmt.Sprintf("%s only allows file extensions: %s", fieldPath, allowedValues)
			default:
				errors[fieldPath] = fmt.Sprintf("%s is invalid", fieldPath)
			}
		}

		return gin.H{"error": errors}
	}

	return gin.H{
		"error":  "Invalid request",
		"detail": err.Error(),
	}
}
