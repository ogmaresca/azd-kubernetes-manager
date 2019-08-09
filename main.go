// main is the main package - golint fix
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/args"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/azuredevops"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/config"
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

	// Initialize
	//azdClient := azuredevops.MakeClient(args.AZD.URL, args.AZD.Token)
	k8sClient, err := kubernetes.MakeClient()
	if err != nil {
		panic(err.Error())
	}

	configFileYaml, err := ioutil.ReadFile(args.ConfigFile)
	if err != nil {
		logging.Logger.Panicf("Error reading config file \"%s\": %s", args.ConfigFile, err.Error())
	}
	configFile, err := config.NewConfigFile(configFileYaml)
	if err != nil {
		logging.Logger.Panicf("Error parsing config file \"%s\": %s", args.ConfigFile, err.Error())
	}

	logging.Logger.Tracef("Parsed config file: %v", configFile)

	func() {
		mux := http.NewServeMux()
		mux.Handle("/serviceHooks", azuredevops.NewServiceHookHandler(args, configFile, k8sClient))

		var healthMux *http.ServeMux
		if args.ServiceHooks.Port != args.Health.Port {
			healthMux = http.NewServeMux()
		} else {
			healthMux = mux
		}

		healthMux.Handle("/healthz", health.LivenessCheck{})
		healthMux.Handle("/metrics", promhttp.Handler())

		go func() {
			err := http.ListenAndServe(fmt.Sprintf(":%d", args.ServiceHooks.Port), mux)
			if err != nil {
				logging.Logger.Panicf("Error serving HTTP requests: %s", err.Error())
			}
		}()

		if args.ServiceHooks.Port != args.Health.Port {
			go func() {
				err = http.ListenAndServe(fmt.Sprintf(":%d", args.Health.Port), healthMux)
				if err != nil {
					logging.Logger.Panicf("Error serving health checks and metrics: %s", err.Error())
				}
			}()
		}
	}()

	for {
		time.Sleep(args.Rate)
	}

	logging.Logger.Info("Exiting azd-kubernetes-manager")
}
