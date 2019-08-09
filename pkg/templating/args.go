package args

// Args holds all of the values available to go templating at runtime
type Args struct {
	EventType string

	PullRequestID int

	BuildID int

	BuildNumber string

	ProjectName string

	ResourceURL string

	ResouceStartTime string

	ResouceEndTime string

	ResourceName string

	ResourceReason string

	CurrentTime string
}
