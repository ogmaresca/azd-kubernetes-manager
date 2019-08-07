package kubernetes

import (
	"fmt"
	"os"

	k8s "k8s.io/client-go/kubernetes"
	k8srest "k8s.io/client-go/rest"
	k8sclientcmd "k8s.io/client-go/tools/clientcmd"
)

// Client is a wrapper around the client-go package for Kubernetes
type Client interface {
}

// ClientImpl is the interface implementation of Client
type ClientImpl struct {
	client *k8s.Clientset
}

// makeClient returns a Client
func makeClient() (Client, error) {
	k8sConfig, err := k8srest.InClusterConfig()
	if err != nil {
		kubeconfigEnv := os.Getenv("KUBECONFIG")
		k8sConfig, err = k8sclientcmd.BuildConfigFromFlags("", kubeconfigEnv)
		if err != nil {
			home := os.Getenv("HOME")
			if home == "" {
				home = os.Getenv("USERPROFILE") // windows
			}
			k8sConfig, err = k8sclientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/.kube/config", home))
			if err != nil {
				return nil, fmt.Errorf("Error initializing Kubernetes config: %s", err.Error())
			}
		}
	}

	clientset, err := k8s.NewForConfig(k8sConfig)
	if err != nil {
		return nil, nil
	}
	return ClientImpl{clientset}, err
}
