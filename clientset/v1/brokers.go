package v1

import (
	"bytes"
	"context"
	"fmt"

	v1 "github.com/diegonayalazo/co/api/types/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/jsonpath"
)

type BrokerInterface interface {
	List(opts metav1.ListOptions) (*v1.BrokerList, error)
	Get(name string, options metav1.GetOptions) (*v1.Broker, error)
	GetPath(name string, options metav1.GetOptions, jsonpath string) (string, error)
	Create(*v1.Broker) (*v1.Broker, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (result *v1.Broker, err error)
	Delete(name string, opts metav1.DeleteOptions) error
}

type brokerClient struct {
	restClient rest.Interface
	ns         string
}

func (c *brokerClient) List(opts metav1.ListOptions) (*v1.BrokerList, error) {
	result := v1.BrokerList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("brokers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *brokerClient) Get(name string, opts metav1.GetOptions) (*v1.Broker, error) {
	result := v1.Broker{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("brokers").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *brokerClient) Create(project *v1.Broker) (*v1.Broker, error) {
	result := v1.Broker{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("brokers").
		Body(project).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *brokerClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.Get().Namespace(c.ns).Resource("brokers").VersionedParams(&opts, scheme.ParameterCodec).Watch(context.TODO())
}

func (c *brokerClient) Patch(name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (result *v1.Broker, err error) {
	result = &v1.Broker{}
	err = c.restClient.Patch(pt).
		Namespace(c.ns).
		Resource("brokers").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(context.TODO()).
		Into(result)
	return
}

// Delete takes name of the broker and deletes it. Returns an error if one occurs.
func (c *brokerClient) Delete(name string, opts metav1.DeleteOptions) error {
	return c.restClient.Delete().
		Namespace(c.ns).
		Resource("brokers").
		Name(name).
		Body(&opts).
		Do(context.TODO()).
		Error()
}

func (c *brokerClient) GetPath(name string, opts metav1.GetOptions, path string) (string, error) {

	result, err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("brokers").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).DoRaw(context.TODO())

	accessToken, parseErr := parseJSONPath(result, "token-key", path)
	if parseErr != nil {
		fmt.Println(fmt.Errorf("error parsing %q from %q: %v", path, string(result), parseErr))
	}

	fmt.Println(accessToken)

	return string(result), err
}

func parseJSONPath(input interface{}, name, template string) (string, error) {
	j := jsonpath.New(name)
	buf := new(bytes.Buffer)
	if err := j.Parse(template); err != nil {
		return "", err
	}
	if err := j.Execute(buf, input); err != nil {
		return "", err
	}
	return buf.String(), nil
}
