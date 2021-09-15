package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"

	v1 "github.com/diegonayalazo/co/api/types/v1"
	clientV1 "github.com/diegonayalazo/co/clientset/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/jsonpath"
)

var kubeconfig string

func init() {
	fmt.Println("conformance officer.")
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
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

func main() {

	var data = []byte(`{
		"kind": "List",
		"items":[
		  {
			"kind":"None",
			"metadata":{
			  "name":"127.0.0.1",
			  "labels":{
				"kubernetes.io/hostname":"127.0.0.1"
			  }
			},
			"status":{
			  "capacity":{"cpu":"4"},
			  "ready": true,
			  "addresses":[{"type": "LegacyHostIP", "address":"127.0.0.1"}]
			}
		  },
		  {
			"kind":"None",
			"metadata":{
			  "name":"127.0.0.2",
			  "labels":{
				"kubernetes.io/hostname":"127.0.0.2"
			  }
			},
			"status":{
			  "capacity":{"cpu":"8"},
			  "ready": false,
			  "addresses":[
				{"type": "LegacyHostIP", "address":"127.0.0.2"},
				{"type": "another", "address":"127.0.0.3"}
			  ]
			}
		  }
		],
		"users":[
		  {
			"name": "myself",
			"user": {}
		  },
		  {
			"name": "e2e",
			"user": {"username": "admin", "password": "secret"}
			}
		]
	  }`)

	accessToken, parseErr := parseJSONPath(data, "token-key", "{.status.address.url}")
	if parseErr != nil {
		fmt.Println(fmt.Errorf("error parsing token-key %q from %q: %v", "{.status.address.url}", string(data), parseErr))
		return
	}

	fmt.Println(accessToken)

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
	//listBroker(*clientSet, namespace)
	createBroker(*clientSet, namespace, brokerName)
	//kubectl get broker conformance-broker -o jsonpath='{.metadata.annotations.eventing\.knative\.dev/broker\.class}'
	getBrokerPath(*clientSet, namespace, brokerName, "{.metadata.annotations.eventing\\.knative\\.dev/broker\\.class}")

	//kubectl patch broker conformance-broker --type merge -p '{"metadata":{"annotations":{"eventing.knative.dev/broker.class":"mutable"}}}'

	patchBroker(*clientSet, []byte(`{"metadata":{"annotations":{"eventing.knative.dev/broker.class":"mutable"}}}`), namespace, brokerName)

	//kubectl patch broker conformance-broker --type merge -p '{"spec":{"config":{"apiVersion":"v1"}}}'

	patchBroker(*clientSet, []byte(`{"spec":{"config":{"apiVersion":"v1"}}}`), namespace, brokerName)

	// kubectl get broker conformance-broker -ojsonpath="{.status.conditions[?(@.type == \"Ready\")].status}"

	getBrokerPath(*clientSet, namespace, brokerName, "{.status.conditions[?(@.type == \"Ready\")].status}")

	//kubectl get broker conformance-broker -ojsonpath="{.status.address.url}"
	getBrokerPath(*clientSet, namespace, brokerName, "{.status.address.url}")
	//kubectl apply -f control-plane/broker-lifecycle/trigger.yaml
	createTrigger(*clientSet, namespace, triggerName, brokerName, uri)

	//kubectl get trigger conformance-trigger -ojsonpath="{.spec.broker}"
	getTriggerPath(*clientSet, namespace, triggerName, "{.spec.broker}")
	//kubectl get trigger conformance-trigger -ojsonpath="{.status.conditions[?(@.type == \"Ready\")].status}"
	getTriggerPath(*clientSet, namespace, triggerName, "{.status.conditions[?(@.type == \"Ready\")].status}")
	//cleanup

	deleteTrigger(*clientSet, namespace, triggerName)
	deleteBroker(*clientSet, namespace, brokerName)

}
