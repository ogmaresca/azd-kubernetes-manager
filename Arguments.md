# Azure Devops Kubernetes Manager Arguments

These arguments are defined in [args.go](pkg/args/args.go)

| Argument    | Description                                                                                                         | Default Value | Required                  |
| ----------- | ------------------------------------------------------------------------------------------------------------------- | ------------- | ------------------------- |
| rate        | How often to query the Azure Devops API.                                                                            | 10s           | If overriden.             |
| token       | The Azure Devops token to call the Azure Devops API with.                                                           |               | If API rules are defined. |
| url         | The Azure Devops organization URL.                                                                                  |               | If API rules are defined. |
| config-file | The path to the config file.                                                                                        |               | Yes                       |
| base-path   | The base path to prepend to every HTTP endpoint.                                                                    |               | No                        |
| port        | The port to listen on for Service Hooks.                                                                            | 10102         | If overridden.            |
| username    | The basic authentication username to use for Service Hooks.                                                         |               | If password is provided.  |
| password    | The basic authentication password to use for Service Hooks.                                                         |               | If username is provided.  |
| healh-port  | The port to listen on for health checks and metrics.                                                                | 10902         | If overridden.            |
| log         | stdlog minimum levels. Allowed values are debug, info, notice, warning, error, critical, alert, emergency and none. | info          | If overridden.            |

