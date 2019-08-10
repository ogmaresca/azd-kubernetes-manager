package config

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// FileSection is an interface for the common methods in all File structs
type FileSection interface {
	Describe() string
	Validate() ([]string, error)
}

// File is the root struct representing the config file
type File struct {
	// Service hook rules
	ServiceHooks []ServiceHook `yaml:"serviceHooks"`
}

// NewConfigFile creates a ConfigFile from YAML
func NewConfigFile(configFileYaml []byte) (File, error) {
	configFile := File{}

	err := yaml.Unmarshal(configFileYaml, &configFile)

	return configFile, err
}

// Describe returns a user-friendly representation of a ConfigFile
func (c File) Describe() string {
	description := "==============\nService Hooks:\n=============="

	var serviceHookDescriptions []string
	for _, serviceHook := range c.ServiceHooks {
		serviceHookDescriptions = append(serviceHookDescriptions, serviceHook.Describe())
	}
	description += joinYAMLSlice(serviceHookDescriptions)

	return description
}

// Validate a Config File. This function returns a slice of warnings and an error.
func (c File) Validate() ([]string, error) {
	if len(c.ServiceHooks) == 0 {
		return []string{"No rules were defined. azd-kubernetes-manager will just log Service Hook requests."}, nil
	}

	var errors []string
	var warnings []string
	for pos, serviceHook := range c.ServiceHooks {
		serviceHookWarnings, err := serviceHook.Validate()
		if len(warnings) > 0 {
			warnings = append(warnings, fmt.Sprintf("Warnings from Service Hook definition %d:%s", pos, joinYAMLSlice(serviceHookWarnings)))
		}
		if err != nil {
			errors = append(errors, fmt.Sprintf("Errors from Service Hook definition %d:\n    %s", pos, strings.ReplaceAll(err.Error(), "\n", "\n  ")))
		}
	}

	var err error
	if len(errors) > 0 {
		err = fmt.Errorf("%s", strings.Join(errors, "\n"))
	}

	return warnings, err
}

func joinYAMLSlice(slice []string) string {
	str := ""
	for _, val := range slice {
		str += fmt.Sprintf("\n- %s", strings.ReplaceAll(val, "\n", "\n  "))
	}
	return str
}
