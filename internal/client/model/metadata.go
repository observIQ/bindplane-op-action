package model

import "time"

type Metadata struct {
	ID              string            `yaml:"id,omitempty" json:"id" mapstructure:"id"`
	Name            string            `yaml:"name,omitempty" json:"name" mapstructure:"name"`
	DisplayName     string            `yaml:"displayName,omitempty" json:"displayName,omitempty" mapstructure:"displayName"`
	Description     string            `yaml:"description,omitempty" json:"description,omitempty" mapstructure:"description"`
	Icon            string            `yaml:"icon,omitempty" json:"icon,omitempty" mapstructure:"icon"`
	Labels          map[string]string `yaml:"labels,omitempty" json:"labels" mapstructure:"labels"`
	Hash            string            `yaml:"hash,omitempty" json:"hash,omitempty" mapstructure:"hash"`
	Version         int               `yaml:"version,omitempty" json:"version,omitempty" mapstructure:"version"`
	DateModified    *time.Time        `yaml:"dateModified,omitempty" json:"dateModified,omitempty" mapstructure:"dateModified"`
	Deprecated      bool              `yaml:"deprecated,omitempty" json:"deprecated,omitempty" mapstructure:"deprecated"`
	AdditionalInfo  *AdditionalInfo   `yaml:"additionalInfo,omitempty" json:"additionalInfo,omitempty" mapstructure:"additionalInfo"`
	ResourceDocLink string            `yaml:"resourceDocLink,omitempty" json:"resourceDocLink,omitempty" mapstructure:"resourceDocLink"`
	Stability       string            `yaml:"stability,omitempty" json:"stability,omitempty" mapstructure:"stability"`
}
