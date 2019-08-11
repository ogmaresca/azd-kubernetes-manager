package config

import (
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

	var fileSections []FileSection
	for _, value := range c.ServiceHooks {
		fileSections = append(fileSections, value)
	}

	return validate(fileSections, "Service Hook definition")
}
