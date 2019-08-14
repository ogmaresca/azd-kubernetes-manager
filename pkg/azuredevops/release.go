package azuredevops

import (
	"net/url"
	"time"
)

// Release holds a Release definition
type Release struct {
	IntDefinition
	ReleaseID                  int                           `json:"releaseId"`
	Status                     string                        `json:"status"`
	CreatedOn                  time.Time                     `json:"createdOn"`
	ModifiedOn                 time.Time                     `json:"modifiedOn"`
	ModifiedBy                 User                          `json:"modifiedBy"`
	CreatedBy                  User                          `json:"createdBy"`
	Environments               []Environment                 `json:"environments"`
	Variables                  map[string]string             `json:"variables"`
	Artifacts                  []ReleaseArtifact             `json:"artifacts"`
	ReleaseDefinition          ServiceHookResourceDefinition `json:"releaseDefinition"`
	Description                string                        `json:"description"`
	Reason                     string                        `json:"reason"`
	ReleaseNameFormat          string                        `json:"releaseNameFormat"`
	KeepForever                bool                          `json:"keepForever"`
	DefinitionSnapshotRevision int                           `json:"definitionSnapshotRevision"`
	Comment                    string                        `json:"comment"`
	LogsContainerURL           *url.URL                      `json:"logsContainerUrl"`
}

// ReleaseArtifact holds an artifact of a Release
type ReleaseArtifact struct {
	SourceID            string                   `json:"sourceId"`
	Type                string                   `json:"type"`
	Alias               string                   `json:"alias"`
	DefinitionReference map[string]StrDefinition `json:"definitionReference"`
	IsPrimary           bool                     `json:"isPrimary"`
}

// ReleaseArtifactType are the artifact types
// https://docs.microsoft.com/en-us/azure/devops/pipelines/artifacts/artifacts-overview?view=azure-devops
type ReleaseArtifactType string

const (
	// ReleaseArtifactTypeBuild is for Build pipeline artifacts
	ReleaseArtifactTypeBuild ReleaseArtifactType = "Build"
	// ReleaseArtifactTypePipeline is for Pipeline pipeline artifacts
	ReleaseArtifactTypePipeline ReleaseArtifactType = "Pipeline"
	// ReleaseArtifactTypeMaven is for Maven pipeline artifacts
	ReleaseArtifactTypeMaven ReleaseArtifactType = "Maven"
	// ReleaseArtifactTypeNPM is for NPM pipeline artifacts
	ReleaseArtifactTypeNPM ReleaseArtifactType = "npm"
	// ReleaseArtifactTypeNuGet is for NuGet pipeline artifacts
	ReleaseArtifactTypeNuGet ReleaseArtifactType = "NuGet"
	// ReleaseArtifactTypePyPI is for PyPI pipeline artifacts
	ReleaseArtifactTypePyPI ReleaseArtifactType = "PyPI"
	// ReleaseArtifactTypeSymbols is for Symbols pipeline artifacts
	ReleaseArtifactTypeSymbols ReleaseArtifactType = "Symbols"
	// ReleaseArtifactTypeUniversal is for Universal pipeline artifacts
	ReleaseArtifactTypeUniversal ReleaseArtifactType = "Universal"
)
