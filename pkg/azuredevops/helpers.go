package azuredevops

import (
	"net/url"
	"strings"
)

// GetOrganizationFromURL grabs the Azure Devops URL
// Ex: https://dev.azure.com/Org/Project/more/path/entries returns Org
func GetOrganizationFromURL(url *url.URL) *string {
	if url != nil {
		splitPath := strings.SplitN(strings.TrimLeft(url.Path, "/"), "/", 2)
		if splitPath != nil && len(splitPath) >= 1 && splitPath[0] != "" {
			return &splitPath[0]
		}
	}

	return nil
}

// GetProjectFromURL grabs the Azure Devops URL
// Ex: https://dev.azure.com/Org/Project/more/path/entries returns Project
func GetProjectFromURL(url *url.URL) *string {
	if url != nil {
		splitPath := strings.SplitN(strings.TrimLeft(url.Path, "/"), "/", 3)
		if splitPath != nil && len(splitPath) >= 2 && splitPath[0] != "" {
			return &splitPath[1]
		}
	}

	return nil
}
