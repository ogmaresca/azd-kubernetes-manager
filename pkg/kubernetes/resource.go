package kubernetes

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Resource represents a Kubernetes resource, which has both Type and Object information
type Resource struct {
	metav1.TypeMeta
	Metadata metav1.ObjectMeta
}
