# Azure Devops Kubernetes Manager

azd-kubernetes-manager manages Kubernetes resources for Azure Devops.

## Azure Devops Abilities

* Taking actions from Service Hooks.

## Kubernetes Abilities

* Deleting resources.

## Installation

First, add this repo to Helm:

``` bash
helm repo add azd-kubernetes-manager https://raw.githubusercontent.com/ggmaresca/azd-kubernetes-manager/master/charts
helm repo update
```

## Configuration

The configuration file is a YAML file. See [Configuration.md](Configuration.md) for more.

## Docker Hub

[Docker Hub link](https://hub.docker.com/r/gmaresca/azd-kubernetes-manager).
