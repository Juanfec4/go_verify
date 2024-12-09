# GoVerify

A minimal, flexible, and extensible tag-based validation and transformation library for Go with zero dependencies.

[![Go Reference](https://pkg.go.dev/badge/github.com/Juanfec4/go_verify.svg)](https://pkg.go.dev/github.com/Juanfec4/go_verify)
[![Go Report Card](https://goreportcard.com/badge/github.com/Juanfec4/go_verify)](https://goreportcard.com/report/github.com/Juanfec4/go_verify)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Features

- Declarative validation using struct tags
- Built-in validation rules and transformations
- Easy-to-use extension system
- Detailed error reporting with JSON support
- Zero external dependencies

## Installation

```bash
go get github.com/Juanfec4/go_verify
```

## Quick Start

```go
type User struct {
    Username string   `validator:"required min=3 max=20" transform:"trim"`
    Email    string   `validator:"required email" transform:"lowercase"`
    Age      int      `validator:"required min_value=18"`
}

func main() {
    user := &User{
        Username: "  john_doe  ",
        Email:    "JOHN@EXAMPLE.COM",
        Age:      15,
    }

    // Transform fields
    if err := goverify.Transform(user); err != nil {
        log.Fatal(err)
    }

    // Validate
    if valid, err := goverify.Validate(user); !valid {
        log.Fatal(err)
    }
}
```

## Built-in Rules

### String Validation

- `required`: Non-empty value
- `min=N`: Minimum length
- `max=N`: Maximum length
- `alphanum`: Letters, numbers, underscores
- `alpha`: Letters only
- `email`: Valid email format
- `url`: Valid URL
- `pattern=regex`: Custom pattern

### Numeric Validation

- `min_value=N`: Minimum value
- `max_value=N`: Maximum value

### Network & Date

- `ipv4`: Valid IPv4 address
- `iso_date`: YYYY-MM-DD format
- `time`: HH:MM:SS format

## Built-in Transformers

- `trim`: Remove whitespace
- `lowercase`: Convert to lowercase
- `uppercase`: Convert to uppercase
- `remove_whitespace`: Remove all whitespace

## Extending

### Custom Validation Rule

```go
goverify.AddRule("uuid", func(v reflect.Value, field reflect.StructField) []string {
    if v.Kind() != reflect.String {
        return nil
    }
    if !isValidUUID(v.String()) {
        return []string{"must be a valid UUID"}
    }
    return nil
})

type Resource struct {
    ID string `validator:"required uuid"`
}
```

### Custom Transformer

```go
goverify.AddTransformer("slugify", func(v reflect.Value) error {
    if v.Kind() != reflect.String {
        return nil
    }
    str := v.String()
    str = strings.ToLower(str)
    str = strings.ReplaceAll(str, " ", "-")
    str = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(str, "")
    v.SetString(str)
    return nil
})

type Post struct {
    Title string `transform:"trim slugify"`
}
```

## Error Handling

```go
if valid, err := goverify.Validate(user); !valid {
    // Get structured error
    fmt.Printf("Error: %v\n", err)

    // Get JSON error
    jsonBytes := goverify.ToJSONErr(err)
    fmt.Printf("JSON: %s\n", jsonBytes)
}
```

Example error output:

```json
{
  "message": "validation failed",
  "fields": {
    "username": ["length must be at least 3"],
    "email": ["invalid email format"],
    "age": ["must be at least 18"]
  }
}
```

## Best Practices

- Add custom rules/transformers during initialization
- Keep transformers and validators thread-safe
- Handle errors appropriately
- Document custom rules and parameters

## Contributing

Pull requests are welcome. For major changes, please open an issue first.

## License

[MIT](https://choosealicense.com/licenses/mit/)
