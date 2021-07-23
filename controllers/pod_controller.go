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

//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;watch;patch
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;patch

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
	// Add TimestampAnnotation to pod if doesn't already exist
	if _, ok := instance.ObjectMeta.Annotations[TimestampAnnotation]; !ok {
		// Annotation the pod with Timestamp
		return r.annotatePodWithTimestamp(instance)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		WithEventFilter(ignoreDeletePredicate()).
		WithEventFilter(ignoreUpdatePredicate()).
		WithEventFilter(ignoreGenericPredicate()).
		WithEventFilter(filterCreatePredicate(mgr.GetClient())).
		Complete(r)
}
