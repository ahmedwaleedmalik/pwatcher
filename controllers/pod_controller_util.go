package controllers

import (
	"time"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
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

// ignoreDelete will suppress Delete events
func ignoreDelete() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
	}
}

// ignoreUpdate will suppress Update events
func ignoreUpdate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return false
		},
	}
}
