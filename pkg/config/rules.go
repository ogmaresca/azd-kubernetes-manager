package config

import (
	"fmt"
	"strings"
)

// Rules lists all of the rules to perform upon an event
type Rules struct {
	// The resources to delete
	Delete []DeleteResourceRule `yaml:"delete"`
}

// DeleteResourceRule lists all of the resources to delete
type DeleteResourceRule struct {
	// The Kubernetes API version of the resource(s) to delete
	APIVersion string `yaml:"apiVersion"`

	// The resource kind
	Kind string `yaml:"kind"`

	// The label selector
	Selector LabelSelector `yaml:"selector"`

	// The maximum resources to delete. If the number of returned resources is < limit, then fail
	Limit *int `yaml:"limit"`
}

///
/// Describe()
///

// Describe returns a user-friendly representation of a Rules
func (r Rules) Describe() string {
	description := "Resource deletion rules:"

	var deletionRuleDescriptions []string
	for _, deletionRule := range r.Delete {
		deletionRuleDescriptions = append(deletionRuleDescriptions, deletionRule.Describe())
	}
	description += joinYAMLSlice(deletionRuleDescriptions)

	return description
}

// Describe returns a user-friendly representation of a DeleteResourceRule
func (r DeleteResourceRule) Describe() string {
	return fmt.Sprintf(
		"API Version: %s\nKinds: %s\nLimit: %d\nLabel Selector:\n  %s",
		r.APIVersion, r.Kind, r.Limit, strings.ReplaceAll(r.Selector.Describe(), "\n", "\n  "),
	)
}

///
/// Validate
///

// Validate a Rules definition. This function returns a slice of warnings and an error.
func (r Rules) Validate() ([]string, error) {
	if len(r.Delete) == 0 {
		return []string{"No rules were defined"}, nil
	}

	var fileSections []FileSection
	for _, value := range r.Delete {
		fileSections = append(fileSections, value)
	}

	return validate(fileSections, "Delete Resource rule definition")
}

// Validate a Delete Kubernetes Resouce rule definition. This function returns a slice of warnings and an error.
func (r DeleteResourceRule) Validate() ([]string, error) {
	var errors []string
	var warnings []string

	if r.APIVersion == "" {
		errors = append(errors, "The Kubernetes API Version `APIVersion` must be defined. Use \"v1\" for the core API.")
	}

	if r.Kind == "" {
		errors = append(errors, "The Kubernetes resource `Kind` must be defined.")
	}

	selectorWarnings, err := r.Selector.Validate()
	if len(selectorWarnings) > 0 {
		warnings = append(warnings, selectorWarnings...)
	}
	if err != nil {
		errors = append(errors, err.Error())
	}

	if len(r.Selector.MatchLabels) == 0 && len(r.Selector.MatchExpressions) == 0 {
		errors = append(errors, "No label `Selector` was defined. See https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/ for defining set-based requirement selectors.")
	}

	if r.Limit != nil && *r.Limit <= 0 {
		errors = append(errors, "If a `Limit` is defined, it must be greater than 0.")
	}

	if len(errors) > 0 {
		err = fmt.Errorf("%s", strings.Join(errors, "\n"))
	}

	return warnings, err
}
