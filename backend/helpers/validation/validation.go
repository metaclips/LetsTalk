package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidateInputs validate inputs
func ValidateInputs(dataSet interface{}) error {
	validate := validator.New()

	err := validate.Struct(dataSet)

	if err != nil {
		// Validation syntax is invalid
		if _, ok := err.(*validator.InvalidValidationError); ok {

			return errors.New("validation syntax is invalid")
		}

		for _, err := range err.(validator.ValidationErrors) {

			reflected := reflect.ValueOf(dataSet)

			// Attempt to find field by name and get json tag name
			field, _ := reflected.Type().FieldByName(err.StructField())

			// If json tag doesn't exist, use lower case of name
			name := field.Tag.Get("json")
			if name == "" {
				name = strings.ToLower(err.StructField())
			}

			switch err.Tag() {
			case "required":
				return fmt.Errorf("%s is required", name)
			case "email":
				return fmt.Errorf("%s should be a valid email", name)
			case "min":
				return fmt.Errorf("provide valid 13 digit number")
			case "max":
				return fmt.Errorf("provide valid 13 digit number")
			case "gt":
				return fmt.Errorf("more than %s %s is required", err.Param(), name)
			case "len":
				return fmt.Errorf("ensure to provide valid phone number")
			default:
				return fmt.Errorf("%s is invalid", name)
			}
		}
	}
	return nil
}
