package config

import (
	newerrors "errors"
	"fmt"
	"strings"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/templating"
	"gopkg.in/yaml.v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Rules lists all of the rules to perform upon an event
type Rules struct {
	// The resources to apply
	Apply []ApplyResourceRule `yaml:"apply"`

	// The resources to delete
	Delete []DeleteResourceRule `yaml:"delete"`
}

// ApplyResourceRule lists a resource to create
type ApplyResourceRule string

// DeleteResourceRule lists a resource to delete
type DeleteResourceRule struct {
	// The Kubernetes API version of the resource(s) to delete
	APIVersion string `yaml:"apiVersion"`

	// The resource kind
	Kind string `yaml:"kind"`

	// The resource namespace
	Namespace string `yaml:"namespace,omitempty"`

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
	if r.IsEmpty() {
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
		err = newerrors.New(strings.Join(errors, "\n"))
	}

	return warnings, err
}

// Validate a Delete Kubernetes Resouce rule definition. This function returns a slice of warnings and an error.
func (r ApplyResourceRule) Validate() ([]string, error) {
	if string(r) == "" {
		return []string{}, newerrors.New("Apply rule must not be empty.")
	}

	strVal := r.String()
	templatedValue, err := templating.Execute("ConfigFileValidation", strVal, sampleTemplatingArgs)
	if err != nil {
		return []string{}, fmt.Errorf("Apply rule error: %s", err.Error())
	} else if logger.LogDebug() && templatedValue != strVal {
		logger.Debugf("Converted Apply rule template:\n  %s\nto:\n  %s", strings.ReplaceAll(strVal, "\n", "\n  "), strings.ReplaceAll(templatedValue, "\n", "\n  "))
	}

	templatedRule := ApplyResourceRule(templatedValue)

	_, err = templatedRule.Parse()
	if err != nil {
		return []string{}, fmt.Errorf("Error parsing Apply rule after templating: %s", err.Error())
	}

	return []string{}, nil
}

// Validate a Delete Kubernetes Resouce rule definition. This function returns a slice of warnings and an error.
func (r DeleteResourceRule) Validate() ([]string, error) {
	var errors []string
	var warnings []string

	if r.APIVersion == "" {
		errors = append(errors, "The Kubernetes API Version `APIVersion` must be defined. Use \"v1\" for the core API.")
	} else {
		split := strings.Split(r.APIVersion, "/")
		if len(split) != 1 && len(split) != 2 {
			errors = append(errors, fmt.Sprintf("Invalid API Version '%s'", r.APIVersion))
		}
	}

	if r.Kind == "" {
		errors = append(errors, "The Kubernetes resource `Kind` must be defined.")
	}

	if r.Namespace != "" {
		templatedNamespace, err := templating.Execute("ConfigFileValidation", r.Namespace, sampleTemplatingArgs)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Delete rule namespace templating error: %s", err.Error()))
		} else if logger.LogDebug() && templatedNamespace != r.Namespace {
			logger.Debugf("Converted Delete rule namespace template:\n  %s\nto:\n  %s", strings.ReplaceAll(r.Namespace, "\n", "\n  "), strings.ReplaceAll(templatedNamespace, "\n", "\n  "))
		}
	}

	selectorWarnings, err := r.Selector.Validate()
	if len(selectorWarnings) > 0 {
		warnings = append(warnings, selectorWarnings...)
	}
	if err != nil {
		errors = append(errors, err.Error())
	}

	if r.Limit != nil && *r.Limit <= 0 {
		errors = append(errors, "If a `Limit` is defined, it must be greater than 0.")
	}

	if len(errors) > 0 {
		err = newerrors.New(strings.Join(errors, "\n"))
	}

	return warnings, err
}

//
// Mappings
//

// IsEmpty returns true if the Rules has doesn't contain any rules
func (r Rules) IsEmpty() bool {
	return len(r.Apply) == 0 && len(r.Delete) == 0
}

// ToTypeMeta maps a DeleteResourceRule to a Kubernetes meta/v1 TypeMeta
func (r DeleteResourceRule) ToTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{Kind: r.Kind, APIVersion: r.APIVersion}
}

// ToGroupVersion maps a DeleteResourceRule to a GroupVersion
func (r DeleteResourceRule) ToGroupVersion() schema.GroupVersion {
	split := strings.Split(r.APIVersion, "/")
	if len(split) == 1 {
		return schema.GroupVersion{Group: "", Version: split[0]}
	} else {
		return schema.GroupVersion{Group: split[0], Version: split[1]}
	}
}

// String returns the rule as a string
func (r ApplyResourceRule) String() string {
	return string(r)
}

// Parse parses the rule as a YAML object
func (r ApplyResourceRule) Parse() (KubernetesResource, error) {
	resource := KubernetesResource{}

	err := yaml.Unmarshal([]byte(r.String()), &resource)

	return resource, err
}

// ParseTemplated templates the rule and then parses it as a YAML object
func (r ApplyResourceRule) ParseTemplated(args templating.Args) (KubernetesResource, error) {
	templatedRule, err := templating.Execute("ApplyResourceRule", r.String(), args)
	if err != nil {
		return KubernetesResource{}, err
	}

	resource := KubernetesResource{}

	err = yaml.Unmarshal([]byte(templatedRule), &resource)

	return resource, err
}
