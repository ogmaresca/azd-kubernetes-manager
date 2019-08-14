package config_test

import (
	"net/url"
	"testing"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/azuredevops"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/config"
)

func TestMatches(t *testing.T) {
	t.Run("test_matches_templating_good", func(t *testing.T) {
		config := config.ServiceHook{
			Event: config.ServiceHookEventTypePullRequestUpdated,
			ResourceFilters: config.ServiceHookResourceFilters{
				Templates: []string{
					"{{ .CreatedBy.DisplayName | title | contains \"Obama\" }}",
					"{{ eq (.Commits | len) 2 }}",
				},
			},
		}

		hook := &azuredevops.ServiceHook{
			EventType: string(config.Event),
			Resource: azuredevops.ServiceHookResource{
				ServiceHookResourcePullRequest: azuredevops.ServiceHookResourcePullRequest{
					PullRequestID: intPtr(1),
					CreatedBy: azuredevops.User{
						ID:          "MockUserId",
						DisplayName: "Barack Obama",
					},
					Title:         strPtr("Mock Pull Request"),
					Description:   strPtr("2 commits"),
					SourceRefName: strPtr("refs/heads/feature/mock"),
					TargetRefName: strPtr("refs/heads/feature/master"),
				},
				Repository: &azuredevops.GitRepository{
					StrDefinition: azuredevops.StrDefinition{
						ID:   "MockRepositoryId",
						Name: "MockRepository",
					},
					URL: forceParseURL("https://dev.azure.com/MockOrganization/MockProject/_apis/repos/git/repositories/MockRepositoryId", t),
					Project: azuredevops.GitProject{
						StrDefinition: azuredevops.StrDefinition{
							ID:   "MockProjectId",
							Name: "MockProject",
						},
						URL:   forceParseURL("https://dev.azure.com/MockOrganization/MockProject/_apis/projects/MockProjectId", t),
						State: "wellFormed",
					},
				},
				Commits: []azuredevops.GitCommit{
					azuredevops.GitCommit{
						CommitID: "SampleCommitId1",
						URL:      forceParseURL("https://dev.azure.com/SampleOrganization/SampleProject/_apis/repos/git/repositories/SampleRepositoryId/commits/SampleCommitId1", t),
					},
					azuredevops.GitCommit{
						CommitID: "SampleCommitId2",
						URL:      forceParseURL("https://dev.azure.com/SampleOrganization/SampleProject/_apis/repos/git/repositories/SampleRepositoryId/commits/SampleCommitId2", t),
					},
				},
				Status: strPtr(string(azuredevops.StatusCompleted)),
			},
		}

		if matches, err := config.Matches(hook); !matches || err != nil {
			t.Errorf("Expected Service Hook to match configuration rule: %+v", err)
		}
	})

	t.Run("test_matches_templating_false", func(t *testing.T) {
		config := config.ServiceHook{
			Event: config.ServiceHookEventTypePullRequestUpdated,
			ResourceFilters: config.ServiceHookResourceFilters{
				Templates: []string{
					"{{ .CreatedBy.DisplayName | title | contains \"Roosevelt\" }}",
					"{{ eq (.Commits | len) 2 }}",
				},
			},
		}

		hook := &azuredevops.ServiceHook{
			EventType: string(config.Event),
			Resource: azuredevops.ServiceHookResource{
				ServiceHookResourcePullRequest: azuredevops.ServiceHookResourcePullRequest{
					PullRequestID: intPtr(1),
					CreatedBy: azuredevops.User{
						ID:          "MockUserId",
						DisplayName: "Barack Obama",
					},
					Title:         strPtr("Mock Pull Request"),
					Description:   strPtr("2 commits"),
					SourceRefName: strPtr("refs/heads/feature/mock"),
					TargetRefName: strPtr("refs/heads/feature/master"),
				},
				Repository: &azuredevops.GitRepository{
					StrDefinition: azuredevops.StrDefinition{
						ID:   "MockRepositoryId",
						Name: "MockRepository",
					},
					URL: forceParseURL("https://dev.azure.com/MockOrganization/MockProject/_apis/repos/git/repositories/MockRepositoryId", t),
					Project: azuredevops.GitProject{
						StrDefinition: azuredevops.StrDefinition{
							ID:   "MockProjectId",
							Name: "MockProject",
						},
						URL:   forceParseURL("https://dev.azure.com/MockOrganization/MockProject/_apis/projects/MockProjectId", t),
						State: "wellFormed",
					},
				},
				Commits: []azuredevops.GitCommit{
					azuredevops.GitCommit{
						CommitID: "SampleCommitId1",
						URL:      forceParseURL("https://dev.azure.com/SampleOrganization/SampleProject/_apis/repos/git/repositories/SampleRepositoryId/commits/SampleCommitId1", t),
					},
					azuredevops.GitCommit{
						CommitID: "SampleCommitId2",
						URL:      forceParseURL("https://dev.azure.com/SampleOrganization/SampleProject/_apis/repos/git/repositories/SampleRepositoryId/commits/SampleCommitId2", t),
					},
				},
				Status: strPtr(string(azuredevops.StatusCompleted)),
			},
		}

		if matches, err := config.Matches(hook); matches || err != nil {
			t.Errorf("Expected Service Hook not to match the configuration rule: %+v", err)
		}
	})

	t.Run("test_matches_templating_bad", func(t *testing.T) {
		config := config.ServiceHook{
			Event: config.ServiceHookEventTypePullRequestUpdated,
			ResourceFilters: config.ServiceHookResourceFilters{
				Templates: []string{
					"{{ .NonexistentField }}",
					"{{ eq (.Commits | len) 2 }}",
				},
			},
		}

		hook := &azuredevops.ServiceHook{
			EventType: string(config.Event),
			Resource: azuredevops.ServiceHookResource{
				ServiceHookResourcePullRequest: azuredevops.ServiceHookResourcePullRequest{
					PullRequestID: intPtr(1),
					CreatedBy: azuredevops.User{
						ID:          "MockUserId",
						DisplayName: "Barack Obama",
					},
					Title:         strPtr("Mock Pull Request"),
					Description:   strPtr("2 commits"),
					SourceRefName: strPtr("refs/heads/feature/mock"),
					TargetRefName: strPtr("refs/heads/feature/master"),
				},
				Repository: &azuredevops.GitRepository{
					StrDefinition: azuredevops.StrDefinition{
						ID:   "MockRepositoryId",
						Name: "MockRepository",
					},
					URL: forceParseURL("https://dev.azure.com/MockOrganization/MockProject/_apis/repos/git/repositories/MockRepositoryId", t),
					Project: azuredevops.GitProject{
						StrDefinition: azuredevops.StrDefinition{
							ID:   "MockProjectId",
							Name: "MockProject",
						},
						URL:   forceParseURL("https://dev.azure.com/MockOrganization/MockProject/_apis/projects/MockProjectId", t),
						State: "wellFormed",
					},
				},
				Commits: []azuredevops.GitCommit{
					azuredevops.GitCommit{
						CommitID: "SampleCommitId1",
						URL:      forceParseURL("https://dev.azure.com/SampleOrganization/SampleProject/_apis/repos/git/repositories/SampleRepositoryId/commits/SampleCommitId1", t),
					},
					azuredevops.GitCommit{
						CommitID: "SampleCommitId2",
						URL:      forceParseURL("https://dev.azure.com/SampleOrganization/SampleProject/_apis/repos/git/repositories/SampleRepositoryId/commits/SampleCommitId2", t),
					},
				},
				Status: strPtr(string(azuredevops.StatusCompleted)),
			},
		}

		if matches, err := config.Matches(hook); matches || err == nil {
			t.Errorf("Expected Service Hook to fail to match configuration rule due to an error")
		}
	})
}

func forceParseURL(strURL string, t *testing.T) url.URL {
	structURL, err := url.Parse(strURL)
	if err != nil {
		t.Errorf("Error parsing URL %s: %s", strURL, err.Error())
	}
	return *structURL
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}
