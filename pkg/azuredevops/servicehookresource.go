package azuredevops

import (
	"net/url"
)

// ServiceHookResource holds fields related to the Service Hook resource
// Release event types:
// - ms.vss-release.release-abandoned-event
// - ms.vss-release.release-created-event
// Release Deployment Approval event types:
// - ms.vss-release.deployment-approval-completed-event
// - ms.vss-release.deployment-approval-pending-event
// Release Deployment event types:
// - ms.vss-release.deployment-completed-event
// - ms.vss-release.deployment-started-event
// Code event types:
// - tfvc.checkin
// - git.push
// Pull Request types:
// - git.pullrequest.created
// - git.pullrequest.merged
// - git.pullrequest.updated
// Work Item types:
// - workitem.created
// - workitem.deleted
// - workitem.restored
// - workitem.updated
// - workitem.commented
type ServiceHookResource struct {
	IntDefinition // Event Types: build.complete (ID), Work Items
	ServiceHookResourceBuildComplete
	ServiceHookResourceCodeCheckedIn
	ServiceHookResourceCodePushed
	ServiceHookResourcePullRequest
	ServiceHookResourceWorkItems
	ServiceHookResourceWorkItemsUpdated

	URL    *url.URL `json:"url"`    // Event types: build.complete, Code, Pull Requests, Work Items
	Reason *string  `json:"reason"` // Event types: build.complete, Releases, Release Deployment Approvals
	Status *string  `json:"status"` // Event types: build.complete, Pull Requests

	Approval    *Approval      `json:"approval"`    // Event types: Approvals
	Release     *Release       `json:"release"`     // Event types: Releases, Release Deployment Approvals
	Environment *Environment   `json:"environment"` // Event types: Release Deployments
	Project     *StrDefinition `json:"project"`     // Event types: Releases, Release Deployment Approvals, Release Deployments
	Repository  *GitRepository `json:"repository"`  // Event types: git.push, Pull Requests
}

// ServiceHookResourceContainers contains the value of the resourceContainers field
type ServiceHookResourceContainers struct {
	Collection ServiceHookResourceContainer `json:"collection"`
	Account    ServiceHookResourceContainer `json:"account"`
	Project    ServiceHookResourceContainer `json:"project"`
}

// ServiceHookResourceContainer contains a Resouce Container defintiion
type ServiceHookResourceContainer struct {
	ID string `json:"id"`
}

// ServiceHookResourceDefinition holds a basic resource definition for Service Hooks
type ServiceHookResourceDefinition struct {
	IntDefinition
	URL string `json:"url"`
}
