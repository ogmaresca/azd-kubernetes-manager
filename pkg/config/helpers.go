package config

import (
	newerrors "errors"
	"fmt"
	"regexp"
	"strings"
)

func joinYAMLSlice(slice []string) string {
	if len(slice) == 0 {
		return "[]"
	}

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
		err = newerrors.New(strings.Join(errors, "\n"))
	}

	return warnings, err
}

func contains(value string, filters []string) bool {
	if len(filters) == 0 {
		return true
	}
	for _, filter := range filters {
		if strings.EqualFold(value, filter) {
			return true
		}
	}
	return false
}

func intersection(values []string, filters []string) bool {
	if len(values) == 0 || len(filters) == 0 {
		return true
	}
	for _, value := range values {
		if contains(value, filters) {
			return true
		}
	}
	return false
}

func containsPOSIXERE(value string, filters []string) (bool, error) {
	if len(filters) == 0 {
		return true, nil
	}
	for _, filter := range filters {
		pattern, err := regexp.CompilePOSIX(filter)
		if err != nil {
			return false, nil
		}
		if pattern.Match([]byte(value)) {
			return true, nil
		}
	}
	return false, nil
}
