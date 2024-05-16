package model

type AgentSelector struct {
	MatchLabels `json:"matchLabels" yaml:"matchLabels" mapstructure:"matchLabels"`
}
