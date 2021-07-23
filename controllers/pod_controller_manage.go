package controllers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	PodTimestampLogFormat string = "\nPod %v/%v - Timestamp %v"
)

func (r *PodReconciler) annotateResource(instance *corev1.Pod) (ctrl.Result, error) {
	log := log.FromContext(context.Background())

	// Base object used for patching
	baseToPatch := client.MergeFrom(instance.DeepCopy())

	// Update annotations and add Timestamp annotation to pod
	instance.ObjectMeta.Annotations = addTimestampAnnotation(instance.ObjectMeta.Annotations)

	// Patch the object using the merge-patch strategy
	err := r.Patch(context.TODO(), instance, baseToPatch)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Log the Pod and Timestamp
	log.Info(fmt.Sprintf(PodTimestampLogFormat, instance.ObjectMeta.Namespace, instance.ObjectMeta.Name, instance.Annotations[TimestampAnnotation]))
	return ctrl.Result{}, nil
}
