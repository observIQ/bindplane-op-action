package model

import (
	"time"

	"k8s.io/apimachinery/pkg/labels"
)

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
	ID              string          `yaml:"id,omitempty" json:"id" mapstructure:"id"`
	Name            string          `yaml:"name,omitempty" json:"name" mapstructure:"name"`
	DisplayName     string          `yaml:"displayName,omitempty" json:"displayName,omitempty" mapstructure:"displayName"`
	Description     string          `yaml:"description,omitempty" json:"description,omitempty" mapstructure:"description"`
	Icon            string          `yaml:"icon,omitempty" json:"icon,omitempty" mapstructure:"icon"`
	Labels          Labels          `yaml:"labels,omitempty" json:"labels" mapstructure:"labels"`
	Hash            string          `yaml:"hash,omitempty" json:"hash,omitempty" mapstructure:"hash"`
	Version         int             `yaml:"version,omitempty" json:"version,omitempty" mapstructure:"version"`
	DateModified    *time.Time      `yaml:"dateModified,omitempty" json:"dateModified,omitempty" mapstructure:"dateModified"`
	Deprecated      bool            `yaml:"deprecated,omitempty" json:"deprecated,omitempty" mapstructure:"deprecated"`
	AdditionalInfo  *AdditionalInfo `yaml:"additionalInfo,omitempty" json:"additionalInfo,omitempty" mapstructure:"additionalInfo"`
	ResourceDocLink string          `yaml:"resourceDocLink,omitempty" json:"resourceDocLink,omitempty" mapstructure:"resourceDocLink"`
	Stability       string          `yaml:"stability,omitempty" json:"stability,omitempty" mapstructure:"stability"`
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

type Labels struct {
	labels.Set `json:"-" yaml:",inline"`
}

type AdditionalInfo struct {
	Message       string              `json:"message" yaml:"message" mapstructure:"message"`
	Documentation []DocumentationLink `json:"documentation" yaml:"documentation"`
}

type DocumentationLink struct {
	Text string `json:"text" yaml:"text"`
	URL  string `json:"url" yaml:"url"`
}

type ConfigurationsResponse struct {
	Configurations []*Configuration `json:"configurations"`
}

// Configuration is a resource that represents a configuration
// NOTE: Modified to only include the configuration name
type Configuration struct {
	// ResourceMeta contains the metadata for this resource
	ResourceMeta `yaml:",inline" mapstructure:",squash"`
	// // Spec contains the spec for the Configuration
	// Spec                            ConfigurationSpec `json:"spec" yaml:"spec" mapstructure:"spec"`
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
	// Rollout contains status for the rollout of this configuration
	Rollout Rollout `json:"rollout,omitempty" yaml:"rollout,omitempty" mapstructure:"rollout"`

	// CurrentVersion is the version of the configuration that has most recently completed a rollout
	CurrentVersion int `json:"currentVersion,omitempty" yaml:"currentVersion,omitempty" mapstructure:"currentVersion"`

	// PendingVersion will be set to the version of a rollout that is in progress. It will be set to 0 when the rollout
	// completes.
	PendingVersion int `json:"pendingVersion,omitempty" yaml:"pendingVersion,omitempty" mapstructure:"pendingVersion"`

	// ----------------------------------------------------------------------
	// transient values set when the configuration is read from the store

	// Latest will be set to true on read if the configuration is the latest version
	Latest bool `json:"latest,omitempty" yaml:"latest,omitempty" mapstructure:"latest"`

	// Pending will be set to true on read if the configuration is the pending version
	Pending bool `json:"pending,omitempty" yaml:"pending,omitempty" mapstructure:"pending"`

	// Current will be set to true on read if the configuration is the current version
	Current bool `json:"current,omitempty" yaml:"current,omitempty" mapstructure:"current"`
}

type RolloutStatus int

