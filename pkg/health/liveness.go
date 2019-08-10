package health

import (
	"net/http"

	"github.com/alexcesaro/log/stdlog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	logger = stdlog.GetFromFlags()

	livenessProbeCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "azd_kubernetes_manager_liveness_probe_count",
		Help: "The total number of liveness probes",
	})
)

// LivenessCheck is an HTTP Handler
type LivenessCheck struct {
}

func (c LivenessCheck) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logger.Debug("Liveness probe")

	livenessProbeCounter.Inc()

	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}
