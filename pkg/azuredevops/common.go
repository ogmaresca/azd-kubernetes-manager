package azuredevops

// IntDefinition is the base type for Azure Devops responses
type IntDefinition struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// StrDefinition is the base type for Azure Devops responses
type StrDefinition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// User fields in the Azure Devops API
type User struct {
	DisplayName string  `json:"displayName"`
	URL         *string `json:"url"`
	ID          string  `json:"id"`
	UniqueName  string  `json:"uniqueName"`
	ImageURL    *string `json:"imageUrl"`
	Descriptor  string  `json:"descriptor"`
}

// Error is returned when an error occurs in the API, such as an invalid ID being used.
type Error struct {
	//ID             string `json:"$id"`
	//InnerException Error  `json:"innerException"`
	Message   string `json:"message"`
	TypeName  string `json:"typeName"`
	TypeKey   string `json:"typeKey"`
	ErrorCode int    `json:"errorCode"`
	EventID   int    `json:"eventId"`
}

// Status are the common Azure Devops statuses
type Status string

const (
	// StatusSucceeded is for resources that have completed successfuly
	StatusSucceeded Status = "succeeded"
	// StatusAbandoned is for resources that have been abandoned
	StatusAbandoned Status = "abandoned"
	// StatusActive is for resources that are active
	StatusActive Status = "active"
	// StatusCompleted is for resources that are completed without a substatus
	StatusCompleted Status = "completed"
	// StatusQueued is for resources that are waiting to run
	StatusQueued Status = "queued"
	// StatusNotSet is for resources that don't have a real status yet
	StatusNotSet Status = "notSet"
)

// Reason are the common Azure Devops reasons
type Reason string

const (
	// ReasonManual is for manually triggered reasources
	ReasonManual Reason = "manual"
	// ReasonContinuousIntegration is for automatically triggered reasources in a CI process
	ReasonContinuousIntegration Reason = "continuousIntegration"
)

// PullRequestStatus lists Azure Devops statuses for Pull Requests
// https://docs.microsoft.com/en-us/rest/api/azure/devops/git/pull%20requests/get%20pull%20requests?view=azure-devops-rest-5.1#pullrequeststatus
type PullRequestStatus string

const (
	// PullRequestStatusCompleted is for completed
	PullRequestStatusCompleted Status = "completed"
	// PullRequestStatusAbandoned is for abandoned PRs
	PullRequestStatusAbandoned Status = "abandoned"
	// PullRequestStatusActive is for PRs that are currently open
	PullRequestStatusActive Status = "active"
	// PullRequestStatusNotSet is the default status of PRs
	PullRequestStatusNotSet Status = "notSet"
)
