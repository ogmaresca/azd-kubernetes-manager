package azuredevops

import "time"

// Environment holds the environment of a Release
type Environment struct {
	IntDefinition
	ReleaseID               int                           `json:"releaseId"`
	Status                  string                        `json:"status"`
	Variables               map[string]string             `json:"variables"`
	PreDeployApprovals      []User                        `json:"preDeployApprovals"`
	PostDeployApprovals     []User                        `json:"postDeployApprovals"`
	PreApprovalsSnapshot    []EnvironmentApprovalSnapshot `json:"preApprovalsSnapshot"`
	PostApprovalsSnapshot   []EnvironmentApprovalSnapshot `json:"postApprovalsSnapshot"`
	DeploySteps             []interface{}                 `json:"deploySteps"`
	Rank                    int                           `json:"rank"`
	DefinitionEnvironmentID int                           `json:"definitionEnvironmentId"`
	QueueID                 int                           `json:"queueId"`
	EnvironmentOptions      map[string]string             `json:"environmentOptions"`
	Demands                 []interface{}                 `json:"demands"`
	Conditions              []interface{}                 `json:"conditions"`
	ModifiedOn              time.Time                     `json:"modifiedOn"`
	WorkflowTasks           []EnvironmentWorkflowTasks    `json:"workflowTasks"`
	DeployPhasesSnapshot    []interface{}                 `json:"deployPhasesSnapshot"`
	Owner                   User                          `json:"owner"`
	ScheduledDeploymentTime time.Time                     `json:"scheduledDeploymentTime"`
	Schedules               []interface{}                 `json:"schedules"`
	Release                 ServiceHookResourceDefinition `json:"release"`
}

// EnvironmentApprovalSnapshot holds the approval snapshot of a Release Environment
type EnvironmentApprovalSnapshot struct {
	Approvals       []User            `json:"approvals"`
	ApprovalOptions map[string]string `json:"approvalOptions"`
}

// EnvironmentWorkflowTasks holds the definition of pipeline environment tasks
type EnvironmentWorkflowTasks struct {
	TaskID           string            `json:"taskId"`
	Version          string            `json:"version"`
	Name             string            `json:"name"`
	Enabled          bool              `json:"enabled"`
	AlwaysRun        bool              `json:"alwaysRun"`
	ContinueOnError  bool              `json:"continueOnError"`
	TimeoutInMinutes int               `json:"timeoutInMinutes"`
	DefinitionType   *string           `json:"definitionType"`
	Inputs           map[string]string `json:"inputs"`
}
