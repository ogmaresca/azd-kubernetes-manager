package azuredevops

import "time"

// Approval is the defintion of release definitions
type Approval struct {
	ID                 int                           `json:"id"`
	Revision           int                           `json:"revision"`
	Approver           User                          `json:"approver"`
	ApprovedBy         User                          `json:"approvedBy"`
	ApprovalType       ApprovalType                  `json:"approvalType"`
	CreatedOn          time.Time                     `json:"createdOn"`
	ModifiedOn         time.Time                     `json:"modifiedOn"`
	Status             ApprovalStatus                `json:"status"`
	Comments           string                        `json:"comments"`
	IsAutomated        bool                          `json:"isAutomated"`
	IsNotificationOn   bool                          `json:"isNotificationOn"`
	TrialNumber        int                           `json:"trialNmber"`
	Attempt            int                           `json:"attempt"`
	Rank               int                           `json:"rank"`
	Release            IntDefinition                 `json:"release"`
	ReleaseDefinition  ServiceHookResourceDefinition `json:"releaseDefinition"`
	ReleaseEnvironment IntDefinition                 `json:"releaseEnvironment"`
}

// ApprovalType are the types of approvals
type ApprovalType string

const (
	// ApprovalTypePreDeploy are approvals before an approval
	ApprovalTypePreDeploy ApprovalType = "preDeploy"
	// ApprovalTypePostDeploy are approvals after an approval
	ApprovalTypePostDeploy ApprovalType = "postDeploy"
)

// ApprovalStatus holds all of the different approval statuses
type ApprovalStatus string

const (
	// ApprovalStatusApproved is for resources that have been approved
	ApprovalStatusApproved ApprovalStatus = "approved"
	// ApprovalStatusRejected is for resources that have been rejected
	ApprovalStatusRejected ApprovalStatus = "rejected"
	// ApprovalStatusPending is for resources that are currently pending
	ApprovalStatusPending ApprovalStatus = "pending"
)
