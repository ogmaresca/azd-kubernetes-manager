package config

// TemplateArgs holds all of the values available to go templating at runtime
type TemplateArgs struct {
	EventType string

	PullRequestId int

	BuildId int

	BuildNumber string

	ProjectName string

	ResourceUrl string

	ResouceStartTime string

	ResouceEndTime string

	ResourceName string

	ResourceReason string

	CurrentTime string
}
