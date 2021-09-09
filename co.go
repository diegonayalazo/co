package main

import (
	"fmt"

	v1 "github.com/diegonayalazo/co/api/types/v1"
	clientV1 "github.com/diegonayalazo/co/clientset/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func listBroker(clientSet clientV1.ExampleV1Client) {
	brokers, err := clientSet.Brokers("default").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("brokers found: %+v\n", brokers)
}

func createBroker(clientSet clientV1.ExampleV1Client) {

	NewBroker := &v1.Broker{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "broker",
			Labels: map[string]string{"mylabel": "test"},
		},
	}
	fmt.Println("creating brokers")
	resp, err := clientSet.Brokers("default").Create(NewBroker)

	if err != nil {
		fmt.Printf("error while creating broker: %v\n", err)
	} else {
		fmt.Printf("object created: %v\n", resp)
	}

}
func patchBroker(clientSet clientV1.ExampleV1Client) {
	//kubectl patch broker conformance-broker --type merge -p '{"metadata":{"annotations":{"eventing.knative.dev/broker.class":"mutable"}}}'
	patch := []byte(`{"metadata":{"annotations":{"eventing.knative.dev/broker.class":"mutable"}}}`)

	/*	payload := []patchUInt32Value{{
			Op:    "replace",
			Path:  "/spec/replicas",
			Value: scale,
		}}

		payloadBytes, _ := json.Marshal(payload)
	*/
	fmt.Println("patching brokers")
	resp, err := clientSet.Brokers("default").Patch("conformance-broker", types.MergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		fmt.Printf("error while patching broker: %v\n", err)
	} else {
		fmt.Printf("object patched: %v\n", resp)
	}
}

//  patchStringValue specifies a patch operation for a string.
type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

//  patchUint32Value specifies a patch operation for a uint32.
type patchUInt32Value struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value uint32 `json:"value"`
}
