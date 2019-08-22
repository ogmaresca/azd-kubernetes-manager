# Azure Devops Kubernetes Manager

azd-kubernetes-manager manages Kubernetes resources for Azure Devops.

## Azure Devops Abilities

* Taking actions from Service Hooks.

## Kubernetes Abilities

* Deleting resources.

## Configuration

The configuration file is a YAML file. See [Configuration.md](Configuration.md) for more.

## Installation

First, add this repo to Helm:

``` bash
helm repo add azd-kubernetes-manager https://raw.githubusercontent.com/ggmaresca/azd-kubernetes-manager/master/charts
helm repo update
```

Then use this command to install it:

``` bash
helm upgrade --install --namespace=kube-public azd-kubernetes-manager azd-kubernetes-manager/azd-kubernetes-manager
```

### Helm Chart Values

| Parameter                           | Description                                                                                                                                                                           | Default                                   |
| ----------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------- |
| `nameOverride`                      | An override value for the name.                                                                                                                                                       |                                           |
| `fullnameOverride`                  | An override value for the full name.                                                                                                                                                  |                                           |
| `minReadySeconds`                   | The deployment's `minReadySeconds`.                                                                                                                                                   | 0                                         |
| `revisionHistoryLimit`              | Number of Deployment versions to keep.                                                                                                                                                | 10                                        |
| `updateStrategy`                    | The Deployment Update Strategy.                                                                                                                                                       | `{ "type": "Recreate" }`                  |
| `image.repository`                  | The Docker Hub repository of the deployment.                                                                                                                                          | docker.io/gmaresca/azd-kubernetes-manager |
| `image.tag`                         | The image tag of the deployment.                                                                                                                                                      | latest version                            |
| `image.pullPolicy`                  | The image pull policy.                                                                                                                                                                | IfNotPresent                              |
| `image.pullSecrets`                 | Image Pull Secrets to use.                                                                                                                                                            | `[]`                                      |
| `logLevel`                          | The log level (debug, info, notice, warning, error, critical, alert, emergency, none)                                                                                                 | info                                      |
| `rate`                              | The period to poll Azure Devops and the Kubernetes API                                                                                                                                | 10s                                       |
| `combinePorts`                      | If true, health and metrics will be exposed on the same port as service hooks.                                                                                                        | `false`                                   |
| `username`                          | The username to use for Service Hook basic authentication.                                                                                                                            |                                           |
| `password`                          | The password to use for Service Hook basic authentication.                                                                                                                            |                                           |
| `configuration`                     | The contents of the [configuration file](Configuration.md).                                                                                                                           | `{ "serviceHooks" : [] }`                 |
| `resources.requests.cpu`            | The CPU requests of the deployment.                                                                                                                                                   | 0.05                                      |
| `resources.requests.memory`         | The memory requests of the deployment.                                                                                                                                                | 16Mi                                      |
| `resources.limits.cpu`              | The CPU limits of the deployment.                                                                                                                                                     | 0.1                                       |
| `resources.limits.memory`           | The memory limits of the deployment.                                                                                                                                                  | 64Mi                                      |
| `service.type`                      | The Service type.                                                                                                                                                                     | LoadBalancer                              |
| `service.externalTrafficPolicy`     | The Service External Traffic Policy (non-ClusterIP services).                                                                                                                         | Local                                     |
| `service.port`                      | The port for service hooks (and health/metrics if combinePorts is true).                                                                                                              | 80                                        |
| `service.nodePort`                  | The node port for service hooks (and health/metrics if combinePorts is true).                                                                                                         |                                           |
| `service.labels`                    | Labels to add to the Service.                                                                                                                                                         | `{}`                                      |
| `service.annotations`               | Annotations to add to the Service.                                                                                                                                                    | `{}`                                      |
| `ingress.enabled`                   | Whether to add an Ingress.                                                                                                                                                            | `false`                                   |
| `ingress.hosts`                     | The Ingress hosts. This is a string array.                                                                                                                                            | `[]`                                      |
| `ingress.basePath`                  | The base path to prepend to all paths (service hooks, metrics, health checks).                                                                                                        |                                           |
| `ingress.labels`                    | Labels to add to the Ingress.                                                                                                                                                         | `{}`                                      |
| `ingress.annotations`               | Annotations to add to the Ingress.                                                                                                                                                    | `{}`                                      |
| `ingress.tls.enabled`               | Whether to enable Ingress TLS on the hosts.                                                                                                                                           | `false`                                   |
| `ingress.tls.secretName`            | The name of the secret for Ingress TLS.                                                                                                                                               |                                           |
| `livenessProbe.failureThreshold`    | The failure threshold for the liveness probe.                                                                                                                                         | 3                                         |
| `livenessProbe.initialDelaySeconds` | The initial delay for the liveness probe.                                                                                                                                             | 1                                         |
| `livenessProbe.periodSeconds`       | The liveness probe period.                                                                                                                                                            | 10                                        |
| `livenessProbe.successThreshold`    | The success threshold for the liveness probe.                                                                                                                                         | 1                                         |
| `livenessProbe.timeoutSeconds`      | The timeout for the liveness probe.                                                                                                                                                   | 1                                         |
| `labels`                            | Labels to add to the Deployment.                                                                                                                                                      | `{}`                                      |
| `annotations`                       | Annotations to add to the Deployment.                                                                                                                                                 | `{}`                                      |
| `podLabels`                         | Labels to add to the Pod.                                                                                                                                                             | `{}`                                      |
| `podAnnotations`                    | Annotations to add to the Pod.                                                                                                                                                        | `{}`                                      |
| `pdb.enabled`                       | Whether to enable a PodDisruptionBudget.                                                                                                                                              | `false`                                   |
| `pdb.minAvailable`                  | The minimum number of pods to keep. Incompatible with `maxUnavailable`.                                                                                                               | 50%                                       |
| `pdb.maxUnavailable`                | The maximum unvailable pods. Incompatible with `minAvailable`.                                                                                                                        | 50%                                       |
| `rbac.create`                       | Whether to create Role Based Access for the deployment.                                                                                                                               | `true`                                    |
| `rbac.clusterRules`                 | The [PolicyRules](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#policyrule-v1-rbac-authorization-k8s-io) for a ClusterRole.                                    | `[]`                                      |
| `rbac.rules`                        | A map with namespaces as the key and a [PolicyRule](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#policyrule-v1-rbac-authorization-k8s-io) array as the value. | `{}`                                      |
| `serviceAccount.create`             | Whether to create a service account for the deployment.                                                                                                                               | `true`                                    |
| `serviceAccount.name`               | The name of an existing SA `serviceAccount.create` is false.                                                                                                                          |                                           |
| `serviceMonitor.enabled`            | Create a `prometheus-operator` ServiceMonitor.                                                                                                                                        | `false`                                   |
| `serviceMonitor.namespace`          | The namespace to install the ServiceMonitor.                                                                                                                                          | Release namespace                         |
| `serviceMonitor.labels`             | Labels to add to the ServiceMonitor.                                                                                                                                                  | `{}`                                      |
| `serviceMonitor.honorLabels`        | Set `honorLabels` on the ServiceMonitor spec.                                                                                                                                         |                                           |
| `serviceMonitor.interval`           | The scrape interval on the ServiceMonitor.                                                                                                                                            | Defaults to `rate`                        |
| `serviceMonitor.metricRelabelings`  | `metricRelabelings` to set on the ServiceMonitor.                                                                                                                                     | `false`                                   |
| `serviceMonitor.relabelings`        | `relabelings` to set on the ServiceMonitor.                                                                                                                                           | `false`                                   |
| `grafanaDashboard.enabled`          | Create a ConfigMap with a Grafana dashboard.                                                                                                                                          | `false`                                   |
| `grafanaDashboard.labels`           | Labels to add to the Grafana dashboard ConfigMap.                                                                                                                                     | `{"grafana_dashboard":"1"}`               |
| `rbac.getConfigmaps`                | Allow getting ConfigMaps, to retrieve the AZP_POOL env value.                                                                                                                         | `false`                                   |
| `rbac.getSecrets`                   | Allow getting Secrets, to retrieve the AZP_POOL env value.                                                                                                                            | `false`                                   |
| `dnsPolicy`                         | The pod DNS policy.                                                                                                                                                                   | `null`                                    |
| `dnsConfig`                         | The pod DNS config.                                                                                                                                                                   | `{}`                                      |
| `restartPolicy`                     | The pod restart policy.                                                                                                                                                               | Always                                    |
| `nodeSelector`                      | The pod node selector.                                                                                                                                                                | `{}`                                      |
| `tolerations`                       | The pod node tolerations.                                                                                                                                                             | `{}`                                      |
| `affinity`                          | The pod node affinity.                                                                                                                                                                | `{}`                                      |
| `securityContext`                   | The pod security context.                                                                                                                                                             | `{}`                                      |
| `hostNetwork`                       | Whether to use the host network of the node.                                                                                                                                          | `false`                                   |
| `initContainers`                    | Init containers to add.                                                                                                                                                               | `[]`                                      |
| `lifecycle`                         | Lifecycle (postStart, preStop) for the pod.                                                                                                                                           | `{}`                                      |
| `sidecars`                          | Additional containers to add.                                                                                                                                                         | `[]`                                      |

## Docker Hub

[Docker Hub link](https://hub.docker.com/r/gmaresca/azd-kubernetes-manager).
