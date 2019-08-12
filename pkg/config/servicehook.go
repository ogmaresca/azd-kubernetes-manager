package config

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/alexcesaro/log/stdlog"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/azuredevops"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/templating"
)

var (
	baseTemplate = templating.New("ServiceHookMatches")
	logger       = stdlog.GetFromFlags()
)

// ServiceHook holds rules for a Service Hook
type ServiceHook struct {
	// The EventType
	// See https://docs.microsoft.com/en-us/azure/devops/service-hooks/events?toc=%2Fazure%2Fdevops%2Fintegrate%2Ftothis.json&bc=%2Fazure%2Fdevops%2Fintegrate%2Fbreadcrumb%2Ftothis.json&view=azure-devops
	// for the event type
	Event ServiceHookEventType `yaml:"event"`

	// Resource filters
	ResourceFilters ServiceHookResourceFilters `yaml:"resourceFilters"`

	// If Continue is true, process any other service hooks that match the Service Hook
	Continue bool `yaml:"continue"`

	// The rules to perform on the Service Hook
	Rules Rules `yaml:"rules"`
}

// ServiceHookResourceFilters holds filters for a Service Hook Resource
type ServiceHookResourceFilters struct {
	// Only execute the rules if {resource.status ∈ this.Statuses}
	Statuses []string `yaml:"statuses"`

	// Only execute the rules if {resource.reason ∈ this.Reasons}
	Reasons []string `yaml:"reasons"`

	// Only execute the rules if {resource.project.name ∈ this.Projects} || {resource.release.project.name ∈ this.Projects}
	Projects []string `yaml:"projects"`

	// Only execute the rules if {resource.release.releaseDefinition.name ∈ this.Releases}
	Releases []string `yaml:"releases"`

	// Only execute the rules if {resource.release.environments[].name ∈ this.Environment}
	Environments []string `yaml:"environments"`

	// Only execute the rules if {resource.approval.type ∈ this.ApprovalTypes}
	ApprovalTypes []azuredevops.ApprovalType `yaml:"approvalTypes"`

	// Only execute the rules if {resource.repository.name ∈ this.Repositories}
	Repositories []string `yaml:"repositories"`

	// Only execute the rules {resource.SourceRefName ∈ this.SourceRefs}
	SourceRefs []string `yaml:"sourceRefs"`

	// Only execute the rules {resource.TargetRefName ∈ this.TargetRefs}
	TargetRefs []string `yaml:"targetRefs"`

	// Filters that are Go templated with the ServiceHook Resource and must return true
	Templates []string `yaml:"templates"`
}

///
/// Describe()
///

// Describe returns a user-friendly representation of a ServiceHook
func (sh ServiceHook) Describe() string {
	return fmt.Sprintf(
		"Event Type: %s\nResource Filters:\n  %s\nContinue: %t\nRules:\n  %s",
		sh.Event, strings.ReplaceAll(sh.ResourceFilters.Describe(), "\n", "\n  "), sh.Continue, strings.ReplaceAll(sh.Rules.Describe(), "\n", "\n  "),
	)
}

// Describe returns a user-friendly representation of a ServiceHook
func (shf ServiceHookResourceFilters) Describe() string {
	return fmt.Sprintf(
		strings.Join([]string{
			"Statuses: %v",
			"Reasons: %v",
			"Projects: %v",
			"Releases: %v",
			"Environments: %v",
			"Approval Types: %v",
			"Repositories: %v",
			"Source Branch Refs: %v",
			"Target Branch Refs: %v",
			"Templates: %s",
		}, "\n"),
		shf.Statuses,
		shf.Reasons,
		shf.Projects,
		shf.Releases,
		shf.Environments,
		shf.ApprovalTypes,
		shf.Repositories,
		shf.SourceRefs,
		shf.TargetRefs,
		joinYAMLSlice(shf.Templates),
	)
}

///
/// Validate()
///

