package azuredevops

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const acceptHeader = "application/json;api-version=5.0-preview.1"

var (
	azdDurations = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "azd_kubernetes_manager_azd_call_duration_seconds",
		Help: "Duration of Azure Devops calls",
	}, []string{"operation"})

	azdCounts = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "azd_kubernetes_manager_azd_call_count",
		Help: "Counts of Azure Devops calls",
	}, []string{"operation"})

	azd429Counts = promauto.NewCounter(prometheus.CounterOpts{
		Name: "azd_kubernetes_manager_azd_call_429_count",
		Help: "Counts of Azure Devops calls returning HTTP 429 (Too Many Requests)",
	})
)

// Client is used to call Azure Devops
type Client interface {
}

// ClientImpl is the interface implementation that calls Azure Devops
type ClientImpl struct {
	baseURL string

	token string
}

func (c ClientImpl) executeGETRequest(endpoint string, response interface{}) error {
	request, err := http.NewRequest("GET", c.baseURL+endpoint, nil)

	if err != nil {
		return err
	}

	request.Header.Set("Accept", acceptHeader)
	request.Header.Set("User-Agent", "go-azd-kubernetes-manager")

	request.SetBasicAuth("user", c.token)

	httpClient := http.Client{}
	httpResponse, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		httpErr := NewHTTPError(httpResponse)
		if httpErr.RetryAfter != nil {
			azd429Counts.Inc()
		}
		return httpErr
	}

	err = json.NewDecoder(httpResponse.Body).Decode(response)
	if err != nil {
		return fmt.Errorf("Error - could not parse JSON response from %s: %s", endpoint, err.Error())
	}

	return nil
}

