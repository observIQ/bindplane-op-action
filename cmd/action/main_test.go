package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_extractName(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		found    bool
		expected string
	}{
		{
			name:     "no name",
			input:    "progress rollout",
			found:    false,
			expected: "",
		},
		{
			name:     "name",
			input:    "progress rollout test",
			found:    true,
			expected: "test",
		},
		{
			name:     "name with prefix and trailing space",
			input:    "this is a commit message progress rollout test  ",
			found:    true,
			expected: "test",
		},
		{
			name:     "name with prefix and leading space",
			input:    "  progress rollout test",
			found:    true,
			expected: "test",
		},
		{
			name:     "name with prefix and leading and trailing space",
			input:    "  this is a commit message   progress rollout test  ",
			found:    true,
			expected: "test",
		},
		{
			name:     "invalid suffix",
			input:    "progress rollout test test",
			found:    false,
			expected: "",
		},
		{
			name:     "Config name with -",
			input:    "progress rollout test-name",
			found:    true,
			expected: "test-name",
		},
		{
			name:     "Config name with - .",
			input:    "progress rollout test-.name",
			found:    true,
			expected: "test-.name",
		},
		{
			name:     "Config name with .",
			input:    "progress rollout test.name",
			found:    true,
			expected: "test.name",
		},
		{
			name:     "Config name with _",
			input:    "progress rollout test_name",
			found:    true,
			expected: "test_name",
		},
		{
			name:     "Readme example",
			input:    "Trigger rollout for dev: progress rollout dev-config",
			found:    true,
			expected: "dev-config",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			name, found := extractConfigName(tc.input)
			require.Equal(t, tc.found, found)
			require.Equal(t, tc.expected, name)
		})
	}

}
