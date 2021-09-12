package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Trigger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TriggerSpec `json:"spec"`
}

type TriggerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Trigger `json:"items"`
}

type TriggerSpec struct {
	Broker            string     `json:"broker"`
	TriggerSubscriber TriggerSub `json:"subscriber"`
}

type TriggerSub struct {
	Uri string `json:"uri"`
}
