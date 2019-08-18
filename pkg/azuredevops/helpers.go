package azuredevops

import (
	"net/url"
	"strings"

	"github.com/alexcesaro/log/stdlog"
)

var logger = stdlog.GetFromFlags()

// GetOrganizationFromURL grabs the Azure Devops URL
// Ex: https://dev.azure.com/Org/Project/more/path/entries returns Org
func GetOrganizationFromURL(urlStr *string) *string {
	if urlStr != nil {
		url, err := url.Parse(*urlStr)
		if err != nil {
			logger.Errorf("Error parsing URL '%s': %s", urlStr, err.Error())
			return nil
		}

		splitPath := strings.SplitN(strings.TrimLeft(url.Path, "/"), "/", 2)
		if splitPath != nil && len(splitPath) >= 1 && splitPath[0] != "" {
			return &splitPath[0]
		}
	}

	return nil
}

// GetProjectFromURL grabs the Azure Devops URL
// Ex: https://dev.azure.com/Org/Project/more/path/entries returns Project
func GetProjectFromURL(urlStr *string) *string {
	if urlStr != nil {
		url, err := url.Parse(*urlStr)
		if err != nil {
			logger.Errorf("Error parsing URL '%s': %s", urlStr, err.Error())
			return nil
		}

		splitPath := strings.SplitN(strings.TrimLeft(url.Path, "/"), "/", 3)
		if splitPath != nil && len(splitPath) >= 2 && splitPath[0] != "" {
			return &splitPath[1]
		}
	}

	return nil
}
