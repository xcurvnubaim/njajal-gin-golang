package CustomValidator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Param   string
	Message string
}

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		return fmt.Sprintf("Minimum length is %s", fe.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Value must be one of %s", fe.Param())
	}
	return fe.Error() // default error
}

// FormatValidationErrors takes an error and returns a formatted error message.
func FormatValidationErrors(err error) string {
	var errMsg []string
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		for _, fe := range ve {
			errMsg = append(errMsg, fmt.Sprintf("%s: %s", fe.Field(), MsgForTag(fe)))
		}
		return strings.Join(errMsg, ", ") // Return the formatted error messages as a single string
	} else {
		return err.Error()
	}
}
