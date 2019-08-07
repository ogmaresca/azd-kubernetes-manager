// main is the main package - golint fix
package main

import(
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/args"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/azuredevops"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/health"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/kubernetes"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/logging"
)

func main() {
	// Parse arguments
	flag.Parse()

	if err := args.ValidateArgs(); err != nil {
		panic(err.Error())
	}
	args := args.ArgsFromFlags()

	logging.Logger.SetLevel(args.Logging.Level)

	// Initialize Azure Devops client
	azdClient := azuredevops.MakeClient(args.AZD.URL, args.AZD.Token)
	k8sClient, err := kubernetes.MakeClient()
	if err != nil {
		panic(err.Error())
	}

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/healthz", health.LivenessCheck{})
		mux.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(fmt.Sprintf(":%d", args.Health.Port), mux)
		if err != nil {
			logging.Logger.Panicf("Error serving health checks and metrics: %s", err.Error())
		}
	}()

	for {
		if(err != nil) {
			// Implement loop
			time.Sleep(args.Rate)
		}
	}

	logging.Logger.Info("Exiting azd-kubernetes-manager")
}