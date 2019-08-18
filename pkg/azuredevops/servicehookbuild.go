package azuredevops

import (
	"time"
)

// ServiceHookResourceBuildComplete holds fields related to the Service Hook resource for build.complete events
type ServiceHookResourceBuildComplete struct {
	URI                *string                             `json:"uri"`
	BuildNumber        *string                             `json:"buildNumber"`
	StartTime          *time.Time                          `json:"startTime"`
	FinishTime         *time.Time                          `json:"finishTime"`
	DropLocation       *string                             `json:"dropLocation"`
	Drop               *ServiceHookResourceBuildDrop       `json:"drop"`
	Log                *ServiceHookResourceBuildLog        `json:"log"`
	SourceGetVersion   *string                             `json:"sourceGetVersion"`
	LastChangedBy      *User                               `json:"lastChangedBy"`
	RetainIndefinitely *bool                               `json:"retainIndefinitely"`
	HasDiagnostics     *bool                               `json:"hasDiagnostics"`
	Definition         *ServiceHookResourceBuildDefinition `json:"definition"`
	Queue              *ServiceHookResourceBuildQueue      `json:"queue"`
	Requests           []ServiceHookResourceBuildRequests  `json:"requests"`
}

// ServiceHookResourceBuildDrop holds the drop field of a Resource
type ServiceHookResourceBuildDrop struct {
	Location    string `json:"location"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	DownloadURL string `json:"downloadUrl"`
}

// ServiceHookResourceBuildLog holds the log field of a Resource
type ServiceHookResourceBuildLog struct {
	Type        string `json:"type"`
	URL         string `json:"url"`
	DownloadURL string `json:"downloadUrl"`
}

// ServiceHookResourceBuildQueue holds the Queue definition of pipeline Service Hooks
type ServiceHookResourceBuildQueue struct {
	ServiceHookResourceDefinition
	QueueType string `json:"queueType"`
}

// ServiceHookResourceBuildRequests holds the Requests definition of pipeline Service Hooks
type ServiceHookResourceBuildRequests struct {
	ID           int    `json:"id"`
	URL          string `json:"url"`
	RequestedFor User   `json:"requestedFor"`
}

// ServiceHookResourceBuildDefinition holds the definition of the Service Hook resource
type ServiceHookResourceBuildDefinition struct {
	ServiceHookResourceDefinition
	BatchSize      int    `json:"batchSize"`
	TriggerType    string `json:"triggerType"`
	DefinitionType string `json:"definitionType"`
}

// ServiceHookResourceBuildDefinitionTriggerType are the Service Hook Resource Definition trigger types
type ServiceHookResourceBuildDefinitionTriggerType string

const (
	// ServiceHookResourceBuildDefinitionTriggerTypeNone is the Service Hook Resource Definition trigger type of none
	ServiceHookResourceBuildDefinitionTriggerTypeNone ServiceHookResourceBuildDefinitionTriggerType = "none"
	// TODO fill out
)

// ServiceHookResourceBuildDefinitionType are the Service Hook Resource Definition types
type ServiceHookResourceBuildDefinitionType string

const (
	// ServiceHookResourceBuildDefinitionTypeXAML is the Service Hook Resource Definition type of xaml
	ServiceHookResourceBuildDefinitionTypeXAML ServiceHookResourceBuildDefinitionType = "xaml"
	// TODO fill out
)

// ServiceHookResourceBuildQueueType are the Service Hook Resource Queue types
type ServiceHookResourceBuildQueueType string

const (
	// ServiceHookResourceBuildQueueTypeBuildController is the Service Hook Resource Queue type of buildController
	ServiceHookResourceBuildQueueTypeBuildController ServiceHookResourceBuildQueueType = "buildController"
	// TODO fill out
)
