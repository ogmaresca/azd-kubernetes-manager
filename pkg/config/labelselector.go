package config

import (
	"fmt"
	"strings"

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

///
/// Describe()
///

// Describe returns a user-friendly representation of a LabelSelector
func (ls LabelSelector) Describe() string {
	description := ""

	if len(ls.MatchLabels) > 0 {
		description += "Match Labels:\n"
		for label, value := range ls.MatchLabels {
			description += fmt.Sprintf("  %s: %s\n", label, value)
		}
	}

	if len(ls.MatchExpressions) > 0 {
		description += "Match Label Expressions:"
		var expressionDescriptions []string
		for _, expression := range ls.MatchExpressions {
			expressionDescriptions = append(expressionDescriptions, expression.Describe())
		}
		description += joinYAMLSlice(expressionDescriptions)
	}

	return description
}

// Describe returns a user-friendly representation of a LabelSelectorRequirement
func (lsr LabelSelectorRequirement) Describe() string {
	switch lsr.Operator {
	case metav1.LabelSelectorOpExists:
	case metav1.LabelSelectorOpDoesNotExist:
		return fmt.Sprintf("%s %s", lsr.Key, lsr.Operator)
	case metav1.LabelSelectorOpIn:
	case metav1.LabelSelectorOpNotIn:
	default:
		break
	}

	return fmt.Sprintf("%s %s %v", lsr.Key, lsr.Operator, lsr.Values)
}

///
/// Validate
///

// Validate a LabelSelector. This function returns a slice of warnings and an error.
func (ls LabelSelector) Validate() ([]string, error) {
	if len(ls.MatchLabels) == 0 && len(ls.MatchExpressions) == 0 {
		return []string{}, fmt.Errorf("%s", "No label `Selector` was defined. See https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/ for defining set-based requirement selectors.")
	}

	var fileSections []FileSection
	for _, value := range ls.MatchExpressions {
		fileSections = append(fileSections, value)
	}

	return validate(fileSections, "Label MatchExpressions")
}

// Validate a LabelSelectorRequirement. This function returns a slice of warnings and an error.
func (lsr LabelSelectorRequirement) Validate() ([]string, error) {
	var errors []string
	var warnings []string

	if lsr.Key == "" {
		errors = append(errors, "The label `Key` must be defined.")
	}

	if lsr.Operator == "" {
		errors = append(errors, "The label expression `Operator` must be defined.")
	} else {
		switch lsr.Operator {
		case metav1.LabelSelectorOpExists:
		case metav1.LabelSelectorOpDoesNotExist:
			if len(lsr.Values) != 0 {
				errors = append(errors, fmt.Sprintf("The label expression `Operator` \"%s\" cannot define `Values`.", lsr.Operator))
			}
		case metav1.LabelSelectorOpIn:
		case metav1.LabelSelectorOpNotIn:
			if len(lsr.Values) == 0 {
				errors = append(errors, fmt.Sprintf("The label expression `Operator` \"%s\" must define `Values`.", lsr.Operator))
			}
		default:
			errors = append(errors, fmt.Sprintf("Invalid label expression `Operator` \"%s\"", lsr.Operator))
		}
	}

	var err error
	if len(errors) > 0 {
		err = fmt.Errorf("%s", strings.Join(errors, "\n"))
	}

	return warnings, err
}

///
/// Mappings
///

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
