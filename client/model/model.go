package model

type UpdateStatus string

const (
	// StatusUnchanged indicates that there were no changes to a modified resource because the existing resource is the same
	StatusUnchanged UpdateStatus = "unchanged"

	// StatusConfigured indicates that changes were applied to an existing resource
	StatusConfigured UpdateStatus = "configured"

	// StatusCreated indicates that a new resource was created
	StatusCreated UpdateStatus = "created"

	// StatusDeleted indicates that a resource was deleted, either from the store or the current filtered view of resources
	StatusDeleted UpdateStatus = "deleted"

	// StatusNotFound indicates a resource was attempted to be deleted but it didn't exist.
	StatusNotFound UpdateStatus = "not-found"

	// StatusInvalid represents an attempt to add or update a resource with an invalid resource
	StatusInvalid UpdateStatus = "invalid"

	// StatusError is used when an individual resource cannot be applied because of an error
	StatusError UpdateStatus = "error"

	// StatusInUse is used when attempting to delete a resource that is being referenced by another
	StatusInUse UpdateStatus = "in-use"

	// StatusForbidden is used when attempting to modify or delete a resource without sufficient permission
	StatusForbidden UpdateStatus = "forbidden"

	// StatusDeprecated is used when attempting to seed a resource that is deprecated that doesn't already exist
	StatusDeprecated UpdateStatus = "deprecated"
)

type ResourceMeta struct {
	APIVersion string   `yaml:"apiVersion,omitempty" json:"apiVersion"`
	Kind       string   `yaml:"kind,omitempty" json:"kind"`
	Metadata   Metadata `yaml:"metadata,omitempty" json:"metadata"`
}

type Metadata struct {
	ID          string `yaml:"id,omitempty" json:"id" mapstructure:"id"`
	Name        string `yaml:"name,omitempty" json:"name" mapstructure:"name"`
	DisplayName string `yaml:"displayName,omitempty" json:"displayName,omitempty" mapstructure:"displayName"`
	Hash        string `yaml:"hash,omitempty" json:"hash,omitempty" mapstructure:"hash"`
	Version     int    `yaml:"version,omitempty" json:"version,omitempty" mapstructure:"version"`
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
