package kubernetes_test

import (
	"testing"

	"github.com/ggmaresca/azd-kubernetes-manager/pkg/kubernetes"
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/kubernetesmock"
)

func TestKubernetesClient(t *testing.T) {
	t.Run("kubernetes_client_interface_test", func(t *testing.T) {
		var client kubernetes.ClientAsync
		client = kubernetes.MakeFromClient(kubernetesmock.MockK8sClient{})

		switch syncClientType := client.Sync().(type) {
		case kubernetesmock.MockK8sClient:
			// Do nothing
		default:
			t.Errorf("Expected a MockK8sClient synchronous client, but found %#v", syncClientType)
		}
	})
}
