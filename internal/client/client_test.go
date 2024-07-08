package client

import (
	"testing"
	"time"

	"github.com/observiq/bindplane-op-action/internal/client/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWithTimeout(t *testing.T) {
	cases := []struct {
		name     string
		timeout  time.Duration
		expected time.Duration
		err      error
	}{
		{
			name:     "nop",
			timeout:  0,
			expected: DefaultTimeout,
			err:      nil,
		},
		{
			name:     "custom",
			timeout:  time.Second * 30,
			expected: time.Second * 30,
			err:      nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			logger := zap.NewNop()
			config := &config.Config{}

			b, err := NewBindPlane(config, logger, WithTimeout(tc.timeout))
			if tc.err != nil {
				require.Error(t, err)
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			actual := b.client.GetClient().Timeout
			require.Equal(t, tc.expected, actual)
		})
	}
}
