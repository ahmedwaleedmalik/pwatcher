# pwatcher

Example Kubernetes Controller built with the `kubebuilder` to showcase different use-cases in a controller.

`pwatcher` is a Kubernetes controller which:

- Listens for new pods
- Annotates those pods with a timestamp
- Logs the pod and timestamp to stdout
- Only respond to pods with a particular annotation (e.g. observable=true)
- Only respond to pods in namespaces with a particular annotation (e.g. observable=true)
- Implement leader election

## Use Cases

### Listens for new pods

To listen for only new pods, we used [predicate](https://stuartleeks.com/posts/kubebuilder-event-filters-part-1-delete/) to supress `Update, Delete, Generic` events.

This filtering only leaves us with `Create` event. Requests against this event will now be passed on to the `reconciler`.

### Annotates those pods with a timestamp

Annotate pod with current timestamp if the annotation doesn't already exist

### Logs the pod and timestamp to stdout

Log the namespaced name of the pod along with the timestamp.

### Only respond to pods with a particular annotation (e.g. observable=true)

To facilitate this, we added a new environment variable `POD_FILTER_KEY`. By using this environment variable we can specify the `key` for filtering. Pod should have the given key in it's annotations.

[predicate](https://stuartleeks.com/posts/kubebuilder-event-filters-part-1-delete/) was used to add this filter. We added the logic to filter out pods against `Create` event in the `predicate`.

### Only respond to pods in namespaces with a particular annotation (e.g. observable=true)

To facilitate this, we added a new environment variable `NAMESPACE_FILTER_KEY`. By using this environment variable we can specify the `key` for filtering. Namespace containing the pod should have the given key in it's annotations.

[predicate](https://stuartleeks.com/posts/kubebuilder-event-filters-part-1-delete/) was used to add this filter. We added the logic to drop namespace, containing the pod, that doesn't have the required key in annotations. `Create` event in the `predicate` was used for this.

## Installation

### Helm Chart

For `helm` view [chart](charts/pwatcher/README.md)

### Kustomize

To install using `kustomize`

- Ensure that [kustomize](https://kubernetes-sigs.github.io/kustomize/installation/) is installed
- Run `VERSION=v0.0.x make deploy` from the root directory to deploy controller with an image tag `v0.0.x`

## Local Development

### Requirements

- golang 1.16
- [kubebuilder](https://book-v1.book.kubebuilder.io/getting_started/installation_and_setup.html)
- [kustomize](https://kubernetes-sigs.github.io/kustomize/installation/)

### Execution

- Ensure that a valid kubeconfig is loaded
- Run `make run` from the root directory to build and execute the binary
