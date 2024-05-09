package state

import (
	"sync"

	"github.com/observiq/bindplane-op-action/client/model"
)

// state can be used to cache data during the
// lifecycle of the action
type State interface {
	// Configurations returns all configuration names
	ConfigurationNames() []string

	// SetConfiguration inserts a configuration into the state
	SetConfiguration(name string, configuration model.AnyResource)
}

// Memory is a state that stores data in memory
type Memory struct {
	mu sync.RWMutex

	// configurations is a map of configurations
	// The key is the name of the configuration
	// and value is the AnyResource representation
	configurations map[string]model.AnyResource
}

var _ State = &Memory{}

// NewMemory creates a new memory state
func NewMemory() *Memory {
	return &Memory{
		configurations: make(map[string]model.AnyResource),
	}
}

// Configurations returns the configurations map
func (m *Memory) ConfigurationNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.configurations))
	for name := range m.configurations {
		names = append(names, name)
	}
	return names
}

// SetConfigurations sets the configurations for a given name. This will overwrite
// any existing configurations for the given name.
func (m *Memory) SetConfiguration(name string, configuration model.AnyResource) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.configurations[name] = configuration
}
