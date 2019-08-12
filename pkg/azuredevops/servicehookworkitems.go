package azuredevops

import (
	"net/url"
	"time"
)

// ServiceHookResourceWorkItems holds fields related to the Service Hook resource for Work Item events
type ServiceHookResourceWorkItems struct {
	Rev    *int              `json:"rev"`
	Fields map[string]string `json:"fields"`
}

// ServiceHookResourceWorkItemsUpdated holds fields related to the Service Hook resource for the workitem.updated
type ServiceHookResourceWorkItemsUpdated struct {
	WorkItemID  *int                                         `json:"workItemId"`
	RevisedBy   *User                                        `json:"revisedBy"`
	RevisedDate *time.Time                                   `json:"revisedDate"`
	Revision    *ServiceHookResourceWorkItemsUpdatedRevision `json:"revision"`
}

// ServiceHookResourceWorkItemsUpdatedRevision holds revision info for a Work Item
type ServiceHookResourceWorkItemsUpdatedRevision struct {
	ServiceHookResourceWorkItems
	ID  int     `json:"id"`
	URL url.URL `json:"url"`
}
