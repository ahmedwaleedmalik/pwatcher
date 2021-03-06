package controllers

import (
	"context"
	"strings"
	"time"

	"github.com/ahmedwaleedmalik/pwatcher/pkg/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	kubeNamespacePrefix string = "kube-"
)

// addTimestampAnnotation adds the Timestamp annotation with the value set to the current time
func addTimestampAnnotation(annotations map[string]string) map[string]string {
	if annotations == nil {
		annotations = make(map[string]string)
	}

	// Add current time stamp
	annotations[TimestampAnnotation] = time.Now().UTC().String()

	return annotations
}

// onlyAllowCreateActionPredicate will suppress Delete, Update and Generic events
func onlyAllowCreateActionPredicate() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return false
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
	}
}

// filterCreatePredicate will apply filters based on pod and namespace annotations
func filterCreatePredicate(client client.Client) predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			pod := e.Object.(*corev1.Pod)

			// Ignore pods from system namespaces
			if isProhibitedNamespace(pod.ObjectMeta.Namespace) {
				return false
			}

			// Drop pre-existing pods, so that they are not queued for reconciliation
			if len(pod.Status.PodIP) != 0 {
				return false
			}

			// Check for filters
			return isObservableNamespace(client, pod) && isObservablePod(pod)
		},
	}
}

// isObservableNamespace ensures that namespace that contains the pod has required annotation, if namespace filter key is defined
func isObservableNamespace(client client.Client, pod *corev1.Pod) bool {
	// Check if Namespace Filter Key exists
	if len(config.NamespaceFilterKey) != 0 {
		// Retrieve namespace
		namespace := &corev1.Namespace{}
		err := client.Get(context.TODO(), types.NamespacedName{Name: pod.ObjectMeta.Namespace}, namespace)
		if err != nil {
			return false
		}

		// Consider resource for reconciliation only if the namespace in which it exists has the required annotation
		if _, ok := namespace.GetAnnotations()[config.NamespaceFilterKey]; ok {
			return true
		}
		return false
	}
	return true
}

// isObservablePod ensures that pod has required annotation, if pod filter key is defined
func isObservablePod(pod *corev1.Pod) bool {
	// Consider resource for reconciliation only if it has the required annotation
	if len(config.PodFilterKey) != 0 {
		if _, ok := pod.GetAnnotations()[config.PodFilterKey]; ok {
			return true
		}
		return false
	}
	return true
}

func isProhibitedNamespace(namespace string) bool {
	return strings.HasPrefix(namespace, kubeNamespacePrefix)
}
