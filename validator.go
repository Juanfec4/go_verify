// Package goverify provides struct validation and transformation functionality using tags.
// It allows for declarative validation rules and data transformations through struct tags.
package goverify

import (
	"reflect"
	"strings"
)

var v = &validator{
	rules: make(map[string]ValidationRule),
}

func init() {
	addRequiredRule()
	addSizeRules()
	addRangeRules()
	addPatternRules()
	addStringRules()
	addNetworkRules()
	addCustomStringRules()
	addDateTimeRules()
}

// Validate validates a struct according to its field tags.
// It returns true if validation passes, false and an error otherwise.
//
// Example:
//
//	type User struct {
//	    Username string `validator:"required min=3 max=20 alphanum"`
//	    Email    string `validator:"required email"`
//	    Age      int    `validator:"required min_value=18 max_value=150"`
//	}
//
//	user := &User{
//	    Username: "john_doe",
//	    Email:    "john@example.com",
//	    Age:      25,
//	}
//
//	valid, err := Validate(user)
//	if !valid {
//	    log.Printf("Validation failed: %v", err)
//	}
func Validate(dto interface{}) (bool, error) {
	if dto == nil {
		return false, NewErr("invalid payload", nil)
	}

	val := reflect.ValueOf(dto)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return false, NewErr("input must be a struct", nil)
	}

	violations := make(map[string][]string)
	t := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		fieldVal := val.Field(i)

		validateTag := field.Tag.Get("validator")
		if validateTag == "" {
			continue
		}

		rules := strings.Fields(validateTag)
		for _, rule := range rules {
			ruleName := strings.Split(rule, "=")[0]
			if ruleFunc, exists := v.rules[ruleName]; exists {
				if errs := ruleFunc(fieldVal, field); len(errs) > 0 {
					violations[field.Name] = append(violations[field.Name], errs...)
				}
			}
		}
	}

	if len(violations) > 0 {
		return false, NewErr("validation failed", violations)
	}

	return true, nil
}

// AddRule adds a new validation rule that can be referenced in struct tags.
// The key parameter is the name used in validator tags.
//
// Example:
//
//	// Add a custom validation rule
//	AddRule("uuid", func(v reflect.Value, field reflect.StructField) []string {
//	    if v.Kind() != reflect.String {
//	        return nil
//	    }
//	    if !isValidUUID(v.String()) {
//	        return []string{"must be a valid UUID"}
//	    }
//	    return nil
//	})
//
//	type Resource struct {
//	    ID string `validator:"required uuid"`
//	}
func AddRule(key string, rule ValidationRule) {
	v.rules[key] = rule
}

func parseParams(tag string) map[string]string {
	params := make(map[string]string)
	pairs := strings.Split(tag, ",")

	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			params[kv[0]] = kv[1]
		}
	}

	return params
}
