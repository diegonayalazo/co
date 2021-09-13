# Testing Knative Eventing Broker lifecycle

This client-go application tests the conformance of the cluster to the broker lifecycle as specified in https://github.com/diegonayalazo/specs/blob/main/specs/eventing/test-plan/broker-lifecycle-conformance.md

## Usage

First of all, start your KNative implementation (using KonK for example). Once your cluster is ready you can execute:

./co --kubeconfig=path-to-your-kube-config file

This will create the needed objects into the cluster using the client-go client and check for the correct results.

This is my first Golang and Client-go app so any feedback is appreciated. I assume there is a huge room for improvement.

Thanks to @salaboy for helping me out!



