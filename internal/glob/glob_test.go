package glob

import "testing"

func TestContainsGlobChars(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "no glob characters",
			path:     "./resources/configurations/config.yaml",
			expected: false,
		},
		{
			name:     "wildcard glob",
			path:     "./resources/configurations/*.yaml",
			expected: true,
		},
		{
			name:     "question mark glob",
			path:     "./resources/configurations/config?.yaml",
			expected: true,
		},
		{
			name:     "bracket glob",
			path:     "./resources/configurations/config[0-9].yaml",
			expected: true,
		},
		{
			name:     "multiple glob characters",
			path:     "./resources/configurations/*.y?ml",
			expected: true,
		},
		{
			name:     "glob at start",
			path:     "*.yaml",
			expected: true,
		},
		{
			name:     "glob at end",
			path:     "./resources/configurations/*",
			expected: true,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
		{
			name:     "path with literal asterisk",
			path:     "./resources/configurations/config*.yaml",
			expected: true,
		},
		{
			name:     "path with literal question mark",
			path:     "./resources/configurations/config?.yaml",
			expected: true,
		},
		{
			name:     "path with literal bracket",
			path:     "./resources/configurations/config[test].yaml",
			expected: true,
		},
		{
			name:     "complex glob pattern",
			path:     "./resources/configurations/config-*.{yaml,yml}",
			expected: true,
		},
		{
			name:     "directory with glob",
			path:     "./resources/*/config.yaml",
			expected: true,
		},
		{
			name:     "simple file path",
			path:     "config.yaml",
			expected: false,
		},
		{
			name:     "relative path without glob",
			path:     "../config.yaml",
			expected: false,
		},
		{
			name:     "absolute path without glob",
			path:     "/home/user/config.yaml",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsGlobChars(tt.path)
			if result != tt.expected {
				t.Errorf("ContainsGlobChars(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}
