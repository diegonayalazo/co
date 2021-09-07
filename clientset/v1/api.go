package v1

import (
	v1 "github.com/diegonayalazo/co/api/types/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ExampleV1Interface interface {
	Brokers(namespace string) BrokerInterface
}

type ExampleV1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*ExampleV1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1.GroupName, Version: v1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ExampleV1Client{restClient: client}, nil
}

func (c *ExampleV1Client) Brokers(namespace string) BrokerInterface {
	return &projectClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
