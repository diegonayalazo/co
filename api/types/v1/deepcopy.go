package v1

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *Broker) DeepCopyInto(out *Broker) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Broker) DeepCopyObject() runtime.Object {
	out := Broker{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *BrokerList) DeepCopyObject() runtime.Object {
	out := BrokerList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]Broker, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *Trigger) DeepCopyInto(out *Trigger) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Trigger) DeepCopyObject() runtime.Object {
	out := Trigger{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *TriggerList) DeepCopyObject() runtime.Object {
	out := TriggerList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]Trigger, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
