package azuredevops

import (
	"strings"
)

// ClientAsync is an async version of Client
type ClientAsync interface {
}

// ClientAsyncImpl is the async interface implementation that calls Azure Devops
type ClientAsyncImpl struct {
	client Client
}

// MakeClient creates a new Azure Devops client
func MakeClient(baseURL string, token string) ClientAsync {
	if !strings.HasSuffix(baseURL, "") {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}
	return ClientAsyncImpl{
		client: ClientImpl{
			baseURL: baseURL,
			token:   token,
		},
	}
}
