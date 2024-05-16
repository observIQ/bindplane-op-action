package model

import "k8s.io/apimachinery/pkg/labels"

type MatchLabels map[string]string

type Labels struct {
	labels.Set `json:"-" yaml:",inline"`
}
