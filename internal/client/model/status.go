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

type AdditionalInfo struct {
	Message       string              `json:"message" yaml:"message" mapstructure:"message"`
	Documentation []DocumentationLink `json:"documentation" yaml:"documentation"`
}

type DocumentationLink struct {
	Text string `json:"text" yaml:"text"`
	URL  string `json:"url" yaml:"url"`
}
