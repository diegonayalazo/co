package main

import (
	"flag"
	"fmt"
	"log"

	v1 "github.com/diegonayalazo/co/api/types/v1"
	clientV1 "github.com/diegonayalazo/co/clientset/v1"
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
	namespace := "default"
	brokerName := "conformance-broker"
	triggerName := "conformance-trigger"
	uri := "http://events-counter-service.default.svc.cluster.local/events"
	listBroker(*clientSet, namespace)
	createBroker(*clientSet, namespace, brokerName)
	//kubectl patch broker conformance-broker --type merge -p '{"metadata":{"annotations":{"eventing.knative.dev/broker.class":"mutable"}}}'

	patchBroker(*clientSet, []byte(`{"metadata":{"annotations":{"eventing.knative.dev/broker.class":"mutable"}}}`), namespace, brokerName)

	//kubectl patch broker conformance-broker --type merge -p '{"spec":{"config":{"apiVersion":"v1"}}}'

	patchBroker(*clientSet, []byte(`{"spec":{"config":{"apiVersion":"v1"}}}`), namespace, brokerName)

	// kubectl get broker conformance-broker -ojsonpath="{.status.conditions[?(@.type == \"Ready\")].status}"

	getBroker(*clientSet, namespace, brokerName)

	//kubectl get broker conformance-broker -ojsonpath="{.status.address.url}"

	//kubectl apply -f control-plane/broker-lifecycle/trigger.yaml
	createTrigger(*clientSet, namespace, triggerName, brokerName, uri)

	//kubectl get trigger conformance-trigger -ojsonpath="{.spec.broker}"
	getTrigger(*clientSet, namespace, triggerName)
	//kubectl get trigger conformance-trigger -ojsonpath="{.status.conditions[?(@.type == \"Ready\")].status}"

	//cleanup
	deleteBroker(*clientSet, namespace, brokerName)
}
