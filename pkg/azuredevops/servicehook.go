package azuredevops

// ServiceHook defines the JSON body of service hooks
// https://docs.microsoft.com/en-us/azure/devops/service-hooks/events?toc=%2Fazure%2Fdevops%2Fintegrate%2Ftoc.json&bc=%2Fazure%2Fdevops%2Fintegrate%2Fbreadcrumb%2Ftoc.json&view=azure-devops
type ServiceHook struct {
	// Fields in all event types
	ID          string `json:"id"`
	EventType   string `json:"eventType"`
	PublisherID string `json:"publisherId"`
	Scope       string `json:"scope"`
	// TODO complete object
}
