#!/bin/bash

AZD_KUBERNETES_MANAGER_VERSION=$(cat version) && \
echo "Uploading azd-kubernetes-manager $AZD_KUBERNETES_MANAGER_VERSION" && \
docker tag azd-kubernetes-manager:dev docker.io/gmaresca/azd-kubernetes-manager:$AZD_KUBERNETES_MANAGER_VERSION && \
docker tag azd-kubernetes-manager:dev docker.io/gmaresca/azd-kubernetes-manager:latest && \
docker push docker.io/gmaresca/azd-kubernetes-manager:$AZD_KUBERNETES_MANAGER_VERSION && \
docker push docker.io/gmaresca/azd-kubernetes-manager:latest && \
docker rmi docker.io/gmaresca/azd-kubernetes-manager:$AZD_KUBERNETES_MANAGER_VERSION && \
docker rmi docker.io/gmaresca/azd-kubernetes-manager:latest && \
echo "Finished uploading azd-kubernetes-manager $AZD_KUBERNETES_MANAGER_VERSION"
