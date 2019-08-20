package azuredevops

import (
	"time"
)

// ServiceHookResourceCodeCheckedIn holds fields related to the Service Hook resource for the tfvc.checkin event
type ServiceHookResourceCodeCheckedIn struct {
	ChangesetID *int       `json:"changesetID"`
	Author      *User      `json:"author"`
	CheckedInBy *User      `json:"checkedInBy"`
	CreatedDate *time.Time `json:"createdDate"`
	Comment     *string    `json:"comment"`
}

// ServiceHookResourceCodePushed holds fields related to the Service Hook resource for the git.push event
type ServiceHookResourceCodePushed struct {
	RefUpdates []GitRefUpdate `json:"refUpdates"`
	PushedBy   *User          `json:"pushedBy"`
	PushID     *int           `json:"pushId"`
	Date       *time.Time     `json:"time"`
}

// ServiceHookResourcePullRequest holds fields related to the Service Hook resource for Pull Request events
type ServiceHookResourcePullRequest struct {
	PullRequestID         *int            `json:"pullRequestId"`
	CreatedBy             *User           `json:"createdBy"`
	CreationDate          *time.Time      `json:"creationDate"`
	ClosedDate            *time.Time      `json:"closedDate"` // Event types: git.pullrequest.merged, git.pullrequest.updated
	Title                 *string         `json:"title"`
	Description           *string         `json:"description"`
	SourceRefName         *string         `json:"sourceRefName"`
	TargetRefName         *string         `json:"targetRefName"`
	MergeStatus           *string         `json:"mergeStatus"`
	MergeID               *string         `json:"mergeId"`
	LastMergeSourceCommit *GitMergeCommit `json:"lastMergeSourceCommit"`
	LastMergeTargetCommit *GitMergeCommit `json:"lastMergeTargetCommit"`
	LastMergeCommit       *GitMergeCommit `json:"lastMergeCommit"`
	Reviewers             []GitReviewer   `json:"reviewers"`
}

// GitCommit represents a single Git commit
type GitCommit struct {
	CommitID string  `json:"commitId"`
	Author   GitUser `json:"author"`
	Comitter GitUser `json:"committer"`
	Comment  string  `json:"comment"`
	URL      string  `json:"url"`
}

// GitUser represents a user in Git
type GitUser struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

// GitRefUpdate holds the old and new Git commit hashes
type GitRefUpdate struct {
	Name        string `json:"name"`
	OldObjectID string `json:"oldObjectId"`
	NewObjectID string `json:"newObjectId"`
}

// GitRepository holds information about a Git repository
type GitRepository struct {
	StrDefinition
	URL           string     `json:"url"`
	Project       GitProject `json:"project"`
	DefaultBranch string     `json:"defaultBranch"`
	RemoteURL     string     `json:"removeURL"`
}

// GitProject holds info about the Azure Devops project of a repository
type GitProject struct {
	StrDefinition
	URL   string `json:"url"`
	State string `json:"string"`
}

// GitMergeCommit holds info about a merge commit
type GitMergeCommit struct {
	CommitID string  `json:"commitId"`
	URL      *string `json:"url"`
}

// GitReviewer holds info about a Pull Request reviewer
type GitReviewer struct {
	User
	ReviewerURL *string `json:"reviewerUrl"`
	Vote        int     `json:"vote"`
	IsContainer bool    `json:"isContainer"`
}
