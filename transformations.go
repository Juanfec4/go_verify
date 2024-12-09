package goverify

import (
	"reflect"
	"strings"
)

func addStringTransformers() {
	// Trim spaces
	AddTransformer("trim", func(v reflect.Value) error {
		if v.Kind() != reflect.String {
			return nil
		}
		v.SetString(strings.TrimSpace(v.String()))
		return nil
	})

	// Convert to lowercase
	AddTransformer("lowercase", func(v reflect.Value) error {
		if v.Kind() != reflect.String {
			return nil
		}
		v.SetString(strings.ToLower(v.String()))
		return nil
	})

	// Convert to uppercase
	AddTransformer("uppercase", func(v reflect.Value) error {
		if v.Kind() != reflect.String {
			return nil
		}
		v.SetString(strings.ToUpper(v.String()))
		return nil
	})

	// Remove all whitespace
	AddTransformer("remove_whitespace", func(v reflect.Value) error {
		if v.Kind() != reflect.String {
			return nil
		}
		s := v.String()
		s = strings.ReplaceAll(s, " ", "")
		s = strings.ReplaceAll(s, "\t", "")
		s = strings.ReplaceAll(s, "\n", "")
		s = strings.ReplaceAll(s, "\r", "")
		v.SetString(s)
		return nil
	})
}
