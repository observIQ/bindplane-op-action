package model

import "time"

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

type RolloutStatus int

type StartRolloutPayload struct {
	Options *RolloutOptions `json:"options"`
}

type Rollout struct {
	Name      string          `json:"name" yaml:"name" mapstructure:"name"`
	Status    RolloutStatus   `json:"status" yaml:"status" mapstructure:"status"`
	Options   RolloutOptions  `json:"options" yaml:"options" mapstructure:"options"`
	Phase     int             `json:"phase" yaml:"phase" mapstructure:"phase"`
	Progress  RolloutProgress `json:"progress" yaml:"progress" mapstructure:"progress"`
	StartedAt *time.Time      `json:"startedAt" yaml:"startedAt" mapstructure:"startedAt"`
	Stages    []RolloutStage  `json:"stages" yaml:"stages" mapstructure:"stages"`
	Stage     int             `json:"stage" yaml:"stage" mapstructure:"stage"`
}

type RolloutOptions struct {
	StartAutomatically bool            `json:"startAutomatically" yaml:"startAutomatically" mapstructure:"startAutomatically"`
	RollbackOnFailure  bool            `json:"rollbackOnFailure" yaml:"rollbackOnFailure" mapstructure:"rollbackOnFailure"`
	PhaseAgentCount    PhaseAgentCount `json:"phaseAgentCount" yaml:"phaseAgentCount" mapstructure:"phaseAgentCount"`
	MaxErrors          int             `json:"maxErrors" yaml:"maxErrors" mapstructure:"maxErrors"`
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
