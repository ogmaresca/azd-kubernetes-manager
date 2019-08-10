package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LabelSelector represents a k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector, because it doesn't have YAML tags
type LabelSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`

	MatchExpressions []LabelSelectorRequirement `yaml:"matchExpressions"`
}

// LabelSelectorRequirement represents a k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelectorRequirement, because it doesn't have YAML tags
type LabelSelectorRequirement struct {
	Key      string                       `yaml:"key"`
	Operator metav1.LabelSelectorOperator `yaml:"operator"`
	Values   []string                     `yaml:"values"`
}

// ToKubernetesLabelSelector maps a LabelSelector to a k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector
func (ls LabelSelector) ToKubernetesLabelSelector() metav1.LabelSelector {
	var expressions []metav1.LabelSelectorRequirement
	for _, expression := range ls.MatchExpressions {
		expressions = append(expressions, expression.ToKubernetesLabelSelectorRequirement())
	}

	return metav1.LabelSelector{
		MatchLabels:      ls.MatchLabels,
		MatchExpressions: expressions,
	}
}

// ToKubernetesLabelSelectorRequirement maps a LabelSelectorRequirement to a k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelectorRequirement
func (lsr LabelSelectorRequirement) ToKubernetesLabelSelectorRequirement() metav1.LabelSelectorRequirement {
	return metav1.LabelSelectorRequirement{
		Key:      lsr.Key,
		Operator: lsr.Operator,
		Values:   lsr.Values,
	}
}
