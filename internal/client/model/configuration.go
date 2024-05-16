package model

type ConfigurationsResponse struct {
	Configurations []*Configuration `json:"configurations"`
}

type Configuration struct {
	ResourceMeta                    `yaml:",inline" mapstructure:",squash"`
	Spec                            ConfigurationSpec `json:"spec" yaml:"spec" mapstructure:"spec"`
	StatusType[ConfigurationStatus] `yaml:",inline" mapstructure:",squash"`
}

type ConfigurationResponse struct {
	Configuration *Configuration `json:"configuration"`
	Raw           string         `json:"raw"`
}

type StatusType[T any] struct {
	Status T `yaml:"status,omitempty" json:"status,omitempty" mapstructure:"status,omitempty"`
}

type ConfigurationStatus struct {
	Rollout        Rollout `json:"rollout,omitempty" yaml:"rollout,omitempty" mapstructure:"rollout"`
	CurrentVersion int     `json:"currentVersion,omitempty" yaml:"currentVersion,omitempty" mapstructure:"currentVersion"`
	PendingVersion int     `json:"pendingVersion,omitempty" yaml:"pendingVersion,omitempty" mapstructure:"pendingVersion"`
	Latest         bool    `json:"latest,omitempty" yaml:"latest,omitempty" mapstructure:"latest"`
	Pending        bool    `json:"pending,omitempty" yaml:"pending,omitempty" mapstructure:"pending"`
	Current        bool    `json:"current,omitempty" yaml:"current,omitempty" mapstructure:"current"`
}

type ConfigurationSpec struct {
	ContentType string `json:"contentType" yaml:"contentType" mapstructure:"contentType"`
	// NOTE: MeasurementInterval is deprecated and will be ignored.
	MeasurementInterval string                  `json:"measurementInterval" yaml:"measurementInterval" mapstructure:"measurementInterval"`
	Raw                 string                  `json:"raw,omitempty" yaml:"raw,omitempty" mapstructure:"raw"`
	Sources             []ResourceConfiguration `json:"sources,omitempty" yaml:"sources,omitempty" mapstructure:"sources"`
	Destinations        []ResourceConfiguration `json:"destinations,omitempty" yaml:"destinations,omitempty" mapstructure:"destinations"`
	Extensions          []ResourceConfiguration `json:"extensions,omitempty" yaml:"extensions,omitempty" mapstructure:"extensions"`
	Selector            AgentSelector           `json:"selector" yaml:"selector" mapstructure:"selector"`
	Rollout             ResourceConfiguration   `json:"rollout,omitempty" yaml:"rollout,omitempty" mapstructure:"rollout"`
}
