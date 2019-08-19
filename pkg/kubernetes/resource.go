package kubernetes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Resource represents a Kubernetes resource, which has both Type and Object information
type Resource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

// DeepCopyInto is required to implement the Kubernetes Object interface
func (in *Resource) DeepCopyInto(out *Resource) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
}

// DeepCopyObject is required to implement the Kubernetes Object interface
func (in *Resource) DeepCopyObject() runtime.Object {
	out := &Resource{}
	in.DeepCopyInto(out)
	return out
}

// ResourceList represents a list of Kubernetes resources
type ResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Resource `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// DeepCopyInto is required to implement the Kubernetes Object interface
func (in *ResourceList) DeepCopyInto(out *ResourceList) {
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	var items []Resource
	for _, resource := range in.Items {
		copy := &Resource{}
		resource.DeepCopyInto(copy)
		items = append(items, *copy)
	}
	out.Items = items
}

// DeepCopyObject is required to implement the Kubernetes Object interface
func (in *ResourceList) DeepCopyObject() runtime.Object {
	out := &ResourceList{}
	in.DeepCopyInto(out)
	return out
}
