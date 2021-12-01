package config

import (
	"time"

	"github.com/ogmaresca/azd-kubernetes-manager/pkg/azuredevops"
	"github.com/ogmaresca/azd-kubernetes-manager/pkg/templating"
)

var (
	sampleServiceHook = func() azuredevops.ServiceHook {
		intPtr := func(i int) *int {
			return &i
		}

		strPtr := func(s string) *string {
			return &s
		}

		sampleTime := time.Now()

		sampleBool := true

		sampleUser := azuredevops.User{
			ID:          "SampleUserId",
			DisplayName: "FirstName Last-Name",
			UniqueName:  "firstnamelast-name@example.com",
			URL:         strPtr("https://dev.azure.com/SampleOrganization/_apis/Identities/d6245f20-2af8-44f4-9451-8107cb2767db"),
			ImageURL:    strPtr("https://dev.azure.com/SampleOrganization/SampleProject/_api/_common/identityImage?id=d6245f20-2af8-44f4-9451-8107cb2767db"),
			Descriptor:  "Sample Descriptor",
		}

		sampleEnvironment := azuredevops.Environment{
			IntDefinition: azuredevops.IntDefinition{ID: 1, Name: "Sample Environment"},
			ReleaseID:     1,
			Status:        string(azuredevops.StatusSucceeded),
			Variables: map[string]string{
				"Sample Variable": "Value",
			},
			PreDeployApprovals:  []azuredevops.User{sampleUser},
			PostDeployApprovals: []azuredevops.User{sampleUser},
			PreApprovalsSnapshot: []azuredevops.EnvironmentApprovalSnapshot{
				azuredevops.EnvironmentApprovalSnapshot{
					Approvals: []azuredevops.User{sampleUser},
					ApprovalOptions: map[string]string{
						"requiredApproverCount":       "0",
						"releaseCreatorCanBeApprover": "true",
					},
				},
			},
			PostApprovalsSnapshot: []azuredevops.EnvironmentApprovalSnapshot{
				azuredevops.EnvironmentApprovalSnapshot{
					Approvals: []azuredevops.User{sampleUser},
					ApprovalOptions: map[string]string{
						"requiredApproverCount":       "0",
						"releaseCreatorCanBeApprover": "true",
					},
				},
			},
			Rank:                    1,
			DefinitionEnvironmentID: 1,
			QueueID:                 1,
			EnvironmentOptions: map[string]string{
				"emailNotificationType": "OnlyOnFailure",
				"emailRecipients":       "release.environment.owner;release.creator",
				"skipArtifactsDownload": "false",
				"timeoutInMinutes":      "0",
				"enableAccessToken":     "false",
			},
			ModifiedOn: sampleTime,
			WorkflowTasks: []azuredevops.EnvironmentWorkflowTasks{
				azuredevops.EnvironmentWorkflowTasks{
					TaskID:           "SampleTaskID",
					Version:          "1.0",
					Name:             "Sample Task",
					Enabled:          true,
					AlwaysRun:        true,
					ContinueOnError:  false,
					TimeoutInMinutes: 10,
					DefinitionType:   strPtr("xaml"),
					Inputs: map[string]string{
						"ConnectedServiceName": "b460b0f8-fe23-4dc2-a99c-fd8b0633fe1c",
						"WebSiteName":          "$(webAppName)",
						"WebSiteLocation":      "Southeast Asia",
						"Slot":                 "",
						"Package":              "$(System.DefaultWorkingDirectory)\\**\\*.zip",
					},
				},
			},
		}

		return azuredevops.ServiceHook{
			ID:          "SampleID",
			PublisherID: "SamplePublisherID",
			Scope:       "all",
			Message: azuredevops.ServiceHookMessage{
				Text:     "Sample Message",
				HTML:     "<strong>Sample Message</strong>",
				Markdown: "Sample Message",
			},
			DetailedMessage: azuredevops.ServiceHookMessage{
				Text:     "Sample Detailed Message",
				HTML:     "<strong>Sample Detailed Message</strong>",
				Markdown: "Sample Detailed Message",
			},
			EventType: string(azuredevops.ServiceHookEventTypeBuildComplete),
			Resource: azuredevops.ServiceHookResource{
				IntDefinition: azuredevops.IntDefinition{ID: 1, Name: "Sample Resource"},
				URL:           strPtr("https://dev.azure.com/SampleOrganization/SampleProject/_apis/sample/SampleID"),
				Reason:        strPtr(string(azuredevops.ReasonContinuousIntegration)),
				Status:        strPtr(string(azuredevops.StatusSucceeded)),

				ServiceHookResourceBuildComplete: azuredevops.ServiceHookResourceBuildComplete{
					URI:          strPtr("vstfs:///Build/Build/2"),
					BuildNumber:  strPtr("190813.01"),
					StartTime:    &sampleTime,
					FinishTime:   &sampleTime,
					DropLocation: strPtr("#/1/drop"),
					Drop: &azuredevops.ServiceHookResourceBuildDrop{
						Location:    "#/1/drop",
						Type:        "container",
						URL:         "https://dev.azure.com/SampleOrganization/SampleProject/_apis/resources/Containers/3/drop",
						DownloadURL: "https://dev.azure.com/SampleOrganization/SampleProject/_apis/resources/Containers/3/drop?api-version=1.0&$format=zip&downloadFileName=ConsumerAddressModule_20150407.1_drop",
					},
					Log: &azuredevops.ServiceHookResourceBuildLog{
						Type:        "container",
						URL:         "https://dev.azure.com/SampleOrganization/SampleProject/_apis/resources/Containers/3/logs",
						DownloadURL: "https://dev.azure.com/SampleOrganization/SampleProject/_apis/resources/Containers/3/logs?api-version=1.0&$format=zip&downloadFileName=ConsumerAddressModule_20150407.1_logs",
					},
					SourceGetVersion:   strPtr("LG:refs/heads/master:aaabcbcajlasjksakhdsakdhsdhkas"),
					LastChangedBy:      &sampleUser,
					RetainIndefinitely: &sampleBool,
					HasDiagnostics:     &sampleBool,
					Definition: &azuredevops.ServiceHookResourceBuildDefinition{
						ServiceHookResourceDefinition: azuredevops.ServiceHookResourceDefinition{
							IntDefinition: azuredevops.IntDefinition{ID: 1, Name: "Sample Build"},
							URL:           "https://dev.azure.com/SampleOrganization/SampleProject/_apis/build-release/Definitions/1",
						},
						BatchSize:      1,
						TriggerType:    string(azuredevops.ServiceHookResourceBuildDefinitionTriggerTypeNone),
						DefinitionType: string(azuredevops.ServiceHookResourceBuildDefinitionTypeXAML),
					},
					Queue: &azuredevops.ServiceHookResourceBuildQueue{
						ServiceHookResourceDefinition: azuredevops.ServiceHookResourceDefinition{
							IntDefinition: azuredevops.IntDefinition{ID: 1, Name: "Sample Queue"},
							URL:           "https://dev.azure.com/SampleOrganization/SampleProject/_apis/build-release/Queues/1",
						},
						QueueType: string(azuredevops.ServiceHookResourceBuildQueueTypeBuildController),
					},
				},

				ServiceHookResourceCodeCheckedIn: azuredevops.ServiceHookResourceCodeCheckedIn{
					ChangesetID: intPtr(1),
					Author:      &sampleUser,
					CheckedInBy: &sampleUser,
					CreatedDate: &sampleTime,
					Comment:     strPtr("Sample Comment"),
				},

				ServiceHookResourceCodePushed: azuredevops.ServiceHookResourceCodePushed{
					RefUpdates: []azuredevops.GitRefUpdate{
						azuredevops.GitRefUpdate{
							Name:        "Sample Ref Update",
							OldObjectID: "ajasjkdkdkadjkdakjdajkdbajksdlka",
							NewObjectID: "laskllsajlaskdnaslkdnalnadksjdaa",
						},
					},
					PushedBy: &sampleUser,
					PushID:   intPtr(1),
					Date:     &sampleTime,
				},

				ServiceHookResourcePullRequest: azuredevops.ServiceHookResourcePullRequest{
					PullRequestID: intPtr(1),
					CreatedBy:     &sampleUser,
					Title:         strPtr("Sample Pull Request"),
					Description:   strPtr("2 commits"),
					SourceRefName: strPtr("refs/heads/feature/Sample"),
					TargetRefName: strPtr("refs/heads/feature/master"),
				},

				ServiceHookResourceWorkItems: azuredevops.ServiceHookResourceWorkItems{
					Rev: intPtr(1),
				},

				ServiceHookResourceWorkItemsUpdated: azuredevops.ServiceHookResourceWorkItemsUpdated{
					WorkItemID:  intPtr(1),
					RevisedBy:   &sampleUser,
					RevisedDate: &sampleTime,
					Revision: &azuredevops.ServiceHookResourceWorkItemsUpdatedRevision{
						Rev: intPtr(1),
						Fields: map[string]string{
							"System.AreaPath":      "SampleProject",
							"System.TeamProject":   "SampleProject",
							"System.IterationPath": "SampleProject\\Release 1\\Sprint 1",
							"System.WorkItemType":  "Bug",
							"System.State":         "New",
							"System.Reason":        "Sample Reason",
							"System.CreatedDate":   sampleTime.String(),
							"System.CreatedBy":     sampleUser.DisplayName,
							"System.ChangedDate":   sampleTime.String(),
							"System.ChangedBy":     sampleUser.DisplayName,
							"System.Title":         "Sample Title",
							"Microsoft.Azure DevOps Services.Common.Severity": "3 - Medium",
						},
						ID:  1,
						URL: "https://dev.azure.com/SampleOrganization/SampleProject/_apis/wit/workItems/1/revisions/1",
					},
				},

				Approval: &azuredevops.Approval{
					ID:               1,
					Revision:         1,
					Approver:         sampleUser,
					ApprovedBy:       sampleUser,
					ApprovalType:     azuredevops.ApprovalTypePostDeploy,
					CreatedOn:        sampleTime,
					ModifiedOn:       sampleTime,
					Status:           azuredevops.ApprovalStatusApproved,
					Comments:         "Sample Comment",
					IsAutomated:      true,
					IsNotificationOn: false,
					TrialNumber:      1,
					Attempt:          1,
					Rank:             1,
					Release:          azuredevops.IntDefinition{ID: 1, Name: "Sample Release"},
					ReleaseDefinition: azuredevops.ServiceHookResourceDefinition{
						IntDefinition: azuredevops.IntDefinition{ID: 1, Name: "Sample Release"},
						URL:           "https://dev.azure.com/SampleOrganization/SampleProject/_apis/Release/definitions/1",
					},
				},

				Release: &azuredevops.Release{
					IntDefinition: azuredevops.IntDefinition{ID: 1, Name: "Sample Release"},
					ReleaseID:     1,
					Status:        string(azuredevops.StatusSucceeded),
					CreatedOn:     sampleTime,
					ModifiedOn:    sampleTime,
					ModifiedBy:    sampleUser,
					CreatedBy:     sampleUser,
					Environments:  []azuredevops.Environment{sampleEnvironment},
					Variables: map[string]string{
						"Configuration": "Release",
					},
					Artifacts: []azuredevops.ReleaseArtifact{
						azuredevops.ReleaseArtifact{
							SourceID: "SampleID",
							Type:     "Build",
							Alias:    "Sample.CI",
							DefinitionReference: map[string]azuredevops.StrDefinition{
								"Definition": azuredevops.StrDefinition{ID: "SampleID", Name: "Sample Artifact"},
							},
							IsPrimary: true,
						},
					},
					ReleaseDefinition: azuredevops.ServiceHookResourceDefinition{
						IntDefinition: azuredevops.IntDefinition{ID: 1, Name: "Sample Release"},
						URL:           "https://dev.azure.com/SampleOrganization/SampleProject/_apis/Release/definitions/1",
					},
					Description:                "Sample Release Description",
					Reason:                     string(azuredevops.ReasonContinuousIntegration),
					ReleaseNameFormat:          "$(yyMMdd)$(.rr)",
					KeepForever:                false,
					DefinitionSnapshotRevision: 1,
					Comment:                    "Sample Release Coment",
					LogsContainerURL:           strPtr("https://dev.azure.com/SampleOrganization/SampleProject/_apis/Release/sample/logs/1"),
				},

				Environment: &sampleEnvironment,

				Project: &azuredevops.StrDefinition{ID: "SampleProjectId", Name: "SampleProject"},

				Repository: &azuredevops.GitRepository{
					StrDefinition: azuredevops.StrDefinition{ID: "SampleRepositoryId", Name: "SampleRepository"},
					URL:           "https://dev.azure.com/SampleOrganization/SampleProject/_apis/repos/git/repositories/SampleRepositoryId",
					Project: azuredevops.GitProject{
						StrDefinition: azuredevops.StrDefinition{ID: "SampleProjectId", Name: "SampleProject"},
						URL:           "https://dev.azure.com/SampleOrganization/SampleProject/_apis/projects/SampleProjectId",
						State:         "wellFormed",
					},
				},

				Commits: []azuredevops.GitCommit{
					azuredevops.GitCommit{
						CommitID: "SampleCommitId1",
						URL:      "https://dev.azure.com/SampleOrganization/SampleProject/_apis/repos/git/repositories/SampleRepositoryId/commits/SampleCommitId1",
					},
					azuredevops.GitCommit{
						CommitID: "SampleCommitId2",
						URL:      "https://dev.azure.com/SampleOrganization/SampleProject/_apis/repos/git/repositories/SampleRepositoryId/commits/SampleCommitId2",
					},
				},
			},
			ResourceVersion: "1.0",
			ResourceContainers: azuredevops.ServiceHookResourceContainers{
				Collection: azuredevops.ServiceHookResourceContainer{
					ID: "SampleID",
				},
				Account: azuredevops.ServiceHookResourceContainer{
					ID: "SampleID",
				},
				Project: azuredevops.ServiceHookResourceContainer{
					ID: "SampleProjectId",
				},
			},
		}
	}()

	sampleTemplatingArgs = templating.NewArgsFromServiceHook(sampleServiceHook)
)
