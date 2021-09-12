package v1

import (
	"context"

	v1 "github.com/diegonayalazo/co/api/types/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type TriggerInterface interface {
	List(opts metav1.ListOptions) (*v1.TriggerList, error)
	Get(name string, options metav1.GetOptions) (*v1.Trigger, error)
	GetPath(name string, options metav1.GetOptions, path string) (string, error)
	Create(*v1.Trigger) (*v1.Trigger, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (result *v1.Trigger, err error)
	Delete(name string, opts metav1.DeleteOptions) error
}
type triggerClient struct {
	restClient rest.Interface
	ns         string
}

func (c *triggerClient) List(opts metav1.ListOptions) (*v1.TriggerList, error) {
	result := v1.TriggerList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("triggers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *triggerClient) Get(name string, opts metav1.GetOptions) (*v1.Trigger, error) {
	result := v1.Trigger{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("triggers").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *triggerClient) GetPath(name string, opts metav1.GetOptions, path string) (string, error) {

	result, err := c.restClient.
		Get().
		AbsPath(path).
		Namespace(c.ns).
		Resource("triggers").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).DoRaw(context.TODO())

	return string(result), err
}

func (c *triggerClient) Create(project *v1.Trigger) (*v1.Trigger, error) {
	result := v1.Trigger{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("triggers").
		Body(project).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *triggerClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.Get().Namespace(c.ns).Resource("triggers").VersionedParams(&opts, scheme.ParameterCodec).Watch(context.TODO())
}

// Delete takes name of the broker and deletes it. Returns an error if one occurs.
func (c *triggerClient) Delete(name string, opts metav1.DeleteOptions) error {
	return c.restClient.Delete().
		Namespace(c.ns).
		Resource("triggers").
		Name(name).
		Body(&opts).
		Do(context.TODO()).
		Error()
}

func (c *triggerClient) Patch(name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (result *v1.Trigger, err error) {
	result = &v1.Trigger{}
	err = c.restClient.Patch(pt).
		Namespace(c.ns).
		Resource("triggers").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(context.TODO()).
		Into(result)
	return
}
