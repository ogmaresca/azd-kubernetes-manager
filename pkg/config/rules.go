package config

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Rules lists all of the rules to perform upon an event
type Rules struct {
	// The resources to apply
	Apply []ApplyResourceRule `yaml:"apply"`

	// The resources to delete
	Delete []DeleteResourceRule `yaml:"delete"`
}

// DeleteResourceRule lists a resource to delete
type DeleteResourceRule struct {
	// The Kubernetes API version of the resource(s) to delete
	APIVersion string `yaml:"apiVersion"`

	// The resource kind
	Kind string `yaml:"kind"`

	// The resource namespace
	Namespace *string `yaml:"namespace"`

	// The label selector
	Selector LabelSelector `yaml:"selector"`

	// The maximum resources to delete. If the number of returned resources is < limit, then fail
	Limit *int `yaml:"limit"`
}

// ApplyResourceRule lists a resource to create
type ApplyResourceRule string

///
/// Describe()
///

// Describe returns a user-friendly representation of a Rules
func (r Rules) Describe() string {
	description := "Resource apply rules:"

	var applyRuleDescriptions []string
	for _, applyRule := range r.Delete {
		applyRuleDescriptions = append(applyRuleDescriptions, applyRule.Describe())
	}
	description += joinYAMLSlice(applyRuleDescriptions)

	description += "\nResource deletion rules:"

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

// Describe returns a user-friendly representation of a ApplyResourceRule
func (r ApplyResourceRule) Describe() string {
	resource, err := r.Parse()
	if err != nil {
		return fmt.Sprintf("Error parsing Kubernetes resource: %s", err.Error())
	}

	name := resource.Metadata.Name
	if name == "" {
		name = resource.Metadata.GenerateName
	}

	return fmt.Sprintf(
		"%s %s/%s",
		resource.APIVersion, resource.Kind, name,
	)
}

///
/// Validate
///

// Validate a Rules definition. This function returns a slice of warnings and an error.
func (r Rules) Validate() ([]string, error) {
	if len(r.Apply) == 0 && len(r.Delete) == 0 {
		return []string{"No rules were defined"}, nil
	}

	var warnings, errors []string

	var applyFileSections []FileSection
	for _, value := range r.Apply {
		applyFileSections = append(applyFileSections, value)
	}

	applyWarnings, applyErr := validate(applyFileSections, "Apply Resource rule definition")
	warnings = append(warnings, applyWarnings...)
	if applyErr != nil {
		errors = append(errors, applyErr.Error())
	}

	var deleteFileSections []FileSection
	for _, value := range r.Delete {
		deleteFileSections = append(deleteFileSections, value)
	}

	deleteWarnings, deleteErr := validate(deleteFileSections, "Delete Resource rule definition")
	warnings = append(warnings, deleteWarnings...)
	if deleteErr != nil {
		errors = append(errors, deleteErr.Error())
	}

	var err error
	if len(errors) > 0 {
		err = fmt.Errorf("%s", strings.Join(errors, "\n"))
	}

	return warnings, err
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

// Validate a Delete Kubernetes Resouce rule definition. This function returns a slice of warnings and an error.
func (r ApplyResourceRule) Validate() ([]string, error) {
	if string(r) == "" {
		return []string{}, fmt.Errorf("%s", "Value must not be empty.")
	}

	template, err := baseTemplate.Clone()
	if err != nil {
		return []string{}, err
	}
	template, err = template.Parse(string(r))
	if err != nil {
		return []string{}, fmt.Errorf("Error parsing Apply rule: %s", err.Error())
	}
	buffer := new(bytes.Buffer)
	err = template.Execute(buffer, sampleTemplatingArgs)
	if err != nil {
		return []string{}, fmt.Errorf("Error executing Apply rule with sample data: %s", err.Error())
	}
	templatedValue := buffer.String()
	if logger.LogDebug() {
		logger.Debugf("Converted Apply rule template:\n  %s\nto:\n  %s", strings.ReplaceAll(string(r), "\n", "\n  "), strings.ReplaceAll(templatedValue, "\n", "\n  "))
	}

	return []string{}, nil
}

//
// Mappings
//

// ToTypeMeta maps a DeleteResourceRule to a Kubernetes meta/v1 TypeMeta
func (r DeleteResourceRule) ToTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{Kind: r.Kind, APIVersion: r.APIVersion}
}

// Parse parses the rule as a YAML object
func (r ApplyResourceRule) Parse() (KubernetesResource, error) {
	resource := KubernetesResource{}

	err := yaml.Unmarshal([]byte(string(r)), &resource)

	return resource, err
}
