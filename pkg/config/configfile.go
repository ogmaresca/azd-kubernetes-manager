package config

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gopkg.in/yaml.v2"
)

// ConfigFile is the root struct representing the config file
type ConfigFile struct {
	// Service hook rules
	ServiceHooks []ServiceHook `yaml:"serviceHooks"`
}

// ServiceHook holds rules for a Service Hook
type ServiceHook struct {
	// The EventType
	// See https://docs.microsoft.com/en-us/azure/devops/service-hooks/events?toc=%2Fazure%2Fdevops%2Fintegrate%2Ftoc.json&bc=%2Fazure%2Fdevops%2Fintegrate%2Fbreadcrumb%2Ftoc.json&view=azure-devops
	// for the event type
	// It must match the API type, ex build.completed
	Event string `yaml:"event"`

	// Only execute the rules if {resource.status ∈ this.ResourceStatuses}
	ResourceStatuses []string `yaml:"resourceStatuses"`

	// Only execute the rules if {resource.reason ∈ this.ResourceReasons}
	ResourceReasons []string `yaml:"resourceReasons"`

	// Only execute the rules if {resource.project.name ∈ this.ProjectNames}
	ProjectNames []string `yaml:"projectNames"`

	// If Continue is true, process any other service hooks that match the Service Hook
	Continue bool `yaml:"continue"`

	// The rules to perform on the Service Hook
	Rules Rules `yaml:"rules"`
}

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
	Selector metav1.LabelSelector `yaml:"selector"`

	// The maximum resources to delete. If the number of returned resources is < limit, then fail
	Limit *int `yaml:"limit"`
}

// NewConfigFile creates a ConfigFile from YAML
func NewConfigFile(configFileYaml []byte) (ConfigFile, error) {
	configFile := ConfigFile{}

	err := yaml.Unmarshal(configFileYaml, &configFile)

	return configFile, err
}

// Describe returns a user-friendly representation of a ConfigFile
func (c ConfigFile) Describe() string {
	description := "==============\nService Hooks:\n=============="

	for _, serviceHook := range c.ServiceHooks {
		description += fmt.Sprintf("\n- %s", strings.ReplaceAll(serviceHook.Describe(), "\n", "\n  "))
	}

	return description
}

// Describe returns a user-friendly representation of a ServiceHook
func (sh ServiceHook) Describe() string {
	return fmt.Sprintf(
		"Event Type: %s\nResource Filters:\n  Statuses: %v\n  Reasons: %+v\n  Projects: %#v\nContinue: %t\nRules:\n  %s",
		sh.Event, sh.ResourceStatuses, sh.ResourceReasons, sh.ProjectNames, sh.Continue, strings.ReplaceAll(sh.Rules.Describe(), "\n", "\n  "),
	)
}

// Describe returns a user-friendly representation of a Rules
func (r Rules) Describe() string {
	description := "Resource deletion rules:"

	for _, serviceHook := range r.Delete {
		description += fmt.Sprintf("\n- %s", strings.ReplaceAll(serviceHook.Describe(), "\n", "\n  "))
	}

	return description
}

// Describe returns a user-friendly representation of a DeleteResourceRule
func (r DeleteResourceRule) Describe() string {
	return fmt.Sprintf(
		"API Version: %s\nKinds: %s\nLimit: %d\nLabel Selector:\n  %s",
		r.APIVersion, r.Kind, r.Limit, strings.ReplaceAll(fmt.Sprintf("%+v", r.Selector), "\n", "\n  "),
	)
}
