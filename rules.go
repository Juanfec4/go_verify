package goverify

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func addSizeRules() {
	// Min length for strings and slices
	AddRule("min", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		rules := strings.Fields(field.Tag.Get("validator"))
		var minLength int

		// Find the min parameter
		for _, rule := range rules {
			if strings.HasPrefix(rule, "min=") {
				valStr := strings.TrimPrefix(rule, "min=")
				val, err := strconv.Atoi(valStr)
				if err != nil {
					return []string{fmt.Sprintf("invalid min: %s", valStr)}
				}
				minLength = val
				break
			}
		}

		switch v.Kind() {
		case reflect.String:
			if len(v.String()) < minLength {
				errs = append(errs, fmt.Sprintf("length must be at least %d", minLength))
			}
		case reflect.Slice, reflect.Array:
			if v.Len() < minLength {
				errs = append(errs, fmt.Sprintf("must have at least %d items", minLength))
			}
		}
		return errs
	})

	// Max length for strings and slices
	AddRule("max", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		rules := strings.Fields(field.Tag.Get("validator"))
		var maxLength int

		// Find the max parameter
		for _, rule := range rules {
			if strings.HasPrefix(rule, "max=") {
				valStr := strings.TrimPrefix(rule, "max=")
				val, err := strconv.Atoi(valStr)
				if err != nil {
					return []string{fmt.Sprintf("invalid max: %s", valStr)}
				}
				maxLength = val
				break
			}
		}

		switch v.Kind() {
		case reflect.String:
			if len(v.String()) > maxLength {
				errs = append(errs, fmt.Sprintf("length must not exceed %d", maxLength))
			}
		case reflect.Slice, reflect.Array:
			if v.Len() > maxLength {
				errs = append(errs, fmt.Sprintf("must not exceed %d items", maxLength))
			}
		}
		return errs
	})
}

func addRangeRules() {
	// Minimum value for numbers
	AddRule("min_value", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		rules := strings.Fields(field.Tag.Get("validator"))
		var minValue float64

		// Find the min_value parameter
		for _, rule := range rules {
			if strings.HasPrefix(rule, "min_value=") {
				valStr := strings.TrimPrefix(rule, "min_value=")
				val, err := strconv.ParseFloat(valStr, 64)
				if err != nil {
					return []string{fmt.Sprintf("invalid min_value: %s", valStr)}
				}
				minValue = val
				break
			}
		}

		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if float64(v.Int()) < minValue {
				errs = append(errs, fmt.Sprintf("must be at least %v", minValue))
			}
		case reflect.Float32, reflect.Float64:
			if v.Float() < minValue {
				errs = append(errs, fmt.Sprintf("must be at least %v", minValue))
			}
		}
		return errs
	})

	// Maximum value for numbers
	AddRule("max_value", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		rules := strings.Fields(field.Tag.Get("validator"))
		var maxValue float64

		// Find the max_value parameter
		for _, rule := range rules {
			if strings.HasPrefix(rule, "max_value=") {
				valStr := strings.TrimPrefix(rule, "max_value=")
				val, err := strconv.ParseFloat(valStr, 64)
				if err != nil {
					return []string{fmt.Sprintf("invalid max_value: %s", valStr)}
				}
				maxValue = val
				break
			}
		}

		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if float64(v.Int()) > maxValue {
				errs = append(errs, fmt.Sprintf("must not exceed %v", maxValue))
			}
		case reflect.Float32, reflect.Float64:
			if v.Float() > maxValue {
				errs = append(errs, fmt.Sprintf("must not exceed %v", maxValue))
			}
		}
		return errs
	})
}

func addRequiredRule() {
	AddRule("required", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		switch v.Kind() {
		case reflect.String:
			if v.String() == "" {
				errs = append(errs, "field is required")
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.Int() == 0 {
				errs = append(errs, "field is required")
			}
		case reflect.Float32, reflect.Float64:
			if v.Float() == 0 {
				errs = append(errs, "field is required")
			}
		case reflect.Slice, reflect.Array:
			if v.Len() == 0 {
				errs = append(errs, "field is required")
			}
		}
		return errs
	})
}

