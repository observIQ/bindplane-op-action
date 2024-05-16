package model

type ResourceConfiguration struct {
	// ID will be generated and is used to uniquely identify the resource
	ID string `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id"`

	// Name must be specified if this is a reference to another resource by name
	Name string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name"`

	// DisplayName is a friendly name of the resource that will be displayed in the UI
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty" mapstructure:"displayName"`

	// ParameterizedSpec contains the definition of an embedded resource if this is not a reference to another resource
	ParameterizedSpec `yaml:",inline" mapstructure:",squash"`
}

type ParameterizedSpec struct {
	Type       string      `yaml:"type,omitempty" json:"type,omitempty" mapstructure:"type"`
	Parameters []Parameter `yaml:"parameters,omitempty" json:"parameters,omitempty" mapstructure:"parameters"`

	Processors []ResourceConfiguration `yaml:"processors,omitempty" json:"processors,omitempty" mapstructure:"processors"`
	Disabled   bool                    `yaml:"disabled,omitempty" json:"disabled,omitempty" mapstructure:"disabled"`
}

type AnyResource struct {
	ResourceMeta `yaml:",inline" mapstructure:",squash"`
	Spec         map[string]any `yaml:"spec" json:"spec" mapstructure:"spec"`
}

type AnyResourceStatus struct {
	Resource AnyResource  `json:"resource" mapstructure:"resource"`
	Status   UpdateStatus `json:"status" mapstructure:"status"`
	Reason   string       `json:"reason" mapstructure:"reason"`
}

type ApplyResponseClientSide struct {
	Updates []*AnyResourceStatus `json:"updates"`
}

type ApplyPayload struct {
	Resources []*AnyResource `json:"resources"`
}

type ResourceMeta struct {
	APIVersion string   `yaml:"apiVersion,omitempty" json:"apiVersion"`
	Kind       string   `yaml:"kind,omitempty" json:"kind"`
	Metadata   Metadata `yaml:"metadata,omitempty" json:"metadata"`
}
