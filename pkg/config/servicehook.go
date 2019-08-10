package config

import (
	"fmt"
	"strings"
)

// ServiceHook holds rules for a Service Hook
type ServiceHook struct {
	// The EventType
	// See https://docs.microsoft.com/en-us/azure/devops/service-hooks/events?toc=%2Fazure%2Fdevops%2Fintegrate%2Ftothis.json&bc=%2Fazure%2Fdevops%2Fintegrate%2Fbreadcrumb%2Ftothis.json&view=azure-devops
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

// Describe returns a user-friendly representation of a ServiceHook
func (sh ServiceHook) Describe() string {
	return fmt.Sprintf(
		"Event Type: %s\nResource Filters:\n  Statuses: %+v\n  Reasons: %+v\n  Projects: %+v\nContinue: %t\nRules:\n  %s",
		sh.Event, sh.ResourceStatuses, sh.ResourceReasons, sh.ProjectNames, sh.Continue, strings.ReplaceAll(sh.Rules.Describe(), "\n", "\n  "),
	)
}

// Validate a Service Hook rule definition. This function returns a slice of warnings and an error.
func (sh ServiceHook) Validate() ([]string, error) {
	var errors []string

	if sh.Event == "" {
		errors = append(errors, "The `event` field must be defined.")
	}

	rulesWarnings, err := sh.Rules.Validate()

	if len(errors) > 0 {
		if err != nil {
			errors = append(errors, err.Error())
		}

		err = fmt.Errorf("%s", joinYAMLSlice(errors))
	}

	return rulesWarnings, err
}