// Validate a Service Hook rule definition. This function returns a slice of warnings and an error.
func (sh ServiceHook) Validate() ([]string, error) {
	var errors []string
	var warnings []string

	if sh.Event == "" {
		errors = append(errors, "The `event` field must be defined.")
	}

	resourceFiltersWarnings, err := sh.ResourceFilters.Validate()
	if len(resourceFiltersWarnings) > 0 {
		warnings = append(warnings, resourceFiltersWarnings...)
	}
	if err != nil {
		errors = append(errors, err.Error())
	}

	rulesWarnings, err := sh.Rules.Validate()
	if len(rulesWarnings) > 0 {
		warnings = append(warnings, rulesWarnings...)
	}
	if err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		err = fmt.Errorf("%s", joinYAMLSlice(errors))
	}

	return rulesWarnings, err
}

// Validate a Service Hook filters definition. This function returns a slice of warnings and an error.
func (shf ServiceHookResourceFilters) Validate() ([]string, error) {
	var errors []string

	gotplRegex, _ := regexp.Compile(".*{{.*}}.*")

	for pos, ref := range shf.SourceRefs {
		_, err := regexp.CompilePOSIX(ref)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error with `SourceRefs` %d: invalid POSIX ERE pattern.", pos))
		}
	}

	for pos, ref := range shf.TargetRefs {
		_, err := regexp.CompilePOSIX(ref)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error with `TargetRefs` %d: invalid POSIX ERE pattern.", pos))
		}
	}

	for pos, templateFilter := range shf.Templates {
		if !gotplRegex.Match([]byte(templateFilter)) {
			errors = append(errors, fmt.Sprintf("Error with template filter %d: invalid templating. Please see https://golang.org/pkg/text/template/.", pos))
		}
	}

	var err error
	if len(errors) > 0 {
		err = fmt.Errorf("%s", joinYAMLSlice(errors))
	}

	return []string{}, err
}

///
/// Other types and methods
///

// Matches determines if a rule should be applied for a Service Hook
func (sh ServiceHook) Matches(serviceHook azuredevops.ServiceHook) (bool, error) {
	status := serviceHook.GetStatus()
	if status != nil && !contains(*status, sh.ResourceFilters.Statuses) {
		return false, nil
	}

	reason := serviceHook.GetReason()
	if reason != nil && !contains(*reason, sh.ResourceFilters.Statuses) {
		return false, nil
	}

	project := serviceHook.GetProjectName()
	if project != nil && !contains(*project, sh.ResourceFilters.Projects) {
		return false, nil
	}

	if serviceHook.Resource.Release != nil {
		if !contains(serviceHook.Resource.Release.ReleaseDefinition.Name, sh.ResourceFilters.Releases) {
			return false, nil
		}
	}

	environments := serviceHook.GetEnvironments()
	if project != nil && !intersection(environments, sh.ResourceFilters.Environments) {
		return false, nil
	}

	if serviceHook.Resource.Approval != nil {
		var approvalTypeFilters []string
		for _, approvalTypeFilter := range sh.ResourceFilters.ApprovalTypes {
			approvalTypeFilters = append(approvalTypeFilters, string(approvalTypeFilter))
		}
		if !contains(string(serviceHook.Resource.Approval.ApprovalType), approvalTypeFilters) {
			return false, nil
		}
	}

	if serviceHook.Resource.Repository != nil {
		if !contains(serviceHook.Resource.Repository.Name, sh.ResourceFilters.Repositories) {
			return false, nil
		}
	}

	if serviceHook.Resource.SourceRefName != nil {
		matches, err := containsPOSIXERE(*serviceHook.Resource.SourceRefName, sh.ResourceFilters.SourceRefs)
		if !matches {
			return false, err
		}
	}

	if serviceHook.Resource.TargetRefName != nil {
		matches, err := containsPOSIXERE(*serviceHook.Resource.TargetRefName, sh.ResourceFilters.TargetRefs)
		if !matches {
			return false, err
		}
	}

	if len(sh.ResourceFilters.Templates) > 0 {
		for pos, filter := range sh.ResourceFilters.Templates {
			template, err := baseTemplate.Clone()
			if err != nil {
				return false, err
			}
			template, err = template.Parse(filter)
			if err != nil {
				return false, fmt.Errorf("Error parsing Service Hook template filter %d: %s", pos, err.Error())
			}
			buffer := new(bytes.Buffer)
			err = template.Execute(buffer, serviceHook.Resource)
			if err != nil {
				return false, fmt.Errorf("Error executing Service Hook template filter %d: %s", pos, err.Error())
			}
			templatedValue := buffer.String()
			if logger.LogDebug() {
				logger.Debugf("Converted Service Hook filter template %d:\n  %s\nto:\n  %s", pos, strings.ReplaceAll(filter, "\n", "\n  "), strings.ReplaceAll(templatedValue, "\n", "\n  "))
			}
			if !strings.EqualFold("true", templatedValue) {
				return false, nil
			}
		}
	}

	return true, nil
}

