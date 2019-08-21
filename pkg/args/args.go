package args

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	rate       = flag.Duration("rate", 10*time.Second, "Duration to check the number of agents.")
	azdToken   = flag.String("token", "", "The Azure Devops token.")
	azdURL     = flag.String("url", "", "The Azure Devops URL. https://dev.azure.com/AccountName")
	configFile = flag.String("config-file", "", "The path to the config file.")
	basePath   = flag.String("base-path", "", "The path to prepend before every path.")
	port       = flag.Int("port", 10102, "The port to serve HTTP requests.")
	username   = flag.String("username", "", "The username to use for Service Hooks basic authentication.")
	password   = flag.String("password", "", "The password to use for Service Hooks basic authentication.")
	healthPort = flag.Int("health-port", 10902, "The port to serve health checks and metrics.")
)

// Args holds all of the program arguments
type Args struct {
	Rate         time.Duration
	ConfigFile   string
	ServiceHooks ServiceHookArgs
	AZD          AzureDevopsArgs
	Health       HealthArgs
}

// ScaleDownArgs holds all of the scale-down related args
type ScaleDownArgs struct {
	Delay time.Duration
	Max   int32
}

// ServiceHookArgs holds all of the service hook related args
type ServiceHookArgs struct {
	BasePath string
	Port     int
	Username string
	Password string
}

// UseBasicAuthentication returns true if the Username and Password are not empty
func (a ServiceHookArgs) UseBasicAuthentication() bool {
	return a.Username != "" && a.Password != ""
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

// FromFlags returns an Args parsed from the program flags
func FromFlags() Args {
	// error should be validated in ValidateArgs()
	return Args{
		Rate:       *rate,
		ConfigFile: *configFile,

		ServiceHooks: ServiceHookArgs{
			BasePath: *basePath,
			Port:     *port,
			Username: *username,
			Password: *password,
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

	if *username != *password && (*username == "" || *password == "") {
		validationErrors = append(validationErrors, "Either the both or neither of the username and password must be provided.")
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("Error(s) with arguments:\n%s", strings.Join(validationErrors, "\n"))
	}
	return nil
}
