package state

import (
	"testing"

	"github.com/observiq/bindplane-op-action/client/model"
	"github.com/stretchr/testify/require"
)

func TestNewMemory(t *testing.T) {
	memory := NewMemory()
	require.NotNil(t, memory)
	require.NotNil(t, memory.configurations)

	c := model.AnyResource{
		ResourceMeta: model.ResourceMeta{
			Metadata: model.Metadata{
				Name: "test",
			},
		},
	}

	memory.SetConfiguration("test", c)

	out := memory.ConfigurationNames()
	require.Len(t, out, 1)
	require.Equal(t, "test", out[0])
}
