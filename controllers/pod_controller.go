/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PodReconciler reconciles a Pod object
type PodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	TimestampAnnotation string = "pwatcher.io/timestamp"
)

//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=pods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("pwatcher-controller", req.NamespacedName)

	// Fetch the Cluster instance
	instance := &corev1.Pod{}

	// Retrieve Pod Object
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Object no longer exists")

			// Do not requeue object for reconciliation
			return ctrl.Result{}, nil
		}

		// Error reading the object
		// Requeue the Object for reconciliation with an error
		return ctrl.Result{}, err
	}

	// Object has been retrieved successfully at this point

	// TODO: Maybe we can move this logic into a separate method; minimal controller
	// Add TimestampAnnotation to pod if doesn't already exist
	if _, ok := instance.ObjectMeta.Annotations[TimestampAnnotation]; !ok {
		// Annotation doesn't exist
		// Base object for patch that patches using the merge-patch strategy with the given object as base.
		baseToPatch := client.MergeFrom(instance.DeepCopy())

		// Update annotations
		instance.ObjectMeta.Annotations = addTimestampAnnotation(instance.ObjectMeta.Annotations)

		// Patch the object
		err := r.Client.Patch(context.TODO(), instance, baseToPatch)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// Log the Pod and Timestamp
	// TODO: Maybe we can just print out pod-namespace/pod-name and timestamp
	log.Info(fmt.Sprintf("\nPod %v/%v - Timestamp %v", instance.ObjectMeta.Namespace, instance.ObjectMeta.Name, instance.Annotations[TimestampAnnotation]))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		WithEventFilter(ignoreDelete()).
		WithEventFilter(ignoreUpdate()).
		Complete(r)
}