// ServiceHookEventType represents all possible Event Type values for a Service Hook configuration
type ServiceHookEventType string

const (
	// ServiceHookEventTypeBuildComplete represents the Build Completed event
	ServiceHookEventTypeBuildComplete ServiceHookEventType = "build.complete"
	// ServiceHookEventTypeBuilds is a sum type for:
	// - build.complete
	ServiceHookEventTypeBuilds ServiceHookEventType = "Builds"

	// ServiceHookEventTypeReleaseAbandoned represents the Release Abandoned event
	ServiceHookEventTypeReleaseAbandoned ServiceHookEventType = "ms.vss-release.release-abandoned-event"
	// ServiceHookEventTypeReleaseCreated represents the Release Created event
	ServiceHookEventTypeReleaseCreated ServiceHookEventType = "ms.vss-release.release-created-event"
	// ServiceHookEventTypeReleases is a sum type for:
	// - ms.vss-release.release-abandoned-event
	// - ms.vss-release.release-created-event
	ServiceHookEventTypeReleases ServiceHookEventType = "Releases"

	// ServiceHookEventTypeReleaseDeploymentApprovalCompleted represents the Release Deployment Approval Completed event
	ServiceHookEventTypeReleaseDeploymentApprovalCompleted ServiceHookEventType = "ms.vss-release.deployment-approval-completed-event"
	// ServiceHookEventTypeReleaseDeploymentApprovalPending represents the Release Deployment Approval Pending event
	ServiceHookEventTypeReleaseDeploymentApprovalPending ServiceHookEventType = "ms.vss-release.deployment-approval-pending-event"
	// ServiceHookEventTypeReleaseDeploymentApprovals is a sum type for:
	// - ms.vss-release.deployment-approval-completed-event
	// - ms.vss-release.deployment-approval-pending-event
	ServiceHookEventTypeReleaseDeploymentApprovals ServiceHookEventType = "Release Deployment Approals"

	// ServiceHookEventTypeReleaseDeploymentCompleted represents the Release Deployment Completed event
	ServiceHookEventTypeReleaseDeploymentCompleted ServiceHookEventType = "ms.vss-release.deployment-completed-event"
	// ServiceHookEventTypeReleaseDeploymentStarted represents the Release Deployment Started event
	ServiceHookEventTypeReleaseDeploymentStarted ServiceHookEventType = "ms.vss-release.deployment-started-event"
	// ServiceHookEventTypeReleaseDeployments is a sum type for:
	// - ms.vss-release.deployment-completed-event
	// - ms.vss-release.deployment-started-event
	ServiceHookEventTypeReleaseDeployments ServiceHookEventType = "Release Deployments"

	// ServiceHookEventTypeCodeCheckedIn represents the Code Checked In event
	ServiceHookEventTypeCodeCheckedIn ServiceHookEventType = "tfvc.checkin"
	// ServiceHookEventTypeCodePushed represents the Code Pushed event
	ServiceHookEventTypeCodePushed ServiceHookEventType = "git.push"
	// ServiceHookEventTypeCode is a sum type for:
	// - tfvc.checkin
	// - git.push
	ServiceHookEventTypeCode ServiceHookEventType = "Code"

	// ServiceHookEventTypePullRequestCreated represents the Pull Request Created event
	ServiceHookEventTypePullRequestCreated ServiceHookEventType = "git.pullrequest.created"
	// ServiceHookEventTypePullRequestMerged represents the Pull Request Merged event
	ServiceHookEventTypePullRequestMerged ServiceHookEventType = "git.pullrequest.merged"
	// ServiceHookEventTypePullRequestUpdated represents the Pull Request Updated event
	ServiceHookEventTypePullRequestUpdated ServiceHookEventType = "git.pullrequest.updated"
	// ServiceHookEventTypePullRequests is a sum type for:
	// - git.pullrequest.created
	// - git.pullrequest.merged
	// - git.pullrequest.updated
	ServiceHookEventTypePullRequests ServiceHookEventType = "Pull Requests"

	// ServiceHookEventTypeWorkItemCreated represents the Work Item Created event
	ServiceHookEventTypeWorkItemCreated ServiceHookEventType = "workitem.created"
	// ServiceHookEventTypeWorkItemDeleted represents the Work Item Deleted event
	ServiceHookEventTypeWorkItemDeleted ServiceHookEventType = "workitem.deleted"
	// ServiceHookEventTypeWorkItemRestored represents the Work Item Restored event
	ServiceHookEventTypeWorkItemRestored ServiceHookEventType = "workitem.restored"
	// ServiceHookEventTypeWorkItemUpdated represents the Work Item Updated event
	ServiceHookEventTypeWorkItemUpdated ServiceHookEventType = "workitem.updated"
	// ServiceHookEventTypeWorkItemCommented represents the Work Item Commented event
	ServiceHookEventTypeWorkItemCommented ServiceHookEventType = "workitem.commented"
	// ServiceHookEventTypeWorkItems is a sum type for:
	// - workitem.created
	// - workitem.deleted
	// - workitem.restored
	// - workitem.updated
	// - workitem.commented
	ServiceHookEventTypeWorkItems ServiceHookEventType = "Work Items"
)

