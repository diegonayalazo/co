package main

import (
	"fmt"

	v1 "github.com/diegonayalazo/co/api/types/v1"
	clientV1 "github.com/diegonayalazo/co/clientset/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func listBroker(clientSet clientV1.ExampleV1Client, namespace string) {
	brokers, err := clientSet.Brokers(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("brokers in namespace : %+v found: %+v\n", namespace, brokers)
}

func createBroker(clientSet clientV1.ExampleV1Client, namespace string, brokerName string) {

	NewBroker := &v1.Broker{
		ObjectMeta: metav1.ObjectMeta{
			Name:   brokerName,
			Labels: map[string]string{"mylabel": "test"},
		},
	}
	fmt.Println("creating brokers")
	//en namespace default
	resp, err := clientSet.Brokers(namespace).Create(NewBroker)

	if err != nil {
		fmt.Printf("error while creating broker: %v\n", err)
	} else {
		fmt.Printf("object created: %v\n", resp)
	}

}
func patchBroker(clientSet clientV1.ExampleV1Client, patch []byte, namespace string, brokerName string) {

	fmt.Printf("patch with %q", patch)
	resp, err := clientSet.Brokers(namespace).Patch(brokerName, types.MergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		fmt.Printf("error while patching broker: %v\n", err)
	} else {
		fmt.Printf("object patched: %v\n", resp)
	}
}

func deleteBroker(clientSet clientV1.ExampleV1Client, namespace string, brokerName string) {
	deleteOptions := metav1.DeleteOptions{}
	fmt.Println("deleting brokers")
	//en namespace default
	err := clientSet.Brokers(namespace).Delete(brokerName, deleteOptions)

	if err != nil {
		fmt.Printf("error while deleting broker: %v\n", err)
	}

}
