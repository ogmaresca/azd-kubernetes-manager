package azuredevops

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/args"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/config"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/kubernetes"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/logging"
)

var (
	serviceHookCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "azd_kubernetes_manager_service_hook_count",
		Help: "The total number of Service Hooks",
	})

	serviceHookDurationHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "azd_kubernetes_manager_service_hook_duration_seconds",
		Help: "The duration of Service Hook requests",
	})
)

// ServiceHookHandler is an HTTP handler for service hooks
type ServiceHookHandler struct {
	args      args.Args
	config    config.ConfigFile
	k8sClient kubernetes.ClientAsync
}

func (h ServiceHookHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	timer := prometheus.NewTimer(serviceHookDurationHistogram)
	defer timer.ObserveDuration()

	if !strings.EqualFold(request.Method, "POST") {
		logging.Logger.Errorf("Service hooks must be POST requests - received %s method", request.Method)
		writer.WriteHeader(405)
		return
	}

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(request.Body)
	if err != nil {
		logging.Logger.Errorf("Error reading request body from service hook: %s", err.Error())
		writer.WriteHeader(500)
		return
	}
	requestStr := string(buffer.Bytes())

	logging.Logger.Debugf("Received service hook: \n%s", requestStr)

	requestObj := new(ServiceHook)
	if err = json.NewDecoder(strings.NewReader(requestStr)).Decode(requestObj); err != nil {
		logging.Logger.Errorf("Error - could not parse JSON from Service hook. Error: %s\nRequest:\n%s", err.Error(), requestStr)
		writer.WriteHeader(400)
		return
	}

	// TODO handle basic authentication

	serviceHookCounter.Inc()

	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}

// NewServiceHookHandler creates a an HTTP handler for Service Hooks
func NewServiceHookHandler(args args.Args, config config.ConfigFile, k8sClient kubernetes.ClientAsync) ServiceHookHandler {
	return ServiceHookHandler{
		args:      args,
		config:    config,
		k8sClient: k8sClient,
	}
}