// GetEventTypes returns all of the Service Hook event types that this value represents
func (et ServiceHookEventType) GetEventTypes() []string {
	switch et {
	case ServiceHookEventTypeBuilds:
		return []string{
			string(ServiceHookEventTypeBuildComplete),
		}
	case ServiceHookEventTypeReleases:
		return []string{
			string(ServiceHookEventTypeReleaseAbandoned),
			string(ServiceHookEventTypeReleaseCreated),
		}
	case ServiceHookEventTypeReleaseDeploymentApprovals:
		return []string{
			string(ServiceHookEventTypeReleaseDeploymentApprovalCompleted),
			string(ServiceHookEventTypeReleaseDeploymentApprovalPending),
		}
	case ServiceHookEventTypeReleaseDeployments:
		return []string{
			string(ServiceHookEventTypeReleaseDeploymentCompleted),
			string(ServiceHookEventTypeReleaseDeploymentStarted),
		}
	case ServiceHookEventTypeCode:
		return []string{
			string(ServiceHookEventTypeCodeCheckedIn),
			string(ServiceHookEventTypeCodePushed),
		}
	case ServiceHookEventTypePullRequests:
		return []string{
			string(ServiceHookEventTypePullRequestCreated),
			string(ServiceHookEventTypePullRequestMerged),
			string(ServiceHookEventTypePullRequestUpdated),
		}
	case ServiceHookEventTypeWorkItems:
		return []string{
			string(ServiceHookEventTypeWorkItemCreated),
			string(ServiceHookEventTypeWorkItemDeleted),
			string(ServiceHookEventTypeWorkItemRestored),
			string(ServiceHookEventTypeWorkItemUpdated),
			string(ServiceHookEventTypeWorkItemCommented),
		}
	}

	return []string{string(et)}
}
