package config

import (
	"fmt"
	"strings"
)

func joinYAMLSlice(slice []string) string {
	str := ""
	for _, val := range slice {
		str += fmt.Sprintf("\n- %s", strings.ReplaceAll(val, "\n", "\n  "))
	}
	return str
}

func validate(slice []FileSection, description string) ([]string, error) {
	var errors []string
	var warnings []string
	for pos, value := range slice {
		valueWarnings, err := value.Validate()
		if len(valueWarnings) > 0 {
			warnings = append(warnings, fmt.Sprintf("Warnings from %s %d:%s", description, pos, joinYAMLSlice(valueWarnings)))
		}
		if err != nil {
			errors = append(errors, fmt.Sprintf("Errors from %s %d:\n    %s", description, pos, strings.ReplaceAll(err.Error(), "\n", "\n  ")))
		}
	}

	var err error
	if len(errors) > 0 {
		err = fmt.Errorf("%s", strings.Join(errors, "\n"))
	}

	return warnings, err
}
