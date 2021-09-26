package v1

import (
	"bytes"
	"context"
	"encoding/json"
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

	fmt.Printf("raw json data:\n %q\n", string(result))
	var data interface{}
	json.Unmarshal(result, &data)

	fmt.Printf("unmarshalled data:\n %q\n", data)

	fmt.Printf("template:\n %q\n", path)

	j := jsonpath.New(name)
	j.AllowMissingKeys(true)
	buf := new(bytes.Buffer)
	if err := j.Parse(path); err != nil {
		fmt.Printf("\nerror parsing template: %q\n", path)
		return "", err
	}
	if err := j.Execute(buf, data); err != nil {
		fmt.Printf("\nerror executing input: %q\n", data)
		return "", err
	}

	fmt.Printf("\njson-path returned: \n %q \n", buf.String())

	return buf.String(), err
}
