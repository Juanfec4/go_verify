package goverify

import "reflect"

type (
	// ValidationRule is a function type that validates a field value against specific rules.
	// It takes a reflect.Value and reflect.StructField as input and returns a slice of validation error messages.
	// If the validation passes, it returns an empty slice.
	ValidationRule func(v reflect.Value, field reflect.StructField) []string

	// TransformFunc is a function type that transforms a field value.
	// It takes a reflect.Value as input and returns an error if the transformation fails.
	// If the transformation succeeds, it returns nil.
	TransformFunc func(reflect.Value) error

	// Err represents a validation or transformation error.
	// It contains a message and a map of field-specific error messages.
	Err struct {
		Msg    string              `json:"message"`
		Fields map[string][]string `json:"fields,omitempty"`
	}

	validator struct {
		rules map[string]ValidationRule
	}
)
