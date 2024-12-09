// validator_test.go
package goverify

import (
	"encoding/json"
	"strings"
	"testing"
)

type UserProfile struct {
	Username   string   `validator:"required min=3 max=20 alphanum" transform:"trim"`
	Email      string   `validator:"required email" transform:"lowercase trim"`
	Age        int      `validator:"required min_value=18 max_value=150"`
	Password   string   `validator:"required min=8 no_whitespace"`
	Interests  []string `validator:"required min=1 max=5"`
	JoinDate   string   `validator:"required iso_date"`
	LastActive string   `validator:"required time"`
}

type ServerConfig struct {
	Hostname   string   `validator:"required alphanum min=3 max=50" transform:"trim lowercase"`
	IPAddress  string   `validator:"required ipv4" transform:"trim"`
	APIKey     string   `validator:"required starts_with=sk_ no_whitespace" transform:"trim"`
	Website    string   `validator:"url" transform:"trim lowercase"`
	Region     string   `validator:"required alpha" transform:"trim uppercase"`
	Tags       []string `validator:"min=1 max=10"`
	SearchTerm string   `validator:"contains=server" transform:"trim lowercase"`
}

func TestValidation(t *testing.T) {
	validUser := &UserProfile{
		Username:   "john_doe123",
		Email:      "john@example.com",
		Age:        25,
		Password:   "securePass123",
		Interests:  []string{"coding"},
		JoinDate:   "2024-03-15",
		LastActive: "14:30:00",
	}

	tests := []struct {
		name        string
		input       interface{}
		wantErr     bool
		errContains []string
	}{
		{
			name:    "Valid complete user",
			input:   validUser,
			wantErr: false,
		},
		{
			name: "String too short",
			input: &UserProfile{
				Username:   "jo",
				Email:      "john@example.com",
				Age:        25,
				Password:   "securePass123",
				Interests:  []string{"coding"},
				JoinDate:   "2024-03-15",
				LastActive: "14:30:00",
			},
			wantErr:     true,
			errContains: []string{"length must be at least 3"},
		},
		{
			name: "Invalid email format",
			input: &UserProfile{
				Username:   "john_doe",
				Email:      "invalid-email",
				Age:        25,
				Password:   "securePass123",
				Interests:  []string{"coding"},
				JoinDate:   "2024-03-15",
				LastActive: "14:30:00",
			},
			wantErr:     true,
			errContains: []string{"invalid email format"},
		},
		{
			name: "Age too low",
			input: &UserProfile{
				Username:   "john_doe",
				Email:      "john@example.com",
				Age:        15,
				Password:   "securePass123",
				Interests:  []string{"coding"},
				JoinDate:   "2024-03-15",
				LastActive: "14:30:00",
			},
			wantErr:     true,
			errContains: []string{"must be at least 18"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := Validate(tt.input)

			if tt.wantErr {
				if valid || err == nil {
					t.Errorf("Validate() error = nil, wantErr = true")
					return
				}

				errStr := err.Error()
				for _, want := range tt.errContains {
					if !strings.Contains(errStr, want) {
						t.Errorf("Error message should contain %q, got %q", want, errStr)
					}
				}
			} else {
				if !valid || err != nil {
					t.Errorf("Validate() error = %v, wantErr = false", err)
				}
			}
		})
	}
}

func TestTransformation(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		check   func(interface{}) bool
		wantErr bool
	}{
		{
			name: "Basic string transformations",
			input: &UserProfile{
				Username:   "  john_doe123  ",
				Email:      "  JOHN@EXAMPLE.COM  ",
				Age:        25,
				Password:   "securePass123",
				Interests:  []string{"coding"},
				JoinDate:   "2024-03-15",
				LastActive: "14:30:00",
			},
			check: func(i interface{}) bool {
				u := i.(*UserProfile)
				return u.Username == "john_doe123" &&
					u.Email == "john@example.com"
			},
			wantErr: false,
		},
		{
			name: "Server config transformations",
			input: &ServerConfig{
				Hostname:   "  SERVER001  ",
				IPAddress:  "192.168.1.1",
				APIKey:     "  sk_test123  ",
				Website:    "  HTTPS://EXAMPLE.COM  ",
				Region:     "north",
				Tags:       []string{"prod"},
				SearchTerm: "  MAIN-SERVER  ",
			},
			check: func(i interface{}) bool {
				s := i.(*ServerConfig)
				result := s.Hostname == "server001"
				result = result && s.APIKey == "sk_test123"
				result = result && s.Region == "NORTH"
				result = result && s.SearchTerm == "main-server"
				result = result && s.Website == "https://example.com"
				if !result {
					t.Logf("Actual values: Hostname='%s', APIKey='%s', Region='%s', SearchTerm='%s', Website='%s'",
						s.Hostname, s.APIKey, s.Region, s.SearchTerm, s.Website)
				}
				return result
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Transform(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Transform() error = nil, wantErr = true")
				}
				return
			}

			if err != nil {
				t.Errorf("Transform() error = %v, wantErr = false", err)
				return
			}

			if !tt.check(tt.input) {
				t.Errorf("Transform() didn't transform as expected")
			}
		})
	}
}

func TestErrorSerialization(t *testing.T) {
	user := &UserProfile{
		Username: "jo", // Should fail validation
	}

	_, err := Validate(user)
	if err == nil {
		t.Fatal("Expected validation error")
	}

	jsonBytes := ToJSONErr(err)
	if len(jsonBytes) == 0 {
		t.Fatal("Expected non-empty JSON error")
	}

	var errMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &errMap); err != nil {
		t.Fatalf("Failed to unmarshal error JSON: %v", err)
	}

	if _, ok := errMap["message"]; !ok {
		t.Error("JSON error should contain 'message' field")
	}
	if _, ok := errMap["fields"]; !ok {
		t.Error("JSON error should contain 'fields' field")
	}
}

func BenchmarkValidation(b *testing.B) {
	user := &UserProfile{
		Username:   "john_doe123",
		Email:      "john@example.com",
		Age:        25,
		Password:   "securePass123",
		Interests:  []string{"coding"},
		JoinDate:   "2024-03-15",
		LastActive: "14:30:00",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Validate(user)
	}
}

func BenchmarkTransformation(b *testing.B) {
	user := &UserProfile{
		Username:   "  JOHN_DOE123  ",
		Email:      "  JOHN@EXAMPLE.COM  ",
		Age:        25,
		Password:   "securePass123",
		Interests:  []string{"coding"},
		JoinDate:   "2024-03-15",
		LastActive: "14:30:00",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Transform(user)
	}
}
