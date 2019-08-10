// main is the main package - golint fix
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/alexcesaro/log/stdlog"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/args"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/azuredevops"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/config"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/health"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/kubernetes"
)

var (
	logger = stdlog.GetFromFlags()
)

func main() {
	// Parse arguments
	flag.Parse()

	if err := args.ValidateArgs(); err != nil {
		panic(err.Error())
	}
	args := args.FromFlags()

	// Initialize
	//azdClient := azuredevops.MakeClient(args.AZD.URL, args.AZD.Token)
	k8sClient, err := kubernetes.MakeClient()
	if err != nil {
		panic(err.Error())
	}

	configFile := getConfigFile(args)

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
				panicf("Error serving HTTP requests: %s", err.Error())
			}
		}()

		if args.ServiceHooks.Port != args.Health.Port {
			go func() {
				err = http.ListenAndServe(fmt.Sprintf(":%d", args.Health.Port), healthMux)
				if err != nil {
					panicf("Error serving health checks and metrics: %s", err.Error())
				}
			}()
		}
	}()

	for {
		time.Sleep(args.Rate)
	}
}

// panicf panics with a formatted message
func panicf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}

func getConfigFile(args args.Args) config.File {
	configFileYaml, err := ioutil.ReadFile(args.ConfigFile)
	if err != nil {
		panicf("Error reading config file \"%s\": %s", args.ConfigFile, err.Error())
	}
	configFile, err := config.NewConfigFile(configFileYaml)
	if err != nil {
		panicf("Error parsing config file \"%s\": %s", args.ConfigFile, err.Error())
	}

	logger.Debugf("Parsed config file:\n%#v", configFile)

	logger.Infof("\n%s", configFile.Describe())

	configFileWarnings, err := configFile.Validate()
	if len(configFileWarnings) > 0 {
		logger.Warningf("Warnings from config file:\n%s", strings.Join(configFileWarnings, "\n"))
	}
	if err != nil {
		panicf("Errors from config file:\n%s", err.Error())
	}
	return configFile
}
