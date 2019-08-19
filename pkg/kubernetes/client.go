package kubernetes

import (
	"fmt"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/alexcesaro/log/stdlog"
)

var logger = stdlog.GetFromFlags()

// Client is a wrapper around the client-go package for Kubernetes
type Client interface {
	List(apiVersion string, kind string, namespace string, labelSelector metav1.LabelSelector) ([]Resource, error)
	Delete(apiVersion string, kind string, namespace string, labelSelector metav1.LabelSelector) error
}

// ClientImpl is the interface implementation of Client
type ClientImpl struct {
	config       *rest.Config
	client       *k8s.Clientset
	apiResources map[string]metav1.APIResourceList
}

// makeClient returns a Client
func makeClient() (Client, error) {
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		kubeconfigEnv := os.Getenv("KUBECONFIG")
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfigEnv)
		if err != nil {
			home := os.Getenv("HOME")
			if home == "" {
				home = os.Getenv("USERPROFILE") // windows
			}
			k8sConfig, err = clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/.kube/config", home))
			if err != nil {
				return nil, fmt.Errorf("Error initializing Kubernetes config: %s", err.Error())
			}
		}
	}

	clientset, err := k8s.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	return ClientImpl{
		config:       k8sConfig,
		client:       clientset,
		apiResources: make(map[string]metav1.APIResourceList),
	}, nil
}

// GetAPIResources retrieves and caches API resources for the given API Version
func (c ClientImpl) GetAPIResources(apiVersion string) (*metav1.APIResourceList, error) {
	resources, exists := c.apiResources[apiVersion]
	if !exists {
		resourcesPtr, err := c.client.Discovery().ServerResourcesForGroupVersion(apiVersion)
		if err != nil {
			return resourcesPtr, err
		}
		c.apiResources[apiVersion] = *resourcesPtr
		return resourcesPtr, nil
	}

	return &resources, nil
}

// GetAPIResource returns a specific API resource
func (c ClientImpl) GetAPIResource(apiVersion string, kind string) (*metav1.APIResource, error) {
	apiResources, err := c.GetAPIResources(apiVersion)
	if err != nil {
		return nil, err
	}

	for _, apiResource := range apiResources.APIResources {
		if strings.EqualFold(apiResource.Kind, kind) {
			return &apiResource, nil
		}
	}

	return nil, fmt.Errorf("Kind '%s' was not found in API Version '%s'", kind, apiVersion)
}

// List a Kubernetes resource
func (c ClientImpl) List(apiVersion string, kind string, namespace string, labelSelector metav1.LabelSelector) ([]Resource, error) {
	apiResource, err := c.GetAPIResource(apiVersion, kind)
	if err != nil {
		return nil, fmt.Errorf("Error getting API Resource %s for API %s: %s", kind, apiVersion, err.Error())
	}

	client, err := c.RESTClient(apiVersion)
	if err != nil {
		return nil, fmt.Errorf("Error getting REST client for API %s: %s", apiVersion, err.Error())
	}

	options := metav1.ListOptions{
		TypeMeta:      metav1.TypeMeta{APIVersion: apiVersion, Kind: kind},
		LabelSelector: metav1.FormatLabelSelector(&labelSelector),
	}

	result := &ResourceList{}

	err = client.
		Get().
		NamespaceIfScoped(namespace, namespace != "").
		Resource(apiResource.Name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)

	if err != nil {
		return nil, fmt.Errorf("Error listing '%s' for API %s: %s", kind, apiVersion, err.Error())
	}
	return result.Items, nil
}

// Delete Kubernetes resource(s)
func (c ClientImpl) Delete(apiVersion string, kind string, namespace string, labelSelector metav1.LabelSelector) error {
	apiResource, err := c.GetAPIResource(apiVersion, kind)
	if err != nil {
		return fmt.Errorf("Error getting API Resource %s for API %s: %s", kind, apiVersion, err.Error())
	}

	client, err := c.RESTClient(apiVersion)
	if err != nil {
		return fmt.Errorf("Error getting REST client for API %s: %s", apiVersion, err.Error())
	}

	resources, err := c.List(apiVersion, kind, namespace, labelSelector)
	if err != nil {
		return err
	}

	var channels []chan error
	for _, resource := range resources {
		channel := make(chan error)
		go func() {
			err := client.Delete().
				NamespaceIfScoped(resource.Namespace, resource.Namespace != "").
				Resource(apiResource.Name).
				Name(resource.Name).
				Do().
				Error()

			if err != nil {
				channel <- fmt.Errorf("Error deleting %s %s %s: %s", apiVersion, kind, resource.Name, err.Error())
			} else {
				logger.Infof("Deleted %s %s %s", apiVersion, kind, resource.Name)
				channel <- nil
			}
		}()
		channels = append(channels, channel)
	}

	var errors []string
	for _, channel := range channels {
		err := <-channel
		if err != nil {
			errors = append(errors, fmt.Sprintf("- %s", strings.ReplaceAll(err.Error(), "\n", "\n  ")))
		}
	}

	if len(errors) > 0 {
		err = fmt.Errorf("Errors deleting resources:\n%s", strings.Join(errors, "\n"))
	}

	return err
}

// RESTClient creates a kubernetes client for the given API version
func (c ClientImpl) RESTClient(apiVersion string) (rest.Interface, error) {
	groupVersion := c.GetGroupVersion(apiVersion)

	config := *c.config
	config.GroupVersion = &groupVersion
	if groupVersion.Group == "" {
		config.APIPath = "/api"
	} else {
		config.APIPath = "/apis"
	}
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	return rest.RESTClientFor(&config)
}

// GetGroupVersion maps an APIVersion to a GroupVersion
func (c ClientImpl) GetGroupVersion(apiVersion string) schema.GroupVersion {
	split := strings.Split(apiVersion, "/")
	if len(split) == 1 {
		return schema.GroupVersion{Group: "", Version: split[0]}
	} else {
		return schema.GroupVersion{Group: split[0], Version: split[1]}
	}
}
