package goconf

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// StructValidator validates a struct's fields against the validation
// rules defined using struct tags. It utilizes the "github.com/go-playground/validator/v10"
// package for validation.
//
// Parameters:
//   - config (interface{}): The struct to be validated. The struct should have validation
//     rules defined using tags such as `validate:"required"`.
//
// Returns:
//   - error: Returns nil if the struct passes validation. If validation fails, it returns
//     an error of type `validator.ValidationErrors` which provides detailed information
//     about validation failures.
//
// Usage Example:
//
//	type Config struct {
//	    Name string `validate:"required"`
//	    Age  int    `validate:"gte=0"`
//	}
//
//	config := Config{
//	    Name: "",
//	    Age: -1,
//	}
//
//	if err := StructValidator(config); err != nil {
//	    // Handle validation errors, e.g., log or return
//	}
//
// Note:
//   - The validator is initialized with the `WithRequiredStructEnabled` option, which ensures
//     that nil struct fields are treated as invalid.
//   - The function will panic if the `config` parameter is not a struct or pointer to a struct.
//
// More validator information https://github.com/go-playground/validator
func StructValidator(config interface{}) error {
	v := validator.New(validator.WithRequiredStructEnabled())
	err := v.Struct(config)

	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return validationErrors
		}

		return errors.Join(err, errors.New("validation failed"))
	}

	return nil
}
