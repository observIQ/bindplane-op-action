package model

type Parameter struct {
	// Name is the name of the parameter
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Value could be any of the following: string, bool, int, enum (string), float, []string, map
	Value interface{} `json:"value" yaml:"value" mapstructure:"value"`

	// Sensitive will be true if the value is sensitive and should be masked when printed.
	Sensitive bool `json:"sensitive,omitempty" yaml:"sensitive,omitempty" mapstructure:"sensitive"`
}
