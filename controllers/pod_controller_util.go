package controllers

import "time"

// AddTimestampAnnotation adds the Timestamp annotation with the value set to the current time
func AddTimestampAnnotation(annotations map[string]string) map[string]string {
	if annotations == nil {
		annotations = make(map[string]string)
	}

	// Add current time stamp
	annotations[TimestampAnnotation] = time.Now().UTC().String()

	return annotations
}