func addPatternRules() {
	// Email format
	AddRule("email", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		if v.Kind() != reflect.String {
			return errs
		}

		emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		re := regexp.MustCompile(emailPattern)
		if !re.MatchString(v.String()) {
			errs = append(errs, "invalid email format")
		}
		return errs
	})

	// Regex pattern matching
	AddRule("pattern", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		if v.Kind() != reflect.String {
			return errs
		}

		tag := field.Tag.Get("validator")
		params := parseParams(tag)
		pattern, ok := params["pattern"]
		if !ok {
			return errs
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			return errs
		}

		if !re.MatchString(v.String()) {
			errs = append(errs, "invalid format")
		}
		return errs
	})
}

func addStringRules() {
	// Alphanumeric and underscore only
	AddRule("alphanum", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		if v.Kind() != reflect.String {
			return errs
		}

		str := v.String()
		for _, char := range str {
			if !unicode.IsLetter(char) && !unicode.IsNumber(char) && char != '_' {
				errs = append(errs, "must contain only letters, numbers, and underscores")
				break
			}
		}
		return errs
	})

	// Letters only
	AddRule("alpha", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		if v.Kind() != reflect.String {
			return errs
		}

		str := v.String()
		for _, char := range str {
			if !unicode.IsLetter(char) {
				errs = append(errs, "must contain only letters")
				break
			}
		}
		return errs
	})

	// No whitespace
	AddRule("no_whitespace", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		if v.Kind() != reflect.String {
			return errs
		}

		if strings.ContainsAny(v.String(), " \t\n\r") {
			errs = append(errs, "must not contain whitespace")
		}
		return errs
	})
}

func addNetworkRules() {
	// URL validation
	AddRule("url", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		if v.Kind() != reflect.String {
			return errs
		}

		str := v.String()
		if str == "" {
			return errs
		}

		_, err := url.ParseRequestURI(str)
		if err != nil {
			errs = append(errs, "must be a valid URL")
		}
		return errs
	})

	// IPv4 validation
	AddRule("ipv4", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		if v.Kind() != reflect.String {
			return errs
		}

		ipv4Pattern := `^(\d{1,3}\.){3}\d{1,3}$`
		if matched, _ := regexp.MatchString(ipv4Pattern, v.String()); !matched {
			errs = append(errs, "must be a valid IPv4 address")
			return errs
		}

		// Validate each octet
		parts := strings.Split(v.String(), ".")
		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil || num < 0 || num > 255 {
				errs = append(errs, "must be a valid IPv4 address")
				break
			}
		}
		return errs
	})
}

func addCustomStringRules() {
	// Contains specific substring
	AddRule("contains", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		if v.Kind() != reflect.String {
			return errs
		}

		rules := strings.Fields(field.Tag.Get("validator"))
		var substring string

		for _, rule := range rules {
			if strings.HasPrefix(rule, "contains=") {
				substring = strings.TrimPrefix(rule, "contains=")
				break
			}
		}

		if !strings.Contains(v.String(), substring) {
			errs = append(errs, fmt.Sprintf("must contain '%s'", substring))
		}
		return errs
	})

	// Starts with prefix
	AddRule("starts_with", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		if v.Kind() != reflect.String {
			return errs
		}

		rules := strings.Fields(field.Tag.Get("validator"))
		var prefix string

		for _, rule := range rules {
			if strings.HasPrefix(rule, "starts_with=") {
				prefix = strings.TrimPrefix(rule, "starts_with=")
				break
			}
		}

		if !strings.HasPrefix(v.String(), prefix) {
			errs = append(errs, fmt.Sprintf("must start with '%s'", prefix))
		}
		return errs
	})
}

func addDateTimeRules() {
	// ISO8601 date validation
	AddRule("iso_date", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		if v.Kind() != reflect.String {
			return errs
		}

		iso8601Pattern := `^\d{4}-\d{2}-\d{2}$`
		if matched, _ := regexp.MatchString(iso8601Pattern, v.String()); !matched {
			errs = append(errs, "must be a valid ISO8601 date (YYYY-MM-DD)")
		}
		return errs
	})

	// Time format validation (HH:MM:SS)
	AddRule("time", func(v reflect.Value, field reflect.StructField) []string {
		var errs []string
		if v.Kind() != reflect.String {
			return errs
		}

		timePattern := `^([01]\d|2[0-3]):([0-5]\d):([0-5]\d)$`
		if matched, _ := regexp.MatchString(timePattern, v.String()); !matched {
			errs = append(errs, "must be a valid time (HH:MM:SS)")
		}
		return errs
	})
}
