package args

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	logLevel   = flag.String("log-level", "info", "Log level (trace, debug, info, warn, error, fatal, panic).")
	rate       = flag.Duration("rate", 10*time.Second, "Duration to check the number of agents.")
	azdToken   = flag.String("token", "", "The Azure Devops token.")
	azdURL     = flag.String("url", "", "The Azure Devops URL. https://dev.azure.com/AccountName")
	configFile = flag.String("config-file", "", "The path to the config file.")
	basePath   = flag.String("base-path", "", "The path to prepend before every path.")
	port       = flag.Int("port", 10102, "The port to serve HTTP requests.")
	username   = flag.String("username", "", "The username to use for Service Hooks basic authentication.")
	password   = flag.String("password", "", "The password to use for Service Hooks basic authentication.")
	healthPort = flag.Int("healh-port", 10902, "The port to serve health checks and metrics.")
)

// Args holds all of the program arguments
type Args struct {
	Rate         time.Duration
	ConfigFile   string
	ServiceHooks ServiceHookArgs
	Logging      LoggingArgs
	AZD          AzureDevopsArgs
	Health       HealthArgs
}

// ScaleDownArgs holds all of the scale-down related args
type ScaleDownArgs struct {
	Delay time.Duration
	Max   int32
}

// LoggingArgs holds all of the logging related args
type LoggingArgs struct {
	Level log.Level
}

// ServiceHookArgs holds all of the service hook related args
type ServiceHookArgs struct {
	BasePath string
	Port     int
	Username string
	Password string
}

// HealthArgs holds all of the healthcheck related args
type HealthArgs struct {
	Port int
}

// AzureDevopsArgs holds all of the Azure Devops related args
type AzureDevopsArgs struct {
	Token string
	URL   string
}

// ArgsFromFlags returns an Args parsed from the program flags
func ArgsFromFlags() Args {
	// error should be validated in ValidateArgs()
	logrusLevel, _ := log.ParseLevel(*logLevel)
	return Args{
		Rate:       *rate,
		ConfigFile: *configFile,

		ServiceHooks: ServiceHookArgs{
			BasePath: *basePath,
			Port:     *port,
			Username: *username,
			Password: *password,
		},

		Logging: LoggingArgs{
			Level: logrusLevel,
		},

		AZD: AzureDevopsArgs{
			Token: *azdToken,
			URL:   *azdURL,
		},

		Health: HealthArgs{
			Port: *healthPort,
		},
	}
}

// ValidateArgs validates all of the command line arguments
func ValidateArgs() error {
	// Validate arguments
	var validationErrors []string
	_, err := log.ParseLevel(*logLevel)
	if err != nil {
		validationErrors = append(validationErrors, err.Error())
	}
	if rate == nil {
		validationErrors = append(validationErrors, "Rate is required.")
	} else if rate.Seconds() <= 1 {
		validationErrors = append(validationErrors, fmt.Sprintf("Rate '%s' is too low.", rate.String()))
	}

	if configFile == nil || *configFile == "" {
		validationErrors = append(validationErrors, "Config File is required.")
	} else {
		configFileInfo, err := os.Stat(*configFile)
		if err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("Error validating config file: %s", err.Error()))
		} else if configFileInfo.IsDir() {
			validationErrors = append(validationErrors, "Configuration file argument points to a directory")
		}
	}

	/*if *azdToken == "" {
		validationErrors = append(validationErrors, "The Azure Devops token is required.")
	}
	if *azdURL == "" {
		validationErrors = append(validationErrors, "The Azure Devops URL is required.")
	}*/

	if *port <= 0 {
		validationErrors = append(validationErrors, "The port must be greater than 0.")
	}
	if *healthPort <= 0 {
		validationErrors = append(validationErrors, "The health port must be greater than 0.")
	}
	if len(validationErrors) > 0 {
		return fmt.Errorf("Error(s) with arguments:\n%s", strings.Join(validationErrors, "\n"))
	}
	return nil
}
