package azuredevops

import (
	"fmt"
	"time"
)

// ServiceHook defines the JSON body of service hooks
// https://docs.microsoft.com/en-us/azure/devops/service-hooks/events?toc=%2Fazure%2Fdevops%2Fintegrate%2Ftoc.json&bc=%2Fazure%2Fdevops%2Fintegrate%2Fbreadcrumb%2Ftoc.json&view=azure-devops
type ServiceHook struct {
	ID                 string                        `json:"id"`
	EventType          string                        `json:"eventType"`
	PublisherID        string                        `json:"publisherId"`
	Scope              string                        `json:"scope"`
	Message            ServiceHookMessage            `json:"message"`
	DetailedMessage    ServiceHookMessage            `json:"detailedMessage"`
	Resource           ServiceHookResource           `json:"resource"`
	ResourceVersion    string                        `json:"resourceVersion"`
	ResourceContainers ServiceHookResourceContainers `json:"resourceContainers"`
	CreatedDate        time.Time                     `json:"createdDate"`
}

// Describe returns a user-friendly descriptor for this Service Hook
func (sh ServiceHook) Describe() string {
	return fmt.Sprintf("(%s %s) %s", sh.EventType, sh.ID, sh.DetailedMessage.Text)
}

// GetStatus returns the status of the project from the Service Hook
func (sh ServiceHook) GetStatus() *string {
	if sh.Resource.Status != nil {
		// build.complete
		// git.pullrequest.created
		// git.pullrequest.merged
		// git.pullrequest.updated
		return sh.Resource.Status
	}

	if sh.Resource.Approval != nil {
		// ms.vss-release.deployment-approval-completed-event
		// ms.vss-release.deployment-approval-pending-event
		status := string(sh.Resource.Approval.Status)
		return &status
	}

	if sh.Resource.Release != nil && sh.Resource.Release.Status != "" {
		// ms.vss-release.release-abandoned-event
		// ms.vss-release.release-created-event
		return &sh.Resource.Release.Status
	}

	if sh.Resource.Environment != nil && sh.Resource.Environment.Status != "" {
		// ms.vss-release.deployment-completed-event
		// ms.vss-release.deployment-started-event
		return &sh.Resource.Environment.Status
	}

	return nil
}

// GetReason returns the reason of the project from the Service Hook
func (sh ServiceHook) GetReason() *string {
	if sh.EventType == string(ServiceHookEventTypeBuildComplete) {
		// build.complete
		return sh.Resource.Reason
	}

	if sh.Resource.Release != nil && sh.Resource.Release.Reason != "" {
		// ms.vss-release.release-abandoned-event
		// ms.vss-release.release-created-event
		// ms.vss-release.deployment-approval-completed-event
		// ms.vss-release.deployment-approval-pending-event
		return &sh.Resource.Release.Reason
	}

	return nil
}

// GetProjectName returns the name of the project from the Service Hook
func (sh ServiceHook) GetProjectName() *string {
	if sh.Resource.Project != nil && sh.Resource.Project.Name != "" {
		// ms.vss-release.release-abandoned-event
		// ms.vss-release.release-created-event
		// ms.vss-release.deployment-approval-completed-event
		// ms.vss-release.deployment-approval-pending-event
		// ms.vss-release.deployment-completed-event
		// ms.vss-release.deployment-started-event
		return &sh.Resource.Project.Name
	}

	if sh.Resource.Repository != nil && sh.Resource.Repository.Project.Name != "" {
		// git.push
		// git.pullrequest.created
		// git.pullrequest.merged
		// git.pullrequest.updated
		return &sh.Resource.Repository.Project.Name
	}

	if sh.Resource.URL != nil {
		// build.complete
		// tfvc.checkin
		// git.push
		// git.pullrequest.created
		// git.pullrequest.merged
		// git.pullrequest.updated
		return GetProjectFromURL(sh.Resource.URL)
	}

	return nil
}

// GetEnvironments returns all of the environment names attached to this ServiceHook
func (sh ServiceHook) GetEnvironments() []string {
	if sh.Resource.Environment != nil {
		// ms.vss-release.deployment-completed-event
		// ms.vss-release.deployment-started-event
		return []string{sh.Resource.Environment.Name}
	}

	if sh.Resource.Release != nil {
		// ms.vss-release.release-abandoned-event
		// ms.vss-release.release-created-event
		// ms.vss-release.deployment-approval-completed-event
		// ms.vss-release.deployment-approval-pending-event
		var environments []string
		for _, environment := range sh.Resource.Release.Environments {
			environments = append(environments, environment.Name)
		}
		return environments
	}

	return nil
}

// ServiceHookMessage is used to hold the message and defaultMessage fields in service hooks
type ServiceHookMessage struct {
	Text     string `json:"text"`
	HTML     string `json:"html"`
	Markdown string `json:"markdown"`
}

// ServiceHookEventType represents all possible Event Type values for a Service Hook configuration
type ServiceHookEventType string

const (
	// ServiceHookEventTypeBuildComplete represents the Build Completed event
	ServiceHookEventTypeBuildComplete ServiceHookEventType = "build.complete"

	// ServiceHookEventTypeReleaseAbandoned represents the Release Abandoned event
	ServiceHookEventTypeReleaseAbandoned ServiceHookEventType = "ms.vss-release.release-abandoned-event"
	// ServiceHookEventTypeReleaseCreated represents the Release Created event
	ServiceHookEventTypeReleaseCreated ServiceHookEventType = "ms.vss-release.release-created-event"

	// ServiceHookEventTypeReleaseDeploymentApprovalCompleted represents the Release Deployment Approval Completed event
	ServiceHookEventTypeReleaseDeploymentApprovalCompleted ServiceHookEventType = "ms.vss-release.deployment-approval-completed-event"
	// ServiceHookEventTypeReleaseDeploymentApprovalPending represents the Release Deployment Approval Pending event
	ServiceHookEventTypeReleaseDeploymentApprovalPending ServiceHookEventType = "ms.vss-release.deployment-approval-pending-event"

	// ServiceHookEventTypeReleaseDeploymentCompleted represents the Release Deployment Completed event
	ServiceHookEventTypeReleaseDeploymentCompleted ServiceHookEventType = "ms.vss-release.deployment-completed-event"
	// ServiceHookEventTypeReleaseDeploymentStarted represents the Release Deployment Started event
	ServiceHookEventTypeReleaseDeploymentStarted ServiceHookEventType = "ms.vss-release.deployment-started-event"

	// ServiceHookEventTypeCodeCheckedIn represents the Code Checked In event
	ServiceHookEventTypeCodeCheckedIn ServiceHookEventType = "tfvc.checkin"
	// ServiceHookEventTypeCodePushed represents the Code Pushed event
	ServiceHookEventTypeCodePushed ServiceHookEventType = "git.push"

	// ServiceHookEventTypePullRequestCreated represents the Pull Request Created event
	ServiceHookEventTypePullRequestCreated ServiceHookEventType = "git.pullrequest.created"
	// ServiceHookEventTypePullRequestMerged represents the Pull Request Merged event
	ServiceHookEventTypePullRequestMerged ServiceHookEventType = "git.pullrequest.merged"
	// ServiceHookEventTypePullRequestUpdated represents the Pull Request Updated event
	ServiceHookEventTypePullRequestUpdated ServiceHookEventType = "git.pullrequest.updated"

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
)
