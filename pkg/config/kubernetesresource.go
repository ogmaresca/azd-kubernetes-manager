package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubernetesResource represents a Kubernetes resource, which has both Type and metadata information
type KubernetesResource struct {
	//
	// TypeMeta
	//
	Kind       string                     `yaml:"kind"`
	APIVersion string                     `yaml:"apiVersion"`
	Metadata   KubernetesResourceMetadata `yaml:"metadata"`
}

// KubernetesResourceMetadata represents a Kubernetes resource metadata
type KubernetesResourceMetadata struct {
	//
	// ObjectMeta
	//
	Name         string            `yaml:"name"`
	GenerateName string            `yaml:"generateName"`
	Namespace    string            `yaml:"namespace"`
	Labels       map[string]string `yaml:"labels"`
	Annotations  map[string]string `yaml:"annotations"`
}

// ToTypeMeta maps a KubernetesResource to a meta/v1 Type
func (r KubernetesResource) ToTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       r.Kind,
		APIVersion: r.APIVersion,
	}
}

// ToObjectMeta maps a KubernetesResourceMetadata to a meta/v1 Object
func (r KubernetesResourceMetadata) ToObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:         r.Name,
		GenerateName: r.GenerateName,
		Namespace:    r.Namespace,
		Labels:       r.Labels,
		Annotations:  r.Annotations,
	}
}
