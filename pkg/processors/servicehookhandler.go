package processors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexcesaro/log/stdlog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/args"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/azuredevops"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/config"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/templating"
)

var (
	logger = stdlog.GetFromFlags()

	serviceHookCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "azd_kubernetes_manager_service_hook_count",
		Help: "The total number of Service Hooks",
	}, []string{"eventType"})

	serviceHookDurationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "azd_kubernetes_manager_service_hook_duration_seconds",
		Help: "The duration of Service Hook requests",
	}, []string{"eventType"})

	serviceHookErrorCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "azd_kubernetes_manager_service_hook_error_count",
		Help: "The total number of Service Hooks",
	}, []string{"eventType", "reason"})
)

// ServiceHookHandler is an HTTP handler for service hooks
type ServiceHookHandler struct {
	args        args.ServiceHookArgs
	config      []config.ServiceHook
	ruleHandler RuleHandler
}

func (h ServiceHookHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	// Assert HTTP method
	if !strings.EqualFold(request.Method, "POST") {
		logger.Errorf("Service hooks must be POST requests - received %s method", request.Method)
		serviceHookCounter.With(prometheus.Labels{"eventType": "unknown"}).Inc()
		serviceHookErrorCounter.With(prometheus.Labels{"eventType": "unknown", "reason": fmt.Sprintf("HTTP %d Method Not Allowed", http.StatusMethodNotAllowed)}).Inc()
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read body into string
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(request.Body)
	if err != nil {
		logger.Errorf("Error reading request body from service hook: %s", err.Error())
		serviceHookCounter.With(prometheus.Labels{"eventType": "unknown"}).Inc()
		serviceHookErrorCounter.With(prometheus.Labels{"eventType": "unknown", "reason": "Error reading body"}).Inc()
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	requestStr := string(buffer.Bytes())
	if logger.LogDebug() {
		logger.Debugf("Received service hook: %s", requestStr)
	}

	// Parse JSON
	requestObj := new(azuredevops.ServiceHook)
	if err = json.NewDecoder(strings.NewReader(requestStr)).Decode(requestObj); err != nil {
		logger.Errorf("Error - could not parse JSON from Service hook. Error: %s\nRequest: %s", err.Error(), requestStr)
		serviceHookCounter.With(prometheus.Labels{"eventType": "unknown"}).Inc()
		serviceHookErrorCounter.With(prometheus.Labels{"eventType": "unknown", "reason": "JSON parse error"}).Inc()
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Prometheus metrics
	serviceHookCounter.With(prometheus.Labels{"eventType": requestObj.EventType}).Inc()
	defer prometheus.NewTimer(serviceHookDurationHistogram.With(prometheus.Labels{"eventType": requestObj.EventType})).ObserveDuration()

	if logger.LogDebug() {
		logger.Debugf("Deserialized response to: %#v", requestObj)
	}

	// Validate basic authentication
	username, password, ok := request.BasicAuth()
	if h.args.UseBasicAuthentication() {
		if ok && username == h.args.Username && password == h.args.Password {
			if logger.LogDebug() {
				logger.Debugf("[%s] Validated basic authentication", requestObj.Describe())
			}
		} else {
			logger.Errorf("[%s] Failed to validate basic authentication", requestObj.Describe())
			serviceHookErrorCounter.With(prometheus.Labels{"eventType": requestObj.EventType, "reason": fmt.Sprintf("HTTP %d Unauthorized", http.StatusUnauthorized)}).Inc()
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else if ok {
		logger.Noticef("[%s] Basic authentication was provided, but basic authentication was not configured.", requestObj.Describe())
	}

	for pos, config := range h.config {
		matches, err := config.Matches(requestObj)
		if err != nil {
			logger.Errorf("[%s] Error determining if Service Hook configuration %d matches request", requestObj.Describe(), pos)
			serviceHookErrorCounter.With(prometheus.Labels{"eventType": requestObj.EventType, "reason": "Error matching configuration"}).Inc()
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if matches {
			logger.Infof("[%s] Processing Service Hook configuration %d", requestObj.Describe(), pos)

			err := h.ruleHandler.Handle(config.Rules, templating.NewArgsFromServiceHook(*requestObj))
			if err != nil {
				logger.Errorf("[%s] Error processing rules: %s", requestObj.Describe(), err.Error())
				serviceHookErrorCounter.With(prometheus.Labels{"eventType": requestObj.EventType, "reason": "Error processing rules"}).Inc()
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !config.Continue {
				break
			}
		}
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

// NewServiceHookHandler creates a an HTTP handler for Service Hooks
func NewServiceHookHandler(args args.ServiceHookArgs, config []config.ServiceHook, ruleHandler RuleHandler) ServiceHookHandler {
	return ServiceHookHandler{
		args:        args,
		config:      config,
		ruleHandler: ruleHandler,
	}
}
