package goverify

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Error implements the error interface for Err.
// It returns a formatted error message including all field-specific errors.
func (e *Err) Error() string {
	if len(e.Fields) == 0 {
		return e.Msg
	}

	var fieldErrors []string
	for field, msgs := range e.Fields {
		fieldErrors = append(fieldErrors, fmt.Sprintf("%s %s", field, strings.Join(msgs, ", ")))
	}

	return fmt.Sprintf("%s - %s", e.Msg, strings.Join(fieldErrors, "; "))
}

// NewErr creates a new validation or transformation error.
// It takes a message and an optional map of field-specific errors.
//
// Example:
//
//	err := NewErr("validation failed", map[string][]string{
//	    "username": {"must be at least 3 characters long"},
//	    "email":    {"invalid email format"},
//	})
func NewErr(msg string, fields map[string][]string) error {
	return &Err{
		Msg:    msg,
		Fields: fields,
	}
}

// ToJSONErr converts a validation or transformation error to JSON format.
// Returns an empty byte slice if the error is nil or not of type *Err.
//
// Example:
//
//	_, err := Validate(user)
//	if err != nil {
//	    jsonBytes := ToJSONErr(err)
//	    fmt.Printf("JSON error: %s\n", string(jsonBytes))
//	}
func ToJSONErr(e error) []byte {
	if e == nil {
		return []byte{}
	}
	err, ok := e.(*Err)
	if !ok {
		return []byte{}
	}
	if json, e := json.Marshal(&err); e == nil {
		return json
	}
	return []byte{}
}
