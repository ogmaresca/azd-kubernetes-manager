package azuredevops_test

import (
	"testing"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/azuredevops"
)

func TestGetOrganizationFromURL(t *testing.T) {
	t.Run("test_getorganizationfromurl_good", func(t *testing.T) {
		url := "https://dev.azure.com/OrganizationName"
		organization := azuredevops.GetOrganizationFromURL(&url)
		if organization == nil || *organization != "OrganizationName" {
			t.Fatalf("Expected %s, but got %#v", "OrganizationName", organization)
		}
	})

	t.Run("test_getorganizationfromurl_withmorepath_good", func(t *testing.T) {
		url := "https://dev.azure.com/OrganizationName/ProjectName/further/url/paths"
		organization := azuredevops.GetOrganizationFromURL(&url)
		if organization == nil || *organization != "OrganizationName" {
			t.Fatalf("Expected %s, but got %#v", "OrganizationName", organization)
		}
	})

	t.Run("test_getorganizationfromurl_bad", func(t *testing.T) {
		url := "https://dev.azure.com/"
		organization := azuredevops.GetOrganizationFromURL(&url)
		if organization != nil {
			t.Fatalf("Expected nil, but got %s", *organization)
		}
	})

	t.Run("test_getorganizationfromurl_nilurl_bad", func(t *testing.T) {
		organization := azuredevops.GetOrganizationFromURL(nil)
		if organization != nil {
			t.Fatalf("Expected nil, but got %s", *organization)
		}
	})
}

func TestGetProjectFromURL(t *testing.T) {
	t.Run("test_getprojectfromurl_good", func(t *testing.T) {
		url := "https://dev.azure.com/OrganizationName/ProjectName"
		project := azuredevops.GetProjectFromURL(&url)
		if project == nil || *project != "ProjectName" {
			t.Fatalf("Expected %s, but got %#v", "ProjectName", project)
		}
	})

	t.Run("test_getprojectfromurl_withmorepath_good", func(t *testing.T) {
		url := "https://dev.azure.com/OrganizationName/ProjectName/further/url/paths"
		project := azuredevops.GetProjectFromURL(&url)
		if project == nil || *project != "ProjectName" {
			t.Fatalf("Expected %s, but got %#v", "ProjectName", project)
		}
	})

	t.Run("test_getprojectfromurl_bad", func(t *testing.T) {
		url := "https://dev.azure.com/"
		project := azuredevops.GetProjectFromURL(&url)
		if project != nil {
			t.Fatalf("Expected nil, but got %s", *project)
		}
	})

	t.Run("test_getprojectfromurl_nilurl_bad", func(t *testing.T) {
		project := azuredevops.GetProjectFromURL(nil)
		if project != nil {
			t.Fatalf("Expected nil, but got %s", *project)
		}
	})
}