type Rollout struct {
	// Name will be set to the Name of the configuration when requested via Configuration.Rollout()
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Status is the status of the rollout
	Status RolloutStatus `json:"status" yaml:"status" mapstructure:"status"`

	// Options are set when the Rollout is created based on the rollout options specified in the configuration
	Options RolloutOptions `json:"options" yaml:"options" mapstructure:"options"`

	// Phase starts at zero and increments until all agents are updated. In each phase, initial*multiplier^phase agents will be updated.
	Phase int `json:"phase" yaml:"phase" mapstructure:"phase"`

	// Progress is the current progress of the rollout
	Progress RolloutProgress `json:"progress" yaml:"progress" mapstructure:"progress"`

	// StartedAt is the time the rollout was started
	StartedAt *time.Time `json:"startedAt" yaml:"startedAt" mapstructure:"startedAt"`

	// Stages are the RolloutProgress (Completed, Errors, Pending, and Waiting agents), Name and Labels for each stage of a progressive Rollout
	Stages []RolloutStage `json:"stages" yaml:"stages" mapstructure:"stages"`

	// Stage is the index of the current stage of the rollout
	Stage int `json:"stage" yaml:"stage" mapstructure:"stage"`
}

type RolloutOptions struct {
	// StartAutomatically determines if this rollout transitions immediately from RolloutStatusPending to
	// RolloutStatusStarted without requiring that it be started manually.
	StartAutomatically bool `json:"startAutomatically" yaml:"startAutomatically" mapstructure:"startAutomatically"`

	// RollbackOnFailure determines if the rollout should be rolled back to the previous configuration if the rollout
	// fails.
	RollbackOnFailure bool `json:"rollbackOnFailure" yaml:"rollbackOnFailure" mapstructure:"rollbackOnFailure"`

	// PhaseAgentCount determines the rate at which agents will be updated during a rollout.
	PhaseAgentCount PhaseAgentCount `json:"phaseAgentCount" yaml:"phaseAgentCount" mapstructure:"phaseAgentCount"`

	// MaxErrors is the maximum number of failed agents before the rollout will be considered an error
	MaxErrors int `json:"maxErrors" yaml:"maxErrors" mapstructure:"maxErrors"`
}

type PhaseAgentCount struct {
	Initial    int     `json:"initial" yaml:"initial" mapstructure:"initial"`
	Multiplier float64 `json:"multiplier" yaml:"multiplier" mapstructure:"multiplier"`
	Maximum    int     `json:"maximum" yaml:"maximum" mapstructure:"maximum"`
}

type RolloutProgress struct {
	// Completed is the number of agents with new version with Connected status
	Completed int `json:"completed" yaml:"completed" mapstructure:"completed"`

	// Errors is the number of agents with new version with Error Status
	Errors int `json:"errors" yaml:"errors" mapstructure:"errors"`

	// Pending is the number of agents that are currently being configured
	Pending int `json:"pending" yaml:"pending" mapstructure:"pending"`

	// Waiting is the number of agents that need to be scheduled for configuration
	Waiting int `json:"waiting" yaml:"waiting" mapstructure:"waiting"`
}

type RolloutStage struct {
	// Name of the stage
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	// Agent labels that will be used to select agents to rollout to this stage
	Labels Labels `json:"labels" yaml:"labels" mapstructure:"labels"`
	// Progress is the current progress of this rollout stage
	Progress RolloutProgress `json:"progress" yaml:"progress" mapstructure:"progress"`
}

const (
	// RolloutStatusPending is created, manual start required
	RolloutStatusPending RolloutStatus = 0

	// RolloutStatusStarted is in progress
	RolloutStatusStarted RolloutStatus = 1

	// RolloutStatusPaused is paused by the user
	RolloutStatusPaused RolloutStatus = 2

	// ----------------------------------------------------------------------
	// terminal states

	// RolloutStatusError is a failed rollout because of too many errors
	RolloutStatusError RolloutStatus = 3

	// RolloutStatusStable is a completed rollout saved for labeled agents connecting
	RolloutStatusStable RolloutStatus = 4

	// RolloutStatusReplaced is an incomplete rollout replaced by another rollout
	RolloutStatusReplaced RolloutStatus = 5
)

type StartRolloutPayload struct {
	Options *RolloutOptions `json:"options"`
}
