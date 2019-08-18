package templating

import (
	"time"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/azuredevops"
)

// Args holds all of the values available to go templating at runtime
type Args struct {
	EventType string

	BuildID *int

	BuildNumber *string

	PullRequestID *int

	ProjectName *string

	ResourceURL *string

	StartTime *time.Time

	FinishTime *time.Time

	ResourceName string

	Status *string

	Reason *string

	ServiceHook azuredevops.ServiceHook
}

// NewArgsFromServiceHook creates an Args from a Service Hook request
func NewArgsFromServiceHook(serviceHook azuredevops.ServiceHook) Args {
	var buildID *int
	var buildNumber *string
	if serviceHook.EventType == string(azuredevops.ServiceHookEventTypeBuildComplete) {
		buildID = &serviceHook.Resource.ID
		buildNumber = serviceHook.Resource.BuildNumber
	}

	return Args{
		EventType:     serviceHook.EventType,
		PullRequestID: serviceHook.Resource.PullRequestID,
		BuildID:       buildID,
		BuildNumber:   buildNumber,
		ProjectName:   serviceHook.GetProjectName(),
		ResourceURL:   serviceHook.Resource.URL,
		StartTime:     serviceHook.Resource.StartTime,
		FinishTime:    serviceHook.Resource.FinishTime,
		ResourceName:  serviceHook.Resource.Name,
		Status:        serviceHook.GetStatus(),
		Reason:        serviceHook.GetReason(),
		ServiceHook:   serviceHook,
	}
}
