package goverify

import (
	"fmt"
	"reflect"
	"strings"
)

var transformers = map[string]TransformFunc{}
var priorityLookup = map[string]int{
	"trim":              1,
	"remove_whitespace": 2,
	"lowercase":         3,
	"uppercase":         4,
}

func init() {
	addStringTransformers()
}

// Transform applies transformations to a struct according to its field tags.
// It returns an error if any transformation fails.
//
// Example:
//
//	type User struct {
//	    Username string `transform:"trim lowercase"`
//	    Email    string `transform:"trim lowercase"`
//	}
//
//	user := &User{
//	    Username: "  JohnDoe  ",
//	    Email:    "  JOHN@EXAMPLE.COM  ",
//	}
//
//	err := Transform(user)
//	if err != nil {
//	    log.Printf("Transform failed: %v", err)
//	}
func Transform(dto interface{}) error {
	if dto == nil {
		return NewErr("invalid payload", nil)
	}

	val := reflect.ValueOf(dto)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return NewErr("input must be a struct", nil)
	}

	return transformStruct(val)
}

// AddTransformer adds a new transformation function that can be referenced in struct tags.
// The name parameter is used in transform tags.
//
// Example:
//
//	// Add a custom transformer
//	AddTransformer("truncate", func(v reflect.Value) error {
//	    if v.Kind() != reflect.String {
//	        return nil
//	    }
//	    str := v.String()
//	    if len(str) > 10 {
//	        v.SetString(str[:10])
//	    }
//	    return nil
//	})
//
//	type Post struct {
//	    Title string `transform:"truncate trim"`
//	}
func AddTransformer(name string, fn TransformFunc) {
	transformers[name] = fn
}

func transformStruct(val reflect.Value) error {
	t := val.Type()
	violations := make(map[string][]string)

	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		fieldVal := val.Field(i)

		// Handle nested structs
		if fieldVal.Kind() == reflect.Struct {
			if err := transformStruct(fieldVal); err != nil {
				if vErr, ok := err.(*Err); ok {
					for k, v := range vErr.Fields {
						violations[field.Name+"."+k] = v
					}
				}
				continue
			}
		}

		// Handle pointers to structs
		if fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() && fieldVal.Elem().Kind() == reflect.Struct {
			if err := transformStruct(fieldVal.Elem()); err != nil {
				if vErr, ok := err.(*Err); ok {
					for k, v := range vErr.Fields {
						violations[field.Name+"."+k] = v
					}
				}
				continue
			}
		}

		// Handle slices of structs
		// TODO maps?
		if fieldVal.Kind() == reflect.Slice {
			for j := 0; j < fieldVal.Len(); j++ {
				elem := fieldVal.Index(j)
				if elem.Kind() == reflect.Struct {
					if err := transformStruct(elem); err != nil {
						if vErr, ok := err.(*Err); ok {
							for k, v := range vErr.Fields {
								violations[fmt.Sprintf("%s[%d].%s", field.Name, j, k)] = v
							}
						}
					}
				}
			}
		}

		// Apply transformations to the field
		if err := applyTransformations(fieldVal, field); err != nil {
			violations[field.Name] = append(violations[field.Name], err.Error())
		}
	}

	if len(violations) > 0 {
		return NewErr("transformation failed", violations)
	}

	return nil
}

func applyTransformations(v reflect.Value, field reflect.StructField) error {
	if !v.CanSet() {
		return nil
	}

	transformTag := field.Tag.Get("transform")
	if transformTag == "" {
		return nil
	}

	transforms := strings.Fields(transformTag)
	orderedTransforms := orderTransforms(transforms)

	for _, t := range orderedTransforms {
		if fn, exists := transformers[t]; exists {
			if err := fn(v); err != nil {
				return err
			}
		}
	}

	return nil
}

// orderTransforms ensures transforms are applied in the correct order
func orderTransforms(transforms []string) []string {
	order := priorityLookup
	// Sort transforms based on predefined order
	ordered := make([]string, len(transforms))
	copy(ordered, transforms)

	// Sort by priority
	for i := 0; i < len(ordered)-1; i++ {
		for j := i + 1; j < len(ordered); j++ {
			priority1 := order[ordered[i]]
			priority2 := order[ordered[j]]
			if priority1 > priority2 {
				ordered[i], ordered[j] = ordered[j], ordered[i]
			}
		}
	}

	return ordered
}
