package main

import (
	"flag"
	"fmt"
	"log"

	v1 "github.com/diegonayalazo/co/api/types/v1"
	clientV1 "github.com/diegonayalazo/co/clientset/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

func init() {
	fmt.Println("conformance officer.")
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {

	var config *rest.Config
	var err error

	log.Printf("using configuration from '%s'", kubeconfig)
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		panic(err)
	}

	v1.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	projects, err := clientSet.Brokers("default").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("brokers found: %+v\n", projects)

	NewBroker := &v1.Broker{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "sslconfigobj",
			Labels: map[string]string{"mylabel": "test"},
		},
	}
	resp, err := clientSet.Brokers("default").Create(NewBroker)

	if err != nil {
		fmt.Printf("error while creating object: %v\n", err)
	} else {
		fmt.Printf("object created: %v\n", resp)
	}

}
