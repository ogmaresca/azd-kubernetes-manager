package processors_test

import (
	"github.com/ggmaresca/azd-kubernetes-manager/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MockKubernetesClient struct {
	listCounts   *map[string]*map[string]uint32
	deleteCounts *map[string]*map[string]uint32
}

func NewMockKubernetesClient() MockKubernetesClient {
	listCounts := make(map[string]*map[string]uint32)
	deleteCounts := make(map[string]*map[string]uint32)
	return MockKubernetesClient{
		listCounts:   &listCounts,
		deleteCounts: &deleteCounts,
	}
}

func (c MockKubernetesClient) DeleteCount(apiVersion string, kind string) *uint32 {
	if kinds, apiVersionExists := (*c.deleteCounts)[apiVersion]; apiVersionExists {
		if count, kindExists := (*kinds)[kind]; kindExists {
			return &count
		}
	}

	return nil
}

func (c MockKubernetesClient) List(apiVersion string, kind string, namespace string, labelSelector metav1.LabelSelector) ([]kubernetes.Resource, error) {
	if kinds, apiVersionExists := (*c.listCounts)[apiVersion]; apiVersionExists {
		if count, kindExists := (*kinds)[kind]; kindExists {
			(*kinds)[kind] = count + 1
		} else {
			(*kinds)[kind] = 1
		}
	} else {
		newKinds := make(map[string]uint32)
		(*c.listCounts)[apiVersion] = &newKinds
		return c.List(apiVersion, kind, namespace, labelSelector)
	}
	return []kubernetes.Resource{}, nil
}

func (c MockKubernetesClient) Delete(apiVersion string, kind string, namespace string, labelSelector metav1.LabelSelector) error {
	if kinds, apiVersionExists := (*c.deleteCounts)[apiVersion]; apiVersionExists {
		if count, kindExists := (*kinds)[kind]; kindExists {
			(*kinds)[kind] = count + 1
		} else {
			(*kinds)[kind] = 1
		}
	} else {
		newKinds := make(map[string]uint32)
		(*c.deleteCounts)[apiVersion] = &newKinds
		return c.Delete(apiVersion, kind, namespace, labelSelector)
	}
	return nil
}
