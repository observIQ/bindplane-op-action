package client

import (
	"testing"
)

func TestBuildBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "URL without trailing slash",
			input:    "https://example.com",
			expected: "https://example.com",
			wantErr:  false,
		},
		{
			name:     "URL with trailing slash",
			input:    "https://example.com/",
			expected: "https://example.com",
			wantErr:  false,
		},
		{
			name:     "URL with multiple trailing slashes",
			input:    "https://example.com///",
			expected: "https://example.com",
			wantErr:  false,
		},
		{
			name:     "URL with path and trailing slash",
			input:    "https://example.com/api/",
			expected: "https://example.com/api",
			wantErr:  false,
		},
		{
			name:     "URL with path without trailing slash",
			input:    "https://example.com/api",
			expected: "https://example.com/api",
			wantErr:  false,
		},
		{
			name:     "URL with query parameters",
			input:    "https://example.com/api/?param=value",
			expected: "https://example.com/api?param=value",
			wantErr:  false,
		},
		{
			name:     "URL with fragment",
			input:    "https://example.com/api/#section",
			expected: "https://example.com/api#section",
			wantErr:  false,
		},
		{
			name:     "URL with port",
			input:    "https://example.com:8080/",
			expected: "https://example.com:8080",
			wantErr:  false,
		},
		{
			name:     "HTTP URL",
			input:    "http://localhost:3000/",
			expected: "http://localhost:3000",
			wantErr:  false,
		},
		{
			name:     "Invalid URL (treated as relative path)",
			input:    "not-a-valid-url",
			expected: "not-a-valid-url",
			wantErr:  false,
		},
		{
			name:     "Empty string (treated as relative path)",
			input:    "",
			expected: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := buildBaseURL(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("buildBaseURL() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("buildBaseURL() unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("buildBaseURL() = %v, want %v", result, tt.expected)
			}

			// Additional check: ensure result doesn't end with trailing slash
			if len(result) > 0 && result[len(result)-1] == '/' {
				t.Errorf("buildBaseURL() result ends with trailing slash: %v", result)
			}
		})
	}
}
