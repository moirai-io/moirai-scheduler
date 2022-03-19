package scheduler

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// GetQueueBinding returns the queue binding of the specified pod
func GetQueueBinding(pod *corev1.Pod) string {
	return ""
}

// SetAnnotation sets the annotation to the specified pod
func SetAnnotation(ctx context.Context, pod *corev1.Pod) {

}

// GetNodeAvaliableResource returns the available resource of the node
func GetNodeAvaliableResource() *framework.Resource {
	return nil
}
